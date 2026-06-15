//go:build integration

// Tests del flujo plan 024 · F4 · S3: la politica de representante de la escuela
// (academic.school_guardian_policy) altera la admision de un ALUMNO.
//
// T1 — on_enrollment + gates_activation (caso feliz end-to-end), sobre S1 con una
//      fila de politica TEST-LOCAL (DML de prueba, NO edita el seed):
//   - admin de San Ignacio crea una invitacion 'student' para 5to A.
//   - un usuario nuevo (el alumno) redime → join-request pending → doble-gate
//     approve. ASSERT: su membership nace 'pending' (gates_activation).
//   - ASSERT: existe una school_invitations de tipo 'guardian' con student_id=el
//     alumno (la auto-generada por on_enrollment).
//   - un usuario nuevo (el representante) redime ESA invitacion guardian →
//     join-request pending → doble-gate approve.
//   - ASSERT: hay guardian_relations (guardian↔student, school) status='active'.
//   - ASSERT: la membership del alumno ahora es 'active' (activada por el vinculo,
//     gating_approver='any').
//
// T2 — default (sin politica) = sin cambio de comportamiento, sobre S1 SIN la fila
//      test-local:
//   - admite un alumno (invitacion student → redencion → doble-gate approve).
//   - ASSERT: su membership nace 'active' de inmediato.
//   - ASSERT: NO se creo ninguna school_invitations de tipo guardian para ese alumno.
//
// DECISION DE APROBADOR (documentada): se usa el camino (ii) del prompt — S1 con una
// fila school_guardian_policy test-local. S1 tiene un school_admin sembrado
// (admin.sanignacio, academic.*) cuyo contexto (S1, unidad de la solicitud) es
// alcanzable via switch-context y firma AMBOS sellos del doble-gate, exactamente como
// ya prueba guardian_invitation_flow. S3 (que SI trae la fila de politica en el seed)
// no tiene un school_admin sembrado con contexto alcanzable, asi que aprobar alli
// seria fragil. El codigo ejercido es el mismo: el resolver de politica lee la fila
// efectiva de la escuela de la solicitud, sea sembrada (S3) o test-local (S1).
//
// Los usuarios alumno/representante se crean frescos via /auth/signup + un user_role
// global con rol 'student' (que tiene acceso al sistema kmp, requisito del login-gate
// MP-08 DEC-C). Al ser nuevos NO tienen membership ni guardian_relations previas, asi
// que el redeem produce ambos gates pending y no colisionan con indices unicos.
//
// Las filas insertadas se limpian con defer (la BD del container es compartida por
// subtests del process).
package guardian_onboarding_policy_flow_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UUIDs/emails sembrados en seeds/playground_v2/base.
const (
	adminSanIgnacioEmail = "admin.sanignacio@edugo.test" // Carmen Valdes, school_admin S1 (academic.*)

	schoolS1ID = "b1000000-0000-0000-0000-000000000001" // Colegio San Ignacio (S1)
	unit5toA   = "ac000000-0000-0000-0000-000000000003" // 5to A (unidad de la admision)

	// Rol con acceso al sistema kmp; se concede a los usuarios frescos para que
	// pasen el login-gate por sistema (MP-08 DEC-C). NO crea membership ni vinculo.
	roleStudentID = "b4000000-0001-0000-0000-000000000001" // L4_ROLE_STUDENT_ID

	// Fila de politica test-local para S1 (Test 1). id fuera del rango sembrado.
	testLocalPolicyS1ID = "9b000000-0000-0000-0000-0000000000f1"

	defaultPassword = "12345678"
)

// --- helpers HTTP ---

// postJSON ejecuta un POST con bearer + body JSON y devuelve (status, body).
func postJSON(t *testing.T, server *httptest.Server, path, bearer string, payload any) (int, []byte) {
	t.Helper()
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, raw
}

// --- tipos del contrato (sub-sets de los DTO) ---

type invitationResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	SchoolID       string `json:"school_id"`
	AcademicUnitID string `json:"academic_unit_id"`
	InvitationType string `json:"invitation_type"`
}

