package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	dbUrl        string
	gitReposPath string
	gitShell     string
	port         int
	timeout      time.Duration

	db *pgxpool.Pool

	query = "select repo.account_id, repo.name from repo inner join account on repo.account_id = account.id and account.ssh_public_key like $1"

	authorizedKeyPrefix string
)

var rootCmd = &cobra.Command{
	Use:   "nf6-git-auth",
	Short: "Nf6 Git Auth server",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "listen" {
			log.Fatal("please pass 'listen' as first argument")
		}

		var err error
		db, err = pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("listening at %v", listener.Addr())

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
	defer conn.Close()
	buffer := make([]byte, 80)
	size, err := conn.Read(buffer)

	pubKey := string(buffer[:size])

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rows, err := db.Query(ctx, query, pubKey+"%")
	if err != nil {
		return
	}

	var authorizedKey strings.Builder
	if _, err := authorizedKey.WriteString(authorizedKeyPrefix); err != nil {
		return
	}
	wroteOne := false
	for rows.Next() {
		var accountId uint64 = 0
		repoName := ""
		err := rows.Scan(&accountId, &repoName)
		if err != nil {
			return
		}
		if match, _ := regexp.MatchString(`^[A-Za-z0-9\-_]+$`, repoName); !match {
			return
		}
		if accountId != 0 && repoName != "" {
			if _, err := authorizedKey.WriteRune(' '); err != nil {
				return
			}
			if _, err := authorizedKey.WriteString(gitReposPath); err != nil {
				return
			}
			if _, err := authorizedKey.WriteRune('/'); err != nil {
				return
			}
			if _, err := authorizedKey.WriteString(fmt.Sprintf("%d", accountId)); err != nil {
				return
			}
			if _, err := authorizedKey.WriteRune('/'); err != nil {
				return
			}
			if _, err := authorizedKey.WriteString(repoName); err != nil {
				return
			}
			if _, err := authorizedKey.WriteString(".git"); err != nil {
				return
			}
		}
		wroteOne = true
	}
	if !wroteOne {
		return
	}
	if _, err := authorizedKey.WriteString(`" `); err != nil {
		return
	}
	if _, err := authorizedKey.WriteString(pubKey); err != nil {
		return
	}

	conn.Write([]byte(authorizedKey.String()))
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/var/lib/nf6-git-auth/config/config.yaml", "config file")
	rootCmd.PersistentFlags().StringVar(&dbUrl, "dbUrl", "dbname=nf6", "url of postgres database")
	rootCmd.PersistentFlags().StringVar(&gitReposPath, "gitReposPath", "/var/lib/nf6-git/repos", "location of git repos")
	rootCmd.PersistentFlags().StringVar(&gitShell, "gitShell", "/bin/nf6-git-shell", "location of git-shell executable")
	rootCmd.PersistentFlags().IntVar(&port, "port", 6970, "listen port")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "timeout for requests")

	viper.BindPFlag("dbUrl", rootCmd.PersistentFlags().Lookup("dbUrl"))
	viper.BindPFlag("gitShell", rootCmd.PersistentFlags().Lookup("gitShell"))
	viper.BindPFlag("gitReposPath", rootCmd.PersistentFlags().Lookup("gitReposPath"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	if err := viper.ReadInConfig(); err == nil {
		dbUrl = viper.GetString("dbUrl")
		gitReposPath = viper.GetString("gitReposPath")
		gitShell = viper.GetString("gitShell")
		port = viper.GetInt("port")
		timeout = viper.GetDuration("timeout")
	}

	authorizedKeyPrefix = `command="` + gitShell
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
