// Package provision da de alta por API los usuarios de un escenario on-top de
// playground_v2 (plan 053 T5.2, paso 3 del flujo).
//
// Por cada usuario:
//  1. POST {identity}/api/v1/auth/signup de edugo-api-identity. Con el sync
//     encendido, identity pide el UUID a identity-core (ensure) y crea la fila
//     de auth.users con ese id, sin credenciales (auth.users ya no guarda contraseña).
//  2. POST {identity-core}/api/v1/auth/signup de identity-core, que adopta esa
//     cuenta (caso B) con la contraseña de prueba, o la reconoce (caso C) si ya
//     la tenía.
//  3. Comprueba que el id de auth.users coincide con el que devuelve
//     identity-core. Si no coincide, el sync de identity estaba apagado y se
//     aborta: seguir crearía cuentas huérfanas en identity-core.
//
// La contraseña de prueba nunca se escribe en el log.
package provision

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"golang.org/x/text/unicode/norm"
)

// Password policy of identity-core (internal/core/domain/password_policy.go):
// after NFKC, at least 12 runes and at most 72 bytes.
const (
	minPasswordRunes = 12
	maxPasswordBytes = 72
)

const (
	signupPath          = "/api/v1/auth/signup"
	codeEmailExists     = "EMAIL_ALREADY_EXISTS"
	maxResponseBodySize = 64 << 10
)

// ErrSyncOff marks the abort when the local id and the identity-core id differ.
var ErrSyncOff = errors.New("el sync de identity con identity-core estaba apagado")

// ValidatePassword applies identity-core's password policy. The error never
// contains the password.
func ValidatePassword(password string) error {
	normalized := norm.NFKC.String(password)
	if utf8.RuneCountInString(normalized) < minPasswordRunes {
		return fmt.Errorf("SCENARIO_TEST_PASSWORD debe tener al menos %d caracteres (tras NFKC)", minPasswordRunes)
	}
	if len(normalized) > maxPasswordBytes {
		return fmt.Errorf("SCENARIO_TEST_PASSWORD no puede pasar de %d bytes (tras NFKC)", maxPasswordBytes)
	}
	return nil
}

// UserIDLookup resolves the id of a live auth.users row by email.
type UserIDLookup interface {
	UserIDByEmail(ctx context.Context, email string) (id string, found bool, err error)
}

// Config holds the provisioner settings.
type Config struct {
	IdentityURL string        // base URL of edugo-api-identity
	CoreURL     string        // base URL of identity-core
	Password    string        // test password; never logged
	Interval    time.Duration // pause between identity-core calls (rate limit per origin)
	RetryWait   time.Duration // wait after a 429 without Retry-After
	MaxRetries  int           // retries after 429
}

// Summary is the final count of the run.
type Summary struct {
	Created  int      // created in edugo-api-identity
	Existed  int      // already in edugo-api-identity (409 EMAIL_ALREADY_EXISTS)
	CoreOK   int      // identity-core answered 201 (case A, B or C: indistinguishable)
	Failed   int      // users that did not finish
	Failures []string // one line per failed user
	Aborted  bool     // stopped early because the sync was off
}

// OK reports whether every user finished.
func (s Summary) OK() bool { return s.Failed == 0 && !s.Aborted }

// Provisioner runs the three steps for each user.
type Provisioner struct {
	cfg    Config
	client *http.Client
	lookup UserIDLookup
	sleep  func(context.Context, time.Duration) error
	logf   func(format string, args ...any)
}

// New builds a Provisioner. sleep and logf may be nil (real sleep, no log).
func New(cfg Config, client *http.Client, lookup UserIDLookup,
	sleep func(context.Context, time.Duration) error, logf func(string, ...any)) *Provisioner {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	if sleep == nil {
		sleep = SleepContext
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}
	cfg.IdentityURL = strings.TrimRight(cfg.IdentityURL, "/")
	cfg.CoreURL = strings.TrimRight(cfg.CoreURL, "/")
	return &Provisioner{cfg: cfg, client: client, lookup: lookup, sleep: sleep, logf: logf}
}

// SleepContext waits d or until ctx ends.
func SleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

type signupBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type apiResponse struct {
	ID      string `json:"id"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r apiResponse) describe() string {
	parts := make([]string, 0, 3)
	for _, p := range []string{r.Code, r.Error, r.Message} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) == 0 {
		return "sin cuerpo"
	}
	return strings.Join(parts, " · ")
}

// Run provisions users in order. It returns an error only for ErrSyncOff or a
// cancelled context; per-user failures are counted in the Summary.
func (p *Provisioner) Run(ctx context.Context, users []playground_v2.ScenarioUser) (Summary, error) {
	var s Summary
	for i, u := range users {
		p.logf("[%d/%d] %s", i+1, len(users), u.Email)
		if err := p.provisionOne(ctx, u, i > 0, &s); err != nil {
			if errors.Is(err, ErrSyncOff) {
				s.Aborted = true
				p.fail(&s, u.Email, err.Error())
				p.logf("⛔ Abortado: no se sigue con los %d usuarios restantes", len(users)-i-1)
				return s, err
			}
			if ctx.Err() != nil {
				return s, ctx.Err()
			}
			p.fail(&s, u.Email, err.Error())
		}
	}
	return s, nil
}

func (p *Provisioner) fail(s *Summary, email, reason string) {
	s.Failed++
	s.Failures = append(s.Failures, email+": "+reason)
	p.logf("   ❌ %s", reason)
}

func (p *Provisioner) provisionOne(ctx context.Context, u playground_v2.ScenarioUser, pause bool, s *Summary) error {
	body := signupBody{Email: u.Email, Password: p.cfg.Password, FirstName: u.FirstName, LastName: u.LastName}

	// Paso a: alta en edugo-api-identity.
	status, resp, err := p.post(ctx, p.cfg.IdentityURL+signupPath, body)
	if err != nil {
		return fmt.Errorf("signup de edugo-api-identity: %w", err)
	}
	switch {
	case status == http.StatusCreated:
		s.Created++
		p.logf("   ✓ creado en edugo-api-identity (id=%s)", resp.ID)
	case status == http.StatusConflict && resp.Code == codeEmailExists:
		s.Existed++
		p.logf("   = ya existía en edugo-api-identity")
	case status == http.StatusBadGateway:
		// No se llama a identity-core: sin fila local, el caso A crearía allí una cuenta con otro id.
		return fmt.Errorf("edugo-api-identity no pudo sincronizar con identity-core (HTTP %d: %s); no se guardó nada",
			status, resp.describe())
	default:
		return fmt.Errorf("signup de edugo-api-identity respondió HTTP %d (%s)", status, resp.describe())
	}

	// Paso b: adopción en identity-core, espaciada por su rate limit por origen.
	if pause {
		if err := p.sleep(ctx, p.cfg.Interval); err != nil {
			return err
		}
	}
	status, resp, err = p.post(ctx, p.cfg.CoreURL+signupPath, body)
	if err != nil {
		return fmt.Errorf("signup de identity-core: %w", err)
	}
	switch status {
	case http.StatusCreated:
		if resp.ID == "" {
			return fmt.Errorf("identity-core respondió 201 sin id")
		}
		s.CoreOK++
		p.logf("   ✓ identity-core respondió 201 (id=%s): cuenta adoptada o ya reconocida", resp.ID)
	case http.StatusConflict:
		// No se reintenta: el caso D suma un fallo al lockout de la cuenta.
		return fmt.Errorf("identity-core respondió 409 (%s): la cuenta tiene otra contraseña (caso D) o está bloqueada o inactiva",
			resp.describe())
	default:
		return fmt.Errorf("signup de identity-core respondió HTTP %d (%s)", status, resp.describe())
	}

	// Paso c: el id local tiene que ser el de identity-core.
	localID, found, err := p.lookup.UserIDByEmail(ctx, u.Email)
	if err != nil {
		return fmt.Errorf("leyendo auth.users: %w", err)
	}
	if !found {
		return fmt.Errorf("auth.users no tiene una fila viva con ese email tras el signup: ¿identity apunta a otra BD?")
	}
	if !strings.EqualFold(localID, resp.ID) {
		return fmt.Errorf("%w: auth.users tiene id=%s e identity-core id=%s. Arranca edugo-api-identity con "+
			"AUTH_IDENTITY_SYNC_ENABLED=true, recrea la BD con base (make docker-recreate) y repite el alta",
			ErrSyncOff, localID, resp.ID)
	}
	p.logf("   ✓ id coincide en auth.users e identity-core")
	return nil
}

// post sends the JSON body and retries on 429, honouring Retry-After (seconds)
// when present and RetryWait otherwise.
func (p *Provisioner) post(ctx context.Context, url string, body signupBody) (int, apiResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return 0, apiResponse{}, fmt.Errorf("serializando el cuerpo: %w", err)
	}
	for attempt := 0; ; attempt++ {
		status, resp, header, err := p.postOnce(ctx, url, payload)
		if err != nil {
			return 0, apiResponse{}, err
		}
		if status != http.StatusTooManyRequests {
			return status, resp, nil
		}
		if attempt >= p.cfg.MaxRetries {
			return status, resp, fmt.Errorf("HTTP 429 tras %d reintentos", p.cfg.MaxRetries)
		}
		wait := p.cfg.RetryWait
		if ra := header.Get("Retry-After"); ra != "" {
			if secs, convErr := strconv.Atoi(strings.TrimSpace(ra)); convErr == nil && secs >= 0 {
				wait = time.Duration(secs) * time.Second
			}
		}
		p.logf("   … HTTP 429 de %s; reintento %d/%d en %s", hostOf(url), attempt+1, p.cfg.MaxRetries, wait)
		if err := p.sleep(ctx, wait); err != nil {
			return 0, apiResponse{}, err
		}
	}
}

func (p *Provisioner) postOnce(ctx context.Context, url string, payload []byte) (int, apiResponse, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return 0, apiResponse{}, nil, fmt.Errorf("armando la petición: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := p.client.Do(req)
	if err != nil {
		return 0, apiResponse{}, nil, fmt.Errorf("llamando a %s: %w", hostOf(url), err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBodySize))
	if err != nil {
		return 0, apiResponse{}, nil, fmt.Errorf("leyendo la respuesta de %s: %w", hostOf(url), err)
	}
	var resp apiResponse
	_ = json.Unmarshal(raw, &resp) // un cuerpo que no es JSON se informa por su código HTTP
	return res.StatusCode, resp, res.Header, nil
}

func hostOf(rawURL string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(rawURL, "https://"), "http://")
	if i := strings.IndexByte(s, '/'); i >= 0 {
		s = s[:i]
	}
	return s
}

// SetTimeout sets the per-request HTTP timeout.
func (p *Provisioner) SetTimeout(d time.Duration) { p.client.Timeout = d }
