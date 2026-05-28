// Package kmp extrae artefactos del repo Kotlin Multiplatform de EduGo.
//
// Su única responsabilidad es leer archivos *.kt y producir un Snapshot
// inmutable con los screenKey, contratos (apiPrefix, basePath, resource),
// permisos explícitos y roles citados que aparecen estáticamente en el
// código. La extracción usa regex pragmáticas adaptadas a las
// convenciones del repo (ver design.md §3.1 y §7.1).
//
// Este paquete es puramente lectura: no toca BD, no resuelve referencias
// runtime, y no muta el código KMP.
package kmp
