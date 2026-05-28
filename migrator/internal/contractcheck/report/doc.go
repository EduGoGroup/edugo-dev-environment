// Package report serializa el resultado del cross-checker FE↔BE como
// JSON + Markdown, computa el diff contra el baseline previo y maneja el
// flag --update-baseline.
//
// La implementación detallada se cablea en tareas posteriores (5.x).
// Este archivo establece el paquete para que el resto del scaffolding
// compile y para reservar el namespace.
package report
