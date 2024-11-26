package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/computerdane/nf6/impl/impl_vip"
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var vipCmd = &cobra.Command{
	Use:    "vip",
	Short:  "Start the Virtual Internet Provider",
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		// initialize

		apiTlsPubKeyData, err := os.ReadFile(apiTlsPubKeyPath)
		if err != nil {
			lib.Crash("failed to open API public key file at ", apiTlsPubKeyPath, ": ", err)
		}
		apiTlsPubKey = string(apiTlsPubKeyData)
		if apiTlsPubKey == "" {
			lib.Crash("API's TLS public key cannot be empty!")
		}
		privKeyData, err := os.ReadFile(wgPrivKeyPath)
		if err != nil {
			lib.Crash("failed to open WireGuard private key file at ", wgPrivKeyPath, ": ", err)
		}
		wgPrivKey, err := wgtypes.ParseKey(string(privKeyData))
		if err != nil {
			lib.Crash(err)
		}
		wg, err := wgctrl.New()
		if err != nil {
			lib.Crash(err)
		}

		// get list of hosts

		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.Vip_ListHosts(ctx, nil)
		if err != nil {
			lib.Crash(err)
		}

		// construct list of peers

		peers := make([]wgtypes.PeerConfig, len(reply.GetHosts()))
		for i, host := range reply.GetHosts() {
			ip := net.ParseIP(host.GetAddr6())
			ipNet := net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(net.IPv6len*8, net.IPv6len*8),
			}
			wgPubKey, err := wgtypes.ParseKey(host.GetWgPubKey())
			if err != nil {
				lib.Warn(err)
				continue
			}
			peers[i] = wgtypes.PeerConfig{
				PublicKey:  wgPubKey,
				AllowedIPs: []net.IPNet{ipNet},
			}
		}

		// configure wg

		if err := wg.ConfigureDevice(wgDeviceName, wgtypes.Config{
			PrivateKey:   &wgPrivKey,
			ListenPort:   &vipWgPort,
			ReplacePeers: true,
			Peers:        peers,
		}); err != nil {
			lib.Crash(err)
		}
		fmt.Printf("added %d peers\n", len(peers))

		// create gRPC listener

		if _, err := os.Stat(tlsCaCertPath); err != nil {
			lib.Crash("ca cert file not found: ", err)
		}
		if _, err := os.Stat(tlsCertPath); err != nil {
			lib.Crash("cert file not found: ", err)
		}
		if _, err := os.Stat(tlsPrivKeyPath); err != nil {
			lib.Crash("private key not found: ", err)
		}

		caCert, err := os.ReadFile(tlsCaCertPath)
		if err != nil {
			lib.Crash("failed to read ca cert: ", err)
		}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(caCert); !ok {
			lib.Crash("failed to append ca cert")
		}
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsPrivKeyPath)
		if err != nil {
			lib.Crash("failed to load x509 keypair: ", err)
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    pool,
			RootCAs:      pool,
		})

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", vipGrpcPort))
		if err != nil {
			lib.Crash("failed to listen: ", err)
		}
		fmt.Printf("listening at %v", lis.Addr())

		server := grpc.NewServer(grpc.Creds(creds))
		nf6.RegisterNf6VipServer(server, &impl_vip.VipServer{
			ApiTlsPubKey: apiTlsPubKey,
			VipWgPort:    vipWgPort,
			Wg:           wg,
			WgDeviceName: wgDeviceName,
			WgPrivKey:    wgPrivKey,
		})
		if err := server.Serve(lis); err != nil {
			lib.Crash("failed to serve: ", err)
		}
	},
}
