package validators

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
)

// knownRefKeys is the single source of truth for slot_data keys treated
// as references the walker must resolve against the seed (Decision D-5).
// Adding a new reference type is a deliberate, explicit change here.
//
// The walker descends through objects and arrays generically; nested
// occurrences (e.g. `field.permission`, `actions[].permission`) are
// detected automatically because the recursion meets the bare key
// `permission` inside the nested object.
var knownRefKeys = map[string]string{
	"permission":   refKindPermission, // string → Permission.Name
	"permissions":  refKindPermission, // []string → each element a Permission.Name
	"requires":     refKindPermission, // []string → each element a Permission.Name
	"resource":     refKindResource,   // string → Resource.Key
	"resource_key": refKindResource,   // string → Resource.Key
}

const (
	refKindPermission = "permission"
	refKindResource   = "resource"
)

// ValidateSlotData covers A-REQ-4: every ScreenInstance.SlotData must
// parse as JSON (SLOT_INVALID_JSON) and any reference under canonical
// keys must resolve against the seed (SLOT_REF_MISSING).
func ValidateSlotData(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	for i := range s.ScreenInstances {
		si := s.ScreenInstances[i]
		raw := si.SlotData
		if len(raw) == 0 {
			continue
		}

		var root interface{}
		if err := json.Unmarshal(raw, &root); err != nil {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeSlotInvalidJSON),
				Code:     report.CodeSlotInvalidJSON,
				Message:  fmt.Sprintf("El slot_data de la pantalla %q no es JSON válido: %s", si.ScreenKey, err.Error()),
				Entity:   "ScreenInstance",
				EntityID: si.ID.String(),
				References: map[string]string{
					"screen_key": si.ScreenKey,
				},
			})
			// Per the spec, do not run the walker on a screen whose
			// slot_data failed to parse.
			continue
		}

		walker := &slotWalker{
			snapshot:   s,
			screenKey:  si.ScreenKey,
			entityID:   si.ID.String(),
			violations: &violations,
		}
		walker.walk("$", root)
	}

	return violations
}

type slotWalker struct {
	snapshot   *loader.SeedSnapshot
	screenKey  string
	entityID   string
	violations *[]report.Violation
}

// walk descends recursively through a JSON value (map / slice / scalar)
// emitting SLOT_REF_MISSING when canonical reference keys point to
// permissions or resources that don't exist in the snapshot.
func (w *slotWalker) walk(path string, node interface{}) {
	switch v := node.(type) {
	case map[string]interface{}:
		w.walkObject(path, v)
	case []interface{}:
		for idx, item := range v {
			w.walk(fmt.Sprintf("%s[%d]", path, idx), item)
		}
	}
}

func (w *slotWalker) walkObject(path string, obj map[string]interface{}) {
	for key, value := range obj {
		childPath := path + "." + key

		if kind, ok := knownRefKeys[key]; ok {
			w.checkRef(childPath, key, value, kind)
		}

		// Continue descending for compound structures (objects/arrays).
		// The check above does not recurse into reference values; their
		// scalars are leaves of the walk.
		switch nested := value.(type) {
		case map[string]interface{}:
			w.walkObject(childPath, nested)
		case []interface{}:
			for idx, item := range nested {
				w.walk(fmt.Sprintf("%s[%d]", childPath, idx), item)
			}
		}
	}
}

// checkRef validates one canonical key/value pair. permissions/requires
// expect arrays of strings; the rest expect a string scalar. Anything
// else is silently ignored — it is the schema's job (out of scope for
// the static auditor) to reject malformed shapes.
func (w *slotWalker) checkRef(path, key string, value interface{}, kind string) {
	switch v := value.(type) {
	case string:
		w.resolveRef(path, v, kind)
	case []interface{}:
		// permissions / requires arrays.
		for idx, item := range v {
			str, ok := item.(string)
			if !ok {
				continue
			}
			w.resolveRef(fmt.Sprintf("%s[%d]", path, idx), str, kind)
		}
	}
	_ = key
}

// resolveRef performs the actual lookup against the snapshot indexes.
func (w *slotWalker) resolveRef(path, value, kind string) {
	if value == "" {
		return
	}

	switch kind {
	case refKindPermission:
		if _, ok := w.snapshot.PermissionByName[value]; ok {
			return
		}
		*w.violations = append(*w.violations, report.Violation{
			Severity: report.SeverityFor(report.CodeSlotRefMissing),
			Code:     report.CodeSlotRefMissing,
			Message:  fmt.Sprintf("La pantalla %q referencia un permiso inexistente %q.", w.screenKey, value),
			Entity:   "ScreenInstance",
			EntityID: w.entityID,
			References: map[string]string{
				"screen_key":      w.screenKey,
				"missing_permission": value,
				"ref_kind":        kind,
			},
			Path: path,
		})
	case refKindResource:
		if _, ok := w.snapshot.ResourceByKey[value]; ok {
			return
		}
		*w.violations = append(*w.violations, report.Violation{
			Severity: report.SeverityFor(report.CodeSlotRefMissing),
			Code:     report.CodeSlotRefMissing,
			Message:  fmt.Sprintf("La pantalla %q referencia un recurso inexistente %q.", w.screenKey, value),
			Entity:   "ScreenInstance",
			EntityID: w.entityID,
			References: map[string]string{
				"screen_key":       w.screenKey,
				"missing_resource": value,
				"ref_kind":         kind,
			},
			Path: path,
		})
	}
}
