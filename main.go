package main

import (
	"log"
	"os"

	"github.com/tattsun/mood-linebot/src"
)

func envvar(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable '%s' is required", key)
	}
	return value
}

func readConfig() *src.Config {
	channelSecret := envvar("LINE_CHANNEL_SECRET")
	channelToken := envvar("LINE_CHANNEL_TOKEN")
	port := envvar("PORT")
	mongoAddr := envvar("MONGODB_ADDR")
	mongoDatabase := envvar("MONGODB_DATABASE")

	config := src.Config{
		Port: port,
		LINE: src.LINEConfig{
			ChannelSecret: channelSecret,
			ChannelToken:  channelToken,
		},
		DB: src.DBConfig{
			MongoAddr:     mongoAddr,
			MongoDatabase: mongoDatabase,
		},
	}
	return &config
}

func main() {
	config := readConfig()
	bot, err := src.NewBot(config)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "check" {
		if err := bot.SendFeelingCheck(); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatal(bot.RunServer())
}
