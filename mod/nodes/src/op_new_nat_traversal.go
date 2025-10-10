package nodes

// NOTE: might  move to mod/nat
import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/shell"
)

type opNewNatTraversal struct {
	Target string
	Out    string `query:"optional"`
}

func (mod *Module) OpNewNatTraversal(ctx *astral.Context, q shell.Query,
	args opNewNatTraversal) (err error) {

	ips := mod.IP.FindIPCandidates()
	if len(ips) == 0 {
		return errors.New("no IP candidates available")
	}

	target, err := mod.Dir.ResolveIdentity(args.Target)
	if err != nil {
		return q.RejectWithCode(4)
	}

	queryArgs := &opStartNatTraversal{}

	var routedQuery = query.New(ctx.Identity(), target,
		methodStartNatTraversal,
		queryArgs)

	_, err = query.Route(ctx, mod.node, routedQuery)
	if err != nil {
		return err
	}

	ch := astral.NewChannelFmt(q.Accept(), "", args.Out)
	defer ch.Close()

	return nil
}
