package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type ConfigFile struct {
	Handles   map[string]string `toml:"handles"`
	Platforms struct {
		Enabled []string `toml:"enabled"`
	} `toml:"platforms"`
}

type Options struct {
	Markdown  bool
	Copy      bool
	Last      int
	Date      string
	Handles   map[string]string
	Platforms []string
}

func LoadConfig() (*ConfigFile, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config dir: %w", err)
	}

	configFilePath := filepath.Join(configDir, "slvd", "config.toml")

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		// Create an empty config file , if it does not exist.
		if errors.Is(err, os.ErrNotExist) {
			return &ConfigFile{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigFile
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid TOML format in config file(%s): %w", configFilePath, err)
	}
	return &config, nil
}

func Parse() *Options {
	opts := &Options{
		Handles:   make(map[string]string),
		Platforms: []string{},
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Print(err)
	} else {
		if cfg.Handles != nil {
			opts.Handles = cfg.Handles
		}
		if cfg.Platforms.Enabled != nil {
			opts.Platforms = cfg.Platforms.Enabled
		}
	}

	// Using StringVar/IntVar to bind terminal inputs directly to the struct's memory addresses
	flag.IntVar(&opts.Last, "last", -1, "Fetches N recent successful submissions")
	flag.StringVar(&opts.Date, "date", "", "Filter by specified date (DD-MM-YYYY)")

	flag.BoolVar(&opts.Markdown, "md", false, "Output table in Markdown format")
	flag.BoolVar(&opts.Copy, "copy", false, "Copy output to clipboard")
	
	var cf, atc, lc bool
	flag.BoolVar(&cf, "cf", false, "Filter by Codeforces submissions")
	flag.BoolVar(&atc, "atc", false, "Filter by AtCoder submissions")
	flag.BoolVar(&lc, "lc", false, "Filter by Leetcode submissions")

	flag.Usage = func() {
		fmt.Printf("USAGE:\n")
		fmt.Printf("  slvd [flags]\n\n")

		fmt.Printf("FILTERING FLAGS:\n")
		fmt.Printf("  -last <int>      Fetch N recent successful submissions (e.g., 10)\n")
		fmt.Printf("  -date <string>   Filter by specified date in DD-MM-YYYY format\n")
		fmt.Printf("  -cf              Show ONLY Codeforces submissions\n")
		fmt.Printf("  -atc             Show ONLY AtCoder submissions\n")
		fmt.Printf("  -lc              Show ONLY LeetCode submissions\n\n")

		fmt.Printf("OUTPUT FLAGS:\n")
		fmt.Printf("  -md              Output table in Markdown format\n")
		fmt.Printf("  -copy            Copy the generated table to system clipboard\n")
		fmt.Printf("  -h, --help       Show this help message\n\n")

		fmt.Printf("EXAMPLES:\n")
		fmt.Printf("  slvd -last 10                # Get last 10 solves from all platforms\n")
		fmt.Printf("  slvd -last 5 -md             # Get last 5 solves in Markdown format\n")
		fmt.Printf("  slvd -date 30-06-2026 -cf    # Get Codeforces solves for a specific date\n")
		fmt.Println()
	}
	flag.Parse()

	if cf || atc || lc {
		opts.Platforms = []string{}

		if cf {
			opts.Platforms = append(opts.Platforms, "codeforces")
		}

		if atc {
			opts.Platforms = append(opts.Platforms, "atcoder")
		}

		if lc {
			opts.Platforms = append(opts.Platforms, "leetcode")
		}
	}

	if len(opts.Platforms) == 0 {
		log.Fatal("Critical: No platforms enabled. Please check config.toml or use flags (-lc, -cf, -atc).")
	}

	if len(opts.Handles) == 0 {
		log.Fatal("Critical: No handles configured. Please add them in config.toml.")
	}

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
