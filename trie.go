package main

import (
	"fmt"
	"strings"
)

// a simple trie implementation for domain matching
// supports wildcard

type Trie struct {
	Content  string
	Children []*Trie
}

func trieNodeContains(s *Trie, e string) (bool, *Trie) {
	for _, a := range s.Children {
		if a.Content == e {
			return true, a
		}
	}
	return false, nil
}

// Insert("www.example.com")
// Insert("www.*.example.com")
func (t *Trie) Insert(domain string) {
	// remove the rightmost dot for standard compliance
	domain = strings.TrimRight(domain, ".")
	parts := strings.Split(domain, ".")

	currentNode := t

	for i := len(parts) - 1; i >= 0; i-- {
		if ok, node := trieNodeContains(currentNode, parts[i]); ok {
			// already exists
			currentNode = node
		} else {
			// create new node
			newNode := Trie{Content: parts[i]}
			currentNode.Children = append(currentNode.Children, &newNode)
			currentNode = &newNode
		}
	}
}

func (t *Trie) RecursiveMatch(reversedDomainParts []string) bool {
	if len(reversedDomainParts) == 0 {
		// log.Printf("Match succeed because domain is shorter than database")
		return true
	}

	if len(t.Children) == 0 {
		// log.Printf("Match failed because domain is shorter than database")
		return false
	}
	for _, subNode := range t.Children {
		if WildcardPatternMatch(reversedDomainParts[0], subNode.Content) {
			// log.Printf("Partial matched domain: %s database: %s\n", reversedDomainParts[0], subNode.Content)
			ret := subNode.RecursiveMatch(reversedDomainParts[1:])
			if ret {
				return true
			}
		} else {
			// log.Printf("Partial match failed domain: %s database: %s\n", reversedDomainParts[0], subNode.Content)
		}
	}

	return false
}

func (t *Trie) Match(domain string) bool {
	// remove the rightmost dot for standard compliance
	domain = strings.TrimRight(domain, ".")
	parts := reverse(strings.Split(domain, "."))

	currentNode := t

	// fmt.Printf("-- %s\n", domain)
	return currentNode.RecursiveMatch(parts)
}

// print the trie - for debug purpose only
func TrieWalk(t *Trie, level int) {
	for i := 0; i < level-2; i++ {
		fmt.Print("  ")
	}
	if level > 1 {
		fmt.Print("└─")
	}

	fmt.Print(t.Content)
	fmt.Print("\n")
	for _, subNode := range t.Children {
		TrieWalk(subNode, level+1)
	}
}
