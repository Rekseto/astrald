package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/astrald/mod/shell"
	"time"
)

type opSignAppContractArgs struct {
	ID       *astral.Identity
	Out      string          `query:"optional"`
	Duration astral.Duration `query:"optional"`
}

func (mod *Module) OpSignAppContract(ctx *astral.Context, q shell.Query, args opSignAppContractArgs) (err error) {
	ch := astral.NewChannelFmt(q.Accept(), "", args.Out)
	defer ch.Close()

	if args.Duration == 0 {
		args.Duration = DefaultAppContractDuration
	}

	// initialize the contract
	var c = &apphost.AppContract{
		AppID:     args.ID,
		HostID:    mod.node.Identity(),
		StartsAt:  astral.Time(time.Now()),
		ExpiresAt: astral.Time(time.Now().Add(time.Duration(args.Duration))),
	}

	// sign the contract
	err = mod.SignAppContract(c)
	if err != nil {
		return ch.Write(astral.NewError(err.Error()))
	}

	contractID, err := objects.Save(ctx, c, mod.Objects.Root())
	if err != nil {
		return ch.Write(astral.NewError(err.Error()))
	}

	mod.log.Infov(1, "signed app contract (%v) with %v", c.AppID, contractID)

	err = mod.Index(ctx, contractID)
	if err != nil {
		return ch.Write(astral.NewError(err.Error()))
	}

	return ch.Write(contractID)
}
