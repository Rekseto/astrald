package profile

import (
	"context"
	"encoding/json"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/mod/discovery"
	"github.com/cryptopunkscc/astrald/mod/profile/proto"
	"github.com/cryptopunkscc/astrald/net"
	"github.com/cryptopunkscc/astrald/node/events"
)

type EventHandler struct {
	*Module
}

func (h *EventHandler) Run(ctx context.Context) error {
	return events.Handle(ctx, h.node.Events(), h.handleDiscovered)
}

func (h *EventHandler) handleDiscovered(e discovery.EventDiscovered) error {
	for _, srv := range e.Info.Services {
		if srv.Identity.IsEqual(h.node.Identity()) {
			continue
		}
		if srv.Type == serviceType {
			return h.updateIdentityProfile(e.Identity, srv.Name)
		}
	}
	return nil
}

func (h *EventHandler) updateIdentityProfile(target id.Identity, serviceName string) error {
	h.log.Infov(2, "updating profile for %s", target)

	conn, err := net.Route(h.ctx, h.node.Router(), net.NewQuery(h.node.Identity(), target, serviceName))
	if err != nil {
		return err
	}
	defer conn.Close()

	var profile proto.Profile
	err = json.NewDecoder(conn).Decode(&profile)
	if err != nil {
		return err
	}

	for _, pep := range profile.Endpoints {
		ep, err := h.node.Infra().Parse(pep.Network, pep.Address)
		if err != nil {
			continue
		}

		_ = h.nodes.AddEndpoint(target, ep)
	}

	h.log.Info("%s profile updated.", target)

	return nil
}
