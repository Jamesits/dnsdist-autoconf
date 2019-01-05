package main

type config struct {
	QuitOnError         bool                     `toml:"quit_on_error"`
	Listen              []string                 `toml:"listen"`
	Upstreams           []string                 `toml:"upstreams"`
	ECS                 bool                     `toml:"enable_ecs"`
	DefaultEcsPrefixV4  int8                     `toml:"default_ecs_prefix_v4"`
	DefaultEcsPrefixV6  int8                     `toml:"default_ecs_prefix_v6"`
	ControlSocketListen string                   `toml:"control_socket"`
	ControlSocketKey    string                   `toml:"control_socket_key"`
	WebServerListen     string                   `toml:"web_server"`
	WebServerPassword   string                   `toml:"web_server_password"`
	AllowDDNSUpdates    bool                     `toml:"allow_ddns_updates"`
	Matches             []map[string]interface{} `toml:"match"`
}
