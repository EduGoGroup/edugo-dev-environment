package kmp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Regex compiladas una sola vez por proceso (B-REQ-10.3 + tarea 2.3).
// Siguen las plantillas declaradas en design.md §3.1, ajustadas a los
// patrones reales observados en kmp-screens/src/commonMain/.
var (
	// "override val screenKey = \"...\"" — el caso dominante (clases que
	// extienden BaseCrudContract / BaseContract).
	reScreenKeyOverride = regexp.MustCompile(`override\s+val\s+screenKey\s*(?::\s*\w+)?\s*=\s*"([^"]+)"`)

	// "val screenKey = \"...\"" — variante para top-level vals o
	// companion objects.
	reScreenKeyPlain = regexp.MustCompile(`(?m)^\s*val\s+screenKey\s*(?::\s*\w+)?\s*=\s*"([^"]+)"`)

	reAPIPrefix = regexp.MustCompile(`apiPrefix\s*=\s*"([^"]+)"`)
	reBasePath  = regexp.MustCompile(`basePath\s*=\s*"([^"]+)"`)
	reResource  = regexp.MustCompile(`(?m)^\s*resource\s*=\s*"([^"]+)"`)

	// "requiredPermission = \"...\"" o "requiredPermission: String? = \"...\"".
	// El segundo caso (con tipo explícito) es el dominante; el primero
	// aparece en handlers donde la herencia ya fija el tipo.
	reRequiredPermission = regexp.MustCompile(`requiredPermission(?:\s*:\s*String\??)?\s*=\s*"([^"]+)"`)

	// "dashboard-<x>" — captura el sufijo (rol). Acepta letras, dígitos
	// y guión bajo en el nombre del rol.
	reDashboardRole = regexp.MustCompile(`"dashboard-([a-zA-Z][a-zA-Z0-9_]*)"`)

	// when (role) | when (role.code) | when (userRole.code) | when (currentRole.code).
	// Restringida a discriminantes cuyo nombre termina en "role" o "Role"
	// (con sufijo opcional ".code") para evitar capturar bloques como
	// when (status.code), when (request.method.code) o when (period.code),
	// que producían ruido masivo en role_phantom (smoke run 2026-05-08).
	reWhenRole = regexp.MustCompile(`(?s)when\s*\(\s*(?:[a-zA-Z]*[Rr]ole)\s*(?:\.code)?\s*\)\s*\{(.*?)\}`)

	// Strings literales simples (sin escapes complejos). Se aplica
	// dentro del bloque when (role.code).
	reStringLiteral = regexp.MustCompile(`"([a-zA-Z][a-zA-Z0-9_-]*)"`)
)

// Directorios que el walker no debe recorrer (tarea 2.2).
var skipDirs = map[string]struct{}{
	"build":           {},
	".gradle":         {},
	"kotlin-js-store": {},
	"node_modules":    {},
	".git":            {},
}

// Extract recorre las raíces dadas, extrae el contrato declarado en cada
// archivo *.kt elegible y devuelve el Snapshot consolidado, junto a los
// errores no-fatales encontrados (B-REQ-9.3).
//
// El segundo error (return) sólo se usa para fallos sistémicos
// (por ejemplo, ningún root existe en disco). Los errores por archivo
// se devuelven en el slice []ExtractError.
func Extract(roots []string) (Snapshot, []ExtractError, error) {
	snap := Snapshot{
		ScreenKeys:  map[string][]Location{},
		Permissions: map[string][]Location{},
		Roles:       map[string][]Location{},
		Contracts:   []ContractDecl{},
	}
	var extractErrs []ExtractError
	visited := 0
	rootsExist := 0

	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			// Path declarado pero inexistente: warning, no abortar.
			extractErrs = append(extractErrs, ExtractError{
				FilePath: root,
				Reason:   fmt.Sprintf("kmp root not found: %v", err),
			})
			continue
		}
		if !info.IsDir() {
			extractErrs = append(extractErrs, ExtractError{
				FilePath: root,
				Reason:   "kmp root is not a directory",
			})
			continue
		}
		rootsExist++

		walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				extractErrs = append(extractErrs, ExtractError{FilePath: path, Reason: walkErr.Error()})
				return nil
			}
			if d.IsDir() {
				if _, skip := skipDirs[d.Name()]; skip {
					return fs.SkipDir
				}
				return nil
			}
			if !isExtractableKotlinFile(path) {
				return nil
			}
			visited++
			data, err := os.ReadFile(path)
			if err != nil {
				extractErrs = append(extractErrs, ExtractError{FilePath: path, Reason: err.Error()})
				return nil
			}
			processFile(path, string(data), &snap)
			return nil
		})
		if walkErr != nil {
			extractErrs = append(extractErrs, ExtractError{FilePath: root, Reason: walkErr.Error()})
		}
	}

	if rootsExist == 0 {
		return snap, extractErrs, fmt.Errorf("ningún kmp-root existe en disco (revisaron: %v)", roots)
	}
	_ = visited // reservado para métricas / debug verboso

	// Determinismo (B-REQ-10.2): ordenamos las locations por (file, line).
	sortSnapshot(&snap)
	return snap, extractErrs, nil
}

