package seed

import (
	"context"
	"fmt"

	seedaudit "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
)

// ProductionLoader bridges the Phase-A loader (which exposes the full
// entity slices) into the slim Snapshot shape the cross-checker needs.
//
// It implements Loader. Phase B was originally written against a
// FixtureLoader stub; this adapter is the real backend.
type ProductionLoader struct {
	source string
}

// NewProductionLoader returns a Loader backed by the production seed.
// The source string is forwarded to seedaudit/loader (today only
// "production" is accepted, see Decision D-1 of Phase A).
func NewProductionLoader(source string) *ProductionLoader {
	if source == "" {
		source = seedaudit.SeedSourceProduction
	}
	return &ProductionLoader{source: source}
}

// Load resolves the seed via Phase-A's loader and projects it into the
// cross-checker Snapshot.
//
// P4-1 (plan B): la tabla iam.role_permissions fue eliminada. El campo
// Snapshot.RolePermissions queda en nil; los validadores que dependen
// de él (cross-checker zombie_permissions) emiten warnings sin detalle
// hasta P4-2, cuando se reescriban contra iam.role_grants.
func (l *ProductionLoader) Load(ctx context.Context) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}

	snap, err := seedaudit.Load(seedaudit.RunOptions{SeedSource: l.source})
	if err != nil {
		return Snapshot{}, fmt.Errorf("seed.ProductionLoader: load: %w", err)
	}

	out := Snapshot{
		Resources:       make([]Resource, 0, len(snap.Resources)),
		Permissions:     make([]Permission, 0, len(snap.Permissions)),
		Roles:           make([]Role, 0, len(snap.Roles)),
		RolePermissions: nil,
		ResourceScreens: make([]ResourceScreen, 0, len(snap.ResourceScreens)),
		ScreenInstances: make([]ScreenInstance, 0, len(snap.ScreenInstances)),
	}

	for i := range snap.Resources {
		r := &snap.Resources[i]
		out.Resources = append(out.Resources, Resource{
			Key:  r.Key,
			Name: r.DisplayName,
		})
	}
	for i := range snap.Permissions {
		p := &snap.Permissions[i]
		out.Permissions = append(out.Permissions, Permission{
			Code: p.Name,
			Name: p.DisplayName,
		})
	}
	for i := range snap.Roles {
		r := &snap.Roles[i]
		out.Roles = append(out.Roles, Role{
			Code:  r.Name,
			Name:  r.DisplayName,
			Scope: r.Scope,
		})
	}
	for i := range snap.ResourceScreens {
		rs := &snap.ResourceScreens[i]
		out.ResourceScreens = append(out.ResourceScreens, ResourceScreen{
			ResourceKey: rs.ResourceKey,
			ScreenKey:   rs.ScreenKey,
			ScreenType:  rs.ScreenType,
			IsDefault:   rs.IsDefault,
		})
	}
	for i := range snap.ScreenInstances {
		si := &snap.ScreenInstances[i]
		out.ScreenInstances = append(out.ScreenInstances, ScreenInstance{
			ScreenKey: si.ScreenKey,
			SlotData:  si.SlotData,
		})
	}
	return out, nil
}
