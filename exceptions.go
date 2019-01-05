package main

import "log"

// legacy funcion to check every fucking err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// if QuitOnError is true, then panic;
// else go on
func softFail(e error) error {
	if e != nil {
		if conf.QuitOnError {
			panic(e)
		} else {
			softErrorCount++
			log.Printf("[ERROR] %s", e)
		}
	}
	return e
}
