// Command provision_users da de alta por API los usuarios de un escenario
// on-top de playground_v2 (plan 053 T5.2, paso 3). Ver internal/provision.
//
// Variables: EDUGO_IDENTITY_URL (por defecto http://localhost:8070),
// IDENTITY_CORE_BASE_URL (obligatoria), SCENARIO_TEST_PASSWORD (obligatoria,
// solo del entorno) y las POSTGRES_* del migrator para leer auth.users.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/provision"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
)

const defaultIdentityURL = "http://localhost:8070"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout))
}

func run(args []string, stdin *os.File, out io.Writer) int {
	logger := log.New(out, "", log.LstdFlags)
	fs := flag.NewFlagSet("provision_users", flag.ContinueOnError)
	scenario := fs.String("scenario", "", fmt.Sprintf("escenario on-top (%s)", strings.Join(playground_v2.OnTopScenarios(), "|")))
	yes := fs.Bool("yes", false, "no pedir confirmación de los destinos")
	interval := fs.Duration("interval", 2*time.Second, "pausa entre llamadas al signup de identity-core (rate limit por origen)")
	retryWait := fs.Duration("retry-wait", 30*time.Second, "espera tras un 429 que no trae Retry-After")
	maxRetries := fs.Int("max-retries", 5, "reintentos ante 429")
	timeout := fs.Duration("timeout", 15*time.Second, "timeout de cada petición HTTP")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *scenario == "" {
		logger.Printf("❌ Falta -scenario=<nombre> (on-top: %v)", playground_v2.OnTopScenarios())
		return 2
	}
	users, err := playground_v2.ScenarioUsers(*scenario)
	if err != nil {
		logger.Printf("❌ %v", err)
		return 2
	}

	identityURL := envOr("EDUGO_IDENTITY_URL", defaultIdentityURL)
	coreURL := os.Getenv("IDENTITY_CORE_BASE_URL")
	if coreURL == "" {
		logger.Printf("❌ Falta IDENTITY_CORE_BASE_URL (sin valor por defecto)")
		return 2
	}
	password := os.Getenv("SCENARIO_TEST_PASSWORD")
	if password == "" {
		logger.Printf("❌ Falta SCENARIO_TEST_PASSWORD (solo del entorno)")
		return 2
	}
	if err := provision.ValidatePassword(password); err != nil {
		logger.Printf("❌ %v", err)
		return 2
	}
	identityHost, err := hostOf(identityURL)
	if err != nil {
		logger.Printf("❌ EDUGO_IDENTITY_URL inválida: %v", err)
		return 2
	}
	coreHost, err := hostOf(coreURL)
	if err != nil {
		logger.Printf("❌ IDENTITY_CORE_BASE_URL inválida: %v", err)
		return 2
	}

	pg := config.Load().Postgres
	logger.Printf("Alta de los %d usuarios del escenario %s", len(users), *scenario)
	logger.Printf("  edugo-api-identity: %s", identityHost)
	logger.Printf("  identity-core:      %s", coreHost)
	logger.Printf("  BD (auth.users):    %s", pg.Target)
	if !*yes {
		if err := confirm(stdin, out, coreHost, pg.Target); err != nil {
			logger.Printf("❌ %v: no se escribió nada", err)
			return 3
		}
	}

	db, err := sql.Open("postgres", pg.ConnStr)
	if err != nil {
		logger.Printf("❌ abriendo la BD: %v", err)
		return 1
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		logger.Printf("❌ no se pudo conectar a la BD %s: %v", pg.Target, err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	p := provision.New(provision.Config{
		IdentityURL: identityURL,
		CoreURL:     coreURL,
		Password:    password,
		Interval:    *interval,
		RetryWait:   *retryWait,
		MaxRetries:  *maxRetries,
	}, nil, provision.SQLLookup{DB: db}, nil, logger.Printf)
	p.SetTimeout(*timeout)

	summary, runErr := p.Run(ctx, users)
	logger.Printf("Resumen %s: creados=%d ya_existían=%d identity-core_201=%d fallidos=%d",
		*scenario, summary.Created, summary.Existed, summary.CoreOK, summary.Failed)
	for _, f := range summary.Failures {
		logger.Printf("  ✗ %s", f)
	}
	if runErr != nil && !errors.Is(runErr, provision.ErrSyncOff) {
		logger.Printf("❌ %v", runErr)
	}
	if !summary.OK() || runErr != nil {
		return 1
	}
	logger.Printf("✅ Usuarios listos. Siguiente paso: make docker-on-top P=%s", *scenario)
	return 0
}

// confirm asks the operator to type the identity-core host and the DB target,
// like scripts/identity/bootstrap.sh does with its hosts.
func confirm(stdin *os.File, out io.Writer, coreHost, dbTarget string) error {
	if fi, err := stdin.Stat(); err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return errors.New("sin terminal interactiva: confirma con -yes")
	}
	r := bufio.NewReader(stdin)
	for _, want := range []struct{ label, value string }{
		{"el host de identity-core", coreHost},
		{"el destino de la BD", dbTarget},
	} {
		fmt.Fprintf(out, "Teclea %s para continuar (%s): ", want.label, want.value)
		line, err := r.ReadString('\n')
		if err != nil && line == "" {
			return fmt.Errorf("leyendo la confirmación: %w", err)
		}
		if strings.TrimSpace(line) != want.value {
			return fmt.Errorf("%s no coincide", want.label)
		}
	}
	return nil
}

func hostOf(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" || u.Host == "" {
		return "", fmt.Errorf("se esperaba http(s)://host[:puerto], llegó %q", raw)
	}
	return u.Host, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
