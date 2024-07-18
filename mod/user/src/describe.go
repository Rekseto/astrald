package user

import (
	"context"
	"github.com/cryptopunkscc/astrald/id"
	"github.com/cryptopunkscc/astrald/lib/desc"
	"github.com/cryptopunkscc/astrald/mod/user"
)

func (mod *Module) Describe(ctx context.Context, identity id.Identity, opts *desc.Opts) []*desc.Desc {
	if identity.IsEqual(mod.UserID()) {
		return []*desc.Desc{{
			Source: mod.node.Identity(),
			Data:   user.UserDesc{Name: mod.dir.DisplayName(identity)},
		}}
	}

	return nil
}
