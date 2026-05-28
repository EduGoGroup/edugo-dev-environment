package loader

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	"github.com/google/uuid"
)

// SeedSourceProduction is the only source supported in v1.
const SeedSourceProduction = "production"

// RunOptions captures the inputs the loader needs from the CLI. Only
// SeedSource is consulted today; future versions may carry alternate
// fixtures (Decision D-1).
type RunOptions struct {
	SeedSource string
}

// SeedSnapshot is the immutable cut of the seed consumed by the
// validators. Slices preserve the seed's declared order; the maps offer
// O(1) lookups so validators can stay linear in the row count.
//
// P4-1 (plan B): RolePermissions fue eliminado del snapshot. La tabla
// iam.role_permissions ya no existe. Los validadores que dependían de
// ese slice (role_permissions, inverse PERMISSION_ZOMBIE) se eliminan
// o quedan en silencio hasta el cierre de P4-2 sobre iam.role_grants.
type SeedSnapshot struct {
	Resources          []entities.Resource
	Permissions        []entities.Permission
	Roles              []entities.Role
	ResourceScreens    []entities.ResourceScreen
	ScreenInstances    []entities.ScreenInstance
	ScreenTemplates    []entities.ScreenTemplate
	ConceptTypes       []entities.ConceptType
	ConceptDefinitions []entities.ConceptDefinition

	ResourceByID     map[uuid.UUID]*entities.Resource
	ResourceByKey    map[string]*entities.Resource
	PermissionByID   map[uuid.UUID]*entities.Permission
	PermissionByName map[string]*entities.Permission
	RoleByID         map[uuid.UUID]*entities.Role
	ScreenByKey      map[string]*entities.ScreenInstance
	ConceptTypeByID  map[uuid.UUID]*entities.ConceptType
}

// Load resolves the seed declared by opts and returns a fully indexed
// snapshot. Per Decision D-1 the only accepted source in v1 is
// `production`.
func Load(opts RunOptions) (*SeedSnapshot, error) {
	source := opts.SeedSource
	if source == "" {
		source = SeedSourceProduction
	}
	if source != SeedSourceProduction {
		return nil, fmt.Errorf("seed source %q not supported (only %q in v1)", source, SeedSourceProduction)
	}

	// El catálogo legacy (archivado pre-Fase-6) fue borrado del disco
	// en el bloque C de Fase 6 (ADR-6 cerrado). Post-Fase-6 el catálogo
	// activo del sistema vive en L0..L4 y todas las capas exponen
	// accessors públicos (`layers.LN*()` y `l4.*()`) que retornan los
	// slices listos para consumo del extractor.
	//
	// TC-5 (cerrado): post-extracción de accessors L0..L3 el loader
	// concatena las 5 capas en orden L0 → L1 → L2 → L3 → L4. Esto
	// elimina los 39 falsos positivos que el cross-checker reportaba
	// pre-fix por entidades de L0..L3 invisibles (announcements,
	// materials, super_admin, announcement_viewer + pantallas y
	// permisos asociados).
	//
	// Out-of-scope del seedaudit (NO se exponen vía accessors):
	// users, user_roles, schools, memberships — el cross-checker no
	// los consume.

	resources, err := loadResources()
	if err != nil {
		return nil, err
	}
	roles, err := loadRoles()
	if err != nil {
		return nil, err
	}
	permissions, err := loadPermissions()
	if err != nil {
		return nil, err
	}
	screenTemplates, err := loadScreenTemplates()
	if err != nil {
		return nil, err
	}
	screenInstances, err := loadScreenInstances()
	if err != nil {
		return nil, err
	}
	resourceScreens, err := loadResourceScreens()
	if err != nil {
		return nil, err
	}
	conceptTypes, err := l4.ConceptTypes()
	if err != nil {
		return nil, fmt.Errorf("load l4 concept_types: %w", err)
	}
	conceptDefinitions, err := l4.ConceptDefinitions()
	if err != nil {
		return nil, fmt.Errorf("load l4 concept_definitions: %w", err)
	}

	snap := &SeedSnapshot{
		Resources:          resources,
		Permissions:        permissions,
		Roles:              roles,
		ResourceScreens:    resourceScreens,
		ScreenInstances:    screenInstances,
		ScreenTemplates:    screenTemplates,
		ConceptTypes:       conceptTypes,
		ConceptDefinitions: conceptDefinitions,
	}
	buildIndexes(snap)
	return snap, nil
}

