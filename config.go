package main

type config struct {
	QuitOnError          bool                     `toml:"quit_on_error"`
	Listen               []string                 `toml:"listen"`
	Upstreams            []string                 `toml:"upstreams"`
	AllowedClientSubnets []string                 `toml:"allowed_client_subnets"`
	AllowDDNSUpdates     bool                     `toml:"allow_ddns_updates"`
	ECS                  ecs                      `toml:"ecs"`
	ControlSocket        controlSocket            `toml:"control_socket"`
	WebServer            webServer                `toml:"web_server"`
	Cache                cache                    `toml:"cache"`
	Pools                []pool                   `toml:"pool"`
	Matches              []map[string]interface{} `toml:"match"`
}

type ecs struct {
	Enabled         bool `toml:"enabled"`
	DefaultPrefixV4 int  `toml:"default_prefix_v4"`
	DefaultPrefixV6 int  `toml:"default_prefix_v6"`
}

type controlSocket struct {
	Listen string `toml:"listen"`
	Key    string `toml:"key"`
}

type webServer struct {
	Listen   string `toml:"listen"`
	Password string `toml:"password"`
	ApiKey   string `toml:"api_key"`
}

type cache struct {
	Enabled                     bool `toml:"enabled"`
	MaxEntries                  int  `toml:"max_entries"`
	MaxLifetime                 int  `toml:"max_lifetime"`
	MinTTL                      int  `toml:"min_ttl"`
	FailureResultTTL            int  `toml:"failure_result_ttl"`
	StaleResultTTL              int  `toml:"stale_result_ttl"`
	AvoidReduceCachedEntriesTTL bool `toml:"avoid_reduce_cached_entries_ttl"`
}

type pool struct {
	Name    string   `toml:"name"`
	Servers []string `toml:"servers"`
}
