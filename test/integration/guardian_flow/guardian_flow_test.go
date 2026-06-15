//go:build integration

// Package guardian_flow valida que el rol L4 guardian
// `tutor.mendoza@edugo.test` puede autenticarse y recibe los patterns
// canónicos del rol en `ActiveContext.Grants.Allow`.
//
// Plan 024·F1 sembró el dato (usuario tutor + guardian_relations) y reescribió
// el contrato de privacidad: el guardián ya NO recibe el wildcard
// `academic.grades.*`; ve a su acudido solo vía `academic.my_wards_*.read:own`
// (sin `academic.calendar.*`). Las assertions de abajo ya reflejan ese modelo.
//
// SIGUE EN SKIP hasta 024·F2: el guardián NO lleva membership (DEC-R-B) y la
// resolución de `ActiveContext` en el login deriva las escuelas de
// `academic.memberships` → un guardián sin membership no resuelve contexto.
// Esa lógica (contexto del representante) es F2 (identity). Al cerrar F2,
// quitar el t.Skip de TestGuardianFlow_Grants: el test debe pasar tal cual.
package guardian_flow_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	roleName  = "guardian"
	userEmail = "tutor.mendoza@edugo.test"
)

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

func TestGuardianFlow_Grants(t *testing.T) {
	t.Skip("024·F2: el guardián no lleva membership (DEC-R-B) y la resolución de " +
		"ActiveContext deriva de academic.memberships; el login del representante se " +
		"resuelve en F2 (identity). F1 ya dejó seed+permisos+contrato listos.")

	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// Patterns extraídos del seed L4 `rolePermissionGrants()` para guardian
	// según el modelo de privacidad del plan 024·F1: el guardián ve las notas,
	// asistencia, anuncios y materiales de su acudido vía
	// `academic.my_wards_*.read:own` (NO el wildcard `academic.grades.*`),
	// y F1 ya no concede `academic.calendar.*`.
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"content.assessments.read",
		"content.assessments.view_results",
		"content.materials.read",
		"reports.progress.read",
		"academic.my_wards_grades.read:own",
		"academic.my_wards_attendance.read:own",
		"academic.my_wards_announcements.read:own",
		"academic.my_wards_materials.read:own",
		"academic.attendance.read",
		"academic.announcements.read",
		"academic.guardian_relations.read",
		"academic.guardian_relations.request",
		"admin.users.read:own",
		"admin.users.update:own",
		"dashboard.view",
		"screens.read",
		"menu.read",
		"notifications.read",
		"reports.read",
		"admin.system_settings.read",
	)

	// Privacidad F1: el guardián NO ve el wildcard de notas; solo a su acudido
	// vía academic.my_wards_*.read:own.
	assert.False(t, roleflow.GrantsAllow(resp.ActiveContext.Grants, "academic.grades.read"),
		"guardian no debe tener academic.grades.* (privacidad plan 024 F1)")

	assert.Empty(t, resp.ActiveContext.Grants.Deny,
		"guardian: grants.deny must be empty")

	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for guardian")
}
