package main

// provides parser for common file types

import (
	"bufio"
	"io"
	"strings"
)

// generate a domain list from the following format:
// server=/example.com/114.114.114.114
// we are making a lot of assumptions here
func generateDomainListFromDnsmasqConfig(i io.Reader) []string {
	var o []string

	scanner := bufio.NewScanner(i)
	for scanner.Scan() {
		// normalize
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if !strings.HasPrefix(line, "server=/") {
			break
		}

		sp := strings.Split(line, "/")
		if len(sp) < 3 {
			break
		}
		o = append(o, sp[1])
	}

	return o
}
