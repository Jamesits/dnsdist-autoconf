package main

import "testing"

func TestTrie(t *testing.T) {
	trie := Trie{}
	trie.Insert("example.com")
	trie.Insert("*.example.com")
	trie.Insert("www.*.example.com")
	trie.Insert("www.example.org")
	trie.Insert("ftp.example.org")

	TrieWalk(&trie, 0)

	if !trie.Match("example.com") {
		t.Error()
	}

	if !trie.Match("a.example.com") {
		t.Error()
	}

	if trie.Match("verylong.a.example.com") {
		t.Error()
	}
}
