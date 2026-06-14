//go:build integration

// Package student_flow valida que el rol L4 student del seed base
// `est.carlos@edugo.test` puede autenticarse y recibe los patterns
// canónicos del rol en `ActiveContext.Grants.Allow`.
package student_flow_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	roleName  = "student"
	userEmail = "est.carlos@edugo.test"
)

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

func TestStudentFlow_Grants(t *testing.T) {
	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// Patterns extraídos del seed L4 `rolePermissionGrants()` para student.
	// Notas del alumno: el rol ya NO recibe `academic.grades.read` (CRUD
	// docente — fuga de privacidad, podado en N3 F4.1); su lectura de notas
	// propias es el feature self `academic.my_grades.read:own`. Los recursos
	// `academic.schedules`/`academic.calendar` se eliminaron del producto.
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"content.assessments.attempt",
		"content.assessments.read",
		"content.assessments_student.read",
		"content.assessments.view_results",
		"content.materials.read",
		"content.materials.download",
		"reports.progress.read:own",
		"academic.my_grades.read:own",
		"academic.my_memberships.read:own",
		"academic.attendance.read",
		"academic.announcements.read",
		"dashboard.view",
		"screens.read",
		"menu.read",
		"notifications.read",
		"admin.system_settings.read",
	)

	// El seed base aplica un deny puntual a est.carlos sobre
	// `academic.grades.read` (override prohibitivo de lectura de notas, el
	// mismo dato que valida user_grants_flow). No es deny vacío.
	assert.Contains(t, resp.ActiveContext.Grants.Deny, "academic.grades.read",
		"student est.carlos: base siembra deny puntual sobre academic.grades.read")

	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for student")
}
