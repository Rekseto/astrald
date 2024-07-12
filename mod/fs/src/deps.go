package fs

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/mod/admin"
	"github.com/cryptopunkscc/astrald/mod/content"
	"github.com/cryptopunkscc/astrald/mod/fs"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/astrald/mod/sets"
	"github.com/cryptopunkscc/astrald/resources"
	"os"
	"path/filepath"
	"time"
)

func (mod *Module) LoadDependencies() error {
	var err error

	// required
	mod.objects, err = core.Load[objects.Module](mod.node, objects.ModuleName)
	if err != nil {
		return err
	}

	mod.content, err = core.Load[content.Module](mod.node, content.ModuleName)
	if err != nil {
		return err
	}

	mod.sets, err = core.Load[sets.Module](mod.node, sets.ModuleName)
	if err != nil {
		return err
	}

	mod.objects.AddOpener(mod, 30)
	mod.objects.AddCreator(mod, 30)
	mod.objects.AddDescriber(mod)
	mod.objects.AddPurger(mod)

	// inject admin command
	if adm, err := core.Load[admin.Module](mod.node, admin.ModuleName); err == nil {
		adm.AddCommand(fs.ModuleName, NewAdmin(mod))
	}

	mod.objects.AddSearcher(NewFinder(mod))
	mod.objects.AddPrototypes(fs.FileDesc{})

	// wait for data module to finish preparing
	ctx, cancel := context.WithTimeoutCause(context.Background(), 15*time.Second, errors.New("data module timed out"))
	defer cancel()
	if err := mod.content.Ready(ctx); err != nil {
		return err
	}

	// if we have file-based resources, use that as writable storage
	fileRes, ok := mod.assets.Res().(*resources.FileResources)
	if ok {
		mod.Watch(filepath.Join(fileRes.Root(), "static_data"))

		dataPath := filepath.Join(fileRes.Root(), "data")
		err = os.MkdirAll(dataPath, 0700)
		if err == nil {
			mod.config.Store = append(mod.config.Store, dataPath)
			mod.Watch(dataPath)
		}
	}

	return nil
}
