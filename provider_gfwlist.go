package main

import "io"

// TODO: GfwList support
// it is bloated and hard to parse
func GfwList(c map[string]interface{}, o io.Writer) {
	const gfwListUrl = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

	// extract a list from gfwlist

}
