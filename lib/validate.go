package lib

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"

	"github.com/manifoldco/promptui"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ValidNixSystems = []string{"x86_64-linux", "aarch64-linux"}

func PromptOrValidate(value *string, prompt *promptui.Prompt) error {
	if *value == "" {
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		*value = result
	} else {
		if err := prompt.Validate(*value); err != nil {
			return err
		}
	}
	return nil
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return status.Error(codes.InvalidArgument, "Invalid email")
	}
	return nil
}

func ValidateHostName(name string) error {
	if match, _ := regexp.MatchString(`^[A-Za-z0-9]([A-Za-z0-9_-]{0,61}[A-Za-z0-9])?$`, name); match {
		return nil
	}
	return status.Error(codes.InvalidArgument, "Invalid host name. Must be alphanumeric with optional underscores in the middle.")
}

func ValidateIpv6Address(addr string) error {
	ip := net.ParseIP(addr)
	if ip == nil || ip.To4() != nil {
		return status.Error(codes.InvalidArgument, "Invalid IPv6 address")
	}
	return nil
}

func ValidateIpv6Prefix(addr string, prefixLen int) error {
	_, ipNet, err := net.ParseCIDR(addr)
	ones, bits := ipNet.Mask.Size()
	if err != nil || ones != prefixLen || bits != 128 {
		return status.Error(codes.InvalidArgument, "Invalid IPv6 prefix")
	}
	return nil
}

func ValidateNixSystem(system string) error {
	for _, validSystem := range ValidNixSystems {
		if system == validSystem {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf("System must be one of: %v", ValidNixSystems))
}

func ValidateRepoName(name string) error {
	if matchValid, _ := regexp.MatchString(`^[A-Za-z0-9]+[A-Za-z0-9\-_]+[A-Za-z0-9]+$`, name); matchValid {
		if matchInvalid, _ := regexp.MatchString(`^.*(\-\-|__|\-_|_\-).*$`, name); !matchInvalid {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, "Repo name must be at least 3 characters. Repo name must only contain characters A-Z, a-z, 0-9, -, and _. Repo name must not start or end with - or _. Repo name must not have two or more consecutive - and/or _.")
}

func ValidateWireguardKey(key string) error {
	_, err := wgtypes.ParseKey(key)
	if err != nil {
		return status.Error(codes.InvalidArgument, "Invalid WireGuard key")
	}
	return nil
}
