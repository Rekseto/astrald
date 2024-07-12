package objects

import (
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/mod/admin"
	"github.com/cryptopunkscc/astrald/mod/content"
	"github.com/cryptopunkscc/astrald/mod/objects"
)

func (mod *Module) LoadDependencies() error {
	// optional
	mod.content, _ = core.Load[content.Module](mod.node, content.ModuleName)

	// inject admin command
	if adm, err := core.Load[admin.Module](mod.node, admin.ModuleName); err == nil {
		adm.AddCommand(objects.ModuleName, NewAdmin(mod))
	}

	return nil
}
