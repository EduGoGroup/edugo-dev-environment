//go:build integration

// Tests del endpoint GET /me/wards/grades de academic (plan 024 · F3 · S1).
//
// Reutiliza los helpers de roleflow (Login/GetJSON/DemoPassword/ActorModeWard)
// apuntándolos a NUESTRO identity (para emitir el token ward) y a NUESTRO
// academic (para pegar el endpoint). El Setup multi-API vive en setup_test.go.
//
// T1 — caso feliz: mendoza (login auto-resuelve a ward sofia) → 200 con la nota
//      de Sofia (Matematicas 5A); NO trae notas de otro alumno.
// T2 — revalidación por-request (el corazón de F3): con el token ward ya emitido,
//      se revoca el vínculo en la BD (status→'revoked') y se re-pega con el MISMO
//      token → 403 GUARDIAN_LINK_NOT_FOUND. Restaura el vínculo al final.
// T3 — gate de autorización: prof.martinez (teacher, contexto válido con unidad
//      pero SIN el permiso academic.my_wards_grades.read:own) → 403
//      INSUFFICIENT_PERMISSIONS (lo corta RequirePermission antes del handler).
package guardian_ward_grades_flow_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UUIDs sembrados en seeds/playground_v2/base. mendoza es guardián de sofia
// (solo S1). La nota de Sofia cuelga de su membership de alumna bb…03 en
// Matematicas dd…01 (grade a0…03, value 6.0). El alumno carlos (…0008) tiene su
// propia nota sobre membership bb…01; el endpoint NUNCA debe filtrarla.
const (
	tutorMendozaEmail = "tutor.mendoza@edugo.test"
	profMartinezEmail = "prof.martinez@edugo.test"

	guardianMendozaID = "00000000-0000-0000-0000-000000000011" // Laura Mendoza (guardián)
	studentSofiaID    = "00000000-0000-0000-0000-000000000009" // Sofia (acudida)

	sofiaMembershipID = "bb000000-0000-0000-0000-000000000003" // membership de alumna de Sofia
	carlosMembershipID = "bb000000-0000-0000-0000-000000000001" // membership de Carlos (NO debe aparecer)
	mateSubjectID     = "dd000000-0000-0000-0000-000000000001" // Matematicas (materia de la nota)
)

// myGradeListBody es el sub-set tipado del 200 (dto.MyGradeListResponse): la
// lista de notas + total. Cada nota trae membership_id/subject_id (lo que usamos
// para confirmar que es la de Sofia y no la de otro alumno).
type myGradeListBody struct {
	Grades []struct {
		ID           string `json:"id"`
		MembershipID string `json:"membership_id"`
		SubjectID    string `json:"subject_id"`
	} `json:"grades"`
	Total int `json:"total"`
}

// errorBody es el dto.ErrorResponse estándar (error/message/code).
type errorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// loginWardMendoza loguea a mendoza contra NUESTRO identity. mendoza es 1
// acudido / 1 escuela → roleflow.Login auto-resuelve a ward sofia, devolviendo
// un access_token con actor_mode=ward + subject_student_id=sofia (precondición
// del endpoint). Devuelve el token y valida la terna del contexto.
func loginWardMendoza(t *testing.T) string {
	t.Helper()
	login := roleflow.Login(t, identityServer, tutorMendozaEmail, roleflow.DemoPassword)
	require.NotNil(t, login.ActiveContext, "login de mendoza debe traer active_context (auto-resuelto)")
	require.Equal(t, roleflow.ActorModeWard, login.ActiveContext.ActorMode,
		"precondición: login de mendoza debe estar en modo ward")
	require.Equal(t, studentSofiaID, login.ActiveContext.SubjectStudentID,
		"precondición: subject_student_id debe ser Sofia")
	return login.AccessToken
}

// TestWardGrades_HappyPath_ReturnsWardGrades (T1): mendoza en contexto ward →
// GET /me/wards/grades contra academic → 200 con ≥1 nota de Sofia (Matematicas);
// ninguna nota de otro alumno (carlos).
func TestWardGrades_HappyPath_ReturnsWardGrades(t *testing.T) {
	token := loginWardMendoza(t)

	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", token)
	require.Equalf(t, http.StatusOK, status,
		"ward feliz: esperaba 200, got %d body=%s", status, string(raw))

	var body myGradeListBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 200 body=%s", string(raw))
	require.GreaterOrEqualf(t, len(body.Grades), 1,
		"el acudido (Sofia) debe tener ≥1 nota body=%s", string(raw))

	// Todas las notas devueltas deben ser de la membership de Sofia (bb…03), nunca
	// de Carlos (bb…01): el endpoint fuerza student_id = subject_student_id del JWT.
	sawSofiaMate := false
	for _, g := range body.Grades {
		assert.Equalf(t, sofiaMembershipID, g.MembershipID,
			"toda nota debe ser de la membership de Sofia; apareció %s (carlos=%s) body=%s",
			g.MembershipID, carlosMembershipID, string(raw))
		assert.NotEqual(t, carlosMembershipID, g.MembershipID,
			"NUNCA debe filtrarse una nota de Carlos (otro alumno)")
		if g.SubjectID == mateSubjectID {
			sawSofiaMate = true
		}
	}
	assert.Truef(t, sawSofiaMate,
		"debe aparecer la nota de Matematicas (%s) de Sofia body=%s", mateSubjectID, string(raw))
}

