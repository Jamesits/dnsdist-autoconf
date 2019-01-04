package main

import "io"

type MatchListProvider func(map[string]interface{}, io.Writer)

var providers = map[string]MatchListProvider{
	"DomainList": DomainList,
	"ActiveDirectory": ActiveDirectory,
	"gfwlist": GfwList,
	"dnsmasq-china-list": DnsmasqChinaList,
}
