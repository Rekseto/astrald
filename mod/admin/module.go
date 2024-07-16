package admin

import (
	"github.com/cryptopunkscc/astrald/id"
	"io"
)

const ModuleName = "admin"
const ActionAccess = "mod.admin.access"
const ActionSudo = "mod.admin.sudo"

type Module interface {
	AddCommand(name string, cmd Command) error
}

type Command interface {
	Exec(term Terminal, args []string) error
}

type Terminal interface {
	UserIdentity() id.Identity
	SetUserIdentity(id.Identity)
	Sprintf(f string, v ...any) string
	Printf(f string, v ...any)
	Println(v ...any)
	Scanf(f string, v ...any)
	ScanLine() (string, error)
	Color() bool
	SetColor(bool)
	io.Writer
}

// Format types are used to format output text on the terminal. Example:
// term.Println("normal %s %s", Keyword("keyword"), Faded("faded"))

type Header string
type Keyword string
type Faded string
type Important string
