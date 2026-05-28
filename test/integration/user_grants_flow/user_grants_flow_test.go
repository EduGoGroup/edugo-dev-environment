//go:build integration

// Package user_grants_flow valida el feature `iam.user_grants` (P4-2):
//
//  1. Que los user_grants del seed demo se materializan en los grants
//     efectivos del usuario consultados vía `GET /api/v1/users/:id/permissions`
//     (deny puntual para un student, allow extra para un teacher).
//  2. Que los endpoints REST bajo `/api/v1/users/:user_id/grants`
//     (list/create/delete) funcionan end-to-end, incluyendo:
//     - rechazo de `permission_pattern = "*"`,
//     - conflicto 409 en duplicados,
//     - grants expirados ignorados por el snapshot de permisos.
//
// Nota de implementación: tanto `LoginResponse.ActiveContext.Grants`
// (rama `Find*Context*`) como `GET /users/:id/permissions` (rama
// `GetUserPermissions`) unen `role_grants ∪ user_grants` activos. Los
// tests cubren ambos caminos: las aserciones de mutación REST consultan
// `/users/:id/permissions` (más directo); los dos demos validan también
// que el login del usuario afectado refleja el override.
//
// El actor admin es `super@edugo.test` (super_admin con pattern `*`, cubre
// `admin.users.grants.manage` y `permissions.mgmt.read` por wildcard).
// El target limpio para las mutaciones es `prof.gonzalez@edugo.test`;
// cada sub-test usa un permission_pattern distinto para evitar interferencia
// entre ejecuciones aunque Go corra los `Test*` en orden alfabético.
package user_grants_flow_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	superAdminEmail = "super@edugo.test"
	studentID       = "00000000-0000-0000-0000-000000000008"
	teacherSeedID   = "00000000-0000-0000-0000-000000000005"
	teacherCleanID  = "00000000-0000-0000-0000-000000000006"
)

