package main

// parse Adblock Plus rules into a domain list
// supports exclusion, supports wildcard
// doesn't support regex for now
// documentation: https://adblockplus.org/filters
// the documentation is not very clear on the syntax (no BNF, etc), lots of assumptions has been made

import (
	"bufio"
	"github.com/asaskevich/govalidator"
	"io"
	"log"
	"strings"
)

// supported (safe to ignore) rule options
var AdblockPlusRuleSupportedOptions = []string{
	"third-party",
}

// generate a domain list from a Adblock Plus rule
// reference: https://github.com/justdomains/ci/blob/master/convertlists.py
func generateDomainListFromAdblockPlusRule(i io.Reader) []string {
	var o []string

	lineCount := 0
	unsupportedCount := 0
	commentCount := 0
	exceptionCount := 0
	emptyLineCount := 0

	exceptionMatchTree := Trie{}

	scanner := bufio.NewScanner(i)
	for scanner.Scan() {
		// normalize
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		if len(line) == 0 {
			emptyLineCount++
			continue
		}

		switch line[0] {
		case '[':
			// the first line will be likely contain [Adblock Plus 2.0]
			fallthrough
		case '!':
			// comment
			commentCount++

		case '|':
			// might be domain name or exact address
			if !strings.HasPrefix(line, "||") {
				// exact address
				unsupportedCount++
			} else {
				// domain name

				// check if there are other options
				splitLine := strings.SplitN(line[2:], "$", 1)
				matchingParts := strings.SplitN(splitLine[0], "^", 1)
				if len(matchingParts) != 2 {
					// no '^' at the end
					unsupportedCount++
					continue
				}

				if len(matchingParts[1]) > 0 {
					// contains extra options
					unsupportedCount++
					continue
				}

				if !govalidator.IsDNSName(matchingParts[0]) {
					// not a valid hostname
					unsupportedCount++
					continue
				}

				if len(splitLine) > 1 {
					// we have to check options
					options := strings.Split(splitLine[1], ",")
					hasUnsupportedOption := false
					for _, option := range options {
						if !contains(AdblockPlusRuleSupportedOptions, option) {
							hasUnsupportedOption = true
							break
						}
					}
					if !hasUnsupportedOption {
						unsupportedCount++
						continue
					}
				}

				// yeah, this is correct!
				o = append(o, matchingParts[0])
			}

		case '@':
			// exception
			if !strings.HasPrefix(line, "@@||") {
				unsupportedCount++
			} else {
				// domain exception

				// strip options
				exceptionRule := strings.SplitN(line[4:], "$", 1)[0]

				// TODO: support regex exceptions
				if strings.HasPrefix(exceptionRule, "/") {
					log.Printf("regex exceptions is not supported for now, ignored. rule: %s\n", exceptionRule)
					unsupportedCount++
					continue
				} else {
					// wildcard rule
					exceptionCount++

					wildcardDomain := strings.TrimRight(strings.SplitN(exceptionRule, "^", 1)[0], ".")
					wildcardDomain = strings.SplitN(wildcardDomain, ":", 1)[0]
					wildcardDomain = strings.SplitN(wildcardDomain, "?", 1)[0]
					wildcardDomain = strings.SplitN(wildcardDomain, "/", 1)[0]

					exceptionMatchTree.Insert(wildcardDomain)
				}
			}

		case ':':
			// extended CSS selectors
			fallthrough
		case '#':
			// element selection
			fallthrough
		case '*':
			// URL matching
			fallthrough
		default:
			unsupportedCount++
		}
	}

	// remove exceptions
	var cleared []string

	for _, domain := range o {
		if !exceptionMatchTree.Match(domain) {
			cleared = append(cleared, domain)
		}
	}

	log.Printf("Total lines: %d\n", lineCount)
	log.Printf("Empty lines: %d\n", emptyLineCount)
	log.Printf("Comments: %d\n", commentCount)
	log.Printf("Unsupported rules: %d\n", unsupportedCount)
	log.Printf("Valid rules: %d\n", len(o))
	log.Printf("Exception rules: %d\n", exceptionCount)
	log.Printf("Final domains: %d\n", len(cleared))

	return cleared
}
