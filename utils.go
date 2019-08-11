package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"reflect"
	"regexp"
)

const OutputCommentPrefix = "--"

type Mode int

const (
	MODE_AUTO Mode = iota
	MODE_IPv4
	MODE_IPv6
)

// struct see: https://dnsdist.org/reference/config.html#newServer
type DnsServer struct {
	address string // "IP:PORT" format
	name    string // for display purpose
}

// code from: https://stackoverflow.com/a/13906031
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// cursed type cast assuming every element inside is a string
// code from: https://gist.github.com/pmn/5374494
func emptyInterfaceToStringArray(i interface{}) []string {
	var o []string

	obj := reflect.ValueOf(i)

	// work around panic: reflect: call of reflect.Value.Len on zero Value
	if !IsZeroOfUnderlyingType(obj) {
		count := obj.Len()
		for index := 0; index < count; index++ {
			elem := obj.Index(index)
			o = append(o, elem.Interface().(string))
		}
	}

	return o
}

const randomLetters = "0123456789ABCDEF"

// generates a random string
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomLetters[rand.Intn(len(randomLetters))]
	}
	return string(b)
}

// contains
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// reverse
func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func getFileHandle(n string) (io.WriteCloser, string) {
	fullPath := path.Join(*configDir, n)
	outputFile, err := os.Create(fullPath)
	check(err)

	return outputFile, fullPath
}

func hashString(s string, len int) string {
	h := md5.New()
	_, err := io.WriteString(h, s)
	check(err)
	var out string
	out = fmt.Sprintf("%x", h.Sum(nil))

	return out[0:len]
}

func stringToIdentifier(s string, maxLen int) string {
	re := regexp.MustCompile(`[^0-9a-zA-Z]`)
	out := re.ReplaceAllString(s, `_`)
	if len(out) <= maxLen {
		return out
	}
	return out[0:maxLen]
}
