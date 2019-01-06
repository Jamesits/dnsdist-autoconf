package main

import (
	"math/rand"
	"reflect"
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