type pendingJoinRequest struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	SchoolID       string `json:"school_id"`
	AcademicUnitID string `json:"academic_unit_id"`
	InvitationType string `json:"invitation_type"`
	Status         string `json:"status"`
}

// adminTokenOn5toA loguea al admin de S1 y deja su contexto activo en (S1, 5to A)
// con grants academic.* — necesario para que el approve firme AMBOS sellos
// (school.<tipo> + unit.<tipo>) sobre la unidad de la solicitud.
func adminTokenOn5toA(t *testing.T) string {
	t.Helper()
	login := roleflow.Login(t, identityServer, adminSanIgnacioEmail, roleflow.DemoPassword)
	require.NotNil(t, login.ActiveContext, "login admin debe traer active_context")

	status, raw := roleflow.SwitchSubject(t, identityServer, login.AccessToken, "self", schoolS1ID, unit5toA)
	require.Equalf(t, http.StatusOK, status,
		"switch-context admin a (S1, 5to A): esperaba 200, got %d body=%s", status, string(raw))

	var sc struct {
		AccessToken string `json:"access_token"`
		Context     *struct {
			SchoolID       string `json:"school_id"`
			AcademicUnitID string `json:"academic_unit_id"`
		} `json:"context"`
	}
	require.NoError(t, json.Unmarshal(raw, &sc), "parse switch-context body=%s", string(raw))
	require.NotNil(t, sc.Context, "switch-context: context nil body=%s", string(raw))
	require.Equal(t, schoolS1ID, sc.Context.SchoolID, "contexto debe quedar en S1")
	require.Equal(t, unit5toA, sc.Context.AcademicUnitID, "contexto debe quedar en 5to A")
	require.NotEmpty(t, sc.AccessToken)
	return sc.AccessToken
}

// signupKmpUser crea un usuario nuevo via /auth/signup y le concede un user_role
// global con rol 'student' (que tiene acceso al sistema kmp). Devuelve (userID,
// token de login). El usuario NO tiene membership ni guardian_relations, asi que el
// redeem produce ambos gates pending.
func signupKmpUser(t *testing.T, label string) (userID, token string) {
	t.Helper()
	email := fmt.Sprintf("onboarding.%s.%s@edugo.test", label, uuid.NewString()[:8])

	status, raw := postJSON(t, identityServer, "/api/v1/auth/signup", "", map[string]any{
		"email":      email,
		"password":   defaultPassword,
		"first_name": "Onboarding",
		"last_name":  label,
	})
	require.Equalf(t, http.StatusCreated, status,
		"signup %s: esperaba 201, got %d body=%s", email, status, string(raw))

	var su struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.Unmarshal(raw, &su), "parse signup body=%s", string(raw))
	require.NotEmpty(t, su.ID, "signup debe devolver id")

	// Concede acceso al sistema kmp via un user_role global (rol student). Sin esto
	// el login devuelve 403 SYSTEM_ACCESS_DENIED (MP-08 DEC-C). school_id NULL ⇒
	// scope_pattern '*', sin membership ni vinculo. granted_at lo exige el schema.
	grantID := uuid.NewString()
	_, err := testSQLDB.Exec(
		`INSERT INTO iam.user_roles (id, user_id, role_id, school_id, academic_unit_id, scope_pattern, is_active, granted_at, created_at, updated_at)
		 VALUES ($1, $2, $3, NULL, NULL, '*', true, now(), now(), now())`,
		grantID, su.ID, roleStudentID)
	require.NoErrorf(t, err, "conceder user_role kmp a %s", email)

	login := roleflow.LoginRaw(t, identityServer, email, defaultPassword)
	require.NotEmpty(t, login.AccessToken, "%s debe poder loguear tras la concesion de rol", email)
	return su.ID, login.AccessToken
}

