package config

import (
	"flag"
)

var (
	// Verbose mode, output more logs.
	Verbose = flag.Bool("v", false, "verbose mode")

	// AlwaysY (yes) as an answer to questions.
	AlwaysY = flag.Bool("y", false, "answer yes to all questions")

	// AlwaysN (no) as an answer to questions.
	AlwaysN = flag.Bool("n", false, "answer no to all questions")

	// Help outputs the usage message.
	Help = flag.Bool("help", false, "show this help message")

	// NoColors remove the colored output.
	NoColors = flag.Bool("no-color", false, "don't use colors in the output")

	// Port for the server tasks
	Port = flag.Int("port", 9810, "server port")
)
