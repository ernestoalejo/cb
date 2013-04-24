package config

import (
	"flag"
)

var (
	Verbose     = flag.Bool("v", false, "verbose mode")
	AlwaysY     = flag.Bool("y", false, "answer yes to all questions")
	AlwaysN     = flag.Bool("n", false, "answer no to all questions")
	Help        = flag.Bool("help", false, "show this help message")
	AngularMode = flag.Bool("angular", false, "use angular tasks")
	ClosureMode = flag.Bool("closure", false, "use closure tasks")
	NoColors    = flag.Bool("no-color", false, "don't use colors in the output")
	ClientOnly  = flag.Bool("client-only", false, "sync only the client items")
)