// admitStudent ejecuta el camino comun de admision de un ALUMNO: el admin crea una
// invitacion 'student' para 5to A, un usuario fresco la redime (join-request pending)
// y el admin la aprueba firmando ambos sellos. Devuelve el userID del alumno admitido.
func admitStudent(t *testing.T, adminToken, label string) string {
	t.Helper()

	status, raw := postJSON(t, academicServer, "/api/v1/schools/invitations", adminToken, map[string]any{
		"academic_unit_id": unit5toA,
		"invitation_type":  "student",
	})
	require.Equalf(t, http.StatusCreated, status,
		"crear invitacion student: esperaba 201, got %d body=%s", status, string(raw))

	var inv invitationResponse
	require.NoError(t, json.Unmarshal(raw, &inv), "parse invitacion body=%s", string(raw))
	require.NotEmpty(t, inv.Code, "la invitacion debe traer code")
	require.Equal(t, "student", inv.InvitationType)
	t.Cleanup(func() { cleanupInvitation(t, inv.ID) })

	studentID, studentToken := signupKmpUser(t, label)

	status, raw = postJSON(t, academicServer, "/api/v1/invitations/redeem", studentToken, map[string]any{
		"code": inv.Code,
	})
	require.Equalf(t, http.StatusCreated, status,
		"redimir code student: esperaba 201, got %d body=%s", status, string(raw))

	jrID := findPendingJoinRequestID(t, adminToken, studentID)
	require.NotEmpty(t, jrID, "debe existir una solicitud pendiente del alumno")
	t.Cleanup(func() { cleanupJoinRequest(t, jrID) })

	status, raw = postJSON(t, academicServer, "/api/v1/schools/join-requests/"+jrID+"/approve", adminToken, map[string]any{})
	require.Equalf(t, http.StatusOK, status,
		"aprobar admision del alumno: esperaba 200, got %d body=%s", status, string(raw))
	var approved struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(raw, &approved), "parse approve body=%s", string(raw))
	require.Equalf(t, "approved", approved.Status,
		"la admision del alumno debe quedar approved tras firmar ambos sellos body=%s", string(raw))

	t.Cleanup(func() { cleanupMembership(t, studentID, schoolS1ID) })
	return studentID
}

// TestOnEnrollmentGatesActivation_HappyPath (T1).
func TestOnEnrollmentGatesActivation_HappyPath(t *testing.T) {
	// Politica test-local en S1: on_enrollment + gates + any (espejo de la de S3).
	insertS1Policy(t)
	t.Cleanup(func() { deleteS1Policy(t) })

	adminToken := adminTokenOn5toA(t)

	// 1-3. Admite al alumno (invitacion → redencion → doble-gate approve).
	studentID := admitStudent(t, adminToken, "student")

	// 3. ASSERT: la membership del alumno nace 'pending' (gates_activation=true).
	assertMembershipStatus(t, studentID, schoolS1ID, "pending",
		"con gates_activation la membership del alumno debe nacer pending")

	// 4. ASSERT: existe una invitacion guardian auto-generada apuntando al alumno
	//    (on_enrollment). Sin email/codigo conocido: se localiza por student_id+tipo.
	guardianInvCode := assertGuardianInvitationForStudent(t, studentID, schoolS1ID)
	t.Cleanup(func() { cleanupInvitationByStudent(t, studentID) })

	// 5. Un usuario nuevo (el representante) redime ESA invitacion guardian →
	//    join-request pending → doble-gate approve.
	guardianID, guardianToken := signupKmpUser(t, "guardian")

	status, raw := postJSON(t, academicServer, "/api/v1/invitations/redeem", guardianToken, map[string]any{
		"code": guardianInvCode,
	})
	require.Equalf(t, http.StatusCreated, status,
		"redimir code guardian: esperaba 201, got %d body=%s", status, string(raw))

	jrID := findPendingJoinRequestID(t, adminToken, guardianID)
	require.NotEmpty(t, jrID, "debe existir una solicitud pendiente del representante")
	t.Cleanup(func() { cleanupJoinRequest(t, jrID) })

	status, raw = postJSON(t, academicServer, "/api/v1/schools/join-requests/"+jrID+"/approve", adminToken, map[string]any{})
	require.Equalf(t, http.StatusOK, status,
		"aprobar admision del representante: esperaba 200, got %d body=%s", status, string(raw))
	var approved struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(raw, &approved), "parse approve guardian body=%s", string(raw))
	require.Equalf(t, "approved", approved.Status,
		"la admision del representante debe quedar approved body=%s", string(raw))
	t.Cleanup(func() { cleanupGuardianRelation(t, guardianID, studentID) })

	// 6. ASSERT: el vinculo guardian↔alumno existe, scoped a S1, status='active'.
	assertGuardianRelationActive(t, guardianID, studentID, schoolS1ID)

	// 7. ASSERT: la membership del alumno ahora es 'active' (activada por el vinculo,
	//    gating_approver='any').
	assertMembershipStatus(t, studentID, schoolS1ID, "active",
		"tras aprobar al representante (any) la membership gateada del alumno debe pasar a active")
}

