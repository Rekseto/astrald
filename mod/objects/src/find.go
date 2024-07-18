package objects

import (
	"context"
	"github.com/cryptopunkscc/astrald/id"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/object"
)

func (mod *Module) AddFinder(finder objects.Finder) error {
	return mod.finders.Add(finder)
}

func (mod *Module) Find(ctx context.Context, objectID object.ID, scope *astral.Scope) (providers []id.Identity) {
	var unique = map[string]struct{}{}

	if scope == nil {
		scope = astral.DefaultScope()
	}

	for _, finder := range mod.finders.Clone() {
		for _, i := range finder.Find(ctx, objectID, scope) {
			if _, found := unique[i.String()]; found {
				continue
			}
			unique[i.String()] = struct{}{}
			providers = append(providers, i)
		}
	}

	return
}
