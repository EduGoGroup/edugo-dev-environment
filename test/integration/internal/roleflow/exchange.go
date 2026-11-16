//go:build integration

package roleflow

// Plan 053 T3.3: edugo-api-identity ya no tiene POST /auth/login, /auth/refresh
// ni /auth/logout. La sesión nace en POST /auth/exchange: identity-core firma un
// Identity Token y identity lo canjea por su Context Token ES256.
//
// Este archivo replica, en mínimo, el banco de los tests de integración de
// identity (internal/testutil/identitytest + test/integration/testutil), que
// este módulo no puede importar por ser `internal/`:
//
//   - un identity-core falso que sirve un JWKS con una clave P-256 generada en
//     memoria y acuña Identity Tokens con identity-shared (el mismo emisor que
//     usa identity-core);
//   - un keyring ES256 del Context Token para identity, generado en memoria;
//   - AUTH_CONTEXT_PUBLIC_KEYS exportada al proceso con la clave pública de ese
//     keyring, para que academic y platform (que la leen de os.LookupEnv en su
//     NewContainer) acepten el Context Token.
//
// Ninguna clave se escribe a disco.

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	identjwt "github.com/EduGoGroup/identity-shared/auth/jwt"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	identityBuilder "github.com/edugo/edugo-api-identity/cmd/builder"
	identityConfig "github.com/edugo/edugo-api-identity/cmd/config"
)

// Literales del contrato con identity-core que el exchange valida. Son copia de
// edugo-api-identity/internal/infra/identitycore/literals.go (no importable).
const (
	identityCoreIssuer = "identity-core"
	identityJWKSPath   = "/.well-known/jwks.json"
	// SystemKMP es el system (y aud) del Identity Token de los clientes KMP.
	SystemKMP = "edu2go.kmp"

	identityCoreKID = "identity-core-test-kid"
	contextKID      = "context-test-kid"

	// identityTokenTTL: el exchange rechaza un Identity Token con menos de un
	// minuto de vida; una hora cubre de sobra cualquier suite.
	identityTokenTTL = time.Hour
)

// Bench es el identity-core falso + la BD de un identity server de test.
type Bench struct {
	identity *httptest.Server
	key      *ecdsa.PrivateKey
	minter   *identjwt.Manager
	db       *gorm.DB
}

var (
	benchesMu sync.Mutex
	benches   = map[string]*Bench{}
)

// benchFor devuelve el Bench del identity server `server` (registrado por
// StartIdentityServer). Falla el test si el server no se levantó con él.
func benchFor(t *testing.T, server *httptest.Server) *Bench {
	t.Helper()
	benchesMu.Lock()
	defer benchesMu.Unlock()
	b, ok := benches[server.URL]
	require.Truef(t, ok, "roleflow: el identity server %s no se levantó con StartIdentityServer", server.URL)
	return b
}

// StartIdentityServer levanta edugo-api-identity in-process sobre db con el
// exchange encendido contra un identity-core falso y un keyring ES256 propio.
//
// EFECTO EN EL PROCESO: exporta AUTH_CONTEXT_PUBLIC_KEYS con la clave pública
// del keyring. Llamarlo ANTES de construir academic/platform, que leen esa
// variable en su NewContainer; si se construyen antes, rechazan el Context
// Token con 401.
func StartIdentityServer(db *gorm.DB, jwtSecret, jwtIssuer string) (*httptest.Server, error) {
	bench, err := newBench(db)
	if err != nil {
		return nil, err
	}

	signing, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("context key: %w", err)
	}
	keyring, err := auth.NewContextKeyring(&auth.ContextSigningKey{ID: contextKID, Private: signing}, nil)
	if err != nil {
		return nil, fmt.Errorf("context keyring: %w", err)
	}
	if err := exportContextPublicKey(&signing.PublicKey); err != nil {
		return nil, err
	}

	cfg := &identityConfig.Config{
		Environment: "development",
		Server: identityConfig.ServerConfig{
			Port: 0,
			Host: "127.0.0.1",
		},
		Auth: identityConfig.AuthConfig{
			JWT: identityConfig.JWTConfig{
				Secret:               jwtSecret,
				Issuer:               jwtIssuer,
				AccessTokenDuration:  accessTokenDuration,
				RefreshTokenDuration: refreshTokenDuration,
			},
			Exchange: identityConfig.ExchangeConfig{Enabled: true, IdentityURL: bench.identity.URL},
		},
		Logging: identityConfig.LoggingConfig{Level: "error", Format: "text"},
		CORS: identityConfig.CORSConfig{
			AllowedOrigins: "*",
			AllowedMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowedHeaders: "Content-Type,Authorization,X-Request-ID",
		},
	}
	if err := identityBuilder.ValidateExchangeSetup(cfg.Auth.Exchange, keyring); err != nil {
		return nil, err
	}

	blacklist := auth.NewInMemoryBlacklist(context.Background())
	log := newNoOpLogger()
	c := identityBuilder.NewContainer(db, log, cfg, blacklist, identityBuilder.WithContextKeyring(keyring))
	server := httptest.NewServer(c.SetupRouter(cfg, log, "dev", "dev"))

	benchesMu.Lock()
	benches[server.URL] = bench
	benchesMu.Unlock()
	return server, nil
}