// TestWardGrades_LinkRevoked_Forbidden (T2): el corazón de F3. Con el token ward
// YA emitido, se revoca el vínculo en la BD del container (status→'revoked') y se
// re-pega el endpoint con el MISMO token → 403 GUARDIAN_LINK_NOT_FOUND, porque la
// revalidación por-request (no la firma del token) gobierna la lectura. Restaura
// el vínculo a 'active' al final para no contaminar otros subtests.
func TestWardGrades_LinkRevoked_Forbidden(t *testing.T) {
	token := loginWardMendoza(t)

	// Sanity: con el vínculo activo el endpoint responde 200 (mismo token que se
	// re-usará tras revocar — así aislamos la causa al estado del vínculo, no al token).
	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", token)
	require.Equalf(t, http.StatusOK, status,
		"precondición T2: con vínculo activo esperaba 200, got %d body=%s", status, string(raw))

	// Revoca el vínculo mendoza↔sofia directamente en la BD. La columna de estado
	// real es academic.guardian_relations.status (check IN
	// 'pending'|'active'|'rejected'|'revoked'); el use case exige status=='active' y
	// FindByGuardianAndStudent solo devuelve filas status IN ('pending','active'),
	// así que 'revoked' la oculta → ErrGuardianLinkNotFound.
	revokeLink(t)
	defer restoreLink(t) // restaura a 'active' pase lo que pase.

	// Re-pega con el MISMO token ward → 403 GUARDIAN_LINK_NOT_FOUND.
	status, raw = roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", token)
	require.Equalf(t, http.StatusForbidden, status,
		"vínculo revocado: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "GUARDIAN_LINK_NOT_FOUND", body.Code,
		"code del 403 debe ser GUARDIAN_LINK_NOT_FOUND body=%s", string(raw))
}

// TestWardGrades_NonGuardian_Forbidden (T3): prof.martinez (teacher) tiene un
// contexto válido con unidad (su login multi-escuela auto-resuelve a una unidad,
// por eso pasa RequireActiveContext) pero NO posee el permiso
// academic.my_wards_grades.read:own → RequirePermission corta con 403
// INSUFFICIENT_PERMISSIONS ANTES de llegar al handler (no es NOT_WARD_CONTEXT,
// que aplicaría a alguien con el permiso pero sin contexto ward).
func TestWardGrades_NonGuardian_Forbidden(t *testing.T) {
	login := roleflow.Login(t, identityServer, profMartinezEmail, roleflow.DemoPassword)
	require.NotNil(t, login.ActiveContext, "login de martinez debe traer active_context")

	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", login.AccessToken)
	require.Equalf(t, http.StatusForbidden, status,
		"teacher sin permiso ward: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", body.Code,
		"code del 403 (teacher sin academic.my_wards_grades.read:own) body=%s", string(raw))
}

// revokeLink marca el vínculo mendoza↔sofia como 'revoked' en la BD del
// container (acceso directo al pool del testcontainer expuesto por el Setup).
func revokeLink(t *testing.T) {
	t.Helper()
	res, err := testSQLDB.Exec(
		`UPDATE academic.guardian_relations SET status = 'revoked', is_active = false
		 WHERE guardian_id = $1 AND student_id = $2`,
		guardianMendozaID, studentSofiaID)
	require.NoError(t, err, "revocar vínculo en BD")
	n, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equalf(t, int64(1), n, "el UPDATE de revocación debió tocar exactamente 1 fila (mendoza↔sofia)")
}

// restoreLink deja el vínculo de vuelta en 'active' (idempotente) para que el
// resto de subtests vea el dato del seed.
func restoreLink(t *testing.T) {
	t.Helper()
	_, err := testSQLDB.Exec(
		`UPDATE academic.guardian_relations SET status = 'active', is_active = true
		 WHERE guardian_id = $1 AND student_id = $2`,
		guardianMendozaID, studentSofiaID)
	require.NoError(t, err, "restaurar vínculo en BD")
}
