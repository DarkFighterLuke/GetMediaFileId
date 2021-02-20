package main

import (
	"encoding/json"
	"fmt"
	"github.com/NicoNex/echotron"
	"log"
	"os"
	"time"
)

const (
	botLogsFolder = "/GetMediaFileIdData/logs/"
)

type bot struct {
	chatId int64
	echotron.Api
}

var TOKEN = os.Getenv("GetMediaFileIdBot")
var logsFolder string

func newBot(chatId int64) echotron.Bot {
	return &bot{
		chatId,
		echotron.NewApi(TOKEN),
	}
}

func initFolders() {
	currentPath, _ := os.Getwd()

	logsFolder = currentPath + botLogsFolder
	_ = os.MkdirAll(logsFolder, 0755)
}

func main() {
	initFolders()

	dsp := echotron.NewDispatcher(TOKEN, newBot)
	dsp.ListenWebhook("https://hiddenfile.tk:443/bot/GetMediaFileId", 40990)
}

func (b *bot) Update(update *echotron.Update) {
	b.logUser(update, logsFolder)
	if update.Message.Text == "/start" {
		b.sendStart(update.Message)
	} else {
		b.sendFileID(update.Message)
	}
}

func (b *bot) sendStart(message *echotron.Message) {
	msg := `Hello, %s!
Send me media and I will tell you their <b>file_id</b> on Telegram servers.`
	msg = fmt.Sprintf(msg, message.User.FirstName)

	b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
}

func (b *bot) sendFileID(message *echotron.Message) {
	if message.Sticker != nil {
		fileID := message.Sticker.FileId
		msg := fmt.Sprintf("The Sticker <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	} else if message.Audio != nil {
		fileID := message.Audio.FileId
		msg := fmt.Sprintf("The Audio <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	} else if message.Document != nil {
		fileID := message.Document.FileId
		msg := fmt.Sprintf("The Document <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	} else if message.Photo != nil {
		photoArray := message.Photo
		for i, v := range photoArray {
			msg := fmt.Sprintf("The Photo number %d <b>file_id</b> is: <code>%s</code>", i, v.FileId)
			b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
		}
	} else if message.Video != nil {
		fileID := message.Video.FileId
		msg := fmt.Sprintf("The Video <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	} else if message.Voice != nil {
		fileID := message.Voice.FileId
		msg := fmt.Sprintf("The Voice <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	} else if message.VideoNote != nil {
		fileID := message.VideoNote.FileId
		msg := fmt.Sprintf("The VideoNote <b>file_id</b> is: <code>%s</code>", fileID)
		b.SendMessage(msg, message.Chat.ID, echotron.PARSE_HTML)
	}
}

func (b *bot) logUser(update *echotron.Update, folder string) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Println("Error marshaling logs: ", err)
		return
	}

	var filename string

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message.Chat.Type == "private" {
			if update.CallbackQuery.Message.Chat.Username == "" {
				filename = folder + update.CallbackQuery.Message.Chat.FirstName + "_" + update.CallbackQuery.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.CallbackQuery.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.Message.Chat.Title + ".txt"
		}

	} else if update.Message != nil {
		if update.Message.Chat.Type == "private" {
			if update.Message.Chat.Username == "" {
				filename = folder + update.Message.Chat.FirstName + "_" + update.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.Message.Chat.Title + ".txt"
		}

	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	dataString := time.Now().Format("2006-01-02T15:04:05") + string(data[:])
	_, err = f.WriteString(dataString + "\n")
	if err != nil {
		log.Println(err)
		return
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
