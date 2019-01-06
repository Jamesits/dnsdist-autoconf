package main

import (
	"fmt"
	"io"
	"strings"
)

func generateServerPool(pool string, servers []DnsServer, o io.Writer) {
	// create newServer() blocks
	for _, server := range servers {
		if len(server.name) > 0 {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", name=\"%s\", pool=\"%s\"})\n",
				server.address,
				server.name,
				pool,
			)
			check(err)
		} else {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", pool=\"%s\"})\n",
				server.address,
				pool,
			)
			check(err)
		}
	}

}

// generate a SuffixMatchNode
// template from https://github.com/PowerDNS/pdns/issues/5433#issuecomment-309860659
// TODO: for large list, we can possibly load it from external text files, like
//		for line in io.lines("/etc/dnsdist/domains.txt") do auto_domain_list_x:add(line) end
// TODO: filter name so it can be a valid Lua variable name
func generateDomainList(name string, domains []string, o io.Writer) string {
	var err error
	// if name is empty, generate a random name
	if len(name) == 0 {
		name = fmt.Sprintf("auto_domain_list_%s", randomString(6))
	}

	_, err = fmt.Fprintf(o, "%s = newSuffixMatchNode()\n", name)
	check(err)

	for _, domain := range domains {
		_, err = fmt.Fprintf(o, "%s:add(\"%s\")\n", name, domain)
		check(err)
	}

	return name
}

// create addAction() blocks
func generateAction(pool string, domainList string, action string, o io.Writer) {
	var err error

	// SuffixMatchNodeRule(..., false) prevents all domains being displayed on the web page
	// as per https://github.com/PowerDNS/pdns/issues/7332#issuecomment-451681228
	_, err = fmt.Fprintf(o, "addAction(SuffixMatchNodeRule(%s, false), ", domainList)
	check(err)

	switch strings.ToLower(action) {
	case "resolve":
		_, err = fmt.Fprintf(o, "PoolAction(\"%s\")", pool)
	case "servfail":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.SERVFAIL)")
	case "block":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.NXDOMAIN)")
	case "nxdomain":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.NXDOMAIN)")
	case "refuse":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.REFUSED)")
	case "drop":
		_, err = fmt.Fprint(o, "DropAction()")
	}
	check(err)

	_, err = fmt.Fprint(o, ")\n")
	check(err)
}

// function to generate a generic domain list / dns server config block
func generateActionFromDomains(pool string, domains []string, action string, o io.Writer) {
	// create SuffixMatchNode
	domainList := generateDomainList("", domains, o)

	// create addAction()
	generateAction(pool, domainList, action, o)
}
