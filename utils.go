package main

import (
	"fmt"
	"io"
	"reflect"
)

const OutputCommentPrefix = "-- "

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

// cursed type cast assuming every element inside is a string
// code from: https://gist.github.com/pmn/5374494
func emptyInterfaceToStringArray(i interface{}) []string {
	var o []string

	obj := reflect.ValueOf(i)
	count := obj.Len()
	for index := 0; index < count; index++ {
		elem := obj.Index(index)
		o = append(o, elem.Interface().(string))
	}

	return o
}

// function to generate a generic domain list / dns server config block
func generateServerPool(pool string, servers []DnsServer, domains []string, o io.Writer) {
	var err error
	// create newServer() blocks
	for _, server := range servers {
		if len(server.name) > 0 {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", name=\"%s\", pool=\"%s\"})\n",
				server.address,
				server.name,
				pool,
			)
			check(err)
		} else {
			_, err := fmt.Fprintf(
				o, "newServer({address=\"%s\", pool=\"%s\"})\n",
				server.address,
				pool,
			)
			check(err)
		}
	}

	// create addAction() blocks
	_, err = fmt.Fprint(o, "addAction({")
	check(err)
	for _, domain := range domains {
		_, err = fmt.Fprintf(
			o,
			"'%s', ",
			domain,
		)
		check(err)
	}
	_, err = fmt.Fprintf(
		o,
		"}, PoolAction(\"%s\"))\n",
		pool,
	)
	check(err)
}
