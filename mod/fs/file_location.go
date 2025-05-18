package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"io"
	"strings"
)

type FileLocation struct {
	NodeID *astral.Identity
	Path   astral.String16
}

// astral

func (fn FileLocation) ObjectType() string {
	return "mod.fs.file_location"
}

func (fn FileLocation) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Struct(fn).WriteTo(w)
}

func (fn *FileLocation) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Struct(fn).ReadFrom(r)
}

// json

func (fn FileLocation) MarshalJSON() ([]byte, error) {
	type alias FileLocation
	return json.Marshal(alias(fn))
}

func (fn *FileLocation) UnmarshalJSON(bytes []byte) error {
	type alias FileLocation
	return json.Unmarshal(bytes, (*alias)(fn))
}

// text

func (fn FileLocation) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%s:%s", fn.NodeID, fn.Path)), nil
}

func (fn *FileLocation) UnmarshalText(text []byte) error {
	split := strings.SplitN(string(text), ":", 2)
	if len(split) != 2 {
		return errors.New("invalid format")
	}

	id, err := astral.IdentityFromString(split[0])
	if err != nil {
		return err
	}

	*fn = FileLocation{
		NodeID: id,
		Path:   astral.String16(split[1]),
	}
	return nil
}

// other

func (fn FileLocation) String() string {
	return string(fn.Path)
}

func init() {
	var fn FileLocation
	astral.DefaultBlueprints.Add(&fn)
}
