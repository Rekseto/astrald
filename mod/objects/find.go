package objects

import (
	"context"
	"github.com/cryptopunkscc/astrald/id"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/object"
)

// Finder is used to figure out which identities can provide access to an object
type Finder interface {
	Find(context.Context, object.ID, *astral.Scope) []id.Identity
}
