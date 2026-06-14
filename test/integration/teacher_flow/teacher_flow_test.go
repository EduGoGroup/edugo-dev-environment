//go:build integration

// Package teacher_flow valida que el rol L4 teacher del seed base
// `prof.martinez@edugo.test` puede autenticarse y recibe los patterns
// canónicos del rol en `ActiveContext.Grants.Allow`.
package teacher_flow_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	roleName  = "teacher"
	userEmail = "prof.martinez@edugo.test"
)

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

func TestTeacherFlow_Grants(t *testing.T) {
	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// Patterns extraídos del seed L4 `rolePermissionGrants()` para teacher.
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"content.assessments.create",
		"content.assessments.read",
		"content.assessments.update",
		"content.assessments.publish",
		"content.assessments.grade",
		"content.assessments.assign",
		"content.assessments.review",
		"content.materials.read",
		"content.materials.create",
		"content.materials.update",
		"content.materials.download",
		"content.materials.publish",
		"academic.grades.read",
		"academic.grades.create",
		"academic.grades.update",
		"academic.grades.finalize",
		"academic.attendance.read",
		"academic.attendance.create",
		"academic.attendance.update",
		"academic.announcements.read",
		"academic.announcements.create",
		"academic.periods.read",
		"reports.progress.read",
		"reports.progress.update",
		"admin.users.read",
		"admin.users.read:own",
		"admin.users.update:own",
		"academic.subjects.read",
		"academic.units.read",
		"reports.stats.unit",
		"dashboard.view",
		"menu.read",
		"notifications.read",
		"admin.system_settings.read",
	)

	assert.Empty(t, resp.ActiveContext.Grants.Deny,
		"teacher: grants.deny must be empty")

	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for teacher")
}
