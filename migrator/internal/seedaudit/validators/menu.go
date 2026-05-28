package validators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/google/uuid"
)

// Three-colour DFS markers for cycle detection (design §3.2):
//   - white = unvisited
//   - gray  = on the current DFS stack
//   - black = fully explored
const (
	menuWhite = 0
	menuGray  = 1
	menuBlack = 2
)

// ValidateMenuHierarchy covers A-REQ-7: every Resource.ParentID must
// resolve to an existing Resource (MENU_PARENT_MISSING) and the parent
// chain must be acyclic (MENU_CYCLE).
//
// Cycles are detected with the standard three-colour DFS. When a gray
// node is revisited on the same DFS path the cycle slice (in order) is
// emitted in Violation.References["cycle"].
func ValidateMenuHierarchy(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	parentMissingFor := make(map[uuid.UUID]uuid.UUID)
	for i := range s.Resources {
		r := s.Resources[i]
		if r.ParentID == nil {
			continue
		}
		if _, ok := s.ResourceByID[*r.ParentID]; !ok {
			parentMissingFor[r.ID] = *r.ParentID
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeMenuParentMissing),
				Code:     report.CodeMenuParentMissing,
				Message:  fmt.Sprintf("El recurso %q tiene parent_id que no existe en el seed.", r.Key),
				Entity:   "Resource",
				EntityID: r.ID.String(),
				References: map[string]string{
					"resource_key":      r.Key,
					"missing_parent_id": r.ParentID.String(),
				},
			})
		}
	}

	color := make(map[uuid.UUID]int, len(s.Resources))
	for i := range s.Resources {
		color[s.Resources[i].ID] = menuWhite
	}

	reportedCycles := make(map[string]struct{})

	// Iterate in deterministic order (input order of Resources slice).
	for i := range s.Resources {
		root := s.Resources[i]
		if color[root.ID] != menuWhite {
			continue
		}
		stack := []uuid.UUID{}
		dfsCycle(s, root.ID, color, &stack, &violations, parentMissingFor, reportedCycles)
	}

	return violations
}

func dfsCycle(
	s *loader.SeedSnapshot,
	id uuid.UUID,
	color map[uuid.UUID]int,
	stack *[]uuid.UUID,
	violations *[]report.Violation,
	parentMissingFor map[uuid.UUID]uuid.UUID,
	reportedCycles map[string]struct{},
) {
	color[id] = menuGray
	*stack = append(*stack, id)
	defer func() {
		*stack = (*stack)[:len(*stack)-1]
		color[id] = menuBlack
	}()

	res, ok := s.ResourceByID[id]
	if !ok {
		return
	}
	if res.ParentID == nil {
		return
	}
	parentID := *res.ParentID

	// If the parent is missing we already reported MENU_PARENT_MISSING;
	// skip recursion to avoid spurious errors.
	if _, missing := parentMissingFor[id]; missing {
		return
	}
	if _, exists := s.ResourceByID[parentID]; !exists {
		return
	}

	switch color[parentID] {
	case menuGray:
		// Cycle detected: the cycle is stack[idx..] + parentID, where
		// idx is the position of parentID in the current stack.
		idx := -1
		for i, v := range *stack {
			if v == parentID {
				idx = i
				break
			}
		}
		if idx == -1 {
			return
		}
		cycle := append([]uuid.UUID(nil), (*stack)[idx:]...)
		key := cycleKey(cycle)
		if _, seen := reportedCycles[key]; seen {
			return
		}
		reportedCycles[key] = struct{}{}

		ids := make([]string, 0, len(cycle))
		keys := make([]string, 0, len(cycle))
		for _, cid := range cycle {
			ids = append(ids, cid.String())
			if r, ok := s.ResourceByID[cid]; ok {
				keys = append(keys, r.Key)
			}
		}
		head := cycle[0]
		var headKey string
		if r, ok := s.ResourceByID[head]; ok {
			headKey = r.Key
		}
		*violations = append(*violations, report.Violation{
			Severity: report.SeverityFor(report.CodeMenuCycle),
			Code:     report.CodeMenuCycle,
			Message:  fmt.Sprintf("Ciclo detectado en la jerarquía de menús (recurso %q).", headKey),
			Entity:   "Resource",
			EntityID: head.String(),
			References: map[string]string{
				"cycle":      strings.Join(ids, ","),
				"cycle_keys": strings.Join(keys, ","),
			},
		})
	case menuWhite:
		dfsCycle(s, parentID, color, stack, violations, parentMissingFor, reportedCycles)
	}
}

// cycleKey returns a canonical string representation of a cycle by
// rotating the slice so it starts with the smallest UUID. This avoids
// reporting the same cycle multiple times depending on the entry node.
func cycleKey(cycle []uuid.UUID) string {
	if len(cycle) == 0 {
		return ""
	}
	minIdx := 0
	for i := 1; i < len(cycle); i++ {
		if cycle[i].String() < cycle[minIdx].String() {
			minIdx = i
		}
	}
	rotated := append([]uuid.UUID{}, cycle[minIdx:]...)
	rotated = append(rotated, cycle[:minIdx]...)
	parts := make([]string, len(rotated))
	for i, v := range rotated {
		parts[i] = v.String()
	}
	return strings.Join(parts, "->")
}

// sort.Strings is not used directly but kept available for future
// helpers that wish to emit deterministic side-data in References.
var _ = sort.Strings
