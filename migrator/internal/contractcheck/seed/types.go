package seed

import (
	"context"
	"encoding/json"
)

// Snapshot agrupa los slices del seed que el cross-checker FE↔BE
// necesita para diagnosticar drift. Coincide con los campos del
// design.md §3.2.
//
// La intención es que esta struct sea compatible con el shape que la
// Fase A exponga eventualmente — cualquier divergencia se resuelve en
// el adapter (TODO.md §3.1) sin tocar el resto del paquete.
type Snapshot struct {
	Resources       []Resource       `json:"resources"`
	Permissions     []Permission     `json:"permissions"`
	Roles           []Role           `json:"roles"`
	RolePermissions []RolePermission `json:"role_permissions"`
	ResourceScreens []ResourceScreen `json:"resource_screens"`
	ScreenInstances []ScreenInstance `json:"screen_instances"`
}

type Resource struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Permission struct {
	Code string `json:"code"`
	Name string `json:"name,omitempty"`
}

type Role struct {
	Code  string `json:"code"`
	Name  string `json:"name,omitempty"`
	Scope string `json:"scope,omitempty"` // system | tenant | school | etc.
}

type RolePermission struct {
	RoleCode       string `json:"role_code"`
	PermissionCode string `json:"permission_code"`
}

type ResourceScreen struct {
	ResourceKey string `json:"resource_key"`
	ScreenKey   string `json:"screen_key"`
	ScreenType  string `json:"screen_type"` // list | form | detail | dashboard | special
	IsDefault   bool   `json:"is_default"`
}

// ScreenInstance representa una fila de screen_instances con su
// slot_data preservado como JSON crudo. El cross-validator escanea ese
// JSON buscando claves "requiredPermission" / "permissions" /
// "permission" para alimentar la heurística de B-REQ-4.
type ScreenInstance struct {
	ScreenKey string          `json:"screen_key"`
	SlotKey   string          `json:"slot_key,omitempty"`
	SlotData  json.RawMessage `json:"slot_data,omitempty"`
}

// Loader es la interfaz que el binario contract-check espera para
// obtener la snapshot del seed. Mantenerla local (en lugar de importar
// internal/seedaudit/loader directamente) permite que Fase B avance en
// paralelo a Fase A.
//
// Cuando Fase A esté lista, escribimos un adapter que la satisfaga y lo
// inyectamos desde main.go (ver TODO.md §3.1).
type Loader interface {
	Load(ctx context.Context) (Snapshot, error)
}
