package report

// Violation codes catalog (v1). The mapping between code and severity
// is the source of truth; validators must use these constants when
// emitting Violations.
//
// Each code traces back to a requirement in the spec
// (system-data-quality-spec/phase-a-static-auditor/requirements.md).
const (
	CodePermResourceMissing     = "PERM_RESOURCE_MISSING"      // A-REQ-1.2
	CodePermDuplicateAction     = "PERM_DUPLICATE_ACTION"      // A-REQ-1.3
	CodeRolePermRoleMissing     = "ROLE_PERM_ROLE_MISSING"     // A-REQ-2.2
	CodeRolePermPermissionMissing = "ROLE_PERM_PERMISSION_MISSING" // A-REQ-2.2
	CodeRolePermDuplicate       = "ROLE_PERM_DUPLICATE"        // A-REQ-2.3
	CodeRSDuplicateDefault      = "RS_DUPLICATE_DEFAULT"       // A-REQ-3.2
	CodeRSScreenMissing         = "RS_SCREEN_MISSING"          // A-REQ-3.3
	CodeRSNoDefault             = "RS_NO_DEFAULT"              // A-REQ-3.4
	CodeSlotInvalidJSON         = "SLOT_INVALID_JSON"          // A-REQ-4.1
	CodeSlotRefMissing          = "SLOT_REF_MISSING"           // A-REQ-4.5
	CodeConceptTypeMissing      = "CONCEPT_TYPE_MISSING"       // A-REQ-5.1
	CodeConceptDuplicateKey     = "CONCEPT_DUPLICATE_KEY"      // A-REQ-5.2
	CodeResourceOrphan          = "RESOURCE_ORPHAN"            // A-REQ-6.1
	CodePermissionZombie        = "PERMISSION_ZOMBIE"          // A-REQ-6.2
	CodeRoleNoDefaultScreen     = "ROLE_NO_DEFAULT_SCREEN"     // A-REQ-6.3
	CodeMenuParentMissing       = "MENU_PARENT_MISSING"        // A-REQ-7.1
	CodeMenuCycle               = "MENU_CYCLE"                 // A-REQ-7.2
	CodeInternalError           = "INTERNAL_ERROR"             // panic recover in RunAll
)

// SeverityFor returns the canonical severity for a known code, or
// SeverityError as a safe default for unknown codes (defensive).
func SeverityFor(code string) Severity {
	switch code {
	case CodeRSNoDefault,
		CodeResourceOrphan,
		CodePermissionZombie,
		CodeRoleNoDefaultScreen:
		return SeverityWarning
	}
	return SeverityError
}
