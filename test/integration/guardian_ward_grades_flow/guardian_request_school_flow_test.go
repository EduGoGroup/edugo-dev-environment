//go:build integration

// Tests de la auto-solicitud de representante bajo el modelo POR-ESCUELA
// (plan 024 · F4 · S1). El vínculo guardian_relations es por-escuela
// (school_id NOT NULL desde F1); este flujo —POST /guardian-relations/request +
// POST /:id/approve— era PRE-F1 y no persistía school_id. S1 lo cierra:
//
//   T1 — feliz: el guardián castro (sin vínculo previo a Sofia) solicita el
//        vínculo con Sofia → 201 PENDING con school_id = San Ignacio (la única
//        escuela activa de Sofia, resuelto automáticamente). Luego el school_admin
//        de San Ignacio (Carmen, contexto activo en S1) aprueba → 200 y el vínculo
//        queda ACTIVE con el school_id correcto.
//   T2 — escuela ajena: un PENDING con school_id = otra escuela (S3) no puede ser
//        aprobado por Carmen (cuyo ActiveContext.SchoolID = S1) → 403
//        APPROVER_SCHOOL_MISMATCH (el aprobador solo gobierna vínculos de SU escuela).
//   T3 — lazo con F2/F3: tras aprobar T1, castro hace switch-context a ward Sofia y
//        GET /me/wards/grades → 200 (el vínculo recién activado habilita la lectura
//        ":own" del acudido).
//
// Reusa el harness multi-API de setup_test.go (identityServer + academicServer
// sobre el mismo testcontainer sembrado con `base`) y el acceso directo al pool
// (testSQLDB) para insertar el PENDING de T2 y limpiar lo creado (defer DELETE).
package guardian_ward_grades_flow_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/test/integration/internal/roleflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// castro es guardián de carlos/diego, NO de Sofia → puede auto-solicitar el
	// vínculo con Sofia sin chocar con un duplicate (mendoza↔sofia ya existe).
	tutorCastroEmail = "tutor.castro@edugo.test"
	// admin de San Ignacio (S1): rol school_admin con academic.* → cubre approve.
	adminSanIgnacioEmail = "admin.sanignacio@edugo.test"

	schoolSanIgnacioID = "b1000000-0000-0000-0000-000000000001" // S1 (escuela de Sofia)
	schoolGlobalID     = "b3000000-0000-0000-0000-000000000003" // S3 (escuela ajena a Sofia)

	sofiaEmail = "est.sofia@edugo.test"
)

// guardianRelationBody es el sub-set tipado de dto.GuardianRelationResponse que
// nos interesa (id + school_id + status) para validar la auto-solicitud/aprobación.
type guardianRelationBody struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	SchoolID  string `json:"school_id"`
	Status    string `json:"status"`
}

// loginCastro loguea a castro y devuelve un access_token EN CONTEXTO. castro es
// guardián de 2 acudidos (carlos/diego) sin membership propia, así que el login NO
// auto-resuelve contexto (DEC-A: >1 sujeto deja active_context nil). Resolvemos con
// un switch-context explícito a ward carlos en San Ignacio; ese token carga el grant
// academic.guardian_relations.request del rol guardian (lo único que el endpoint de
// auto-solicitud exige, junto con el user_id del JWT = castro).
func loginCastro(t *testing.T) string {
	t.Helper()
	raw := roleflow.LoginRaw(t, identityServer, tutorCastroEmail, roleflow.DemoPassword)
	require.NotEmpty(t, raw.AccessToken, "castro: token vacío")

	status, body := roleflow.SwitchSubject(t, identityServer, raw.AccessToken,
		"ward:"+studentCarlosID, schoolSanIgnacioID, "")
	require.Equalf(t, http.StatusOK, status,
		"castro switch a ward Carlos: esperaba 200, got %d body=%s", status, string(body))

	var sw struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.Unmarshal(body, &sw), "parse switch body=%s", string(body))
	require.NotEmpty(t, sw.AccessToken, "castro switch: access_token vacío")
	return sw.AccessToken
}

// doPost ejecuta POST con bearer y body JSON, devolviendo (status, raw). Helper
// local porque roleflow solo expone GetJSON (los flujos previos eran de lectura).
func doPost(t *testing.T, baseURL, path, bearer string, payload any) (int, []byte) {
	t.Helper()
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, body)
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

