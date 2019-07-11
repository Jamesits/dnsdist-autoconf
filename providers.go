package main

import "io"

// Provider invocation function
// args:
// 	index int: a number that won't collide with other [[match]] blocks, can be used for naming
// 	c map[string]interface{}: the user-supplied config
//  o io.Writer: a binary buffer you should put your part of config to
type MatchListProvider func(int, map[string]interface{}, io.Writer)

// define a list of provider functions
var providers = map[string]MatchListProvider{
	"DomainList":      DomainList,
	"ActiveDirectory": ActiveDirectory,
	// "gfwlist":            GfwList,
	"dnsmasq-china-list": DnsmasqChinaList,
	"DomainListUrl":      DomainListUrl,
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
