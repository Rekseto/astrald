package objects

import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/objects"
)

func (mod *Module) AddReceiver(receiver objects.Receiver) error {
	if receiver == nil {
		return errors.New("receiver is nil")
	}

	return mod.receivers.Add(receiver)
}

func (mod *Module) Receive(obj astral.Object, source *astral.Identity) (err error) {
	if source.IsZero() {
		source = mod.node.Identity()
	}

	ok := mod.receive(mod.node.Identity(), obj)

	if !ok {
		err = errors.New("rejected")
	}

	return
}

func (mod *Module) receive(senderID *astral.Identity, object astral.Object) (ok bool) {
	drop := &Drop{
		senderID: senderID,
		object:   object,
		repo:     mod.Root(),
	}
	for _, r := range mod.receivers.Clone() {
		err := r.ReceiveObject(drop)
		if err != nil {
			mod.log.Errorv(1, "receiver %v errored: %v", r, err)
		}
	}

	return drop.accepted.Get()
}
