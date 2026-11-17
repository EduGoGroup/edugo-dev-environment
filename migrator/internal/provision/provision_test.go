package provision

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
)

const testPassword = "una-clave-de-prueba-larga"

// fakeLookup simula auth.users: email -> id.
type fakeLookup map[string]string

func (f fakeLookup) UserIDByEmail(_ context.Context, email string) (string, bool, error) {
	id, ok := f[email]
	return id, ok, nil
}

// reply es una respuesta enlatada de un doble HTTP.
type reply struct {
	status     int
	body       string
	retryAfter string
}

// fakeServer devuelve, por email, la cola de respuestas indicada y registra las llamadas.
type fakeServer struct {
	mu      sync.Mutex
	replies map[string][]reply
	calls   map[string]int
	bodies  []map[string]string
}

func newFakeServer(t *testing.T, replies map[string][]reply) (*httptest.Server, *fakeServer) {
	t.Helper()
	fs := &fakeServer{replies: replies, calls: map[string]int{}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/auth/signup" {
			t.Errorf("petición inesperada %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("cuerpo no JSON: %v", err)
		}
		fs.mu.Lock()
		fs.bodies = append(fs.bodies, body)
		email := body["email"]
		n := fs.calls[email]
		fs.calls[email]++
		queue := fs.replies[email]
		fs.mu.Unlock()
		if n >= len(queue) {
			t.Errorf("llamada %d no prevista para %s", n+1, email)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rep := queue[n]
		if rep.retryAfter != "" {
			w.Header().Set("Retry-After", rep.retryAfter)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rep.status)
		_, _ = w.Write([]byte(rep.body))
	}))
	t.Cleanup(srv.Close)
	return srv, fs
}

func created(id string) reply {
	return reply{status: http.StatusCreated, body: fmt.Sprintf(`{"id":%q,"email":"x"}`, id)}
}

var (
	alreadyExists = reply{status: http.StatusConflict, body: `{"error":"conflict","message":"email already exists","code":"EMAIL_ALREADY_EXISTS"}`}
	coreConflict  = reply{status: http.StatusConflict, body: `{"error":"email already registered","code":"CONFLICT"}`}
	tooMany       = reply{status: http.StatusTooManyRequests, body: `{"error":"too many requests","code":"TOO_MANY_REQUESTS"}`}
)

type harness struct {
	p      *Provisioner
	id     *fakeServer
	core   *fakeServer
	sleeps []time.Duration
	log    *strings.Builder
}

func newHarness(t *testing.T, idReplies, coreReplies map[string][]reply, lookup fakeLookup) *harness {
	t.Helper()
	idSrv, idFake := newFakeServer(t, idReplies)
	coreSrv, coreFake := newFakeServer(t, coreReplies)
	h := &harness{id: idFake, core: coreFake, log: &strings.Builder{}}
	h.p = New(Config{
		IdentityURL: idSrv.URL + "/",
		CoreURL:     coreSrv.URL,
		Password:    testPassword,
		Interval:    2 * time.Second,
		RetryWait:   30 * time.Second,
		MaxRetries:  3,
	}, nil, lookup,
		func(_ context.Context, d time.Duration) error { h.sleeps = append(h.sleeps, d); return nil },
		func(format string, args ...any) { fmt.Fprintf(h.log, format+"\n", args...) })
	return h
}

func user(email string) playground_v2.ScenarioUser {
	return playground_v2.ScenarioUser{Email: email, FirstName: "Francisco", LastName: "I"}
}

func (h *harness) assertNoPasswordInLog(t *testing.T) {
	t.Helper()
	if strings.Contains(h.log.String(), testPassword) {
		t.Fatalf("la contraseña aparece en el log:\n%s", h.log.String())
	}
}

