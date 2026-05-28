package seed

import (
	"encoding/json"
	"sort"
)

// PermissionsReferencedInSlots inspecciona el slot_data de cada
// ScreenInstance del Snapshot y devuelve, en orden alfabético sin
// duplicados, los códigos de permiso que aparecen referenciados bajo las
// claves "requiredPermission", "permissions" o "permission".
//
// Es la heurística que B-REQ-4 necesita para distinguir un permiso
// "zombie" de uno consumido por configuración dinámica del backend.
//
// Estructura aceptada del slot_data (todos opcionales):
//
//	{
//	  "requiredPermission": "schools:read",
//	  "permissions": ["schools:read", "schools:update"],
//	  "permission": "schools:delete",
//	  "actions": [
//	    {"requiredPermission": "schools:create"}
//	  ],
//	  "items": [...]
//	}
//
// El recorrido es genérico: cualquier nivel de anidamiento se inspecciona
// vía json.Unmarshal sobre interface{}. Se descartan valores no-string
// y strings vacíos.
func (s Snapshot) PermissionsReferencedInSlots() []string {
	set := map[string]struct{}{}
	for _, inst := range s.ScreenInstances {
		if len(inst.SlotData) == 0 {
			continue
		}
		var any interface{}
		if err := json.Unmarshal(inst.SlotData, &any); err != nil {
			continue // slot inválido: ignorar, no abortar el run.
		}
		collectPermissionRefs(any, set)
	}
	out := make([]string, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// collectPermissionRefs recorre recursivamente el árbol JSON acumulando
// los strings encontrados bajo claves de interés.
func collectPermissionRefs(node interface{}, set map[string]struct{}) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, val := range v {
			switch key {
			case "requiredPermission", "permission":
				if s, ok := val.(string); ok && s != "" {
					set[s] = struct{}{}
				}
			case "permissions":
				if list, ok := val.([]interface{}); ok {
					for _, item := range list {
						if s, ok := item.(string); ok && s != "" {
							set[s] = struct{}{}
						}
					}
				} else if s, ok := val.(string); ok && s != "" {
					// Tolerancia: a veces "permissions" viene como string.
					set[s] = struct{}{}
				}
			}
			// Siempre seguir profundizando: el valor puede ser otro mapa
			// o una lista de mapas.
			collectPermissionRefs(val, set)
		}
	case []interface{}:
		for _, item := range v {
			collectPermissionRefs(item, set)
		}
	}
}
