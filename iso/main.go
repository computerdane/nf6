package iso

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	AccountSshPubKey string
	HostAddr6        string
	HostSystem       string
	HostWgPrivKey    string
	VipWgEndpoint    string
	VipWgPubKey      string
}

//go:embed flake.nix
var flakeNix []byte

func Generate(dir string, config *Config) (isoPath string, err error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", status.Error(codes.Unknown, fmt.Sprintf("failed to create directory %s", dir))
	}
	if err := os.WriteFile(dir+"/flake.nix", flakeNix, os.ModePerm); err != nil {
		return "", status.Error(codes.Unknown, "failed to write flake.nix")
	}
	configJson, err := json.Marshal(config)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "failed to marshal iso config as json")
	}
	if err := os.WriteFile(dir+"/config.json", configJson, os.ModePerm); err != nil {
		return "", status.Error(codes.Unknown, "failed to write config.json")
	}
	cmd := exec.Command("nix", "build", ".#nixosConfigurations.nf6.config.system.build.isoImage")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", status.Error(codes.Unknown, "failed to build iso")
	}
	outDir, err := filepath.EvalSymlinks(dir + "/result")
	if err != nil {
		return "", status.Error(codes.Unknown, "failed to resolve symlink for result")
	}
	err = filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() && strings.HasSuffix(path, ".iso") {
			isoPath = path
			return nil
		}
		return err
	})
	return isoPath, nil
}
