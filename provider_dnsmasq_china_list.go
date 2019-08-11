package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func DnsmasqChinaList(index int, c map[string]interface{}, o io.Writer) {
	var servers []DnsServer
	const poolName = "dnsmasq-china-list"

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}
	generateServerPoolInline(poolName, servers, o)

	for index, url := range localDomainListUrls {
		log.Printf("Downloading rule %s...\n", url)

		resp, err := http.Get(url)
		// currently if download fail then go on
		// TODO: retry
		if softFail(err) != nil {
			continue
		}
		defer resp.Body.Close()

		domains := generateDomainListFromDnsmasqConfig(resp.Body)
		_, err = fmt.Fprintf(o, "%s Domain list generated from %s\n", OutputCommentPrefix, url)
		check(err)

		generateActionFromDomains(fmt.Sprintf("%s-%d", poolName, index), poolName, domains, c["action"].(string), o)

		generateDefaultProviderTasks(poolName, c, o)

		_, err = fmt.Fprintf(o, "\n")
		check(err)

		log.Printf("Generated %d rules\n", len(domains))
	}
}
