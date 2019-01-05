package main

import "io"

type MatchListProvider func(map[string]interface{}, io.Writer)

// define a list of provider functions
var providers = map[string]MatchListProvider{
	"DomainList":         DomainList,
	"ActiveDirectory":    ActiveDirectory,
	"gfwlist":            GfwList,
	"dnsmasq-china-list": DnsmasqChinaList,
}
