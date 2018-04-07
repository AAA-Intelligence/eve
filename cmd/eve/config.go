package main

import "flag"

// Config holds information to configure the web server
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
// program usage e.g.: eve -host eve.de -http 80 -https 443
func loadConfig() *Config {
	var config Config
	flag.StringVar(&config.Host, "host", "", "hostname")
	flag.IntVar(&config.HTTP, "http", 80, "HTTP port")
	flag.IntVar(&config.HTTPS, "https", 443, "HTTPS port")
	flag.Parse()
	return &config
}