func TestRun_CreatedAndAdopted(t *testing.T) {
	const id = "11111111-1111-1111-1111-111111111111"
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {created(id)}},
		map[string][]reply{"a@edugo.local": {created(id)}},
		fakeLookup{"a@edugo.local": id})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local")})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if s.Created != 1 || s.Existed != 0 || s.CoreOK != 1 || !s.OK() {
		t.Fatalf("resumen inesperado: %+v", s)
	}
	// Los dos signup reciben los mismos cuatro campos, incluido el apellido de 1 carácter.
	for _, b := range append(h.id.bodies, h.core.bodies...) {
		if b["email"] != "a@edugo.local" || b["password"] != testPassword || b["first_name"] != "Francisco" || b["last_name"] != "I" {
			t.Errorf("cuerpo inesperado: %v", b)
		}
	}
	if len(h.sleeps) != 0 {
		t.Errorf("el primer usuario no debía esperar: %v", h.sleeps)
	}
	h.assertNoPasswordInLog(t)
}

func TestRun_AlreadyExistedIsNotAFailure(t *testing.T) {
	const id = "22222222-2222-2222-2222-222222222222"
	h := newHarness(t,
		map[string][]reply{"b@edugo.local": {alreadyExists}},
		map[string][]reply{"b@edugo.local": {created(id)}},
		fakeLookup{"b@edugo.local": id})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("b@edugo.local")})
	if err != nil || !s.OK() || s.Existed != 1 || s.Created != 0 || s.CoreOK != 1 {
		t.Fatalf("resumen %+v, err %v", s, err)
	}
	h.assertNoPasswordInLog(t)
}

func TestRun_Core409IsReportedAndRunContinues(t *testing.T) {
	const idA, idB = "33333333-3333-3333-3333-333333333333", "44444444-4444-4444-4444-444444444444"
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {alreadyExists}, "b@edugo.local": {created(idB)}},
		map[string][]reply{"a@edugo.local": {coreConflict}, "b@edugo.local": {created(idB)}},
		fakeLookup{"a@edugo.local": idA, "b@edugo.local": idB})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local"), user("b@edugo.local")})
	if err != nil {
		t.Fatalf("un 409 no aborta: %v", err)
	}
	if s.Failed != 1 || s.CoreOK != 1 || s.OK() {
		t.Fatalf("resumen inesperado: %+v", s)
	}
	if !strings.Contains(s.Failures[0], "a@edugo.local") || !strings.Contains(s.Failures[0], "caso D") {
		t.Errorf("fallo mal descrito: %q", s.Failures[0])
	}
	if h.core.calls["a@edugo.local"] != 1 {
		t.Errorf("un 409 de identity-core no se reintenta: %d llamadas", h.core.calls["a@edugo.local"])
	}
	// Entre el primer y el segundo signup de identity-core hay una pausa de Interval.
	if len(h.sleeps) != 1 || h.sleeps[0] != 2*time.Second {
		t.Errorf("pausas inesperadas: %v", h.sleeps)
	}
	h.assertNoPasswordInLog(t)
}

func TestRun_429RetriesHonouringRetryAfter(t *testing.T) {
	const id = "55555555-5555-5555-5555-555555555555"
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {created(id)}},
		map[string][]reply{"a@edugo.local": {
			{status: http.StatusTooManyRequests, body: tooMany.body, retryAfter: "7"},
			tooMany, // sin Retry-After: usa RetryWait
			created(id),
		}},
		fakeLookup{"a@edugo.local": id})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local")})
	if err != nil || !s.OK() {
		t.Fatalf("resumen %+v, err %v", s, err)
	}
	if h.core.calls["a@edugo.local"] != 3 {
		t.Errorf("llamadas a identity-core = %d, esperado 3", h.core.calls["a@edugo.local"])
	}
	want := []time.Duration{7 * time.Second, 30 * time.Second}
	if fmt.Sprint(h.sleeps) != fmt.Sprint(want) {
		t.Errorf("esperas = %v, esperado %v", h.sleeps, want)
	}
	h.assertNoPasswordInLog(t)
}