// loadResources concatena los recursos sembrados por L0 → L3 → L4 en
// orden de capas. L1 y L2 no siembran recursos.
func loadResources() ([]entities.Resource, error) {
	l0, err := layers.L0Resources()
	if err != nil {
		return nil, fmt.Errorf("load l0 resources: %w", err)
	}
	l3, err := layers.L3Resources()
	if err != nil {
		return nil, fmt.Errorf("load l3 resources: %w", err)
	}
	l4Res, err := l4.Resources()
	if err != nil {
		return nil, fmt.Errorf("load l4 resources: %w", err)
	}
	out := make([]entities.Resource, 0, len(l0)+len(l3)+len(l4Res))
	out = append(out, l0...)
	out = append(out, l3...)
	out = append(out, l4Res...)
	return out, nil
}

// loadRoles concatena los roles sembrados por L0 → L1 → L4 en orden
// de capas. L2 y L3 no siembran roles.
func loadRoles() ([]entities.Role, error) {
	l0, err := layers.L0Roles()
	if err != nil {
		return nil, fmt.Errorf("load l0 roles: %w", err)
	}
	l1, err := layers.L1Roles()
	if err != nil {
		return nil, fmt.Errorf("load l1 roles: %w", err)
	}
	l4Roles, err := l4.Roles()
	if err != nil {
		return nil, fmt.Errorf("load l4 roles: %w", err)
	}
	out := make([]entities.Role, 0, len(l0)+len(l1)+len(l4Roles))
	out = append(out, l0...)
	out = append(out, l1...)
	out = append(out, l4Roles...)
	return out, nil
}

// loadPermissions concatena los permisos sembrados por L0 → L3 → L4.
// L1 y L2 no siembran permisos (reusan announcements:read de L0).
func loadPermissions() ([]entities.Permission, error) {
	l0, err := layers.L0Permissions()
	if err != nil {
		return nil, fmt.Errorf("load l0 permissions: %w", err)
	}
	l3, err := layers.L3Permissions()
	if err != nil {
		return nil, fmt.Errorf("load l3 permissions: %w", err)
	}
	l4Perms, err := l4.Permissions()
	if err != nil {
		return nil, fmt.Errorf("load l4 permissions: %w", err)
	}
	out := make([]entities.Permission, 0, len(l0)+len(l3)+len(l4Perms))
	out = append(out, l0...)
	out = append(out, l3...)
	out = append(out, l4Perms...)
	return out, nil
}

// loadScreenTemplates retorna las screen_templates sembradas por L0 → L4.
// L1..L3 no siembran templates (reutilizan list/detail/form-basic-v1 de L0).
func loadScreenTemplates() ([]entities.ScreenTemplate, error) {
	l0, err := layers.L0ScreenTemplates()
	if err != nil {
		return nil, fmt.Errorf("load l0 screen_templates: %w", err)
	}
	l4ST, err := l4.ScreenTemplates()
	if err != nil {
		return nil, fmt.Errorf("load l4 screen_templates: %w", err)
	}
	out := make([]entities.ScreenTemplate, 0, len(l0)+len(l4ST))
	out = append(out, l0...)
	out = append(out, l4ST...)
	return out, nil
}

// loadScreenInstances concatena las screen_instances sembradas por
// L0 → L2 → L3 → L4. L1 no siembra instances.
func loadScreenInstances() ([]entities.ScreenInstance, error) {
	l0, err := layers.L0ScreenInstances()
	if err != nil {
		return nil, fmt.Errorf("load l0 screen_instances: %w", err)
	}
	l2, err := layers.L2ScreenInstances()
	if err != nil {
		return nil, fmt.Errorf("load l2 screen_instances: %w", err)
	}
	l3, err := layers.L3ScreenInstances()
	if err != nil {
		return nil, fmt.Errorf("load l3 screen_instances: %w", err)
	}
	l4SI, err := l4.ScreenInstances()
	if err != nil {
		return nil, fmt.Errorf("load l4 screen_instances: %w", err)
	}
	out := make([]entities.ScreenInstance, 0, len(l0)+len(l2)+len(l3)+len(l4SI))
	out = append(out, l0...)
	out = append(out, l2...)
	out = append(out, l3...)
	out = append(out, l4SI...)
	return out, nil
}

