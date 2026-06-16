//go:build integration

// Tests del endpoint GET /me/wards/attendance de academic (plan 024 · F3 · S2).
//
// Comparte el Setup multi-API de setup_test.go (mismo paquete que el flow de
// grades): UN testcontainer postgres con identity + academic in-process sobre la
// MISMA gorm.DB y el MISMO AUTH_JWT_SECRET, sembrado con migrations + base. El
// token ward lo emite IDENTITY (login de mendoza auto-resuelve a su acudida
// Sofia) y se pega contra ACADEMIC, igual que el flow de grades.
//
// T1 — caso feliz: mendoza (login auto-resuelve a ward sofia) → 200 con ≥1
//      registro de asistencia de Sofia; ninguno de otro alumno (carlos). El seed
//      base SÍ trae asistencia para Sofia (3 filas en Matematicas 5A), así que
//      tomamos el camino "SÍ hay": verificamos por SELECT al pool y afirmamos
//      Total>=1 + records de la membership de Sofia. (El camino "NO hay" inserta 1
//      fila con defer de limpieza; queda como red de seguridad si el seed cambia.)
// T2 — revalidación por-request (el corazón de F3): con el token ward ya emitido,
//      se revoca el vínculo en la BD (status→'revoked') y se re-pega con el MISMO
//      token → 403 GUARDIAN_LINK_NOT_FOUND. Restaura el vínculo al final.
// T3 — gate de autorización: prof.martinez (teacher, contexto válido con unidad
//      pero SIN el permiso academic.my_wards_attendance.read:own) → 403
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

// UUIDs sembrados en seeds/playground_v2/base relevantes a la asistencia de
// Sofia. Su asistencia cuelga de la misma membership de alumna (bb…03) en la
// materia Matematicas (dd…01), unidad 5to A. El alumno Carlos (bb…01) también
// tiene asistencia sembrada en la misma sección; el endpoint NUNCA debe filtrarla.
// El recorded_by de las filas sembradas es la admin (…0005), válido como autor.
const (
	wardAttRecorderID = "00000000-0000-0000-0000-000000000005" // admin que registró la asistencia sembrada
	// wardAttInsertID/Date solo se usan en el camino "NO hay" (fallback): id y
	// fecha que no colisionan con la UNIQUE (membership_id, subject_id, date) del
	// seed (Sofia trae 2026-03-17/18/19 en Matematicas).
	wardAttInsertID   = "a1000000-0000-0000-0000-0000000000ff"
	wardAttInsertDate = "2026-03-25"
)

