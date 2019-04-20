package main

// parse dnsmasq config file into a domain list

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// generate a domain list from the following format:
// server=/example.com/114.114.114.114
// we are making a lot of assumptions here
func generateDomainListFromDnsmasqConfig(i io.Reader) []string {
	var o []string

	countSuccessLine := 0
	countPrefixIncorrectLine := 0
	countCannotSplit := 0

	scanner := bufio.NewScanner(i)
	for scanner.Scan() {
		// normalize
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if !strings.HasPrefix(line, "server=/") {
			countPrefixIncorrectLine += 1
			log.Printf("PrefixIncorrect: %s\n", line)
			continue
		}

		sp := strings.Split(line, "/")
		if len(sp) < 3 {
			countCannotSplit += 1
			log.Printf("CannotSplit: %s\n", line)
			continue
		}

		countSuccessLine += 1
		o = append(o, sp[1])
	}

	log.Printf("Success: %d, PrefixIncorrect: %d, CannotSplit: %d\n", countSuccessLine, countPrefixIncorrectLine, countCannotSplit)
	return o
}
