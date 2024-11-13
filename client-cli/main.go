package main

import (
	"flag"
	"log"
	"os"

	openssl "github.com/computerdane/nf6/lib"
)

var (
	baseDir = flag.String("base-dir", "", "location of api data")
	sslDir  = flag.String("ssl-dir", "", "location of ssl data")

	ssl *openssl.Openssl
)

func mkdirAll(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if *baseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not find user home dir: %v", err)
			os.Exit(1)
		}
		*baseDir = homeDir + "/.nf6/client-cli"
	}
	if *sslDir == "" {
		*sslDir = *baseDir + "/ssl"
	}

	mkdirAll(*baseDir)
	mkdirAll(*sslDir)

	ssl = &openssl.Openssl{Dir: *sslDir}
	err := ssl.GenConfigFile()
	if err != nil {
		log.Fatalf("failed to generate ssl config file: %v", err)
		os.Exit(1)
	}
}
