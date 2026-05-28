// Package loader resolves the production seed into an immutable
// SeedSnapshot consumed by the static auditor (Phase A). It imports the
// canonical seed package and materializes its rows into the
// `entities.*` structs, then precomputes the indexes the validators
// rely on to stay linear-time.
package loader
