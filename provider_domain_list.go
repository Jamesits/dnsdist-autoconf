package main

import (
	"fmt"
	"io"
)

func DomainList(c map[string]interface{}, o io.Writer) {
	randomSuffix := randomString(6)
	poolName := fmt.Sprintf("DomainList-%s", randomSuffix)
	domainListName := fmt.Sprintf("domain_list_%s", randomSuffix)
	var servers []DnsServer

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}
	generateServerPool(poolName, servers, o)
	generateDomainList(domainListName, emptyInterfaceToStringArray(c["domains"]), o)
	generateAction(poolName, domainListName, c["action"].(string), o)

	// cache
	if conf.Cache.Enabled {
		assignCache(poolName, globalPacketCache, o)
	}
}
