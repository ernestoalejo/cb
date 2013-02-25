package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

var ()

func main() {
	flag.Parse()

	if *config.Help {
		usage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	q := &registry.Queue{}
	for _, task := range args {
		q.AddTask(task)
	}

	if err := q.RunWithTimer(config); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("\n Usage: cb [target] [options...]")
	flag.PrintDefaults()
	registry.PrintTasks()
}
