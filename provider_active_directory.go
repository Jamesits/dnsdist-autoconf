package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// ActiveDirectory match list
// automatic query DC IP address assuming all DC have DNS role configured
func ActiveDirectory(c map[string]interface{}, o io.Writer) {
	// prepare
	ctx := context.Background()
	primaryDomain := c["primary_domain"].(string)
	var resolver *net.Resolver
	log.Printf("Querying domain %s...\n", primaryDomain)
	if val, ok := c["bootstrap_server"]; ok {
		resolver = newCustomResolver(MODE_AUTO, val.(string))
	} else {
		resolver = newResolver()
	}

	// gather a list of DCs
	// TODO: set current site name in config so we can query current site first
	_, records, err := resolver.LookupSRV(ctx, "ldap", "tcp", primaryDomain)
	check(err)

	var servers []DnsServer

	// TODO: sort by Priority and Weight?
	// not needed currently since AD DNS server seldom use them, they sort the DNS records
	for _, record := range records {
		// get IPs via address
		addresses, err := resolver.LookupHost(ctx, record.Target)
		check(err)
		log.Printf("Address: %s (%s, total %d), Priority: %d, Weight: %d\n", record.Target, addresses[0], len(addresses), record.Priority, record.Weight)
		for _, address := range addresses {
			name := fmt.Sprintf("%s (%s)", record.Target, address)

			var addressWithPort string
			if strings.Contains(address, ":") {
				// IPv6
				addressWithPort = fmt.Sprintf("[%s]:53", address)
			} else {
				// IPv4
				addressWithPort = fmt.Sprintf("%s:53", address)
			}
			servers = append(servers, DnsServer{
				name:    name,
				address: addressWithPort,
			})
		}

	}

	poolName := fmt.Sprintf("AD-%s", primaryDomain)
	var domains []string
	if val, ok := c["domains"]; ok {
		domains = emptyInterfaceToStringArray(val)
	}
	domains = append([]string{primaryDomain}, domains...)

	generateServerPool(poolName, servers, domains, c["action"].(string), o)

	// emptyInterfaceToStringArray(c["domains"])
}
