package kmp

// Snapshot agrupa todo el contrato estático declarado en archivos *.kt
// del repo KMP. Las claves de los maps son los identificadores
// canónicos (screenKey, código de permiso, código de rol); los values
// llevan todas las ubicaciones donde el literal aparece.
//
// Es seguro de mutar antes de pasarlo al cross-validator, pero los
// callers deberían tratarlo como inmutable: la implementación de Extract
// devuelve la estructura ya construida y no la reusa.
type Snapshot struct {
	// ScreenKeys -> ubicaciones donde aparece "override val screenKey =
	// \"...\"" o "val screenKey = \"...\"". Cubre B-REQ-1.
	ScreenKeys map[string][]Location

	// Permissions -> ubicaciones donde aparece un literal
	// requiredPermission = "...". La inferencia canónica
	// (<resource>:{read,create,update,delete}) NO se hace aquí, vive en
	// validate. Cubre B-REQ-3.2.
	Permissions map[string][]Location

	// Roles -> ubicaciones donde aparece "dashboard-<X>" o un string
	// comparado en when (role.code). Cubre B-REQ-5.1, B-REQ-5.2.
	Roles map[string][]Location

	// Contracts agrupa por archivo el screenKey + apiPrefix + basePath +
	// resource declarados en BaseCrudContract(...) o equivalente. Cubre
	// B-REQ-3.1 y B-REQ-7.1.
	Contracts []ContractDecl
}

// Location identifica un literal extraído del código KMP. El Snippet se
// recorta a 120 chars (B-REQ-1.3 espera ubicaciones detalladas pero no
// chorizos completos).
type Location struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Snippet  string `json:"snippet"`
}

// ContractDecl agrupa los cuatro literales que permiten razonar sobre
// el ruteo de servicio (B-REQ-7) y la inferencia de permisos (B-REQ-3.1).
//
// Cualquier campo puede quedar vacío si el archivo declara solo un
// subset (por ejemplo un ScreenContract custom que no extiende
// BaseCrudContract). El cross-validator es responsable de tratar los
// vacíos como "no aplica".
type ContractDecl struct {
	ScreenKey string   `json:"screen_key"`
	APIPrefix string   `json:"api_prefix"`
	BasePath  string   `json:"base_path"`
	Resource  string   `json:"resource"`
	File      Location `json:"file"`
}

// ExtractError reporta un archivo que el extractor no pudo leer o
// procesar parcialmente. NO aborta el run (B-REQ design §5): los
// errores se acumulan y se exponen al final.
type ExtractError struct {
	FilePath string `json:"file_path"`
	Reason   string `json:"reason"`
}

func (e ExtractError) Error() string {
	return e.FilePath + ": " + e.Reason
}

// snippetMaxLen define el truncamiento del Snippet de Location.
const snippetMaxLen = 120
