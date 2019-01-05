package main

import (
	"context"
	"net"
)

// force Dialer to connect via IPv4/v6 and a specific server
// code from https://github.com/Jamesits/myip/blob/master/src/opendns-dns.go
func dialContextFactory(mode Mode, server string) func(context.Context, string, string) (net.Conn, error) {
	var ret func(context.Context, string, string) (net.Conn, error)
	switch mode {
	case MODE_IPv4:
		ret = func(ctx context.Context, network, address string) (net.Conn, error) {
			switch network {
			case "tcp", "tcp4", "tcp6":
				network = "tcp4"
			case "udp", "udp4", "udp6":
				network = "udp4"
			}

			return (&net.Dialer{}).DialContext(ctx, network, server)
		}
	case MODE_IPv6:
		ret = func(ctx context.Context, network, address string) (net.Conn, error) {
			switch network {
			case "tcp", "tcp4", "tcp6":
				network = "tcp6"
			case "udp", "udp4", "udp6":
				network = "udp6"
			}
			return (&net.Dialer{}).DialContext(ctx, network, server)
		}
	default:
		ret = func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{DualStack: true}).DialContext(ctx, network, server)
		}
	}

	return ret
}

func newResolver() *net.Resolver {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	return resolver
}

func newCustomResolver(mode Mode, server string) *net.Resolver {
	resolver := newResolver()
	resolver.Dial = dialContextFactory(mode, server)
	return resolver
}
