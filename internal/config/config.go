package config

import (
	"flag"
	"log"
)

type Options struct {
	Handle string
	Last   int
	Date   string
}

func Parse() *Options {
	lastFlag := flag.Int("last", -1, "Fetches N recent successful submissions")
	dateFlag := flag.String("date", "", "Filter by specified date (DD-MM-YYYY)")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Usage: slvd [--last N] [--date DD-MM-YYYY] <codeforces-handle>")
	}

	return &Options{Handle: args[0], Last: *lastFlag, Date: *dateFlag}
}
