package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func DomainListUrl(index int, c map[string]interface{}, o io.Writer) {
	randomSuffix := randomString(6)
	poolName := fmt.Sprintf("DomainListUrl-%s", randomSuffix)
	// domainListName := fmt.Sprintf("domain_list_%s", randomSuffix)
	var servers []DnsServer

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}
	generateServerPoolInline(poolName, servers, o)

	for _, url := range localDomainListUrls {
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

		generateActionFromDomains(poolName, poolName, domains, c["action"].(string), o)

		log.Printf("Generated %d rules\n", len(domains))
	}

	generateDefaultProviderTasks(poolName, c, o)
}
