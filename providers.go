package main

import "io"

type MatchListProvider func(map[string]interface{}, io.Writer)

// define a list of provider functions
var providers = map[string]MatchListProvider{
	"DomainList":      DomainList,
	"ActiveDirectory": ActiveDirectory,
	// "gfwlist":            GfwList,
	"dnsmasq-china-list": DnsmasqChinaList,
}

// provide a convenient helper to generate the same set of config for a provider that:
// * only generate a single pool
// it currently does:
// * assign an packet cache
func generateDefaultProviderTasks(poolName string, c map[string]interface{}, o io.Writer) {
	// cache
	if conf.Cache.Enabled {
		assignCache(poolName, globalPacketCache, o)
	}
}
