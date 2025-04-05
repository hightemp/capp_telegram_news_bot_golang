package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/sym01/htmlsanitizer"
)

func PrepareString(input string) string {
	re := regexp.MustCompile(`\s+`)
	output := re.ReplaceAllString(input, " ")
	output = strings.TrimSpace(output)
	return output
}

func main() {
	cd, _ := os.Getwd()
	fmt.Println("Current working directory:", cd)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env")
	}

	RSS_LINK := os.Getenv("RSS_LINK")
	// TELEGRAM_BOT_URL := os.Getenv("TELEGRAM_BOT_URL")
	// TELEGRAM_BOT_NAME := os.Getenv("TELEGRAM_BOT_NAME")
	TELEGRAM_BOT_KEY := os.Getenv("TELEGRAM_BOT_KEY")
	TELEGRAM_BOT_CHANNEL_ID := os.Getenv("TELEGRAM_BOT_CHANNEL_ID")
	TELEGRAM_BOT_UPDATE_TIMEOUT := os.Getenv("TELEGRAM_BOT_UPDATE_TIMEOUT")

	channel_id, _ := strconv.ParseInt(TELEGRAM_BOT_CHANNEL_ID, 10, 64)
	sleep_hours, _ := strconv.Atoi(TELEGRAM_BOT_UPDATE_TIMEOUT)

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_KEY)
	if err != nil {
		log.Panic(err)
	}

	LastPublishedDate, err := os.ReadFile("lastdate.txt")
	if err != nil {
		log.Printf("[WARNING] Can't read file: lastdate.txt\n")
	}
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	LastPublishedDateTime, err := time.Parse(layout, string(LastPublishedDate))
	if err != nil {
		LastPublishedDateTime = time.Now()
	}

	DefaultAllowList := &htmlsanitizer.AllowList{
		Tags: []*htmlsanitizer.Tag{
			{"a", []string{"href"}, []string{}},
			{"b", []string{}, []string{}},
			{"i", []string{}, []string{}},
			{"u", []string{}, []string{}},
			{"code", []string{}, []string{}},
			{"pre", []string{}, []string{}},
		},
		GlobalAttr: []string{
			"class",
			"id",
		},
	}

	sanhtml := htmlsanitizer.NewHTMLSanitizer()
	sanhtml.AllowList = DefaultAllowList

	for {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(RSS_LINK)
		log.Println(feed.Title)
		log.Println(LastPublishedDateTime)

		for i, j := 0, len(feed.Items)-1; i < j; i, j = i+1, j-1 {
			feed.Items[i], feed.Items[j] = feed.Items[j], feed.Items[i]
		}

		for _, item := range feed.Items {
			if (*item.PublishedParsed).Before(LastPublishedDateTime) || (*item.PublishedParsed).Equal(LastPublishedDateTime) {
				continue
			}

			log.Println("Заголовок статьи:", item.Title)
			log.Println("Ссылка на статью:", item.Link)
			log.Println("Дата:", item.PublishedParsed)
			log.Println("Дата:", item.Published)
			// log.Println("Description:", item.Description)
			// log.Println("Content:", item.Content)
			log.Println("=============================")

			// log.Println(channel_id)
			// log.Println(bot)
			description, _ := sanhtml.SanitizeString(item.Description)
			description = PrepareString(description)
			messageText := html.UnescapeString("<b>" + item.Title + "</b>" + "\n\n" + description + "\n\n" + item.Link)

			log.Println("messageText:", messageText)

			message := tgbotapi.NewMessage(channel_id, messageText)
			message.DisableWebPagePreview = true
			message.ParseMode = "HTML"

			_, err = bot.Send(message)
			if err != nil {
				log.Println("Error sending bot message:", err)
			}

			LastPublishedDateTime = *item.PublishedParsed
			err = os.WriteFile("lastdate.txt", []byte(item.Published), 0644)

			if err != nil {
				log.Println("Error writing file 'lastdate.txt':", err)
			}
		}

		time.Sleep(time.Duration(sleep_hours) * time.Hour)
	}
}