func newBench(db *gorm.DB) (*Bench, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("identity-core key: %w", err)
	}
	minter, err := identjwt.NewManager(key, identityCoreIssuer, identityCoreKID)
	if err != nil {
		return nil, fmt.Errorf("identity-core minter: %w", err)
	}
	b := &Bench{key: key, minter: minter, db: db}
	mux := http.NewServeMux()
	mux.HandleFunc(identityJWKSPath, b.serveJWKS)
	// Vive hasta el fin del proceso, como los servers de TestMain.
	b.identity = httptest.NewServer(mux)
	return b, nil
}

func (b *Bench) serveJWKS(w http.ResponseWriter, _ *http.Request) {
	enc := base64.RawURLEncoding
	// Punto sin comprimir: 0x04 || X (32 bytes) || Y (32 bytes).
	point, err := b.key.PublicKey.Bytes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	doc := map[string]any{"keys": []map[string]string{{
		"kty": "EC", "crv": "P-256", "kid": identityCoreKID, "use": "sig", "alg": "ES256",
		"x": enc.EncodeToString(point[1:33]),
		"y": enc.EncodeToString(point[33:65]),
	}}}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(doc)
}

func exportContextPublicKey(pub *ecdsa.PublicKey) error {
	pemKey, err := auth.MarshalECPublicKeyPEM(pub)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(map[string]string{contextKID: pemKey})
	if err != nil {
		return err
	}
	return os.Setenv(auth.EnvContextPublicKeys, string(raw))
}

// MintIdentityToken acuña un Identity Token de identity-core para el usuario
// con ese email (sub = su UUID en auth.users) y el system KMP.
func MintIdentityToken(t *testing.T, server *httptest.Server, email string) string {
	t.Helper()
	b := benchFor(t, server)
	token, _, err := b.minter.GenerateIdentityToken(identjwt.IdentityTokenInput{
		UserID: UserIDByEmail(t, b.db, email),
		System: SystemKMP,
		Email:  email,
		TTL:    identityTokenTTL,
	})
	require.NoError(t, err, "roleflow: mint identity token")
	return token
}

// UserIDByEmail resuelve el UUID de auth.users para email. Falla el test si no
// existe.
func UserIDByEmail(t *testing.T, db *gorm.DB, email string) string {
	t.Helper()
	var id string
	require.NoError(t, db.Raw(
		"SELECT id::text FROM auth.users WHERE email = ? AND deleted_at IS NULL", email,
	).Scan(&id).Error, "roleflow: lookup user %s", email)
	require.NotEmptyf(t, id, "roleflow: no existe el usuario %s", email)
	return id
}

// DesiredContext es el `desired_context` opcional del exchange: el contexto que
// el cliente pide reconstruir al renovar la sesión (re-exchange).
type DesiredContext struct {
	SchoolID       string `json:"school_id,omitempty"`
	AcademicUnitID string `json:"academic_unit_id,omitempty"`
	RoleID         string `json:"role_id,omitempty"`
	Subject        string `json:"subject,omitempty"`
}

// Exchange ejecuta POST /api/v1/auth/exchange con un Identity Token recién
// acuñado para email y, si no es nil, el desired_context. Devuelve (status,
// body) crudos.
func Exchange(t *testing.T, server *httptest.Server, email string, desired *DesiredContext) (int, []byte) {
	t.Helper()
	payload := map[string]any{"identity_token": MintIdentityToken(t, server, email)}
	if desired != nil {
		payload["desired_context"] = desired
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/api/v1/auth/exchange", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, raw
}
