package media

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/mod/admin"
	"github.com/cryptopunkscc/astrald/mod/content"
	"github.com/cryptopunkscc/astrald/mod/media"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"time"
)

func (mod *Module) LoadDependencies() error {
	var err error

	if adm, err := core.Load[admin.Module](mod.node, admin.ModuleName); err == nil {
		adm.AddCommand(media.ModuleName, NewAdmin(mod))
	}

	mod.content, err = core.Load[content.Module](mod.node, content.ModuleName)
	if err != nil {
		return err
	}

	mod.objects, err = core.Load[objects.Module](mod.node, objects.ModuleName)
	if err != nil {
		return err
	}

	mod.objects.AddDescriber(mod)

	// wait for data module to finish preparing
	ctx, cancel := context.WithTimeoutCause(context.Background(), 15*time.Second, errors.New("data module timed out"))
	defer cancel()
	if err := mod.content.Ready(ctx); err != nil {
		return err
	}

	mod.objects.AddSearcher(mod)
	mod.objects.AddPrototypes(&media.Audio{})

	return nil
}
