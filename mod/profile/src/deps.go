package profile

import (
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/mod/dir"
	"github.com/cryptopunkscc/astrald/mod/exonet"
	"github.com/cryptopunkscc/astrald/mod/nodes"
	"github.com/cryptopunkscc/astrald/mod/objects"
)

type Deps struct {
	Dir     dir.Module
	Exonet  exonet.Module
	Nodes   nodes.Module
	Objects objects.Module
}

func (mod *Module) LoadDependencies() (err error) {
	return core.Inject(mod.node, &mod.Deps)
}
