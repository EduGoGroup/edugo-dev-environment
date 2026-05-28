// Package report owns the wire-format model of the audit (Violation,
// Stats, Summary, AuditReport) and the JSON / Markdown renderers. It is
// versioned via SchemaVersion so consumers (CI, dashboards, humans) can
// rely on a stable contract across auditor releases.
package report
