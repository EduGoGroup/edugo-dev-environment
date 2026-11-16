//go:build integration

// Package guardian_flow (archivo switch) valida el switch-context POR SUJETO del
// representante (plan 024·F2·T4, ADR 0026): el cliente nunca fija el sujeto a mano;
// pide `subject:"ward:<student_id>"` y el server valida el vínculo guardián↔acudido
// e impersona al hijo emitiendo un contexto con `actor_mode=ward` +
// `subject_student_id`. Casos cubiertos:
//
//   - ward con >1 escuela y sin school_id → 409 CONTEXT_SCHOOL_REQUIRED + candidatos.
//   - ward con school_id válido → 200, contexto del guardián impersonando al hijo.
//   - ward con 1 sola escuela y sin school_id → 200 (auto-resuelve la escuela).
//   - student no vinculado al usuario → 403 GUARDIAN_LINK_NOT_FOUND.
//   - la renovación (re-exchange con desired_context) preserva la terna
//     sujeto/actor (viaja firmada en el access_token).
//   - las rutas retiradas /auth/login, /auth/refresh y /auth/logout dan 404.
//
// Reusa el TestMain/roleflow.Setup del paquete (guardian_flow_test.go).
package guardian_flow_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UUIDs sembrados en `seeds/playground_v2/base` (edugo-infrastructure). El
// representante castro tiene 2 acudidos; carlos está en S1 y S3 (multi-escuela),
// diego solo en S1. mendoza es guardián de sofia (solo S1), NO de carlos.
const (
	tutorCastroEmail  = "tutor.castro@edugo.test"
	tutorMendozaEmail = "tutor.mendoza@edugo.test"

	studentCarlosID = "00000000-0000-0000-0000-000000000008" // acudido de castro, S1 + S3
	studentDiegoID  = "00000000-0000-0000-0000-000000000010" // acudido de castro, solo S1

	schoolS1ID = "b1000000-0000-0000-0000-000000000001"
	schoolS3ID = "b3000000-0000-0000-0000-000000000003"
)

// switchContextBody es el sub-set tipado del 200 de switch-context: la clave de
// respuesta es `context` (igual que el helper doSwitchContext interno).
type switchContextBody struct {
	AccessToken string `json:"access_token"`
	Context     *struct {
		RoleName         string `json:"role_name"`
		SchoolID         string `json:"school_id"`
		SubjectStudentID string `json:"subject_student_id"`
		ActorMode        string `json:"actor_mode"`
	} `json:"context"`
}

