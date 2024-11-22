package lib

import (
	"crypto/rand"
	"fmt"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RandomIpv6Prefix(from *net.IPNet, length int) (*net.IPNet, error) {
	if length <= 0 {
		return nil, status.Error(codes.InvalidArgument, "length must be a positive number")
	}
	if from.IP.To4() != nil {
		return nil, status.Error(codes.InvalidArgument, "from address must be an IPv6 address")
	}
	{
		ones, _ := from.Mask.Size()
		if ones >= length {
			return nil, status.Error(codes.InvalidArgument, "length must be larger than from address prefix length")
		}
	}

	ip := make([]byte, net.IPv6len)
	_, err := rand.Read(ip)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate random IPv6 prefix")
	}
	mask := net.CIDRMask(length, net.IPv6len*8)
	for i := range net.IPv6len {
		ip[i] &= ^from.Mask[i]
		ip[i] |= from.IP[i]
		ip[i] &= mask[i]
	}

	return &net.IPNet{
		IP:   ip,
		Mask: mask,
	}, nil
}

func RandomIpv6Addr(from *net.IPNet) (net.IP, error) {
	if from.IP.To4() != nil {
		return nil, status.Error(codes.InvalidArgument, "from address must be an IPv6 address")
	}
	ipNet, err := RandomIpv6Prefix(from, net.IPv6len*8)
	if err != nil {
		return nil, err
	}
	return ipNet.IP, nil
}

func EnsureIpv6PrefixContainsAddr(prefix *net.IPNet, addr net.IP) error {
	if prefix.IP.To4() != nil || addr.To4() != nil {
		return status.Error(codes.InvalidArgument, "prefix and address must both be IPv6")
	}
	for i := range net.IPv6len {
		if addr[i]&prefix.Mask[i] != prefix.IP[i] {
			return status.Error(codes.PermissionDenied, fmt.Sprintf("ip must be in the prefix %s", prefix.String()))
		}
	}
	return nil
}
