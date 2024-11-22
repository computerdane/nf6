package lib

import (
	"crypto/rand"
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

	randIp := make([]byte, net.IPv6len)
	_, err := rand.Read(randIp)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate random IPv6 address")
	}
	ip := make([]byte, net.IPv6len)
	mask := net.CIDRMask(length, net.IPv6len*8)
	for i := range net.IPv6len {
		ip[i] = from.IP[i] | ((mask[i] ^ from.Mask[i]) & randIp[i])
	}

	return &net.IPNet{
		IP:   ip,
		Mask: mask,
	}, nil
}
