package main

import "io"

func ActiveDirectory(c map[string]interface{}, o io.Writer){
	io.WriteString(o, "test")
}