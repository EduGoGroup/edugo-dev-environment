//go:build integration

// Tests del flujo plan 024 · F4 · S2 (invitación de representante con selección
// de alumno → vínculo al aprobar).
//
// T1 — caso feliz end-to-end:
//   - admin de San Ignacio (con contexto S1 + unidad de Sofia) crea una
//     invitación tipo "guardian" apuntando a Sofia → obtiene el code.
//   - un representante SIN vínculo a Sofia (tutor.castro) redime el code →
//     join-request; el doble-gate (school+unit) se firma → approved.
//   - ASERCIÓN CENTRAL: existe una fila en academic.guardian_relations con
//     guardian_id=castro, student_id=Sofia, school_id=S1, status='active'.
//   - lazo F2/F3: el representante hace switch-context al hijo (Sofia) →
//     GET /me/wards/grades devuelve 200 (el vínculo recién creado lo habilita).
//
// T2 — caso negativo: crear invitación "guardian" SIN student_id → 400
//      STUDENT_REQUIRED_FOR_GUARDIAN.
//
// Las filas insertadas se limpian con defer (la BD del container es compartida
// por subtests del process).
package guardian_invitation_flow_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UUIDs/emails sembrados en seeds/playground_v2/base.
const (
	adminSanIgnacioEmail = "admin.sanignacio@edugo.test" // Carmen Valdes, school_admin S1 (academic.*)
	tutorCastroEmail     = "tutor.castro@edugo.test"     // Miguel Castro, guardián de Carlos/Diego, NO de Sofia

	schoolS1ID = "b1000000-0000-0000-0000-000000000001" // Colegio San Ignacio (S1)
	sofiaUnit  = "ac000000-0000-0000-0000-000000000003" // 5to A (unidad de Sofia)

	guardianCastroID = "00000000-0000-0000-0000-000000000012" // Miguel Castro (redime → representante)
	studentSofiaID   = "00000000-0000-0000-0000-000000000009" // Sofia (alumna apuntada)
)

// postJSON ejecuta un POST con bearer + body JSON y devuelve (status, body). El
// paquete roleflow no expone un helper genérico de POST-con-body (solo Login/
// SwitchSubject/GetJSON), así que lo declaramos local.
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

// adminTokenOnSofiaUnit loguea al admin de S1 y deja su contexto activo en
// (S1, unidad de Sofia) con grants academic.* — necesario para que el approve
// firme AMBOS sellos (school.guardian + unit.guardian) sobre la unidad de la
// solicitud (la unidad de la invitación guardian = la de Sofia).
func adminTokenOnSofiaUnit(t *testing.T) string {
	t.Helper()
	login := roleflow.Login(t, identityServer, adminSanIgnacioEmail, roleflow.DemoPassword)
	require.NotNil(t, login.ActiveContext, "login admin debe traer active_context")

	// Fija el contexto en la unidad de Sofia (5to A). switch-context permite a un
	// school_admin (rol school-scoped) activar cualquier unidad activa de su
	// escuela; los grants academic.* se conservan.
	status, raw := roleflow.SwitchSubject(t, identityServer, login.AccessToken, "self", schoolS1ID, sofiaUnit)
	require.Equalf(t, http.StatusOK, status,
		"switch-context admin a (S1, unidad Sofia): esperaba 200, got %d body=%s", status, string(raw))

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
	require.Equal(t, sofiaUnit, sc.Context.AcademicUnitID, "contexto debe quedar en la unidad de Sofia")
	require.NotEmpty(t, sc.AccessToken)
	return sc.AccessToken
}

// invitationResponse es el sub-set tipado de dto.InvitationResponse.
type invitationResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	SchoolID       string `json:"school_id"`
	AcademicUnitID string `json:"academic_unit_id"`
	InvitationType string `json:"invitation_type"`
}

// errorBody es el dto.ErrorResponse estándar.
type errorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// pendingJoinRequest es el sub-set tipado de dto.JoinRequestResponse.
type pendingJoinRequest struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	SchoolID       string `json:"school_id"`
	AcademicUnitID string `json:"academic_unit_id"`
	InvitationType string `json:"invitation_type"`
	Status         string `json:"status"`
}