func TestRun_429ExhaustsRetries(t *testing.T) {
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {created("x")}},
		map[string][]reply{"a@edugo.local": {tooMany, tooMany, tooMany, tooMany}},
		fakeLookup{})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local")})
	if err != nil {
		t.Fatalf("agotar los reintentos es un fallo del usuario, no un aborto: %v", err)
	}
	if s.Failed != 1 || !strings.Contains(s.Failures[0], "429") {
		t.Fatalf("resumen inesperado: %+v", s)
	}
	if h.core.calls["a@edugo.local"] != 4 {
		t.Errorf("llamadas = %d, esperado 1 + 3 reintentos", h.core.calls["a@edugo.local"])
	}
}

func TestRun_MismatchedIDsAbort(t *testing.T) {
	h := newHarness(t,
		// Sync apagado: identity crea con su propio UUID.
		map[string][]reply{"a@edugo.local": {created("66666666-6666-6666-6666-666666666666")}},
		map[string][]reply{"a@edugo.local": {created("77777777-7777-7777-7777-777777777777")}},
		fakeLookup{"a@edugo.local": "66666666-6666-6666-6666-666666666666"})

	users := []playground_v2.ScenarioUser{user("a@edugo.local"), user("b@edugo.local")}
	s, err := h.p.Run(context.Background(), users)
	if !errors.Is(err, ErrSyncOff) {
		t.Fatalf("esperado ErrSyncOff, llegó %v", err)
	}
	if !s.Aborted || s.OK() || !strings.Contains(s.Failures[0], "AUTH_IDENTITY_SYNC_ENABLED=true") {
		t.Fatalf("resumen inesperado: %+v", s)
	}
	if h.id.calls["b@edugo.local"] != 0 || h.core.calls["b@edugo.local"] != 0 {
		t.Errorf("tras el aborto no se debía tocar a b@edugo.local")
	}
	h.assertNoPasswordInLog(t)
}

func TestRun_IdentitySyncFailureSkipsCore(t *testing.T) {
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {{status: http.StatusBadGateway, body: `{"code":"IDENTITY_SYNC_FAILED"}`}}},
		map[string][]reply{},
		fakeLookup{})

	s, err := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local")})
	if err != nil || s.Failed != 1 || !strings.Contains(s.Failures[0], "IDENTITY_SYNC_FAILED") {
		t.Fatalf("resumen %+v, err %v", s, err)
	}
	if h.core.calls["a@edugo.local"] != 0 {
		t.Errorf("sin fila local no se llama a identity-core")
	}
}

func TestRun_IdentityMismatch409IsAFailure(t *testing.T) {
	h := newHarness(t,
		map[string][]reply{"a@edugo.local": {{status: http.StatusConflict, body: `{"code":"IDENTITY_USER_MISMATCH"}`}}},
		map[string][]reply{},
		fakeLookup{})

	s, _ := h.p.Run(context.Background(), []playground_v2.ScenarioUser{user("a@edugo.local")})
	if s.Failed != 1 || s.Existed != 0 || !strings.Contains(s.Failures[0], "IDENTITY_USER_MISMATCH") {
		t.Fatalf("resumen inesperado: %+v", s)
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name string
		pw   string
		ok   bool
	}{
		{"corta", "once-chars!", false},
		{"doce", "doce-chars!!", true},
		{"72 bytes", strings.Repeat("a", 72), true},
		{"73 bytes", strings.Repeat("a", 73), false},
		// 12 runas de 2 bytes: válida por longitud en caracteres.
		{"multibyte", strings.Repeat("ñ", 12), true},
		// 37 runas de 2 bytes = 74 bytes: rechazada por bytes.
		{"multibyte larga", strings.Repeat("ñ", 37), false},
	}
	for _, c := range cases {
		err := ValidatePassword(c.pw)
		if (err == nil) != c.ok {
			t.Errorf("%s: err=%v, ok esperado %v", c.name, err, c.ok)
		}
		if err != nil && strings.Contains(err.Error(), c.pw) {
			t.Errorf("%s: el error contiene la contraseña", c.name)
		}
	}
}