// TestDefaultPolicy_NoChange (T2): S1 SIN fila de politica = comportamiento de hoy.
func TestDefaultPolicy_NoChange(t *testing.T) {
	// Defensa: asegura que NO queda la fila test-local del T1 (orden de tests no
	// garantizado). deleteS1Policy es idempotente.
	deleteS1Policy(t)

	adminToken := adminTokenOn5toA(t)

	// 1. Admite un alumno en S1 (defaults).
	studentID := admitStudent(t, adminToken, "default")

	// 2. ASSERT: la membership nace 'active' de inmediato (sin gating).
	assertMembershipStatus(t, studentID, schoolS1ID, "active",
		"sin politica la membership del alumno debe nacer active (comportamiento de hoy)")

	// 3. ASSERT: NO se creo ninguna invitacion guardian para ese alumno
	//    (invitation_mode default = 'manual' ≠ 'on_enrollment').
	assertNoGuardianInvitationForStudent(t, studentID, schoolS1ID)
}

// --- helpers de query/lectura ---

// findPendingJoinRequestID lista GET /schools/join-requests/pending (contexto del
// admin) y devuelve el ID de la solicitud cuyo user_id == userID.
func findPendingJoinRequestID(t *testing.T, adminToken, userID string) string {
	t.Helper()
	status, raw := roleflow.GetJSON(t, academicServer, "/api/v1/schools/join-requests/pending", adminToken)
	require.Equalf(t, http.StatusOK, status,
		"listar pendientes: esperaba 200, got %d body=%s", status, string(raw))

	var list struct {
		Requests []pendingJoinRequest `json:"requests"`
	}
	require.NoError(t, json.Unmarshal(raw, &list), "parse pendientes body=%s", string(raw))
	for _, r := range list.Requests {
		if r.UserID == userID {
			return r.ID
		}
	}
	return ""
}

