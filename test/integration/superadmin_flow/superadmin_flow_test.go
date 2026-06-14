//go:build integration

package superadmin_flow_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSuperAdminGlobalFlow_CrossAPI ejercita el flujo
// SchoolSelector → switchContext → UnitSelector → Dashboard del super_admin
// global (rol L0 super_admin, school_id=NULL, sin membership) contra los
// 3 AppServer reales (identity, academic, platform).
//
// Mecánica MP-08 DEC-A (login → switch-context): el super_admin global no
// tiene membership, así que el login NO auto-selecciona contexto. En su
// lugar `FindUserSchools` (vía hasActiveSystemRole→findAllActiveSchools)
// devuelve TODAS las escuelas activas del sistema y `active_context` queda
// nil. Para operar contra los handlers scope=school del backend hay que
// resolver un contexto primero con switch-context, que rota el JWT con el
// school_id/academic_unit_id elegidos y los grants del rol (`*`).
//
// El test es secuencial — cada paso depende del anterior:
//
//	step 1 (login)              → schools[] poblado con las escuelas del
//	                              sistema y active_context nil (DEC-A).
//	step 2 (switch-context)     → resuelve contexto sobre una escuela del
//	                              seed (multi-unidad → 409 CONTEXT_UNIT_REQUIRED
//	                              → reintento con unidad). El rol global
//	                              injecta school_id/name/academic_unit_id/name
//	                              y los grants del super_admin; el JWT se rota.
//	step 3 (academic /schools)  → con el token rotado, academic.routes_school
//	                              acepta context.browse_schools (wildcard `*`).
//	step 4 (identity /units)    → con el token rotado, identity lista las
//	                              unidades de la escuela elegida.
//	step 5 (platform /sync)     → el bundle del Dashboard no llega vacío con
//	                              el token que ya trae contexto.
func TestSuperAdminGlobalFlow_CrossAPI(t *testing.T) {
	// ----------------------------------------------------------------
	// Step 1: POST identity /auth/login con global-super@e2e.test.
	// MP-08 DEC-C: el login exige `system`; el super_admin tiene acceso
	// a "kmp" vía iam.system_roles (L4).
	// ----------------------------------------------------------------
	loginBody := mustJSON(t, map[string]string{
		"email":    globalUserEmail,
		"password": globalUserPassword,
		"system":   "kmp",
	})

	loginResp := mustPost(t, identityServer.URL+"/api/v1/auth/login", "", loginBody)
	require.Equalf(t, http.StatusOK, loginResp.status,
		"step 1 login: expected 200, body=%s", string(loginResp.body))

	var login struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Schools      []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"schools"`
		ActiveContext *struct {
			RoleID   string `json:"role_id"`
			RoleName string `json:"role_name"`
			Grants   struct {
				Allow []string `json:"allow"`
				Deny  []string `json:"deny"`
			} `json:"grants"`
		} `json:"active_context"`
	}
	require.NoError(t, json.Unmarshal(loginResp.body, &login),
		"step 1 login: parse response body=%s", string(loginResp.body))

	require.NotEmpty(t, login.AccessToken, "step 1: access_token must not be empty")

	// MP-08 DEC-A: con >1 escuela el login NO auto-selecciona contexto, así que
	// active_context queda nil y la cascada se completa con switch-context.
	require.Nil(t, login.ActiveContext,
		"step 1 (DEC-A): active_context must be nil for the multi-school global super_admin")

	// El super_admin global no tiene membership, pero hasActiveSystemRole→
	// findAllActiveSchools le devuelve TODAS las escuelas activas del sistema.
	require.GreaterOrEqual(t, len(login.Schools), 1,
		"step 1 (DEC-A): schools[] must list the system schools for the global super_admin")

	loginToken := login.AccessToken

	// ----------------------------------------------------------------
	// Step 2: POST identity /api/v1/auth/switch-context.
	// El SchoolSelector real prueba escuela por escuela hasta resolver un
	// contexto a nivel unidad. Las escuelas del seed base son multi-unidad,
	// así que el switch sin unidad devuelve 409 CONTEXT_UNIT_REQUIRED y hay
	// que reintentar con una unidad de la escuela. Replicamos esa semántica
	// (idéntica a roleflow.switchContext) para que el test sea estable ante
	// el orden de listado.
	// ----------------------------------------------------------------
	switched, pickedSchoolID, pickedUnitID := resolveUnitContext(t, loginToken, login.Schools)

	require.NotEmpty(t, switched.AccessToken, "step 2: new access_token must not be empty")
	require.NotNil(t, switched.Context, "step 2: context must not be nil")

	// El rol global injecta school_id/school_name/academic_unit_id/name aunque
	// el user_role matched tenga school_id IS NULL (super_admin).
	assert.NotEmpty(t, switched.Context.SchoolID,
		"step 2: switch_context must inject school_id for global role")
	assert.Equal(t, pickedSchoolID, switched.Context.SchoolID,
		"step 2: returned school_id must equal the one requested")
	assert.NotEmpty(t, switched.Context.SchoolName,
		"step 2: switch_context must inject school_name for global role")
	assert.NotEmpty(t, switched.Context.AcademicUnitID,
		"step 2: switch_context must inject academic_unit_id for global role")
	assert.Equal(t, pickedUnitID, switched.Context.AcademicUnitID,
		"step 2: returned academic_unit_id must equal the one requested")
	assert.NotEmpty(t, switched.Context.AcademicUnitName,
		"step 2: switch_context must inject academic_unit_name for global role")

	// El super_admin resuelve por wildcard: el contexto rotado cubre
	// context.browse_schools/units (el rol global tiene `*`).
	roleflow.AssertGrantsContains(t, roleflow.Grants{
		Allow: switched.Context.Grants.Allow,
		Deny:  switched.Context.Grants.Deny,
	}, "context.browse_schools", "context.browse_units")

	rotatedToken := switched.AccessToken

	// ----------------------------------------------------------------
	// Step 3: GET academic /api/v1/schools con el token rotado.
	// El handler usa RequireAnyPermission(schools:read | context:browse_schools);
	// el token de login (sin contexto) no trae grants y daría 403 — el rotado sí.
	// ----------------------------------------------------------------
	schoolsResp := mustGet(t, academicServer.URL+"/api/v1/schools", rotatedToken)
	require.Equalf(t, http.StatusOK, schoolsResp.status,
		"step 3 academic /schools: expected 200, body=%s", string(schoolsResp.body))

	var schoolsBody struct {
		Schools []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"schools"`
		Total int `json:"total"`
	}
	require.NoError(t, json.Unmarshal(schoolsResp.body, &schoolsBody),
		"step 3 schools: parse body=%s", string(schoolsResp.body))

	require.GreaterOrEqual(t, len(schoolsBody.Schools), 1,
		"step 3: academic /schools must list at least 1 school (base seed)")

	// ----------------------------------------------------------------
	// Step 4: GET identity /api/v1/auth/contexts/schools/:id/units con el
	// token rotado, sobre la escuela ya elegida. Confirma el lado units del
	// SchoolSelector.
	// ----------------------------------------------------------------
	unitsURL := fmt.Sprintf("%s/api/v1/auth/contexts/schools/%s/units",
		identityServer.URL, pickedSchoolID)
	unitsResp := mustGet(t, unitsURL, rotatedToken)
	require.Equalf(t, http.StatusOK, unitsResp.status,
		"step 4 identity /units (school=%s): expected 200, body=%s",
		pickedSchoolID, string(unitsResp.body))

	var unitsBody struct {
		Units []schoolUnit `json:"units"`
		Total int          `json:"total"`
	}
	require.NoError(t, json.Unmarshal(unitsResp.body, &unitsBody),
		"step 4 units: parse body=%s", string(unitsResp.body))
	require.GreaterOrEqual(t, len(unitsBody.Units), 1,
		"step 4: picked school must have units")

	// ----------------------------------------------------------------
	// Step 5: GET platform /api/v1/sync/bundle con el token rotado.
	// ----------------------------------------------------------------
	bundleResp := mustGet(t, platformServer.URL+"/api/v1/sync/bundle", rotatedToken)
	require.Equalf(t, http.StatusOK, bundleResp.status,
		"step 5 platform /sync/bundle: expected 200, body=%s", string(bundleResp.body))

	var bundle struct {
		Menu    []map[string]any `json:"menu"`
		Screens map[string]any   `json:"screens"`
	}
	require.NoError(t, json.Unmarshal(bundleResp.body, &bundle),
		"step 5 bundle: parse body=%s", string(bundleResp.body))

	assert.NotEmpty(t, bundle.Menu, "step 5: bundle.menu must not be empty for super_admin")
	assert.NotEmpty(t, bundle.Screens, "step 5: bundle.screens must not be empty for super_admin")
}

