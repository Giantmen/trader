package main

import (
	"flag"
	stdlog "log"

	"github.com/Giantmen/trader/api"
	"github.com/Giantmen/trader/config"
	"github.com/Giantmen/trader/log"

	"github.com/BurntSushi/toml"
	"github.com/solomoner/gozilla"
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

	api.Register(&cfg)
	log.Info("register over")
	gozilla.DefaultLogOpt.Format += " {{.Body}}"
	stdlog.Fatal(gozilla.ListenAndServe(cfg.Listen))
}
