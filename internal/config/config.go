package config

import (
	"flag"
	"log"
)

type Options struct {
	Handle string
	Last   int
}

func Parse() *Options {
	lastFlag := flag.Int("last", 0, "Fetches N recent successful submissions")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Usage: slvd [--last N] <codeforces-handle>")
	}

	return &Options{Handle: args[0], Last: *lastFlag}
}