// schoolUnit es el sub-set tipado de una unidad listada por
// GET /auth/contexts/schools/:id/units.
type schoolUnit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// switchContextResult es el sub-set tipado del payload
// POST /auth/switch-context que el test consume.
type switchContextResult struct {
	AccessToken string `json:"access_token"`
	Context     *struct {
		RoleID           string `json:"role_id"`
		RoleName         string `json:"role_name"`
		SchoolID         string `json:"school_id"`
		SchoolName       string `json:"school_name"`
		AcademicUnitID   string `json:"academic_unit_id"`
		AcademicUnitName string `json:"academic_unit_name"`
		Grants           struct {
			Allow []string `json:"allow"`
			Deny  []string `json:"deny"`
		} `json:"grants"`
	} `json:"context"`
}

// resolveUnitContext replica la cascada del SchoolSelector real: recorre las
// escuelas disponibles y para cada una intenta switch-context. Si la escuela es
// multi-unidad (409 CONTEXT_UNIT_REQUIRED) resuelve su primera unidad vía
// GET /auth/contexts/schools/:id/units y reintenta. Devuelve el primer contexto
// resuelto a nivel unidad (con academic_unit_id), junto a la escuela y unidad
// elegidas. Falla el test si ninguna escuela resuelve a una unidad.
func resolveUnitContext(t *testing.T, loginToken string, schools []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}) (switchContextResult, string, string) {
	t.Helper()
	for _, candidate := range schools {
		unitID := firstSchoolUnit(t, loginToken, candidate.ID)
		if unitID == "" {
			continue
		}
		switchBody := mustJSON(t, map[string]string{
			"school_id":        candidate.ID,
			"academic_unit_id": unitID,
		})
		switchResp := mustPost(t, identityServer.URL+"/api/v1/auth/switch-context",
			loginToken, switchBody)
		require.Equalf(t, http.StatusOK, switchResp.status,
			"switch-context school=%s unit=%s: expected 200, body=%s",
			candidate.ID, unitID, string(switchResp.body))

		var switched switchContextResult
		require.NoError(t, json.Unmarshal(switchResp.body, &switched),
			"switch-context: parse body=%s", string(switchResp.body))
		require.NotNil(t, switched.Context,
			"switch-context: context nil body=%s", string(switchResp.body))
		if switched.Context.AcademicUnitID != "" {
			return switched, candidate.ID, unitID
		}
	}
	require.FailNow(t, "no school resolved a unit-level context for the global super_admin")
	return switchContextResult{}, "", ""
}

