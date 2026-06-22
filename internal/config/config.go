package config

import (
	"flag"
	"log"
	"os"
	"time"
)

type Options struct {
	Handle string
	Last   int
	Date   string
}

func Parse() *Options {
	opts := &Options{}

	// Using StringVar/IntVar to bind terminal inputs directly to the struct's memory addresses
	flag.IntVar(&opts.Last, "last", -1, "Fetches N recent successful submissions")
	flag.StringVar(&opts.Date, "date", "", "Filter by specified date (DD-MM-YYYY)")

	flag.Usage = func() {
		log.Printf("Usage: %s [flags] <platform-handle>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	opts.Handle = args[0]
	return opts

}

func (opts *Options) GetAtCoderSecond() int64 {
	now := time.Now()

	// Date Flag
	if opts.Date != "" {
		t, err := time.ParseInLocation("02-01-2006", opts.Date, now.Location())
		if err != nil {
			log.Fatal("Critical: Invalid date format. Please use DD-MM-YYYY")
		}
		return t.Unix()
	}

	// Last Flag
	if opts.Last != -1 {
		return now.AddDate(-1, 0, 0).Unix()
	}

	// Default Today
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return today.Unix()
}
