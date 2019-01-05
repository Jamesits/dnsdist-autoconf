package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// generate ruleset for https://github.com/felixonmars/dnsmasq-china-list
// rules:
//   resolved using local dns upstreams:
//      accelerated-domains.china.conf
//      apple.china.conf
//      google.china.conf
//   TODO: blacklisted result:
//      bogus-nxdomain.china.conf

var localDomainListUrls = []string{
	"https://github.com/felixonmars/dnsmasq-china-list/raw/master/accelerated-domains.china.conf",
	"https://github.com/felixonmars/dnsmasq-china-list/raw/master/apple.china.conf",
	"https://github.com/felixonmars/dnsmasq-china-list/raw/master/google.china.conf",
}

func DnsmasqChinaList(c map[string]interface{}, o io.Writer) {
	var servers []DnsServer
	const poolName = "dnsmasq-china-list"

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}
	generateServerPool(poolName, servers, o)

	for _, url := range localDomainListUrls {
		resp, err := http.Get(url)
		check(err)
		defer resp.Body.Close()

		domains := generateDomainListFromDnsmasqConfig(resp.Body)
		_, err = fmt.Fprintf(o, "%s Domain list generated from %s\n", OutputCommentPrefix, url)
		check(err)

		generateActions(poolName, domains, c["action"].(string), o)
	}
}

// generate a domain list from the following format:
// server=/example.com/114.114.114.114
// we are making a lot of assumptions here
func generateDomainListFromDnsmasqConfig(i io.Reader) []string {
	var o []string

	scanner := bufio.NewScanner(i)
	for scanner.Scan() {
		// normalize
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if !strings.HasPrefix(line, "server=/") {
			break
		}

		sp := strings.Split(line, "/")
		if len(sp) < 3 {
			break
		}
		o = append(o, sp[1])
	}

	return o
}