// firstSchoolUnit lista las unidades de la escuela y devuelve el id de la
// primera, o "" si la escuela no tiene unidades.
func firstSchoolUnit(t *testing.T, bearer, schoolID string) string {
	t.Helper()
	unitsURL := fmt.Sprintf("%s/api/v1/auth/contexts/schools/%s/units",
		identityServer.URL, schoolID)
	unitsResp := mustGet(t, unitsURL, bearer)
	require.Equalf(t, http.StatusOK, unitsResp.status,
		"list units (school=%s): expected 200, body=%s",
		schoolID, string(unitsResp.body))
	var body struct {
		Units []schoolUnit `json:"units"`
		Total int          `json:"total"`
	}
	require.NoError(t, json.Unmarshal(unitsResp.body, &body),
		"list units: parse body=%s", string(unitsResp.body))
	if len(body.Units) == 0 {
		return ""
	}
	return body.Units[0].ID
}

// ---------------------------------------------------------------------
// HTTP helpers (compactos, sin dependencias del testutil de cada API).
// ---------------------------------------------------------------------

type httpResult struct {
	status int
	body   []byte
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

func mustPost(t *testing.T, url, bearer string, body []byte) httpResult {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	return doRequest(t, req)
}

func mustGet(t *testing.T, url, bearer string) httpResult {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	return doRequest(t, req)
}

func doRequest(t *testing.T, req *http.Request) httpResult {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return httpResult{status: resp.StatusCode, body: body}
}
