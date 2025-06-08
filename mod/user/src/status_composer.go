package user

import "github.com/cryptopunkscc/astrald/mod/nearby"

var _ nearby.Composer = &Module{}

func (mod *Module) ComposeStatus(a nearby.Composition) {
	if mod.config.Public {
		c := mod.ActiveContract()
		if c != nil {
			a.Attach(c)
		} else {
			a.Attach(nearby.NewFlag("claimable"))
		}
	}
}
