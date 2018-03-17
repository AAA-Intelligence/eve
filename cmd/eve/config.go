package main

import "flag"

// Config configures web server
// used in cmd package
type Config struct {

	// Host e.g. google.de, mypage.com, localhost
	Host string

	// HTTP port
	HTTP int

	// HTTPS port
	HTTPS int
}

// loads config data from program arguments
// defaults are:
// 		host: "" (empty)
//		http: 80
//		https: 443
func loadConfig() *Config {
	var config Config
	flag.StringVar(&config.Host, "host", "", "hostname")
	flag.IntVar(&config.HTTP, "http", 80, "HTTP port")
	flag.IntVar(&config.HTTPS, "https", 443, "HTTPS port")
	flag.Parse()
	return &config
}
