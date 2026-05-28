// Package validators houses the pure, side-effect-free validators that
// inspect a SeedSnapshot and return a slice of report.Violation. Each
// validator covers one requirement family (permissions, role/permission
// matrix, resource screens, slot data, concepts, inverse coverage, and
// menu hierarchy) and is composed by the registry into RunAll.
package validators
