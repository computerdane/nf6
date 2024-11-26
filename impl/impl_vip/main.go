package impl_vip

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

type VipServer struct {
	nf6.UnimplementedNf6VipServer
	ApiTlsPubKey string
	Wg           *wgctrl.Client
	WgDeviceName string
	VipWgPort    int
	WgPrivKey    wgtypes.Key
}

func (s *VipServer) CreatePeer(ctx context.Context, in *nf6.CreatePeer_Request) (*nf6.None, error) {
	pubKey, err := lib.TlsGetGrpcPubKey(ctx)
	if err != nil {
		return nil, err
	}
	if pubKey != s.ApiTlsPubKey {
		lib.Warn("attempt made with unknown public key: ", pubKey)
		return nil, status.Error(codes.Unauthenticated, "access denied")
	}

	ip := net.ParseIP(in.GetAddr6())
	ipNet := net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(net.IPv6len*8, net.IPv6len*8),
	}
	wgPubKey, err := wgtypes.ParseKey(in.GetWgPubKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse wg pub key")
	}
	peer := wgtypes.PeerConfig{
		PublicKey:  wgPubKey,
		AllowedIPs: []net.IPNet{ipNet},
	}
	if err := s.Wg.ConfigureDevice(s.WgDeviceName, wgtypes.Config{
		PrivateKey: &s.WgPrivKey,
		ListenPort: &s.VipWgPort,
		Peers:      []wgtypes.PeerConfig{peer},
	}); err != nil {
		lib.Warn("failed to configure wg device: ", err)
		return nil, status.Error(codes.Internal, "failed to configure wg device")
	}
	return nil, nil
}

func (s *VipServer) DeletePeer(ctx context.Context, in *nf6.DeletePeer_Request) (*nf6.None, error) {
	pubKey, err := lib.TlsGetGrpcPubKey(ctx)
	if err != nil {
		return nil, err
	}
	if pubKey != s.ApiTlsPubKey {
		lib.Warn("attempt made with unknown public key: ", pubKey)
		return nil, status.Error(codes.Unauthenticated, "access denied")
	}

	wgPubKey, err := wgtypes.ParseKey(in.GetWgPubKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse wg pub key")
	}
	peer := wgtypes.PeerConfig{
		Remove:    true,
		PublicKey: wgPubKey,
	}
	if err := s.Wg.ConfigureDevice(s.WgDeviceName, wgtypes.Config{
		PrivateKey: &s.WgPrivKey,
		ListenPort: &s.VipWgPort,
		Peers:      []wgtypes.PeerConfig{peer},
	}); err != nil {
		lib.Warn("failed to configure wg device: ", err)
		return nil, status.Error(codes.Internal, "failed to configure wg device")
	}
	return nil, nil
}