// assertMembershipStatus exige una membership del usuario en la escuela con el
// status esperado.
func assertMembershipStatus(t *testing.T, userID, schoolID, want, msg string) {
	t.Helper()
	var got string
	err := testSQLDB.QueryRow(
		`SELECT status FROM academic.memberships WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID).Scan(&got)
	require.NoErrorf(t, err, "debe existir UNA membership user=%s school=%s", userID, schoolID)
	assert.Equalf(t, want, got, "%s (got status=%q)", msg, got)
}

// assertGuardianInvitationForStudent exige una school_invitations de tipo guardian
// con student_id=studentID, activa, y devuelve su code (para que el representante la
// redima). La auto-generada por on_enrollment no trae email; se localiza por
// (student_id, tipo, school).
func assertGuardianInvitationForStudent(t *testing.T, studentID, schoolID string) string {
	t.Helper()
	var code string
	err := testSQLDB.QueryRow(
		`SELECT si.code
		   FROM academic.school_invitations si
		   JOIN academic.invitation_types it ON it.id = si.invitation_type_id
		  WHERE si.student_id = $1 AND si.school_id = $2 AND it.key = 'guardian' AND si.is_active = true`,
		studentID, schoolID).Scan(&code)
	require.NoErrorf(t, err,
		"on_enrollment: debe existir UNA invitacion guardian auto-generada con student_id=%s en %s", studentID, schoolID)
	require.NotEmpty(t, code, "la invitacion guardian auto-generada debe traer code")
	return code
}

// assertNoGuardianInvitationForStudent exige que NO exista ninguna school_invitations
// de tipo guardian apuntando a ese alumno (default = sin on_enrollment).
func assertNoGuardianInvitationForStudent(t *testing.T, studentID, schoolID string) {
	t.Helper()
	var count int
	err := testSQLDB.QueryRow(
		`SELECT COUNT(*)
		   FROM academic.school_invitations si
		   JOIN academic.invitation_types it ON it.id = si.invitation_type_id
		  WHERE si.student_id = $1 AND si.school_id = $2 AND it.key = 'guardian'`,
		studentID, schoolID).Scan(&count)
	require.NoError(t, err, "contar invitaciones guardian del alumno")
	assert.Equalf(t, 0, count,
		"sin on_enrollment NO debe auto-generarse invitacion guardian para el alumno %s (got %d)", studentID, count)
}

// assertGuardianRelationActive lee academic.guardian_relations y exige una fila
// guardian↔student scoped a school con status='active' e is_active=true.
func assertGuardianRelationActive(t *testing.T, guardianID, studentID, schoolID string) {
	t.Helper()
	var status string
	var isActive bool
	err := testSQLDB.QueryRow(
		`SELECT status, is_active FROM academic.guardian_relations
		 WHERE guardian_id = $1 AND student_id = $2 AND school_id = $3
		 AND relationship_type = 'guardian'`,
		guardianID, studentID, schoolID).Scan(&status, &isActive)
	require.NoErrorf(t, err,
		"debe existir UNA fila guardian_relations guardian↔alumno en %s (creada al aprobar)", schoolID)
	assert.Equal(t, "active", status, "el vinculo debe nacer active (la aprobacion del colegio ES la autorizacion)")
	assert.True(t, isActive, "is_active debe ser true para un vinculo active")
}

// --- DML de prueba (politica test-local) ---

// insertS1Policy inserta la fila school_guardian_policy test-local de S1 (default de
// escuela, academic_unit_id NULL): on_enrollment + gates + any. Es DML de prueba, NO
// edita el archivo seed.
func insertS1Policy(t *testing.T) {
	t.Helper()
	_, err := testSQLDB.Exec(
		`INSERT INTO academic.school_guardian_policy
		   (id, school_id, academic_unit_id, invitation_mode, gates_activation, gating_approver, link_scope, created_at, updated_at)
		 VALUES ($1, $2, NULL, 'on_enrollment', true, 'any', 'school', now(), now())
		 ON CONFLICT (id) DO NOTHING`,
		testLocalPolicyS1ID, schoolS1ID)
	require.NoError(t, err, "insertar politica test-local de S1")
}

// deleteS1Policy elimina la fila test-local de S1. Idempotente.
func deleteS1Policy(t *testing.T) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_guardian_policy WHERE id = $1`, testLocalPolicyS1ID)
	require.NoError(t, err, "borrar politica test-local de S1")
}

// --- limpieza (defer) ---

func cleanupInvitation(t *testing.T, invitationID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_invitations WHERE id = $1`, invitationID)
	assert.NoError(t, err, "cleanup invitacion")
}

func cleanupInvitationByStudent(t *testing.T, studentID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_invitations WHERE student_id = $1`, studentID)
	assert.NoError(t, err, "cleanup invitacion auto-generada del alumno")
}

func cleanupJoinRequest(t *testing.T, joinRequestID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_join_requests WHERE id = $1`, joinRequestID)
	assert.NoError(t, err, "cleanup join-request")
}

func cleanupMembership(t *testing.T, userID, schoolID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.memberships WHERE user_id = $1 AND school_id = $2`, userID, schoolID)
	assert.NoError(t, err, "cleanup membership")
}

func cleanupGuardianRelation(t *testing.T, guardianID, studentID string) {
	t.Helper()
	_, err := testSQLDB.Exec(
		`DELETE FROM academic.guardian_relations
		 WHERE guardian_id = $1 AND student_id = $2 AND relationship_type = 'guardian'`,
		guardianID, studentID)
	assert.NoError(t, err, "cleanup guardian_relation")
}
