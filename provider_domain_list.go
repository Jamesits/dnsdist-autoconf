package main

import (
	"fmt"
	"io"
)

func DomainList(index int, c map[string]interface{}, o io.Writer) {
	//randomSuffix := randomString(6)
	//poolName := fmt.Sprintf("DomainList-%s", randomSuffix)
	//domainListName := fmt.Sprintf("domain_list_%s", randomSuffix)
	poolName := fmt.Sprintf("DomainList-%d", index+1)
	domainListName := fmt.Sprintf("domain_list_%d", index+1)
	var servers []DnsServer

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}

	generateServerPoolInline(poolName, servers, o)
	generateDomainList(domainListName, emptyInterfaceToStringArray(c["domains"]), o)
	generateAction(poolName, domainListName, c["action"].(string), o)

	generateDefaultProviderTasks(poolName, c, o)
}
