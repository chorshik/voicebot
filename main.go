package main

import (
	"github.com/ebladrocher/voicebot/bot"
	conf "github.com/ebladrocher/voicebot/system"
	model "github.com/ebladrocher/voicebot/system/model"
	"log"
	//"net/http"
)

const (
	token string = ""
)

func main() {
	cfg, err := conf.ReadConfig()
	if err != nil {
		log.Println(err)
		tmpCfg := &model.Config{
			Debug: false,
		}
		bot, err := bot.NewBot(tmpCfg)
		if err != nil {
			log.Print(err)
		}

		bot.Start()

	} else {
		bot, err := bot.NewBot(cfg)
		if err != nil {
			log.Print(err)
		}

		bot.Start()
	}

}
