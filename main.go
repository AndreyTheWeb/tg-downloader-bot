package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"tg-bot-insta-downloader/utils"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

var (
	bot            *telebot.Bot
	instagramRegex = regexp.MustCompile(`^(https:\/\/(\w+\.)?instagram.com\/)`)
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	rapidApiKey := os.Getenv("RAPID_API_KEY")
	if rapidApiKey == "" {
		log.Fatal("RAPID_API_KEY is required")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL is required")
	}

	bot, err = telebot.NewBot(telebot.Settings{
		Token: botToken,
		Poller: &telebot.Webhook{
			Listen:   ":3000",
			Endpoint: &telebot.WebhookEndpoint{PublicURL: webhookURL},
		},
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
			url, err := utils.FetchInstagramVideo(message, rapidApiKey)
			if err != nil {
				return err
			}
			return c.Send(&telebot.Video{File: telebot.FromURL(url)})
		}
		return nil
	})
}

// Handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Bot is running")
	update := &telebot.Update{}
	bot.ProcessUpdate(*update)
}
