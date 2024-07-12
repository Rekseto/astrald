package policy

import (
	"context"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/mod/admin"
	"github.com/cryptopunkscc/astrald/mod/policy"
)

func (mod *Module) Prepare(ctx context.Context) error {
	// inject admin command
	if adm, err := core.Load[admin.Module](mod.node, admin.ModuleName); err == nil {
		adm.AddCommand(policy.ModuleName, NewAdmin(mod))
	}

	return nil
}
