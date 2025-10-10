package nodes

import (
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/shell"
	"github.com/cryptopunkscc/astrald/mod/utp"
)

// NOTE: might  move to mod/nat
import (
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/shell"
	"github.com/cryptopunkscc/astrald/mod/utp"
)

type opStartNatTraversal struct {
	Out string `query:"optional"`
}

func (mod *Module) OpStartNatTraversal(ctx *astral.Context, q shell.Query,
	args opStartNatTraversal) error {
	ips := mod.IP.FindIPCandidates()
	if len(ips) == 0 {
		return fmt.Errorf(``)
	}

	localPort := mod.UTP.ListenPort()

	ch := astral.NewChannelFmt(q.Accept(), "", args.Out)
	defer ch.Close()

	localEndpoint := &utp.Endpoint{
		IP:   ips[0],
		Port: astral.Uint16(localPort),
	}

	err := ch.Write(localEndpoint)
	if err != nil {
		return err
	}

	initiatorEndpoint, err := ch.ReadPayload((&utp.Endpoint{}).ObjectType())
	if err != nil {
		return err
	}

	mod.log.Info("NAT traversal info exchanged: local=%v, remote=%v", localEndpoint, initiatorEndpoint)
	// FIXME: hole punching (by sending UTP packets to each other)

	// FIXME: return established pair of endpoints
	return nil
}
