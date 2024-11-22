package impl_wgserver

import (
	"context"
	"net"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WgServer struct {
	nf6.UnimplementedNf6WgServer
	ApiTlsPubKey   string
	Wg             *wgctrl.Client
	WgDeviceName   string
	WgServerWgPort int
	WgPrivKey      wgtypes.Key
}

func (s *WgServer) CreateRoute(ctx context.Context, in *nf6.CreateRoute_Request) (*nf6.None, error) {
	pubKey, err := lib.TlsGetGrpcPubKey(ctx)
	if err != nil {
		return nil, err
	}
	if pubKey != s.ApiTlsPubKey {
		lib.Warn("attempt made with unknown public key: ", pubKey)
		return nil, status.Error(codes.Unauthenticated, "access denied")
	}

	_, ipNet, err := net.ParseCIDR(in.GetAddr6())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse addr6")
	}
	wgPubKey, err := wgtypes.ParseKey(in.GetWgPubKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse wg pub key")
	}
	peer := wgtypes.PeerConfig{
		PublicKey:  wgPubKey,
		AllowedIPs: []net.IPNet{*ipNet},
	}
	if err := s.Wg.ConfigureDevice(s.WgDeviceName, wgtypes.Config{
		PrivateKey:   &s.WgPrivKey,
		ListenPort:   &s.WgServerWgPort,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{peer},
	}); err != nil {
		lib.Warn("failed to configure wg device: ", err)
		return nil, status.Error(codes.Internal, "failed to configure wg device")
	}
	return nil, nil
}
