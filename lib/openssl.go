package lib

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

type Openssl struct {
	Dir string
}

func (ssl *Openssl) GetPath(name string) string {
	return ssl.Dir + "/" + name
}

func (ssl *Openssl) GetConfigPath() string {
	return ssl.GetPath("openssl.cnf")
}

func (ssl *Openssl) GetExtConfigPath() string {
	return ssl.GetPath("openssl_ext.cnf")
}

func (ssl *Openssl) InitCmd(cmd *exec.Cmd) {
	cmd.Dir = ssl.Dir
	cmd.Env = []string{"OPENSSL_CONF=" + ssl.GetConfigPath()}
}

func (ssl *Openssl) FileExists(name string) (bool, error) {
	_, err := os.Stat(ssl.Dir + "/" + name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (ssl *Openssl) RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	ssl.InitCmd(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run command: %s", string(output))
		return err
	}
	return nil
}

func (ssl *Openssl) RunCommandInMemory(cmd *exec.Cmd, in string) ([]byte, error) {
	ssl.InitCmd(cmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, in)
	if err != nil {
		return nil, err
	}
	return cmd.Output()
}

func (ssl *Openssl) GenConfigFiles() error {
	file, err := os.Create(ssl.GetConfigPath())
	if err != nil {
		return err
	}
	extFile, err := os.Create(ssl.GetExtConfigPath())
	if err != nil {
		return err
	}

	file.WriteString(`
[ req ] 
encrypt_key = no
default_md = sha256
prompt = no
utf8 = yes
distinguished_name = my_req_distinguished_name

[ my_req_distinguished_name ]
C = AA
ST = a
L = a
O  = a
CN = a
`)
	extFile.WriteString(`
subjectAltName = DNS:a
`)
	return nil
}

func (ssl *Openssl) GenKey(outPath string) error {
	exists, err := ssl.FileExists(outPath)
	if !exists {
		return ssl.RunCommand("openssl", "genpkey", "-algorithm", "ED25519", "-out", outPath)
	}
	return err
}

func (ssl *Openssl) GenCert(keyPath string, outPath string) error {
	exists, err := ssl.FileExists(outPath)
	if !exists {
		return ssl.RunCommand("openssl", "req", "-addext", "subjectAltName = DNS:a", "-new", "-x509", "-key", keyPath, "-out", outPath)
	}
	return err
}

func (ssl *Openssl) GenCsr(keyPath string, outPath string) error {
	exists, err := ssl.FileExists(outPath)
	if !exists {
		return ssl.RunCommand("openssl", "req", "-addext", "subjectAltName = DNS:a", "-new", "-key", keyPath, "-out", outPath)
	}
	return err
}

func (ssl *Openssl) GenCertFromCsr(csrPath string, caKeyPath string, caCertPath string, outPath string) error {
	exists, err := ssl.FileExists(outPath)
	if !exists {
		return ssl.RunCommand("openssl", "x509", "-extfile", ssl.GetExtConfigPath(), "-req", "-CA", caCertPath, "-CAkey", caKeyPath, "-in", csrPath, "-out", outPath)
	}
	return err
}

func (ssl *Openssl) GenCertFromCsrInMemory(csr string, caKeyPath string, caCertPath string) ([]byte, error) {
	cmd := exec.Command("openssl", "x509", "-req", "-CA", caCertPath, "-CAkey", caKeyPath)
	return ssl.RunCommandInMemory(cmd, csr)
}

func (ssl *Openssl) GetPublicKeyInMemory(cert string) ([]byte, error) {
	cmd := exec.Command("openssl", "x509", "-pubkey", "-noout")
	return ssl.RunCommandInMemory(cmd, cert)
}
