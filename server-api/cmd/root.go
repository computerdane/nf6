package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/computerdane/nf6/nf6"
	"github.com/computerdane/nf6/server-api/server_insecure"
	"github.com/computerdane/nf6/server-api/server_secure"
	"github.com/computerdane/nf6/server-api/ssl_util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	cfgFile string

	dataDir      string
	dbUrl        string
	portInsecure int
	portSecure   int

	sslDir string
)

var rootCmd = &cobra.Command{
	Use:   "nfapi",
	Short: "Nf6 API server",
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("loading certs...")

		ssl := &ssl_util.SslUtil{Dir: sslDir}
		err := ssl.GenCaFiles("ca")
		if err != nil {
			log.Fatalf("failed to generate ca files: %v", err)
		}
		err = ssl.GenCertFiles("ca", "server")
		if err != nil {
			log.Fatalf("failed to generate server cert files: %v", err)
		}
		caCert, err := os.ReadFile(sslDir + "/ca.crt")
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatal("failed to append certs from pem")
		}
		cert, err := tls.LoadX509KeyPair(sslDir+"/server.crt", sslDir+"/server.key")
		if err != nil {
			log.Fatal(err)
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
			RootCAs:      caCertPool,
		})

		log.Println("connecting to db...")

		db, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		numAccounts := 0
		if err := db.QueryRow(context.Background(), "select count(*) from account").Scan(&numAccounts); err != nil {
			log.Fatal(err)
		}
		log.Printf("connected to database with %d accounts", numAccounts)

		log.Println("creating grpc servers...")

		serverInsecure := grpc.NewServer()
		serverSecure := grpc.NewServer(grpc.Creds(creds))
		nf6.RegisterNf6InsecureServer(serverInsecure, server_insecure.NewServer(db, caCert, ssl))
		nf6.RegisterNf6SecureServer(serverSecure, server_secure.NewServer(db))

		log.Println("opening listeners...")

		listenerInsecure, err := net.Listen("tcp", fmt.Sprintf(":%d", portInsecure))
		if err != nil {
			log.Fatal(err)
		}
		listenerSecure, err := net.Listen("tcp", fmt.Sprintf(":%d", portSecure))
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			log.Printf("insecure server listening at %v", listenerInsecure.Addr())
			if err := serverInsecure.Serve(listenerInsecure); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()
		log.Printf("secure server listening at %v", listenerSecure.Addr())
		if err := serverSecure.Serve(listenerSecure); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig, initDataDir)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/var/lib/nfapi/config/config.yaml", "config file")
	rootCmd.PersistentFlags().StringVar(&dataDir, "dataDir", "/var/lib/nfapi/data", "where to store persistent data")
	rootCmd.PersistentFlags().StringVar(&dbUrl, "dbUrl", "dbname=nf6", "url of postgres database")
	rootCmd.PersistentFlags().IntVar(&portInsecure, "portInsecure", 6968, "port for insecure connections")
	rootCmd.PersistentFlags().IntVar(&portSecure, "portSecure", 6969, "port for secure connections")

	viper.BindPFlag("dataDir", rootCmd.PersistentFlags().Lookup("dataDir"))
	viper.BindPFlag("dbUrl", rootCmd.PersistentFlags().Lookup("dbUrl"))
	viper.BindPFlag("portInsecure", rootCmd.PersistentFlags().Lookup("portInsecure"))
	viper.BindPFlag("portSecure", rootCmd.PersistentFlags().Lookup("portSecure"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("could not read config file: %v", err)
	} else {
		dataDir = viper.GetString("dataDir")
		dbUrl = viper.GetString("dbUrl")
		portInsecure = viper.GetInt("portInsecure")
		portSecure = viper.GetInt("portSecure")
	}
}

func initDataDir() {
	sslDir = dataDir + "/ssl"

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(sslDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