// TestGuardianInvitation_HappyPath (T1): el flujo completo invitación→redención→
// aprobación→vínculo→lectura de notas del acudido.
func TestGuardianInvitation_HappyPath(t *testing.T) {
	adminToken := adminTokenOnSofiaUnit(t)

	// 1. El admin crea una invitación tipo "guardian" apuntando a Sofia (en la
	//    unidad de Sofia, que es el contexto activo del admin).
	sofia := studentSofiaID
	status, raw := postJSON(t, academicServer, "/api/v1/schools/invitations", adminToken, map[string]any{
		"academic_unit_id": sofiaUnit,
		"invitation_type":  "guardian",
		"student_id":       sofia,
	})
	require.Equalf(t, http.StatusCreated, status,
		"crear invitación guardian: esperaba 201, got %d body=%s", status, string(raw))

	var inv invitationResponse
	require.NoError(t, json.Unmarshal(raw, &inv), "parse invitación body=%s", string(raw))
	require.NotEmpty(t, inv.Code, "la invitación debe traer code")
	require.Equal(t, "guardian", inv.InvitationType)
	defer cleanupInvitation(t, inv.ID)

	// 2. El representante (castro) redime el code → join-request pending. Castro
	//    es guardián en S1 (membresía activa) → el sello de colegio nace firmado
	//    (auto); queda pendiente el de unidad.
	castroLogin := roleflow.LoginRaw(t, identityServer, tutorCastroEmail, roleflow.DemoPassword)
	require.NotEmpty(t, castroLogin.AccessToken, "castro debe poder loguear")

	status, raw = postJSON(t, academicServer, "/api/v1/invitations/redeem", castroLogin.AccessToken, map[string]any{
		"code": inv.Code,
	})
	require.Equalf(t, http.StatusCreated, status,
		"redimir code: esperaba 201, got %d body=%s", status, string(raw))

	// 3. El admin lista las solicitudes pendientes de S1 y localiza la de castro.
	jrID := findPendingJoinRequestID(t, adminToken, guardianCastroID)
	require.NotEmpty(t, jrID, "debe existir una solicitud pendiente de castro")
	defer cleanupJoinRequest(t, jrID)

	// 4. El admin aprueba: con contexto (S1, unidad de Sofia) + grants academic.*
	//    firma el/los sello(s) restante(s) → approved. Al aprobar se crea el
	//    vínculo guardian_relations (el corazón de S2).
	status, raw = postJSON(t, academicServer, "/api/v1/schools/join-requests/"+jrID+"/approve", adminToken, map[string]any{})
	require.Equalf(t, http.StatusOK, status,
		"aprobar join-request: esperaba 200, got %d body=%s", status, string(raw))

	var approved struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(raw, &approved), "parse approve body=%s", string(raw))
	require.Equalf(t, "approved", approved.Status,
		"la solicitud debe quedar approved tras firmar ambos sellos body=%s", string(raw))

	defer cleanupGuardianRelation(t, guardianCastroID, studentSofiaID)

	// 5. ASERCIÓN CENTRAL: el vínculo castro↔Sofia existe en BD, scoped a S1, active.
	assertGuardianRelation(t, guardianCastroID, studentSofiaID, schoolS1ID)

	// 6. Lazo F2/F3: el representante hace switch-context al acudido (Sofia) y lee
	//    sus notas. El vínculo recién creado es lo que habilita el contexto ward y
	//    la revalidación por-request del endpoint.
	status, raw = roleflow.SwitchSubject(t, identityServer, castroLogin.AccessToken, "ward:"+studentSofiaID, schoolS1ID, "")
	require.Equalf(t, http.StatusOK, status,
		"switch-context de castro a ward Sofia: esperaba 200, got %d body=%s", status, string(raw))

	var wardCtx struct {
		AccessToken string `json:"access_token"`
		Context     *struct {
			ActorMode        string `json:"actor_mode"`
			SubjectStudentID string `json:"subject_student_id"`
		} `json:"context"`
	}
	require.NoError(t, json.Unmarshal(raw, &wardCtx), "parse ward switch body=%s", string(raw))
	require.NotNil(t, wardCtx.Context, "switch ward: context nil body=%s", string(raw))
	require.Equal(t, roleflow.ActorModeWard, wardCtx.Context.ActorMode, "el contexto debe ser ward")
	require.Equal(t, studentSofiaID, wardCtx.Context.SubjectStudentID, "el acudido activo debe ser Sofia")

	status, raw = roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", wardCtx.AccessToken)
	require.Equalf(t, http.StatusOK, status,
		"GET /me/wards/grades como representante de Sofia: esperaba 200, got %d body=%s", status, string(raw))
}

