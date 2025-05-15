package gateway

import (
	"bytes"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/exonet"
	"github.com/cryptopunkscc/astrald/mod/gateway"
)

func (mod *Module) Unpack(network string, data []byte) (exonet.Endpoint, error) {
	if network != NetworkName {
		return nil, exonet.ErrUnsupportedNetwork
	}
	return Unpack(data)
}

// Unpack converts a binary representation of the address to a struct
func Unpack(data []byte) (addr *gateway.Endpoint, err error) {
	addr = &gateway.Endpoint{}
	_, err = astral.Struct(addr).ReadFrom(bytes.NewReader(data))
	return
}
