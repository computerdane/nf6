package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	cfgFile          string
	shouldSaveConfig bool

	dataDir      string
	dbUrl        string
	port         int
	wgPrivKeyStr string
	timeout      time.Duration
	deviceName   string

	db     *pgxpool.Pool
	socket string
)

var rootCmd = &cobra.Command{
	Use:   "nf6-router",
	Short: "Nf6 Router",
}

var listenCmd = &cobra.Command{
	Use: "listen",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		db, err = pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// initialize wgctrl
		wg, err := wgctrl.New()
		if err != nil {
			log.Fatalf("failed to initialize wgctrl: %v", err)
		}
		privKey, err := wgtypes.ParseKey(wgPrivKeyStr)
		if err != nil {
			log.Fatalf("filed to parse private key: %v", err)
		}

		// get pubkeys and addrs of all machines
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		query := "select wg_public_key, text(addr_ipv6) from machine"
		rows, err := db.Query(ctx, query)
		if err != nil {
			log.Fatalf("query failed: %v", err)
		}

		// construct list of peers
		var peers []wgtypes.PeerConfig
		for rows.Next() {
			var (
				wgPubKeyStr string
				addrIpv6    string
			)
			err := rows.Scan(&wgPubKeyStr, &addrIpv6)
			if err != nil {
				log.Fatal(err)
			}

			_, ipNet, err := net.ParseCIDR(addrIpv6)
			pubKey, err := wgtypes.ParseKey(wgPubKeyStr)

			peer := wgtypes.PeerConfig{
				PublicKey:  pubKey,
				AllowedIPs: []net.IPNet{*ipNet},
			}

			peers = append(peers, peer)
		}

		// configure wg
		if err := wg.ConfigureDevice(deviceName, wgtypes.Config{
			PrivateKey:   &privKey,
			ListenPort:   &port,
			ReplacePeers: true,
			Peers:        peers,
		}); err != nil {
			log.Fatalf("failed to configure device: %v", err)
		}
		log.Printf("added %d peers", len(peers))

		// create socket
		listener, err := net.Listen("unix", socket)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("listening at %v", listener.Addr())
		if err := os.Chmod(socket, 0666); err != nil {
			log.Fatal(err)
		}

		// clean up socket file
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Remove(socket)
			os.Exit(1)
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("failed to accept connection: ", err)
				continue
			}

			go handle(conn)
		}
	},
}

func handle(conn net.Conn) {
}

func init() {
	cobra.OnInitialize(initConfig, initDataDir)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/var/lib/nf6-router/config/config.yaml", "config file")
	rootCmd.PersistentFlags().BoolVar(&shouldSaveConfig, "save-config", false, "save to the config file with the provided flags")

	lib.AddOption(rootCmd, lib.Option{P: &dataDir, Name: "dataDir", Shorthand: "", Value: "/var/lib/nf6-router/data", Usage: "where to store persistent data"})
	lib.AddOption(rootCmd, lib.Option{P: &dbUrl, Name: "dbUrl", Shorthand: "", Value: "dbname=nf6", Usage: "url of postgres database"})
	lib.AddOption(rootCmd, lib.Option{P: &port, Name: "port", Shorthand: "", Value: 51820, Usage: "wireguard listen port"})
	lib.AddOption(rootCmd, lib.Option{P: &wgPrivKeyStr, Name: "privateKey", Shorthand: "", Value: "", Usage: "wireguard private key"})
	lib.AddOption(rootCmd, lib.Option{P: &timeout, Name: "timeout", Shorthand: "", Value: 5 * time.Second, Usage: "timeout for requests"})
	lib.AddOption(rootCmd, lib.Option{P: &deviceName, Name: "deviceName", Shorthand: "", Value: "wg", Usage: "name of wireguard device"})

	rootCmd.AddCommand(listenCmd)
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		cfgFile = home + "/.config/nf6-router/config.yaml"
	}
	viper.SetConfigFile(cfgFile)

	if _, err := os.Stat(cfgFile); err != nil {
		genConfig()
	}

	if err := viper.ReadInConfig(); err == nil {
		lib.LoadOptions()
	}

	if shouldSaveConfig {
		genConfig()
	}
}

func genConfig() {
	cfgFileDir := path.Dir(cfgFile)
	if err := os.MkdirAll(cfgFileDir, os.ModePerm); err != nil {
		log.Println("failed to make config directory: ", err)
	}
	if _, err := os.OpenFile(cfgFile, os.O_CREATE|os.O_RDONLY, 0600); err != nil {
		log.Println("failed to create config file: ", err)
	}
	if err := viper.WriteConfig(); err != nil {
		log.Println("failed to generate config: ", err)
	}
}

func initDataDir() {
	socket = dataDir + "/socket.sock"

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
