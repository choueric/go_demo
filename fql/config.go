package main

import (
	"fmt"
	"os"
)

const (
	CONFIG_FILE = "config.json"
)

const DefaultConfig = `{
	"chrome": "/opt/google/chrome/chrome",
	"url_prefix": "url_prefix"
} `

type Config struct {
	Chrome    string `json:"chrome"`
	UrlPrefix string `json:"url_prefix"`
}

func getConfig() *Config {
	jc := New(CONFIG_FILE, Config{})

	if _, err := jc.Load(DefaultConfig); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v", err)
		os.Exit(1)
	}

	config := jc.Data().(*Config)
	fmt.Println(config)
	return config
}
