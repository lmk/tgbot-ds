package main

// go get -u github.com/go-telegram-bot-api/telegram-bot-api

import (
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	botToken := os.Getenv("BOTTOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("[ERROR] Fail create bot [%s]", botToken)
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("[READY] Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		respMsg := parsingBotMessage(bot, update.Message)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, respMsg)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)

	}
}

// parsingBotMessage parsing bot message
// download torrent file to destnation
// create magent uri to destnation
func parsingBotMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) string {

	myTgName := os.Getenv("MYTGNAME")
	destination := os.Getenv("DESTINATION")

	if myTgName != message.From.UserName {
		log.Printf("[WARNING] Is not you [%s]", message.From.UserName)
		return "I am not yours"
	}

	if message.Document != nil { // torrent file

		if message.Document.MimeType != "application/x-bittorrent" {
			log.Printf("[WARNING] Is not torrent file [%s]", message.Document.MimeType)
			return "Not torrent file"
		}

		url, err := bot.GetFileDirectURL(message.Document.FileID)
		if err != nil {
			log.Printf("[ERROR] Does not get File URL - %s", err)
			return "Does not get File URL"
		}

		err = downloadFile(destination, message.Document.FileName, url)
		if err != nil {
			log.Printf("[ERROR] Download fail - %s", err)
			return "Fail download"
		}

	} else { // magnet

		//log.Printf("[%s] %s", message.From.UserName, message.Text)
		err := CreateMagnet(destination, message.Text)
		if err != nil {
			log.Printf("[ERROR] Create magnet fail - %s", err)
			return "Fail create magent"
		}
	}

	return OkMesssge()
}

// downloadFile  download torrent file to destnation
func downloadFile(dst string, name string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(dst + string(os.PathSeparator) + name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