// TestGuardianInvitation_MissingStudent_BadRequest (T2): crear invitación
// "guardian" sin student_id → 400 STUDENT_REQUIRED_FOR_GUARDIAN.
func TestGuardianInvitation_MissingStudent_BadRequest(t *testing.T) {
	adminToken := adminTokenOnSofiaUnit(t)

	status, raw := postJSON(t, academicServer, "/api/v1/schools/invitations", adminToken, map[string]any{
		"academic_unit_id": sofiaUnit,
		"invitation_type":  "guardian",
		// student_id ausente a propósito.
	})
	require.Equalf(t, http.StatusBadRequest, status,
		"guardian sin student_id: esperaba 400, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 400 body=%s", string(raw))
	assert.Equal(t, "STUDENT_REQUIRED_FOR_GUARDIAN", body.Code,
		"code del 400 debe ser STUDENT_REQUIRED_FOR_GUARDIAN body=%s", string(raw))
}

// findPendingJoinRequestID lista GET /schools/join-requests/pending (contexto del
// admin) y devuelve el ID de la solicitud cuyo user_id == el representante.
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
			require.Equal(t, "guardian", r.InvitationType, "la solicitud de castro debe ser tipo guardian")
			return r.ID
		}
	}
	return ""
}

// assertGuardianRelation lee directamente academic.guardian_relations y exige una
// fila guardian↔student scoped a school con status='active' e is_active=true.
func assertGuardianRelation(t *testing.T, guardianID, studentID, schoolID string) {
	t.Helper()
	var status string
	var isActive bool
	err := testSQLDB.QueryRow(
		`SELECT status, is_active FROM academic.guardian_relations
		 WHERE guardian_id = $1 AND student_id = $2 AND school_id = $3
		 AND relationship_type = 'guardian'`,
		guardianID, studentID, schoolID).Scan(&status, &isActive)
	require.NoErrorf(t, err,
		"debe existir UNA fila guardian_relations castro↔sofia en S1 (creada al aprobar)")
	assert.Equal(t, "active", status, "el vínculo debe nacer active (la aprobación del colegio ES la autorización)")
	assert.True(t, isActive, "is_active debe ser true para un vínculo active")
}

// --- limpieza (defer) ---

func cleanupInvitation(t *testing.T, invitationID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_invitations WHERE id = $1`, invitationID)
	assert.NoError(t, err, "cleanup invitación")
}

func cleanupJoinRequest(t *testing.T, joinRequestID string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.school_join_requests WHERE id = $1`, joinRequestID)
	assert.NoError(t, err, "cleanup join-request")
}

func cleanupGuardianRelation(t *testing.T, guardianID, studentID string) {
	t.Helper()
	// Borra SOLO el vínculo creado por el test (guardian+student+tipo guardian);
	// los vínculos sembrados (castro↔carlos/diego) usan otros student_id.
	_, err := testSQLDB.Exec(
		`DELETE FROM academic.guardian_relations
		 WHERE guardian_id = $1 AND student_id = $2 AND relationship_type = 'guardian'`,
		guardianID, studentID)
	assert.NoError(t, err, "cleanup guardian_relation")
}