// TestRequestApprove_PerSchool_HappyPath (T1): castro solicita vínculo con Sofia
// (school_id resuelto a S1 automáticamente) y Carmen (admin de S1) lo aprueba.
func TestRequestApprove_PerSchool_HappyPath(t *testing.T) {
	// 1. castro loguea (switch a ward carlos en S1): el token carga el grant
	//    academic.guardian_relations.request (seed L4 guardian).
	castroToken := loginCastro(t)

	// 2. POST /guardian-relations/request por Sofia (por email). Sofia tiene UNA sola
	//    escuela activa (S1) → school_id se resuelve solo, sin enviarlo en el body.
	status, raw := doPost(t, academicServer.URL, "/api/v1/guardian-relations/request", castroToken, map[string]string{
		"identifier":        sofiaEmail,
		"relationship_type": "guardian",
	})
	require.Equalf(t, http.StatusCreated, status,
		"request feliz: esperaba 201, got %d body=%s", status, string(raw))

	var created guardianRelationBody
	require.NoError(t, json.Unmarshal(raw, &created), "parse 201 body=%s", string(raw))
	defer deleteRelation(t, created.ID) // limpia el vínculo creado pase lo que pase.

	require.NotEmpty(t, created.ID, "el pending creado debe traer id body=%s", string(raw))
	assert.Equal(t, "pending", created.Status, "auto-solicitud nace en pending body=%s", string(raw))
	assert.Equalf(t, schoolSanIgnacioID, created.SchoolID,
		"school_id del pending debe ser San Ignacio (única escuela de Sofia) body=%s", string(raw))

	// 3. Carmen (admin de San Ignacio) loguea: single-school → contexto activo en S1.
	admin := roleflow.Login(t, identityServer, adminSanIgnacioEmail, roleflow.DemoPassword)
	require.NotNil(t, admin.ActiveContext, "admin: active_context debe estar resuelto")
	require.Equalf(t, schoolSanIgnacioID, admin.ActiveContext.SchoolID,
		"precondición: el admin debe estar en contexto de San Ignacio, got %s", admin.ActiveContext.SchoolID)

	// 4. POST /:id/approve en contexto de S1 → 200 (escuela del aprobador == escuela
	//    del pending).
	status, raw = doPost(t, academicServer.URL, "/api/v1/guardian-relations/"+created.ID+"/approve", admin.AccessToken, nil)
	require.Equalf(t, http.StatusOK, status,
		"approve feliz: esperaba 200, got %d body=%s", status, string(raw))

	// 5. Verifica que el vínculo quedó ACTIVE conservando el school_id.
	status, raw = roleflow.GetJSON(t, academicServer, "/api/v1/guardian-relations/"+created.ID, admin.AccessToken)
	require.Equalf(t, http.StatusOK, status,
		"GET del vínculo tras aprobar: esperaba 200, got %d body=%s", status, string(raw))

	var after guardianRelationBody
	require.NoError(t, json.Unmarshal(raw, &after), "parse GET body=%s", string(raw))
	assert.Equal(t, "active", after.Status, "tras approve el vínculo debe estar active body=%s", string(raw))
	assert.Equal(t, schoolSanIgnacioID, after.SchoolID, "school_id se preserva tras aprobar body=%s", string(raw))
	assert.Equal(t, studentSofiaID, after.StudentID, "el vínculo debe ser con Sofia body=%s", string(raw))
}

// TestApprove_ForeignSchool_Forbidden (T2): un PENDING cuya escuela es S3 no puede
// ser aprobado por Carmen (ActiveContext.SchoolID = S1) → 403 APPROVER_SCHOOL_MISMATCH.
func TestApprove_ForeignSchool_Forbidden(t *testing.T) {
	// Inserta un PENDING en S3 directamente en la BD del container. Usamos un par
	// guardian/student que NO colisiona con el índice único (guardian_id, student_id,
	// school_id) del seed: castro(…0012) ↔ Sofia(…0009) en S3 no existe en `base`.
	relID := insertPendingRelation(t, guardianCastroID, studentSofiaID, schoolGlobalID)
	defer deleteRelation(t, relID)

	// Carmen (admin de San Ignacio, contexto en S1) intenta aprobar el pending de S3.
	admin := roleflow.Login(t, identityServer, adminSanIgnacioEmail, roleflow.DemoPassword)
	require.Equal(t, schoolSanIgnacioID, admin.ActiveContext.SchoolID,
		"precondición: admin en contexto de San Ignacio")

	status, raw := doPost(t, academicServer.URL, "/api/v1/guardian-relations/"+relID+"/approve", admin.AccessToken, nil)
	require.Equalf(t, http.StatusForbidden, status,
		"approve de escuela ajena: esperaba 403, got %d body=%s", status, string(raw))

	var body errorBody
	require.NoError(t, json.Unmarshal(raw, &body), "parse 403 body=%s", string(raw))
	assert.Equal(t, "APPROVER_SCHOOL_MISMATCH", body.Code,
		"code del 403 debe ser APPROVER_SCHOOL_MISMATCH body=%s", string(raw))

	// El pending NO debe haber cambiado de estado.
	var status0 string
	require.NoError(t, testSQLDB.QueryRow(
		`SELECT status FROM academic.guardian_relations WHERE id = $1`, relID).Scan(&status0))
	assert.Equal(t, "pending", status0, "el pending de S3 no debió activarse")
}