// isExtractableKotlinFile decide si un archivo es candidato. Se ignoran
// los *Test.kt y los archivos no .kt.
func isExtractableKotlinFile(path string) bool {
	if !strings.HasSuffix(path, ".kt") {
		return false
	}
	base := filepath.Base(path)
	if strings.HasSuffix(base, "Test.kt") || strings.HasSuffix(base, "Tests.kt") {
		return false
	}
	// Ignorar archivos en directorios de test (defensa adicional para
	// estructuras como src/commonTest/, src/desktopTest/, etc.).
	if strings.Contains(path, "/src/") {
		idx := strings.Index(path, "/src/")
		rest := path[idx+len("/src/"):]
		if seg := strings.SplitN(rest, "/", 2); len(seg) > 0 {
			if strings.HasSuffix(seg[0], "Test") {
				return false
			}
		}
	}
	return true
}

// processFile aplica las regex sobre el contenido ya saneado de un
// archivo y lo agrega al Snapshot.
func processFile(path string, raw string, snap *Snapshot) {
	cleaned, lineMap := stripComments(raw)
	lines := strings.Split(cleaned, "\n")

	extractScreenKeys(path, lines, lineMap, snap)
	extractContractDecls(path, raw, cleaned, lines, lineMap, snap)
	extractExplicitPermissions(path, lines, lineMap, snap)
	extractRoles(path, cleaned, lines, lineMap, snap)
}

// extractScreenKeys captura override val screenKey = "..." y la variante
// val screenKey = "..." (B-REQ-1.1, tarea 2.4).
func extractScreenKeys(path string, lines []string, lineMap []int, snap *Snapshot) {
	for i, line := range lines {
		if m := reScreenKeyOverride.FindStringSubmatch(line); m != nil {
			addLocation(snap.ScreenKeys, m[1], path, lineMap[i], line)
			continue
		}
		if m := reScreenKeyPlain.FindStringSubmatch(line); m != nil {
			addLocation(snap.ScreenKeys, m[1], path, lineMap[i], line)
		}
	}
}

// extractContractDecls compone un ContractDecl por archivo a partir del
// primer match de cada literal. Si el archivo declara múltiples
// contratos (raro), sólo se conserva el primero — el cross-validator
// puede inspeccionar Snapshot.ScreenKeys para ver duplicados.
func extractContractDecls(path string, raw, cleaned string, lines []string, lineMap []int, snap *Snapshot) {
	decl := ContractDecl{}
	declLine := 0

	if m := reAPIPrefix.FindStringSubmatchIndex(cleaned); m != nil {
		decl.APIPrefix = cleaned[m[2]:m[3]]
		declLine = lineForOffset(cleaned, m[2], lineMap)
	}
	if m := reBasePath.FindStringSubmatchIndex(cleaned); m != nil {
		decl.BasePath = cleaned[m[2]:m[3]]
		if declLine == 0 {
			declLine = lineForOffset(cleaned, m[2], lineMap)
		}
	}
	if m := reResource.FindStringSubmatchIndex(cleaned); m != nil {
		decl.Resource = cleaned[m[2]:m[3]]
		if declLine == 0 {
			declLine = lineForOffset(cleaned, m[2], lineMap)
		}
	}

	// Resolvemos screenKey al primer match en este archivo.
	for i, line := range lines {
		if m := reScreenKeyOverride.FindStringSubmatch(line); m != nil {
			decl.ScreenKey = m[1]
			if declLine == 0 {
				declLine = lineMap[i]
			}
			break
		}
		if m := reScreenKeyPlain.FindStringSubmatch(line); m != nil {
			decl.ScreenKey = m[1]
			if declLine == 0 {
				declLine = lineMap[i]
			}
			break
		}
	}

	// Sólo guardamos un ContractDecl si hay al menos un literal útil
	// para el cross-validator.
	if decl.ScreenKey == "" && decl.APIPrefix == "" && decl.BasePath == "" && decl.Resource == "" {
		return
	}

	snippet := truncateSnippet(firstNonEmptyLine(lines, declLine))
	decl.File = Location{FilePath: path, Line: declLine, Snippet: snippet}
	snap.Contracts = append(snap.Contracts, decl)
	_ = raw
}

