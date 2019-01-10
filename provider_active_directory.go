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
	if !conf.AllowDDNSUpdates {
		log.Println("Warning: if \"allow_ddns_updates\" is set to false, Active Directory DNS auto registration might not work.")
	}

	// prepare
	ctx := context.Background()
	primaryDomain := c["primary_domain"].(string)
	var resolver *net.Resolver
	var err error
	log.Printf("Querying domain %s...\n", primaryDomain)

	// gather a list of DCs
	// TODO: set current site name in config so we can query current site first
	var domain string
	var records []*net.SRV
	var bootstrapServers = emptyInterfaceToStringArray(c["bootstrap_servers"])

	// try each pre-defined one
	for _, server := range bootstrapServers {
		log.Printf("Trying server %s...\n", server)
		resolver = newCustomResolver(MODE_AUTO, server)
		domain, records, err = resolver.LookupSRV(ctx, "ldap", "tcp", primaryDomain)
		if softFail(err) != nil || len(records) == 0 {
			log.Println("FAILED")
			continue
		} else {
			log.Println("SUCCESS")
			break
		}
	}

	// try system server , if failed then softFail
	if len(records) == 0 {
		log.Println("Trying system DNS server...")
		resolver = newResolver()
		domain, records, err = resolver.LookupSRV(ctx, "ldap", "tcp", primaryDomain)
		if softFail(err) != nil {
			return
		}
	}

	// if no valid SRV record is found, that's an error
	if len(records) == 0 {
		if softFail(fmt.Errorf("no valid SRV records for %s", domain)) != nil {
			return
		}
	}

	var servers []DnsServer

	// TODO: sort by Priority and Weight?
	// not needed currently since AD DNS server seldom use them, they sort the DNS records
	for _, record := range records {
		// get IPs via address
		addresses, err := resolver.LookupHost(ctx, record.Target)
		if softFail(err) != nil {
			return
		}
		log.Printf("Address: %s (%s, total %d), Priority: %d, Weight: %d\n", record.Target, addresses[0], len(addresses), record.Priority, record.Weight)
		for _, address := range addresses {
			name := record.Target

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

	generateServerPool(poolName, servers, o)
	generateActionFromDomains(poolName, domains, c["action"].(string), o)

	// cache
	if conf.Cache.Enabled {
		assignCache(poolName, globalPacketCache, o)
	}

	// emptyInterfaceToStringArray(c["domains"])
}