// loadResourceScreens concatena los mappings resource↔screen sembrados
// por L0 → L2 → L3 → L4. L1 no siembra resource_screens.
func loadResourceScreens() ([]entities.ResourceScreen, error) {
	l0, err := layers.L0ResourceScreens()
	if err != nil {
		return nil, fmt.Errorf("load l0 resource_screens: %w", err)
	}
	l2, err := layers.L2ResourceScreens()
	if err != nil {
		return nil, fmt.Errorf("load l2 resource_screens: %w", err)
	}
	l3, err := layers.L3ResourceScreens()
	if err != nil {
		return nil, fmt.Errorf("load l3 resource_screens: %w", err)
	}
	l4RS, err := l4.ResourceScreens()
	if err != nil {
		return nil, fmt.Errorf("load l4 resource_screens: %w", err)
	}
	out := make([]entities.ResourceScreen, 0, len(l0)+len(l2)+len(l3)+len(l4RS))
	out = append(out, l0...)
	out = append(out, l2...)
	out = append(out, l3...)
	out = append(out, l4RS...)
	return out, nil
}

// NewSnapshot is exposed for tests that need to feed synthetic data
// directly. It runs the same indexing pipeline as Load.
func NewSnapshot(
	resources []entities.Resource,
	permissions []entities.Permission,
	roles []entities.Role,
	resourceScreens []entities.ResourceScreen,
	screenInstances []entities.ScreenInstance,
	screenTemplates []entities.ScreenTemplate,
	conceptTypes []entities.ConceptType,
	conceptDefinitions []entities.ConceptDefinition,
) *SeedSnapshot {
	snap := &SeedSnapshot{
		Resources:          resources,
		Permissions:        permissions,
		Roles:              roles,
		ResourceScreens:    resourceScreens,
		ScreenInstances:    screenInstances,
		ScreenTemplates:    screenTemplates,
		ConceptTypes:       conceptTypes,
		ConceptDefinitions: conceptDefinitions,
	}
	buildIndexes(snap)
	return snap
}

// buildIndexes populates every map declared on SeedSnapshot. Pointers
// returned by the maps alias the slice elements, so callers must not
// mutate them.
func buildIndexes(s *SeedSnapshot) {
	s.ResourceByID = make(map[uuid.UUID]*entities.Resource, len(s.Resources))
	s.ResourceByKey = make(map[string]*entities.Resource, len(s.Resources))
	for i := range s.Resources {
		r := &s.Resources[i]
		s.ResourceByID[r.ID] = r
		s.ResourceByKey[r.Key] = r
	}

	s.PermissionByID = make(map[uuid.UUID]*entities.Permission, len(s.Permissions))
	s.PermissionByName = make(map[string]*entities.Permission, len(s.Permissions))
	for i := range s.Permissions {
		p := &s.Permissions[i]
		s.PermissionByID[p.ID] = p
		s.PermissionByName[p.Name] = p
	}

	s.RoleByID = make(map[uuid.UUID]*entities.Role, len(s.Roles))
	for i := range s.Roles {
		r := &s.Roles[i]
		s.RoleByID[r.ID] = r
	}

	s.ScreenByKey = make(map[string]*entities.ScreenInstance, len(s.ScreenInstances))
	for i := range s.ScreenInstances {
		si := &s.ScreenInstances[i]
		s.ScreenByKey[si.ScreenKey] = si
	}

	s.ConceptTypeByID = make(map[uuid.UUID]*entities.ConceptType, len(s.ConceptTypes))
	for i := range s.ConceptTypes {
		ct := &s.ConceptTypes[i]
		s.ConceptTypeByID[ct.ID] = ct
	}
}