// attendanceListBody es el sub-set tipado del 200 (dto.AttendanceListResponse):
// la lista de registros + total/page/limit. Cada record trae membership_id (lo
// que usamos para confirmar que es de Sofia y no de otro alumno).
type attendanceListBody struct {
	Records []struct {
		ID           string `json:"id"`
		MembershipID string `json:"membership_id"`
		SubjectID    string `json:"subject_id"`
		Date         string `json:"date"`
		Status       string `json:"status"`
	} `json:"records"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// TestWardAttendance_HappyPath_ReturnsWardAttendance (T1): mendoza en contexto
// ward → GET /me/wards/attendance contra academic → 200 con ≥1 registro de Sofia
// (membership bb…03); ninguno de otro alumno (carlos, bb…01). El endpoint fuerza
// student_id = subject_student_id del JWT, así que el filtrado es del backend.
func TestWardAttendance_HappyPath_ReturnsWardAttendance(t *testing.T) {
	token := loginWardMendoza(t)

	// PRIMERO decidimos el camino según el dato real del seed: ¿tiene Sofia
	// asistencia? Contamos sus filas en academic.attendance vía la membership de
	// alumna. (En base hay 3; el bloque de inserción es la red de seguridad por si
	// el seed cambiara y el acudido quedara sin asistencia.)
	if sofiaAttendanceCount(t) == 0 {
		// Camino "NO hay": el seed base no tiene asistencia para el acudido →
		// insertamos 1 fila para Sofia (membership bb…03 + subject Matematicas dd…01
		// + recorded_by admin …0005), con fecha que no colisiona con la UNIQUE del
		// seed. defer limpia para no contaminar otros subtests.
		insertSofiaAttendance(t)
		defer deleteSofiaAttendance(t)
	}

	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/attendance", token)
	require.Equalf(t, http.StatusOK, status,
		"ward feliz: esperaba 200, got %d body=%s", status, string(raw))

	var body attendanceListBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 200 body=%s", string(raw))
	require.GreaterOrEqualf(t, body.Total, 1,
		"el acudido (Sofia) debe tener ≥1 registro de asistencia body=%s", string(raw))
	require.GreaterOrEqualf(t, len(body.Records), 1,
		"el acudido (Sofia) debe traer ≥1 record body=%s", string(raw))

	// Todos los records devueltos deben ser de la membership de Sofia (bb…03),
	// nunca de Carlos (bb…01): el endpoint fuerza student_id = subject_student_id
	// del JWT, no acepta filtro por query.
	for _, r := range body.Records {
		assert.Equalf(t, sofiaMembershipID, r.MembershipID,
			"todo registro debe ser de la membership de Sofia; apareció %s (carlos=%s) body=%s",
			r.MembershipID, carlosMembershipID, string(raw))
		assert.NotEqual(t, carlosMembershipID, r.MembershipID,
			"NUNCA debe filtrarse asistencia de Carlos (otro alumno)")
	}
}

// TestWardAttendance_LinkRevoked_Forbidden (T2): el corazón de F3. Con el token
// ward YA emitido, se revoca el vínculo en la BD del container (status→'revoked')
// y se re-pega el endpoint con el MISMO token → 403 GUARDIAN_LINK_NOT_FOUND,
// porque la revalidación por-request (no la firma del token) gobierna la lectura.
// Restaura el vínculo a 'active' al final para no contaminar otros subtests.
func TestWardAttendance_LinkRevoked_Forbidden(t *testing.T) {
	token := loginWardMendoza(t)

	// Sanity: con el vínculo activo el endpoint responde 200 (mismo token que se
	// re-usará tras revocar — así aislamos la causa al estado del vínculo, no al token).
	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/attendance", token)
	require.Equalf(t, http.StatusOK, status,
		"precondición T2: con vínculo activo esperaba 200, got %d body=%s", status, string(raw))

	// Revoca el vínculo mendoza↔sofia directamente en la BD. FindByGuardianAndStudent
	// solo devuelve filas status IN ('pending','active') y el use case exige
	// status=='active', así que 'revoked' la oculta → ErrGuardianLinkNotFound.
	// revokeLink/restoreLink viven en guardian_ward_grades_flow_test.go (mismo paquete).
	revokeLink(t)
	defer restoreLink(t) // restaura a 'active' pase lo que pase.

	// Re-pega con el MISMO token ward → 403 GUARDIAN_LINK_NOT_FOUND.
	status, raw = roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/attendance", token)
	require.Equalf(t, http.StatusForbidden, status,
		"vínculo revocado: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "GUARDIAN_LINK_NOT_FOUND", body.Code,
		"code del 403 debe ser GUARDIAN_LINK_NOT_FOUND body=%s", string(raw))
}

// TestWardAttendance_NonGuardian_Forbidden (T3): prof.martinez (teacher) tiene un
// contexto válido con unidad (su login multi-escuela auto-resuelve a una unidad,
// por eso pasa RequireActiveContext) pero NO posee el permiso
// academic.my_wards_attendance.read:own → RequirePermission corta con 403
// INSUFFICIENT_PERMISSIONS ANTES de llegar al handler (no es NOT_WARD_CONTEXT,
// que aplicaría a alguien con el permiso pero sin contexto ward).
func TestWardAttendance_NonGuardian_Forbidden(t *testing.T) {
	login := roleflow.Login(t, identityServer, profMartinezEmail, roleflow.DemoPassword)
	require.NotNil(t, login.ActiveContext, "login de martinez debe traer active_context")

	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/attendance", login.AccessToken)
	require.Equalf(t, http.StatusForbidden, status,
		"teacher sin permiso ward: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", body.Code,
		"code del 403 (teacher sin academic.my_wards_attendance.read:own) body=%s", string(raw))
}

// sofiaAttendanceCount cuenta las filas de asistencia del acudido (Sofia) vía su
// membership de alumna, consultando directamente el pool del testcontainer
// expuesto por el Setup. Es el SELECT que decide el camino de T1.
func sofiaAttendanceCount(t *testing.T) int {
	t.Helper()
	var n int
	err := testSQLDB.QueryRow(
		`SELECT count(*) FROM academic.attendance a
		 JOIN academic.memberships m ON m.id = a.membership_id
		 WHERE m.user_id = $1`,
		studentSofiaID).Scan(&n)
	require.NoError(t, err, "contar asistencia de Sofia en BD")
	return n
}

// insertSofiaAttendance inserta 1 fila de asistencia para Sofia en la BD del
// container. Solo se invoca en el camino "NO hay" (fallback) de T1: el seed base
// hoy SÍ trae asistencia para el acudido, pero si dejara de traerla este INSERT
// garantiza el ≥1 que el endpoint debe devolver. La fecha 2026-03-25 evita
// colisionar con la UNIQUE (membership_id, subject_id, date) del seed.
func insertSofiaAttendance(t *testing.T) {
	t.Helper()
	res, err := testSQLDB.Exec(
		`INSERT INTO academic.attendance (id, membership_id, subject_id, date, status, recorded_by)
		 VALUES ($1, $2, $3, $4, 'present', $5)`,
		wardAttInsertID, sofiaMembershipID, mateSubjectID, wardAttInsertDate, wardAttRecorderID)
	require.NoError(t, err, "insertar asistencia de Sofia en BD")
	n, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equalf(t, int64(1), n, "el INSERT de asistencia debió tocar exactamente 1 fila")
}

// deleteSofiaAttendance borra la fila insertada por el fallback de T1 (idempotente
// por id) para no contaminar otros subtests.
func deleteSofiaAttendance(t *testing.T) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.attendance WHERE id = $1`, wardAttInsertID)
	require.NoError(t, err, "limpiar asistencia insertada de Sofia")
}