// TestRequestApprove_ThenWardGradesReadable (T3): lazo con F2/F3. Tras aprobar el
// vínculo castro↔Sofia en S1, castro impersona a Sofia (switch-context ward) y lee
// sus notas vía GET /me/wards/grades → 200. Confirma que la auto-solicitud aprobada
// habilita la lectura ":own" del acudido (la revalidación por-request encuentra el
// vínculo ACTIVE).
func TestRequestApprove_ThenWardGradesReadable(t *testing.T) {
	// 1. castro solicita el vínculo con Sofia (S1 auto).
	castroToken := loginCastro(t)
	status, raw := doPost(t, academicServer.URL, "/api/v1/guardian-relations/request", castroToken, map[string]string{
		"identifier":        sofiaEmail,
		"relationship_type": "guardian",
	})
	require.Equalf(t, http.StatusCreated, status, "request: esperaba 201, got %d body=%s", status, string(raw))
	var created guardianRelationBody
	require.NoError(t, json.Unmarshal(raw, &created), "parse 201 body=%s", string(raw))
	defer deleteRelation(t, created.ID)

	// 2. Carmen aprueba en S1.
	admin := roleflow.Login(t, identityServer, adminSanIgnacioEmail, roleflow.DemoPassword)
	status, raw = doPost(t, academicServer.URL, "/api/v1/guardian-relations/"+created.ID+"/approve", admin.AccessToken, nil)
	require.Equalf(t, http.StatusOK, status, "approve: esperaba 200, got %d body=%s", status, string(raw))

	// 3. castro impersona a Sofia (ward) y lee sus notas. Re-login para que el token
	//    refleje el vínculo recién activado (los wards se calculan al firmar).
	castroRaw := roleflow.LoginRaw(t, identityServer, tutorCastroEmail, roleflow.DemoPassword)
	sw, swRaw := roleflow.SwitchSubject(t, identityServer, castroRaw.AccessToken,
		"ward:"+studentSofiaID, schoolSanIgnacioID, "")
	require.Equalf(t, http.StatusOK, sw,
		"switch a ward Sofia: esperaba 200, got %d body=%s", sw, string(swRaw))

	var swBody struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.Unmarshal(swRaw, &swBody), "parse switch body=%s", string(swRaw))
	require.NotEmpty(t, swBody.AccessToken, "switch a ward debe devolver access_token")

	status, raw = roleflow.GetJSON(t, academicServer, "/api/v1/me/wards/grades", swBody.AccessToken)
	require.Equalf(t, http.StatusOK, status,
		"ward grades tras aprobar el vínculo: esperaba 200, got %d body=%s", status, string(raw))
}

// insertPendingRelation inserta una fila guardian_relations PENDING con el school_id
// dado, devolviendo su id. Usado por T2 para fabricar un pending de escuela ajena
// sin pasar por el endpoint (que siempre resuelve la escuela del alumno).
func insertPendingRelation(t *testing.T, guardianID, studentID, schoolID string) string {
	t.Helper()
	var id string
	err := testSQLDB.QueryRow(
		`INSERT INTO academic.guardian_relations
		    (guardian_id, student_id, school_id, relationship_type, status, is_primary, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, 'guardian', 'pending', false, false, NOW(), NOW())
		 RETURNING id::text`,
		guardianID, studentID, schoolID).Scan(&id)
	require.NoError(t, err, "insertar pending de prueba")
	require.NotEmpty(t, id)
	return id
}

// deleteRelation borra una fila guardian_relations por id (idempotente). Limpieza
// de los vínculos que crean/insertan los subtests.
func deleteRelation(t *testing.T, id string) {
	t.Helper()
	_, err := testSQLDB.Exec(`DELETE FROM academic.guardian_relations WHERE id = $1`, id)
	require.NoError(t, err, "borrar vínculo de prueba")
}

const (
	guardianCastroID = "00000000-0000-0000-0000-000000000012" // Miguel Castro (guardián)
	studentCarlosID  = "00000000-0000-0000-0000-000000000008" // Carlos (acudido de castro en S1)
)
