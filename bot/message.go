package bot

import (
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func message(Message *tgbotapi.Message) {
	// CLB - Made for custom craft
	if Message.ForwardFrom != nil {
		if Message.ForwardFrom.ID == 280391978 {
			if Message.Document != nil {
				fileName := "C-" + strconv.Itoa(Message.Date)
				fileID := Message.Document.FileID
				fileURL, _ := bot.GetFileDirectURL(fileID)

				if fileURL != "" {
					donwloadStatus := DownloadFileFromURL(fileName, fileURL)
					if donwloadStatus == true {
						log.Println("File Creato correttamente")

						msg := tgbotapi.NewMessage(Message.Chat.ID, "Lista caricata correttamente")
						msg.ReplyToMessageID = Message.MessageID

						replyInputMessage := "Custom Craft 🔨 " + fileName + ":1"
						msg.ReplyMarkup = SetterSwitchCLBInlineKeyboard("Start", replyInputMessage)

						// Send message start craft
						if _, err := bot.Send(msg); err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}
}
