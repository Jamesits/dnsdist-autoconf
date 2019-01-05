package main

import (
	"fmt"
	"io"
)

func DomainList(c map[string]interface{}, o io.Writer) {
	poolName := fmt.Sprintf("DomainList-%s", randomString(6))
	var servers []DnsServer

	for _, server := range emptyInterfaceToStringArray(c["upstreams"]) {
		servers = append(servers, DnsServer{
			address: server,
		})
	}
	generateServerPool(poolName, servers, o)
	generateActions(poolName, emptyInterfaceToStringArray(c["domains"]), c["action"].(string), o)
}