// Espejo de `dto.UserGrantDTO` / `dto.ListUserGrantsResponse` / `dto.CreateUserGrantResponse`
// del módulo identity. Sólo los campos que los tests consumen.
type userGrantDTO struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	ScopePattern      string  `json:"scope_pattern"`
	PermissionPattern string  `json:"permission_pattern"`
	Effect            string  `json:"effect"`
	ExpiresAt         *string `json:"expires_at,omitempty"`
	GrantedBy         *string `json:"granted_by,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

type listUserGrantsResponse struct {
	Items []userGrantDTO `json:"items"`
}

type createUserGrantResponse struct {
	Grant userGrantDTO `json:"grant"`
}

// userPermissionsResponse — mirror de `dto.UserPermissionsResponse`.
type userPermissionsResponse struct {
	UserID string          `json:"user_id"`
	Grants roleflow.Grants `json:"grants"`
}

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

// TestUserGrants_DemoSeedDenyOverride — verifica que el seed demo emite un
// deny puntual sobre `academic.grades.read` para `est.carlos`, y que el
// matcher aplica deny > allow al consultar `GrantsAllow` sobre el snapshot
// efectivo (role_grants ∪ user_grants).
func TestUserGrants_DemoSeedDenyOverride(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)
	grants := fetchEffectiveGrants(t, env.Server, studentID, super.AccessToken)

	assert.Contains(t, grants.Deny, "academic.grades.read",
		"seed demo: student debe tener deny puntual sobre academic.grades.read")

	assert.False(t, roleflow.GrantsAllow(grants, "academic.grades.read"),
		"deny > allow: GrantsAllow debe retornar false sobre el permiso negado")

	// Sanity: otros permisos del rol student no se ven afectados.
	assert.True(t, roleflow.GrantsAllow(grants, "content.materials.read"),
		"student conserva el resto de sus permisos por rol")

	// El login del propio student también refleja el override:
	// ActiveContext.Grants debe unir role_grants ∪ user_grants.
	studentLogin := roleflow.Login(t, env.Server, "est.carlos@edugo.test", roleflow.DemoPassword)
	require.NotNil(t, studentLogin.ActiveContext)
	assert.Contains(t, studentLogin.ActiveContext.Grants.Deny, "academic.grades.read",
		"login.ActiveContext.Grants debe incluir el deny puntual del seed")
}

// TestUserGrants_DemoSeedAllowOverride — verifica que el seed demo concede
// `admin.users.create` (que NO está en el rol teacher) a `prof.martinez`
// vía user_grant con expires_at futuro.
func TestUserGrants_DemoSeedAllowOverride(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)
	grants := fetchEffectiveGrants(t, env.Server, teacherSeedID, super.AccessToken)

	assert.Contains(t, grants.Allow, "admin.users.create",
		"seed demo: teacher debe recibir allow extra sobre admin.users.create")
	assert.True(t, roleflow.GrantsAllow(grants, "admin.users.create"),
		"GrantsAllow debe permitir el permiso concedido por user_grant")

	// El login del propio teacher también refleja el override.
	teacherLogin := roleflow.Login(t, env.Server, "prof.martinez@edugo.test", roleflow.DemoPassword)
	require.NotNil(t, teacherLogin.ActiveContext)
	assert.Contains(t, teacherLogin.ActiveContext.Grants.Allow, "admin.users.create",
		"login.ActiveContext.Grants debe incluir el allow extra del seed")
}

// TestUserGrants_API_CreateListDelete — flujo CRUD completo desde super_admin
// sobre prof.gonzalez. Verifica que los cambios se reflejan en los grants
// efectivos vía /users/:id/permissions tras cada mutación.
func TestUserGrants_API_CreateListDelete(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)
	bearer := super.AccessToken

	// Baseline: no debe existir el deny pre-existente.
	pre := fetchEffectiveGrants(t, env.Server, teacherCleanID, bearer)
	assert.NotContains(t, pre.Deny, "academic.attendance.update",
		"baseline: prof.gonzalez no debe tener deny pre-existente sobre academic.attendance.update")

	createPath := "/api/v1/users/" + teacherCleanID + "/grants"

	status, body := postJSON(t, env.Server, createPath, bearer, map[string]any{
		"permission_pattern": "academic.attendance.update",
		"effect":             "deny",
	})
	require.Truef(t, status == http.StatusCreated || status == http.StatusOK,
		"POST grant: expected 2xx, got %d body=%s", status, string(body))

	// GET list — busca el grant recién creado y captura su id.
	status, body = roleflow.GetJSON(t, env.Server, createPath, bearer)
	require.Equalf(t, http.StatusOK, status, "GET grants: status %d body=%s", status, string(body))

	var list listUserGrantsResponse
	require.NoError(t, json.Unmarshal(body, &list), "GET grants: parse body=%s", string(body))

	grantID := findGrantID(list.Items, "academic.attendance.update", "deny")
	require.NotEmptyf(t, grantID,
		"GET grants: debe contener el grant recién creado, items=%+v", list.Items)

	// Tras el POST: el deny debe aparecer en los grants efectivos.
	post := fetchEffectiveGrants(t, env.Server, teacherCleanID, bearer)
	assert.Contains(t, post.Deny, "academic.attendance.update",
		"prof.gonzalez debe tener deny sobre academic.attendance.update tras el POST")

	// DELETE — limpia el grant.
	deletePath := createPath + "/" + grantID
	status, body = deleteRequest(t, env.Server, deletePath, bearer)
	require.Truef(t, status == http.StatusOK || status == http.StatusNoContent,
		"DELETE grant: expected 200/204, got %d body=%s", status, string(body))

	// Tras el DELETE: el deny ya no debe estar.
	after := fetchEffectiveGrants(t, env.Server, teacherCleanID, bearer)
	assert.NotContains(t, after.Deny, "academic.attendance.update",
		"prof.gonzalez no debe tener deny tras el DELETE")
}

// TestUserGrants_API_RejectWildcardStar — el usecase rechaza pattern "*"
// porque concedería/negaría TODO el catálogo, equivalente a un bypass.
func TestUserGrants_API_RejectWildcardStar(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)

	status, body := postJSON(t, env.Server,
		"/api/v1/users/"+teacherCleanID+"/grants",
		super.AccessToken,
		map[string]any{
			"permission_pattern": "*",
			"effect":             "deny",
		})
	require.Equalf(t, http.StatusBadRequest, status,
		"POST '*' grant: expected 400, got %d body=%s", status, string(body))
}

// TestUserGrants_API_Duplicate — un (user_id, scope_pattern, permission_pattern, effect)
// duplicado retorna 409. Se hace cleanup best-effort al final.
func TestUserGrants_API_Duplicate(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)
	bearer := super.AccessToken
	path := "/api/v1/users/" + teacherCleanID + "/grants"

	body := map[string]any{
		"permission_pattern": "content.materials.delete",
		"effect":             "deny",
	}

	status, raw := postJSON(t, env.Server, path, bearer, body)
	require.Truef(t, status == http.StatusCreated || status == http.StatusOK,
		"POST primer grant: expected 2xx, got %d body=%s", status, string(raw))

	var created createUserGrantResponse
	_ = json.Unmarshal(raw, &created) // best-effort para cleanup

	// POST idéntico → 409.
	status, raw = postJSON(t, env.Server, path, bearer, body)
	assert.Equalf(t, http.StatusConflict, status,
		"POST duplicado: expected 409, got %d body=%s", status, string(raw))

	// Cleanup best-effort: no aborta el test si el grant no existe.
	if created.Grant.ID != "" {
		_, _ = deleteRequest(t, env.Server, path+"/"+created.Grant.ID, bearer)
	}
}

// TestUserGrants_API_ExpiredIgnored — un grant con expires_at en el pasado
// se acepta al crearlo (la regla "active" vive en el consumidor), pero NO
// aparece en los grants efectivos porque GetUserPermissions filtra los
// expirados antes de unirlos a role_grants.
func TestUserGrants_API_ExpiredIgnored(t *testing.T) {
	env := roleflow.Get()

	super := roleflow.Login(t, env.Server, superAdminEmail, roleflow.DemoPassword)
	bearer := super.AccessToken
	path := "/api/v1/users/" + teacherCleanID + "/grants"

	expiredAt := "2020-01-01T00:00:00Z"
	status, raw := postJSON(t, env.Server, path, bearer, map[string]any{
		"permission_pattern": "academic.subjects.delete",
		"effect":             "deny",
		"expires_at":         expiredAt,
	})
	require.Truef(t, status == http.StatusCreated || status == http.StatusOK,
		"POST grant expirado: expected 2xx, got %d body=%s", status, string(raw))

	var created createUserGrantResponse
	_ = json.Unmarshal(raw, &created)

	// Snapshot efectivo debe ignorar el grant expirado.
	grants := fetchEffectiveGrants(t, env.Server, teacherCleanID, bearer)
	assert.NotContains(t, grants.Deny, "academic.subjects.delete",
		"grant expirado no debe aparecer en grants.deny efectivos")

	// Cleanup best-effort.
	if created.Grant.ID != "" {
		_, _ = deleteRequest(t, env.Server, path+"/"+created.Grant.ID, bearer)
	}
}

// fetchEffectiveGrants consulta GET /api/v1/users/:id/permissions con bearer
// y retorna los grants efectivos (role_grants ∪ user_grants activos).
func fetchEffectiveGrants(t *testing.T, server *httptest.Server, userID, bearer string) roleflow.Grants {
	t.Helper()
	status, body := roleflow.GetJSON(t, server,
		"/api/v1/users/"+userID+"/permissions", bearer)
	require.Equalf(t, http.StatusOK, status,
		"GET /users/%s/permissions: status %d body=%s", userID, status, string(body))

	var resp userPermissionsResponse
	require.NoError(t, json.Unmarshal(body, &resp),
		"parse permissions response: body=%s", string(body))
	return resp.Grants
}

// findGrantID busca el primer grant con el (permission_pattern, effect)
// indicado y retorna su id, o "" si no existe.
func findGrantID(items []userGrantDTO, pattern, effect string) string {
	for _, it := range items {
		if it.PermissionPattern == pattern && it.Effect == effect {
			return it.ID
		}
	}
	return ""
}

// postJSON ejecuta un POST JSON con bearer y retorna (status, body bytes).
func postJSON(t *testing.T, server *httptest.Server, path, bearer string, body any) (int, []byte) {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, out
}

// deleteRequest ejecuta un DELETE con bearer y retorna (status, body bytes).
func deleteRequest(t *testing.T, server *httptest.Server, path, bearer string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodDelete, server.URL+path, nil)
	require.NoError(t, err)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, out
}
