package main

type config struct {
	QuitOnError      bool                     `toml:"quit_on_error"`
	Listen           []string                 `toml:"listen"`
	Upstreams        []string                 `toml:"upstreams"`
	ECS              ecs                      `toml:"ecs"`
	ControlSocket    controlSocket            `toml:"control_socket"`
	WebServer        webServer                `toml:"web_server"`
	AllowDDNSUpdates bool                     `toml:"allow_ddns_updates"`
	Matches          []map[string]interface{} `toml:"match"`
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
