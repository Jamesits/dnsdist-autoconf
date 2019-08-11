package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func setServerPolicy(policy string, o io.Writer) {
	_, err := fmt.Fprintf(
		o, "\nsetServerPolicy(%s)\n",
		policy,
	)
	check(err)
}

// generate newServer() blocks from a servers=[] inside a [[match]]
func generateServerPoolInline(pool string, servers []DnsServer, o io.Writer) {
	// create newServer() blocks
	for _, server := range servers {
		if len(server.name) > 0 {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", name=\"%s\", pool=\"%s\", useClientSubnet=%t, checkInterval=%d})\n",
				server.address,
				server.name,
				pool,
				conf.ECS.Enabled,
				conf.HealthCheckInterval,
			)
			check(err)
		} else {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", pool=\"%s\", useClientSubnet=%t, checkInterval=%d})\n",
				server.address,
				pool,
				conf.ECS.Enabled,
				conf.HealthCheckInterval,
			)
			check(err)
		}
	}
}

// generate newServer() blocks from a [[pool]] definition
func generateServerPool(pool pool, o io.Writer) {
	for _, server := range pool.Servers {
		_, err := fmt.Fprintf(
			o, "newServer({address=\"%s\", pool=\"%s\", useClientSubnet=%t, checkInterval=%d})\n",
			server,
			pool.Name,
			conf.ECS.Enabled,
			conf.HealthCheckInterval,
		)
		check(err)
	}
}

// generate a SuffixMatchNode
// TODO: filter name so it can be a valid Lua variable name
func generateDomainList(name string, domains []string, o io.Writer) string {
	var err error
	// if name is empty, generate a random name
	if len(name) == 0 {
		name = fmt.Sprintf("auto_domain_list_%s", randomString(6))
	}

	_, err = fmt.Fprintf(o, "%s = newSuffixMatchNode()\n", name)
	check(err)

	listFile, fullPath := getFileHandle(name + ".list")
	defer listFile.Close()

	for _, domain := range domains {
		_, err = fmt.Fprintln(listFile, domain)
		check(err)
	}

	_, err = fmt.Fprintf(o, "for line in io.lines(\"%s\") do %s:add(newDNSName(line)) end\n", fullPath, name)
	check(err)

	return name
}

// create addAction() blocks
func generateAction(pool string, domainList string, action string, o io.Writer) {
	var err error

	hideDomainListOnWebPortal := false
	if len(domainList) > 10 {
		hideDomainListOnWebPortal = true
	}

	// SuffixMatchNodeRule(..., true) prevents all domains being displayed on the web page
	// as per https://github.com/PowerDNS/pdns/issues/7332#issuecomment-451681228
	_, err = fmt.Fprintf(o, "addAction(SuffixMatchNodeRule(%s, %s), ", domainList, strconv.FormatBool(hideDomainListOnWebPortal))
	check(err)

	switch strings.ToLower(action) {
	case "resolv": // compatibility
		fallthrough
	case "resolve":
		_, err = fmt.Fprintf(o, "PoolAction(\"%s\")", pool)
	case "servfail":
		_, err = fmt.Fprint(o, "RCodeAction(DNSRCode.SERVFAIL)")
	case "block":
		// return NXDOMAIN so the request would be cached by client (thank @m13253)
		fallthrough
	case "nxdomain":
		_, err = fmt.Fprint(o, "RCodeAction(DNSRCode.NXDOMAIN)")
	case "refuse":
		_, err = fmt.Fprint(o, "RCodeAction(DNSRCode.REFUSED)")
	case "drop":
		_, err = fmt.Fprint(o, "DropAction()")
	}
	check(err)

	_, err = fmt.Fprint(o, ")\n")
	check(err)
}

// function to generate a generic domain list / dns server config block
func generateActionFromDomains(domainListName string, pool string, domains []string, action string, o io.Writer) {
	// create SuffixMatchNode
	domainList := generateDomainList(stringToIdentifier(domainListName, 32), domains, o)

	// create addAction()
	generateAction(pool, domainList, action, o)
}

// Create a packet cache
// https://dnsdist.org/guides/cache.html
func createCache(name string, cacheConfig cache, o io.Writer) string {
	// if name is empty, generate a random name
	if len(name) == 0 {
		name = fmt.Sprintf("cache_%s", randomString(6))
	}
	_, err := fmt.Fprintf(o, "%s = newPacketCache(%d,{maxTTL=%d, minTTL=%d, temporaryFailureTTL=%d, staleTTL=%d})\n",
		name,
		cacheConfig.MaxEntries,
		cacheConfig.MaxLifetime,
		cacheConfig.MinTTL,
		cacheConfig.FailureResultTTL,
		cacheConfig.StaleResultTTL,
	)
	check(err)

	return name
}

// Assign a packet cache to a pool
func assignCache(poolName string, cacheName string, o io.Writer) {
	_, err := fmt.Fprintf(o, "getPool(\"%s\"):setCache(%s)\n", poolName, cacheName)
	check(err)
}
