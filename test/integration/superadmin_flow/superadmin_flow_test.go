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

// TestSuperAdminGlobalFlow_CrossAPI ejercita los 5 pasos del flujo
// SchoolSelector → switchContext → UnitSelector → Dashboard contra los
// 3 AppServer reales (identity, academic, platform).
//
// Cada paso cubre una clase de bug (numerada por la sesión 2026-05-12):
//
//	step 1 (login)              → cubre bug-1: super_admin debe recibir
//	                              context:browse_schools y context:browse_units
//	                              en active_context.permissions.
//	step 2 (academic /schools)  → cubre bug-2: academic.routes_school
//	                              acepta context:browse_schools.
//	step 3 (identity /units)    → cubre bug-2 lado units (routes_unit /
//	                              auth/contexts/schools/:id/units).
//	step 4 (switch-context)     → cubre bug-4: switch_context inyecta
//	                              school_id/school_name/academic_unit_id/
//	                              academic_unit_name aunque el rol matched
//	                              sea global.
//	step 5 (platform /sync)     → cubre regresiones del bundle del
//	                              Dashboard (menu/screens no vacíos con el
//	                              nuevo token que ya tiene contexto).
//
// El test es secuencial — cada paso depende del anterior (el token de
// login se reusa hasta switch-context, y luego se rota al token nuevo).
func TestSuperAdminGlobalFlow_CrossAPI(t *testing.T) {
	// ----------------------------------------------------------------
	// Step 1: POST identity /auth/login con global-super@e2e.test.
	// ----------------------------------------------------------------
	loginBody := mustJSON(t, map[string]string{
		"email":    globalUserEmail,
		"password": globalUserPassword,
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
	require.NotNil(t, login.ActiveContext, "step 1: active_context must not be nil")

	// Schools[] DEBE estar vacío — el global super_admin no tiene
	// membership, no debería listarse ninguna escuela en el login.
	assert.Empty(t, login.Schools, "step 1: schools[] must be empty for global super_admin")

	// Bug-class 1: super_admin cubre context.browse_schools/units en
	// grants (Pass 3 wildcard-first: el rol global tiene `*`).
	roleflow.AssertGrantsContains(t, login.ActiveContext.Grants,
		"context.browse_schools", "context.browse_units")

	loginToken := login.AccessToken

	// ----------------------------------------------------------------
	// Step 2: GET academic /api/v1/schools con loginToken.
	// El handler usa RequireAnyPermission(schools:read | context:browse_schools).
	// ----------------------------------------------------------------
	schoolsResp := mustGet(t, academicServer.URL+"/api/v1/schools", loginToken)
	require.Equalf(t, http.StatusOK, schoolsResp.status,
		"step 2 academic /schools: expected 200, body=%s", string(schoolsResp.body))

	var schoolsBody struct {
		Schools []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"schools"`
		Total int `json:"total"`
	}
	require.NoError(t, json.Unmarshal(schoolsResp.body, &schoolsBody),
		"step 2 schools: parse body=%s", string(schoolsResp.body))

	require.GreaterOrEqual(t, len(schoolsBody.Schools), 1,
		"step 2 bug-2: academic /schools must list at least 1 school (demo seed)")

	// ----------------------------------------------------------------
	// Step 3: GET identity /api/v1/auth/contexts/schools/:id/units.
	// El SchoolSelector real prueba escuela por escuela hasta encontrar
	// una con unidades. Replicamos esa semántica para que el test sea
	// estable ante el orden de listado (algunas escuelas del seed e2e
	// no tienen unidades; las del demo seed sí).
	// ----------------------------------------------------------------
	type schoolUnit struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var unitsBody struct {
		Units []schoolUnit `json:"units"`
		Total int          `json:"total"`
	}
	var pickedSchoolID, pickedUnitID string
	for _, candidate := range schoolsBody.Schools {
		unitsURL := fmt.Sprintf("%s/api/v1/auth/contexts/schools/%s/units",
			identityServer.URL, candidate.ID)
		unitsResp := mustGet(t, unitsURL, loginToken)
		require.Equalf(t, http.StatusOK, unitsResp.status,
			"step 3 identity /units (school=%s): expected 200, body=%s",
			candidate.ID, string(unitsResp.body))
		var body struct {
			Units []schoolUnit `json:"units"`
			Total int          `json:"total"`
		}
		require.NoError(t, json.Unmarshal(unitsResp.body, &body))
		if len(body.Units) >= 1 {
			pickedSchoolID = candidate.ID
			pickedUnitID = body.Units[0].ID
			unitsBody = body
			break
		}
	}
	require.NotEmpty(t, pickedSchoolID,
		"step 3 bug-2 (units side): at least one listed school must have units")
	require.NotEmpty(t, pickedUnitID,
		"step 3: picked academic_unit_id must not be empty")
	require.GreaterOrEqual(t, len(unitsBody.Units), 1,
		"step 3: picked school must have units")

	// ----------------------------------------------------------------
	// Step 4: POST identity /api/v1/auth/switch-context.
	// ----------------------------------------------------------------
	switchBody := mustJSON(t, map[string]string{
		"school_id":        pickedSchoolID,
		"academic_unit_id": pickedUnitID,
	})
	switchResp := mustPost(t, identityServer.URL+"/api/v1/auth/switch-context",
		loginToken, switchBody)
	require.Equalf(t, http.StatusOK, switchResp.status,
		"step 4 switch-context: expected 200, body=%s", string(switchResp.body))

	var switched struct {
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
	require.NoError(t, json.Unmarshal(switchResp.body, &switched),
		"step 4 switch-context: parse body=%s", string(switchResp.body))

	require.NotEmpty(t, switched.AccessToken, "step 4: new access_token must not be empty")
	require.NotNil(t, switched.Context, "step 4: context must not be nil")

	// Bug-class 4: el rol global injecta school_id/school_name/academic_unit_*.
	assert.NotEmpty(t, switched.Context.SchoolID,
		"step 4 bug-4: switch_context must inject school_id for global role")
	assert.Equal(t, pickedSchoolID, switched.Context.SchoolID,
		"step 4: returned school_id must equal the one requested")
	assert.NotEmpty(t, switched.Context.SchoolName,
		"step 4 bug-4: switch_context must inject school_name for global role")
	assert.NotEmpty(t, switched.Context.AcademicUnitID,
		"step 4 bug-4: switch_context must inject academic_unit_id for global role")
	assert.Equal(t, pickedUnitID, switched.Context.AcademicUnitID,
		"step 4: returned academic_unit_id must equal the one requested")
	assert.NotEmpty(t, switched.Context.AcademicUnitName,
		"step 4 bug-4: switch_context must inject academic_unit_name for global role")

	rotatedToken := switched.AccessToken

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
