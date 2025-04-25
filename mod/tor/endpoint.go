package tor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/term"
	"io"
	"strconv"
	"strings"
)

var _ astral.Object = &Endpoint{}

// Endpoint is an astral.Object holding information about a Tor endpoint (digest and port).
// Supports JSON and text.
type Endpoint struct {
	Digest Digest
	Port   astral.Uint16
}

func (Endpoint) ObjectType() string {
	return "mod.tor.endpoint"
}

func (e Endpoint) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Struct(e).WriteTo(w)
}

func (e *Endpoint) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Struct(e).ReadFrom(r)
}

// text support

func (e *Endpoint) UnmarshalText(text []byte) (err error) {
	var hp = strings.SplitN(string(text), ":", 2)
	if len(hp) != 2 {
		return fmt.Errorf("invalid format")
	}

	err = e.Digest.UnmarshalText([]byte(hp[0]))
	if err != nil {
		return err
	}

	p, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	e.Port = astral.Uint16(p)

	return nil
}

func (e Endpoint) MarshalText() (text []byte, err error) {
	s := fmt.Sprintf("%s:%d", e.Digest, e.Port)

	return []byte(s), nil
}

// JSON marshaling

func (e *Endpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Address())
}

func (e *Endpoint) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	return e.UnmarshalText([]byte(str))
}

// exonet.Endpoint

// Network returns the name of the network the address belongs to
func (e *Endpoint) Network() string {
	return ModuleName
}

// Address returns a human-readable representation of the address
func (e *Endpoint) Address() string {
	if e == nil || e.IsZero() {
		return "none"
	}

	return fmt.Sprintf("%s:%d", e.Digest, e.Port)
}

// other

// Pack returns binary representation of the address
func (e *Endpoint) Pack() []byte {
	var b = &bytes.Buffer{}

	_, err := e.WriteTo(b)
	if err != nil {
		return nil
	}

	return b.Bytes()
}

// other

// IsZero returns true if the address has zero-value
func (e *Endpoint) IsZero() bool {
	return e == nil || len(e.Digest) == 0
}

func (e *Endpoint) String() string {
	return e.Network() + ":" + e.Address()
}

func init() {
	_ = astral.DefaultBlueprints.Add(&Endpoint{})

	term.SetTranslateFunc(func(o *Endpoint) astral.Object {
		return &term.ColorString{
			Color: term.HighlightColor,
			Text:  astral.String32(o.String()),
		}
	})
}
