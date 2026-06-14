//go:build integration

// Package school_admin_flow valida que el rol L4 school_admin del seed
// base `admin.sanignacio@edugo.test` puede autenticarse y recibe los
// patterns canónicos en `ActiveContext.Grants.Allow`.
//
// Pass 2 (single-path Grants): ya no existe path "legacy"; el test
// se limita a verificar el formato Grants emitido por el identity
// server contra el seed L4.
package school_admin_flow_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	roleName = "school_admin"
	// admin.sanignacio@edugo.test → user_role role_id = L4_ROLE_SCHOOL_ADMIN_ID,
	// school_id = b1000000-0000-0000-0000-000000000001 (San Ignacio).
	userEmail = "admin.sanignacio@edugo.test"
)

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

func TestSchoolAdminFlow_Grants(t *testing.T) {
	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// Patterns extraídos del seed L4 `rolePermissionGrants()` para
	// school_admin: dominio CRUD de la institución (sample
	// representativo — la lista completa supera 60 patterns).
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"admin.users.read",
		"admin.users.create",
		"admin.users.update",
		"admin.schools.read",
		"admin.schools.update",
		"admin.schools.manage",
		"academic.units.create",
		"academic.units.read",
		"academic.units.update",
		"academic.units.delete",
		"academic.subjects.create",
		"academic.subjects.read",
		"academic.periods.create",
		"academic.periods.activate",
		"academic.grades.read",
		"academic.grades.finalize",
		"academic.memberships.create",
		"academic.memberships.read",
		"academic.announcements.read",
		"academic.announcements.create",
		"academic.announcements.update",
		"academic.announcements.delete",
		"content.materials.delete",
		"content.assessments.create",
		"content.assessments.delete",
		"admin.roles.read",
		"admin.roles.create",
		"admin.roles.update",
		"admin.roles.delete",
		"admin.system_settings.read",
		"admin.system_settings.update",
		"dashboard.view",
		"menu.read",
		"notifications.read",
	)

	// Sanity: deny vacío (school_admin no recibe denies en seed).
	assert.Empty(t, resp.ActiveContext.Grants.Deny,
		"school_admin: grants.deny must be empty")

	// E2E: el endpoint autenticado debe responder 200 con el token.
	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for school_admin")
}
