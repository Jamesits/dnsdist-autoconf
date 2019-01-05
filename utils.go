package main

import (
	"fmt"
	"io"
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

// check every fucking err
func check(e error) {
	if e != nil {
		panic(e)
	}
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

func generateServerPool(pool string, servers []DnsServer, o io.Writer) {
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

}

// function to generate a generic domain list / dns server config block
func generateActions(pool string, domains []string, action string, o io.Writer) {
	var err error
	// create addAction() blocks
	_, err = fmt.Fprint(o, "addAction({\n")
	check(err)
	for index, domain := range domains {
		if index > 0 {
			_, err = fmt.Fprint(o, ", \n")
			check(err)
		}
		_, err = fmt.Fprintf(o, "    '%s'", domain)
		check(err)
	}
	_, err = fmt.Fprint(o, "\n}, ")
	check(err)

	switch action {
	case "resolve":
		_, err = fmt.Fprintf(o, "PoolAction(\"%s\")", pool)
	case "servfail":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.SERVFAIL)")
	case "block":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.NXDOMAIN)")
	case "nxdomain":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.NXDOMAIN)")
	case "refuse":
		_, err = fmt.Fprint(o, "RCodeAction(dnsdist.REFUSED)")
	case "drop":
		_, err = fmt.Fprint(o, "DropAction()")
	}
	check(err)

	_, err = fmt.Fprint(o, ")\n")
	check(err)
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
