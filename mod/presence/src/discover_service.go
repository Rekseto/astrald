package presence

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/cslq"
	"github.com/cryptopunkscc/astrald/mod/presence"
	"github.com/cryptopunkscc/astrald/mod/presence/proto"
	"net"
	"strconv"
	"time"
)

type DiscoverService struct {
	*Module
	cache map[string]*Ad
}

func NewDiscoverService(module *Module) *DiscoverService {
	return &DiscoverService{
		Module: module,
		cache:  make(map[string]*Ad),
	}
}

func (srv *DiscoverService) Run(ctx context.Context) error {
	for {
		ad, err := srv.readAd()
		if err != nil {
			return err
		}

		srv.save(ad)

		srv.log.Logv(
			2,
			"received an ad from %v endpoint %v expires %v",
			ad.Identity,
			ad.Endpoint,
			ad.ExpiresAt,
		)

		srv.events.Emit(EventAdReceived{ad})

		if ad.Has(presence.DiscoverFlag) {
			err = srv.announce.sendWithFlags(ad.UDPAddr)
			if err != nil {
				srv.log.Errorv(2, "error responding to discover request: %v", err)
			}
		}

		if srv.config.AutoAdd {
			_ = srv.Nodes.AddEndpoint(ad.Identity, ad.Endpoint)
		}

		if !srv.config.TrustAliases || ad.Alias == "" {
			continue
		}
		if _, err := srv.Dir.GetAlias(ad.Identity); err == nil {
			continue
		}
		if _, err := srv.Dir.Resolve(ad.Alias); err == nil {
			continue
		}

		err = srv.Dir.SetAlias(ad.Identity, ad.Alias)
		if err != nil {
			srv.log.Error("error setting alias '%v' for %v: %v", ad.Alias, ad.Identity.Fingerprint(), err)
		} else {
			srv.log.Info("alias set for %v (%v)", ad.Identity, ad.Identity.Fingerprint())
		}
	}
}

func (srv *DiscoverService) RecentAds() []*Ad {
	var res = make([]*Ad, 0, len(srv.cache))
	for _, p := range srv.cache {
		if p.ExpiresAt.After(time.Now()) {
			res = append(res, p)
		}
	}
	return res
}

func (srv *DiscoverService) readAd() (*Ad, error) {
	for {
		buf := make([]byte, 1024)

		n, srcAddr, err := srv.socket.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		}

		var msg proto.Ad
		err = cslq.Decode(bytes.NewReader(buf[:n]), "v", &msg)
		if err != nil {
			srv.log.Errorv(1, "error decoding ad from %v: %v", srcAddr, err)
			fmt.Println(buf[:n])
			continue
		}

		// ignore our own ad
		if msg.Identity.IsEqual(srv.node.Identity()) {
			continue
		}

		// verify signature
		if !ecdsa.VerifyASN1(msg.Identity.PublicKey().ToECDSA(), msg.Hash(), msg.Sig) {
			return nil, errors.New("invalid ad signature")
		}

		hostPort := net.JoinHostPort(srcAddr.IP.String(), strconv.Itoa(msg.Port))

		endpoint, err := srv.TCP.Parse("tcp", hostPort)
		if err != nil {
			panic(err)
		}

		return &Ad{
			UDPAddr:   srcAddr,
			Identity:  msg.Identity,
			Alias:     msg.Alias,
			Endpoint:  endpoint,
			ExpiresAt: msg.ExpiresAt,
			Flags:     msg.Flags,
		}, nil
	}
}

func (srv *DiscoverService) save(ad *Ad) {
	hexID := ad.Identity.String()
	srv.cache[hexID] = ad
	srv.clean()
}

func (srv *DiscoverService) clean() {
	for hexID, p := range srv.cache {
		if p.ExpiresAt.Before(time.Now()) {
			delete(srv.cache, hexID)
		}
	}
}
