package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"tg-bot-insta-downloader/utils"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

var (
	instagramRegex = regexp.MustCompile(`^(https:\/\/(\w+\.)?instagram.com\/)`)
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	rapidApiKey := os.Getenv("RAPID_API_KEY")
	if rapidApiKey == "" {
		log.Fatal("RAPID_API_KEY is required")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Добро пожаловать! Отправьте ссылку на Instagram Reels для загрузки.")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		message := c.Message().Text
		if instagramRegex.MatchString(message) {
			if err != nil {
				return err
			}

			url, err := utils.FetchInstagramVideo(message, rapidApiKey)
			if err != nil {
				return err
			}

			c.Send(&telebot.Video{File: telebot.FromURL(url)})
		}
		return nil
	})

	bot.Start()
}
