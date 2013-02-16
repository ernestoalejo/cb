package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	help = flag.Bool("help", false, "show this help message")
)

func main() {
	flag.Parse()

	if *help {
		usage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	queue = args
	if err := runQueue(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("Usage: cb [target] [options...]")
	flag.PrintDefaults()
}