// extractExplicitPermissions captura "requiredPermission = \"...\"" y
// la variante con tipo (B-REQ-3.2, tarea 2.6).
func extractExplicitPermissions(path string, lines []string, lineMap []int, snap *Snapshot) {
	for i, line := range lines {
		if m := reRequiredPermission.FindStringSubmatch(line); m != nil {
			addLocation(snap.Permissions, m[1], path, lineMap[i], line)
		}
	}
}

// extractRoles captura literales "dashboard-<X>" y los strings dentro
// de when (role.code) { ... } (B-REQ-5.1, B-REQ-5.2, tarea 2.7).
func extractRoles(path string, cleaned string, lines []string, lineMap []int, snap *Snapshot) {
	for i, line := range lines {
		matches := reDashboardRole.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			addLocation(snap.Roles, m[1], path, lineMap[i], line)
		}
	}
	// Bloques when (role.code) { ... }: capturamos cada string literal
	// del bloque como candidato a rol. Esto es una sobre-aproximación
	// (cualquier string en el when entra), pero el cross-validator
	// filtra contra el seed y los falsos positivos quedan acotados.
	whenBlocks := reWhenRole.FindAllStringSubmatchIndex(cleaned, -1)
	for _, idx := range whenBlocks {
		// idx[2:4] delimita el contenido (capturing group 1).
		blockStart, blockEnd := idx[2], idx[3]
		block := cleaned[blockStart:blockEnd]
		for _, sm := range reStringLiteral.FindAllStringSubmatchIndex(block, -1) {
			str := block[sm[2]:sm[3]]
			line := lineForOffset(cleaned, blockStart+sm[2], lineMap)
			snippet := lineSnippet(lines, line, lineMap)
			addLocation(snap.Roles, str, path, line, snippet)
		}
	}
}

// addLocation agrega una Location al map, deduplicando por
// (FilePath, Line, Snippet).
func addLocation(m map[string][]Location, key, file string, line int, rawLine string) {
	if key == "" {
		return
	}
	loc := Location{FilePath: file, Line: line, Snippet: truncateSnippet(rawLine)}
	for _, existing := range m[key] {
		if existing.FilePath == loc.FilePath && existing.Line == loc.Line && existing.Snippet == loc.Snippet {
			return
		}
	}
	m[key] = append(m[key], loc)
}

// truncateSnippet recorta a snippetMaxLen y normaliza whitespace.
func truncateSnippet(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= snippetMaxLen {
		return s
	}
	return s[:snippetMaxLen]
}

// firstNonEmptyLine devuelve la línea no vacía empezando en startLine
// (1-indexed); si startLine está fuera de rango, devuelve "".
func firstNonEmptyLine(lines []string, startLine int) string {
	if startLine <= 0 || startLine > len(lines) {
		return ""
	}
	for i := startLine - 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			return lines[i]
		}
	}
	return ""
}

// lineSnippet devuelve la línea original (limpia de comentarios) en la
// posición line (1-indexed contra el archivo original).
func lineSnippet(lines []string, line int, lineMap []int) string {
	for i, orig := range lineMap {
		if orig == line && i < len(lines) {
			return lines[i]
		}
	}
	return ""
}

// sortSnapshot ordena las Location de cada bucket por (FilePath, Line).
func sortSnapshot(snap *Snapshot) {
	for _, locs := range snap.ScreenKeys {
		sortLocations(locs)
	}
	for _, locs := range snap.Permissions {
		sortLocations(locs)
	}
	for _, locs := range snap.Roles {
		sortLocations(locs)
	}
	sort.SliceStable(snap.Contracts, func(i, j int) bool {
		if snap.Contracts[i].File.FilePath != snap.Contracts[j].File.FilePath {
			return snap.Contracts[i].File.FilePath < snap.Contracts[j].File.FilePath
		}
		return snap.Contracts[i].File.Line < snap.Contracts[j].File.Line
	})
}

func sortLocations(locs []Location) {
	sort.SliceStable(locs, func(i, j int) bool {
		if locs[i].FilePath != locs[j].FilePath {
			return locs[i].FilePath < locs[j].FilePath
		}
		return locs[i].Line < locs[j].Line
	})
}
