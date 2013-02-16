package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ernestokarim/cb/registry"
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

	q := &registry.Queue{}
	for _, task := range args {
		q.AddTask(task)
	}

	if err := q.Run(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("Usage: cb [target] [options...]")
	flag.PrintDefaults()
}
