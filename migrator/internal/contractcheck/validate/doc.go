// Package validate cruza el Snapshot del frontend KMP contra el Snapshot
// del seed de producción y emite los Drift detectados, agrupados en las 7
// categorías del design (screen_key_phantom, screen_key_dead,
// permission_phantom, permission_zombie, role_phantom, role_unused,
// service_prefix_mismatch).
//
// La implementación detallada se cablea en tareas posteriores (4.x).
// Este archivo establece el paquete para que el resto del scaffolding
// compile y para reservar el namespace.
package validate
