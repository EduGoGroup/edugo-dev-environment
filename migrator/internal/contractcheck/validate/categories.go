package validate

// CategoryMeta describes the canonical severity envelope of a category.
// Detectors may downgrade or escalate within their bounds based on the
// extra rules listed in design.md §4 (e.g. screen_key_dead is normally
// `warning` but escalates to `error` when the dead screen is a
// dashboard or `is_default=true`).
type CategoryMeta struct {
	Default Severity
	Notes   string
}

// Catalog is the source of truth for category → default severity. The
// reporter and any future linter should consult it instead of hardcoding.
var Catalog = map[string]CategoryMeta{
	CategoryScreenKeyPhantom: {
		Default: SeverityError,
		Notes:   "FE referencia un screenKey que el seed no declara; el FE quedaría sin pantalla en runtime.",
	},
	CategoryScreenKeyDead: {
		Default: SeverityWarning,
		Notes:   "El seed declara un screen_key sin implementación KMP. Escala a error si screen_type=dashboard o is_default=true.",
	},
	CategoryPermissionPhantom: {
		Default: SeverityError,
		Notes:   "FE consume un permiso (literal o inferido) que el seed no contiene. Degrada a warning si el resource existe pero la acción inferida no.",
	},
	CategoryPermissionZombie: {
		Default: SeverityWarning,
		Notes:   "Permiso seedado que ningún role_permission asigna y que ningún slot_data ni FE referencia.",
	},
	CategoryRolePhantom: {
		Default: SeverityError,
		Notes:   "FE menciona un role.code que el seed no declara.",
	},
	CategoryRoleUnused: {
		Default: SeverityWarning,
		Notes:   "Rol seedado que el FE nunca atiende. Escala a error si scope=system.",
	},
	CategoryServicePrefixMismatch: {
		Default: SeverityError,
		Notes:   "El apiPrefix declarado por el FE no coincide con la tabla canónica resource→servicio.",
	},
}

// SeverityFor returns the default severity for a known category, or
// SeverityWarning as a safe fallback. Detectors with conditional
// escalation/downgrade should set Drift.Severity directly.
func SeverityFor(category string) Severity {
	if meta, ok := Catalog[category]; ok {
		return meta.Default
	}
	return SeverityWarning
}
