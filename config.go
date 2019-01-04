package main

type config struct {
	Listen           []string `toml:"listen"`
	Upstreams        []string `toml:"upstreams"`
	ECS				 bool     `toml:"enable_ecs"`
	DefaultEcsPrefixV4 int8 `toml:"default_ecs_prefix_v4"`
	DefaultEcsPrefixV6 int8 `toml:"default_ecs_prefix_v6"`
	Matches            []map[string]interface {} `toml:"match"`
}
