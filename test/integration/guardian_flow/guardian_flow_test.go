//go:build integration

// Package guardian_flow valida que el rol L4 guardian del seed demo
// `tutor.mendoza@edugo.test` puede autenticarse y recibe los patterns
// canónicos del rol en `ActiveContext.Grants.Allow`.
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
	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// Patterns extraídos del seed L4 `rolePermissionGrants()` para guardian.
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"content.assessments.read",
		"content.assessments.view_results",
		"content.materials.read",
		"reports.progress.read",
		"academic.grades.read",
		"academic.attendance.read",
		"academic.announcements.read",
		"academic.calendar.read",
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

	assert.Empty(t, resp.ActiveContext.Grants.Deny,
		"guardian: grants.deny must be empty")

	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for guardian")
}