// ambiguousSchoolBody es el cuerpo del 409 CONTEXT_SCHOOL_REQUIRED: code +
// la lista de escuelas candidatas para que el cliente pinte el selector.
type ambiguousSchoolBody struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Candidates []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"candidates"`
}

// errorBody es el dto.ErrorResponse estándar (error/message/code).
type errorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// TestGuardianSwitch_WardMultiSchool_RequiresSelection: castro pide impersonar a
// carlos (vinculado en S1 y S3) sin school_id → 409 CONTEXT_SCHOOL_REQUIRED con
// las 2 escuelas candidatas en el body.
func TestGuardianSwitch_WardMultiSchool_RequiresSelection(t *testing.T) {
	env := roleflow.Get()
	login := roleflow.LoginRaw(t, env.Server, tutorCastroEmail)

	status, raw := roleflow.SwitchSubject(t, env.Server, login.AccessToken,
		"ward:"+studentCarlosID, "", "")
	require.Equalf(t, http.StatusConflict, status,
		"ward multi-escuela sin school_id: esperaba 409, got %d body=%s", status, string(raw))

	var body ambiguousSchoolBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 409 body=%s", string(raw))
	assert.Equal(t, "CONTEXT_SCHOOL_REQUIRED", body.Code,
		"code del 409 body=%s", string(raw))
	require.Lenf(t, body.Candidates, 2,
		"el 409 debe traer 2 escuelas candidatas (S1 y S3) body=%s", string(raw))

	ids := map[string]bool{}
	for _, c := range body.Candidates {
		ids[c.ID] = true
	}
	assert.Truef(t, ids[schoolS1ID], "S1 debe estar entre los candidatos body=%s", string(raw))
	assert.Truef(t, ids[schoolS3ID], "S3 debe estar entre los candidatos body=%s", string(raw))
}

// TestGuardianSwitch_WardWithSchool_Resolves: castro pide a carlos en S1 → 200,
// contexto del guardián impersonando al hijo (role=guardian, actor_mode=ward,
// subject=carlos, school=S1).
func TestGuardianSwitch_WardWithSchool_Resolves(t *testing.T) {
	env := roleflow.Get()
	login := roleflow.LoginRaw(t, env.Server, tutorCastroEmail)

	status, raw := roleflow.SwitchSubject(t, env.Server, login.AccessToken,
		"ward:"+studentCarlosID, schoolS1ID, "")
	require.Equalf(t, http.StatusOK, status,
		"ward con school_id válido: esperaba 200, got %d body=%s", status, string(raw))

	var body switchContextBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 200 body=%s", string(raw))
	require.NotNilf(t, body.Context, "context nil body=%s", string(raw))
	assert.Equal(t, "guardian", body.Context.RoleName, "role_name del contexto resuelto")
	assert.Equal(t, studentCarlosID, body.Context.SubjectStudentID, "subject_student_id = carlos")
	assert.Equal(t, roleflow.ActorModeWard, body.Context.ActorMode, "actor_mode = ward")
	assert.Equal(t, schoolS1ID, body.Context.SchoolID, "school_id = S1")
}

// TestGuardianSwitch_WardSingleSchool_Auto: castro pide a diego (solo S1) sin
// school_id → 200, la escuela se auto-resuelve (caso inequívoco).
func TestGuardianSwitch_WardSingleSchool_Auto(t *testing.T) {
	env := roleflow.Get()
	login := roleflow.LoginRaw(t, env.Server, tutorCastroEmail)

	status, raw := roleflow.SwitchSubject(t, env.Server, login.AccessToken,
		"ward:"+studentDiegoID, "", "")
	require.Equalf(t, http.StatusOK, status,
		"ward 1-escuela sin school_id: esperaba 200 (auto-resuelve), got %d body=%s",
		status, string(raw))

	var body switchContextBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 200 body=%s", string(raw))
	require.NotNilf(t, body.Context, "context nil body=%s", string(raw))
	assert.Equal(t, studentDiegoID, body.Context.SubjectStudentID, "subject_student_id = diego")
	assert.Equal(t, roleflow.ActorModeWard, body.Context.ActorMode, "actor_mode = ward")
	assert.Equal(t, schoolS1ID, body.Context.SchoolID, "school_id auto-resuelto = S1")
}

// TestGuardianSwitch_NoLink_Forbidden: mendoza (guardián de sofia, NO de carlos)
// pide impersonar a carlos → 403 GUARDIAN_LINK_NOT_FOUND.
func TestGuardianSwitch_NoLink_Forbidden(t *testing.T) {
	env := roleflow.Get()
	login := roleflow.LoginRaw(t, env.Server, tutorMendozaEmail)

	status, raw := roleflow.SwitchSubject(t, env.Server, login.AccessToken,
		"ward:"+studentCarlosID, schoolS1ID, "")
	require.Equalf(t, http.StatusForbidden, status,
		"ward no vinculado: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "GUARDIAN_LINK_NOT_FOUND", body.Code, "code del 403 body=%s", string(raw))
}

// TestGuardianSwitch_Renewal_PreservesSubject: la renovación de sesión conserva
// la terna sujeto/actor. Hasta el plan 053 T3.3 se renovaba con POST
// /auth/refresh; ya no existe (404): el refresh es de identity-core y el
// cliente RE-CANJEA un Identity Token nuevo con `desired_context` = el contexto
// que tenía. Se usa castro (no mendoza) para que la aserción discrimine: castro
// tiene 2 acudidos, así que su exchange inicial NO auto-resuelve a ningún ward;
// si el re-exchange ignorara desired_context, el contexto sería otro (o nil).
//
// Pasos: exchange → switch-context a ward carlos en S3 → re-exchange con el
// desired_context de ese token → el token nuevo (claims firmados, decodificados
// sin verificar) y el active_context de la respuesta conservan
// subject_student_id=carlos, actor_mode=ward, school_id=S3 y el rol.
func TestGuardianSwitch_Renewal_PreservesSubject(t *testing.T) {
	env := roleflow.Get()

	login := roleflow.LoginRaw(t, env.Server, tutorCastroEmail)
	require.Nil(t, login.ActiveContext,
		"precondición: castro (2 acudidos) no auto-resuelve contexto en el exchange")

	status, raw := roleflow.SwitchSubject(t, env.Server, login.AccessToken,
		"ward:"+studentCarlosID, schoolS3ID, "")
	require.Equalf(t, http.StatusOK, status,
		"switch a ward carlos en S3: esperaba 200, got %d body=%s", status, string(raw))
	var switched switchContextBody
	require.NoError(t, json.Unmarshal(raw, &switched), "parse switch body=%s", string(raw))
	before, ok := decodeJWTPayload(t, switched.AccessToken)["active_context"].(map[string]any)
	require.Truef(t, ok, "el token del switch debe traer active_context")
	require.Equal(t, studentCarlosID, before["subject_student_id"], "precondición: sujeto carlos")
	roleID, _ := before["role_id"].(string)
	require.NotEmpty(t, roleID, "precondición: role_id en el token del switch")

	// Renovación: re-exchange con un Identity Token nuevo + desired_context.
	status, raw = roleflow.Exchange(t, env.Server, tutorCastroEmail, &roleflow.DesiredContext{
		SchoolID: schoolS3ID,
		RoleID:   roleID,
		Subject:  "ward:" + studentCarlosID,
	})
	require.Equalf(t, http.StatusOK, status,
		"re-exchange con desired_context: esperaba 200, got %d body=%s", status, string(raw))
	var renewed roleflow.LoginResponse
	require.NoError(t, json.Unmarshal(raw, &renewed), "parse re-exchange body=%s", string(raw))
	require.NotEmpty(t, renewed.AccessToken, "re-exchange debe devolver un access_token nuevo")
	require.NotNil(t, renewed.ActiveContext, "re-exchange debe devolver active_context")
	assert.Equal(t, studentCarlosID, renewed.ActiveContext.SubjectStudentID,
		"active_context de la respuesta: subject_student_id = carlos")
	assert.Equal(t, roleflow.ActorModeWard, renewed.ActiveContext.ActorMode,
		"active_context de la respuesta: actor_mode = ward")

	ac, ok := decodeJWTPayload(t, renewed.AccessToken)["active_context"].(map[string]any)
	require.Truef(t, ok, "el access_token renovado debe traer active_context")
	assert.Equal(t, studentCarlosID, ac["subject_student_id"],
		"la renovación debe preservar subject_student_id (carlos)")
	assert.Equal(t, roleflow.ActorModeWard, ac["actor_mode"],
		"la renovación debe preservar actor_mode=ward")
	assert.Equal(t, schoolS3ID, ac["school_id"], "la renovación debe preservar school_id = S3")
	assert.Equal(t, roleID, ac["role_id"], "la renovación debe preservar el rol")
}

// TestLegacyAuthRoutes_Gone: plan 053 T3.3 retiró el login, el refresh y el
// logout propios de identity; las rutas ya no existen.
func TestLegacyAuthRoutes_Gone(t *testing.T) {
	env := roleflow.Get()
	for _, path := range []string{"/api/v1/auth/login", "/api/v1/auth/refresh", "/api/v1/auth/logout"} {
		resp, err := http.Post(env.Server.URL+path, "application/json", bytes.NewReader([]byte("{}")))
		require.NoError(t, err)
		_ = resp.Body.Close()
		assert.Equalf(t, http.StatusNotFound, resp.StatusCode, "POST %s debe dar 404", path)
	}
}

// decodeJWTPayload decodifica el segmento central (payload) de un JWT como
// base64url sin verificar la firma — solo nos interesan los claims, no la
// autenticidad (el server ya validó/firmó). Devuelve el mapa de claims.
func decodeJWTPayload(t *testing.T, token string) map[string]any {
	t.Helper()
	parts := strings.Split(token, ".")
	require.Lenf(t, parts, 3, "JWT mal formado: %s", token)

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoErrorf(t, err, "decode base64url del payload JWT: %s", parts[1])

	var claims map[string]any
	require.NoError(t, json.Unmarshal(payload, &claims), "parse claims JWT: %s", string(payload))
	return claims
}
