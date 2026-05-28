// Package seed expone la snapshot del seed de producción que el
// cross-checker FE↔BE necesita para diagnosticar drift.
//
// Mientras la Fase A (internal/seedaudit/loader) madura en paralelo, este
// paquete declara una interfaz local Loader y provee un mock alimentado
// por fixtures JSON (testdata/seed/*.json). Cuando Fase A esté lista,
// basta con escribir un adapter que implemente seed.Loader envolviendo el
// loader real (ver TODO.md §3.1).
package seed
