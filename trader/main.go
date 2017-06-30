package main

import (
	"flag"
	stdlog "log"
	"net/http"

	"github.com/Giantmen/trader/config"
	"github.com/Giantmen/trader/log"

	"github.com/BurntSushi/toml"
)

var (
	cfgPath = flag.String("config", "config.toml", "config file path")
)

func initLog(cfg *config.Config) {
	log.SetLevelByString(cfg.LogLevel)
	if !cfg.Debug {
		log.SetHighlighting(false)
		err := log.SetOutputByName(cfg.LogPath)
		if err != nil {
			log.Fatal(err)
		}
		log.SetRotateByDay()
	}
}

func main() {
	flag.Parse()
	var cfg config.Config
	_, err := toml.DecodeFile(*cfgPath, &cfg)
	if err != nil {
		stdlog.Fatal("DecodeConfigFile error: ", err)
	}
	initLog(&cfg)
	stdlog.Fatal(http.ListenAndServe(cfg.Listen, nil))
}
