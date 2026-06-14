//go:build integration

// Package readonly_auditor_flow valida que el rol alias L4
// `readonly_auditor` del seed base `readonly@edugo.test` hereda los
// patterns de teacher pero filtrados a sólo acciones de lectura. Los
// patterns con verbos mutativos (create/update/delete/publish/grade/
// finalize/assign/review/manage/request/activate/attempt) NO deben
// aparecer en `Grants.Allow`.
package readonly_auditor_flow_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	roleName  = "readonly_auditor"
	userEmail = "readonly@edugo.test"
)

func TestMain(m *testing.M) {
	os.Exit(roleflow.Setup(m))
}

func TestReadonlyAuditorFlow_Grants(t *testing.T) {
	env := roleflow.Get()

	resp := roleflow.Login(t, env.Server, userEmail, roleflow.DemoPassword)
	require.NotNil(t, resp.ActiveContext, "active_context must be present")
	assert.Equal(t, roleName, resp.ActiveContext.RoleName)

	// readonly_auditor hereda de teacher con filterMutationGrants(): sólo
	// quedan acciones de lectura, view, download, etc.
	roleflow.AssertGrantsContains(t, resp.ActiveContext.Grants,
		"content.assessments.read",
		"content.materials.read",
		"content.materials.download",
		"academic.grades.read",
		"academic.attendance.read",
		"academic.schedules.read",
		"academic.announcements.read",
		"academic.periods.read",
		"reports.progress.read",
		"admin.users.read",
		"admin.users.read:own",
		"academic.subjects.read",
		"academic.units.read",
		"reports.stats.unit",
		"dashboard.view",
		"menu.read",
		"notifications.read",
		"reports.read",
		"academic.calendar.read",
	)

	// Validación negativa: el rol no debe poder mutar nada. Pass 3
	// wildcard-first: `readonly_auditor` recibe allow amplio
	// (`academic.*`, `content.*`, ...) + denies `*.suffix` para verbos
	// mutativos. La aserción usa el matcher real (deny > allow).
	mutationProbes := []string{
		"academic.announcements.create",
		"academic.grades.update",
		"academic.attendance.delete",
		"content.materials.publish",
		"content.assessments.grade",
		"academic.guardian_relations.approve",
		"academic.guardian_relations.request",
		"academic.guardian_relations.manage",
		"content.assessments.review",
		"content.assessments.assign",
		"content.assessments.attempt",
		"academic.periods.activate",
		"academic.grades.finalize",
	}
	for _, probe := range mutationProbes {
		assert.Falsef(t, roleflow.GrantsAllow(resp.ActiveContext.Grants, probe),
			"readonly_auditor must NOT be granted mutation %q", probe)
	}

	status, _ := roleflow.GetJSON(t, env.Server,
		"/api/v1/auth/contexts", resp.AccessToken)
	assert.Equal(t, 200, status,
		"GET /auth/contexts must return 200 for readonly_auditor")
}
