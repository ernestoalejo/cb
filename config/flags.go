package config

import (
	"flag"
)

var (
	Verbose  = flag.Bool("v", false, "verbose mode")
	AlwaysY  = flag.Bool("y", false, "answer yes to all questions")
	AlwaysN  = flag.Bool("n", false, "answer no to all questions")
	Help     = flag.Bool("help", false, "show this help message")
	NoColors = flag.Bool("no-color", false, "don't use colors in the output")
)
