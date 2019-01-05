package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var configPath = flag.String("config", "config.toml", "config file")
	var outputPath = flag.String("output", "-", "output file")
	flag.Parse()

	var outputFile = os.Stdout
	if *outputPath != "-" {
		outputFile, err := os.Create(*outputPath)
		check(err)
		defer outputFile.Close()
	}

	conf := &config{}
	_, err := toml.DecodeFile(*configPath, conf)
	check(err)

	// default
	if len(conf.Listen) == 0 {
		conf.Listen = []string{"127.0.0.1:53", "[::1]:53"}
	}

	_, err = fmt.Fprintf(outputFile, "%s Auto generated by dnsconf \n\n", OutputCommentPrefix)
	check(err)

	// TODO: generate global config

	// TODO: ECS https://dnsdist.org/advanced/ecs.html

	log.Printf("Match list count: %d", len(conf.Matches))
	for index, m := range conf.Matches {
		// normalize options and process default options
		if m["provider"] == nil {
			m["provider"] = "DomainList"
		}
		if m["action"] == nil {
			m["action"] = "resolve"
		} else {
			m["action"] = strings.ToLower(m["action"].(string))
		}

		var o bytes.Buffer

		// find the appropriate provider
		providerName := strings.ToLower(m["provider"].(string))
		found := false
		for key, value := range providers {
			if strings.ToLower(key) == providerName {
				// got a match
				log.Printf("Processing match #%d, type %s, action %s\n", index+1, key, m["action"])
				found = true
				value(m, &o)
				_, err = fmt.Fprintf(outputFile, "\n%s match #%d [%s] -> %s\n", OutputCommentPrefix, index+1, m["provider"], m["action"])
				check(err)
				_, err = outputFile.WriteString(o.String())
				check(err)
				_, err = fmt.Fprintf(outputFile, "\n%s end match #%d\n\n", OutputCommentPrefix, index+1)
				check(err)
				break
			}
		}

		if !found {
			log.Fatalf("Unknown provider %s at match #%d\n", m["provider"], index+1)
		}

	}
}
