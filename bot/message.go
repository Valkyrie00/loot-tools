package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func message(Message *tgbotapi.Message) {
	if Message.ForwardFrom != nil {
		// CLB - Made for custom craft
		if Message.ForwardFrom.ID == 280391978 {
			clbCraftList(Message)
		}

		// Loot Bot - Forward
		if Message.ForwardFrom.ID == 171514820 {
			if strings.Contains(Message.Text, "migliorare la postazione") {
				lootPlatformCraftParser(Message)
			}
		}
	}
}

// Loot - Platflorm craft parser
func lootPlatformCraftParser(Message *tgbotapi.Message) {
	lines := strings.Split(Message.Text, "\n")
	var stringResult string
	stringResult = "/craft "

	for _, line := range lines {
		if strings.Contains(line, "🚫") {

			itemName := GetStringInBetween(line, "> ", " (")
			todoItems := strings.Split(GetStringInBetween(line, ") ", " 🚫"), "/")
			qToCraft, _ := strconv.Atoi(todoItems[0])
			qHaveItem, _ := strconv.Atoi(todoItems[1])

			partialResult := itemName + ":" + strconv.Itoa(qHaveItem-qToCraft) + ","
			stringResult = stringResult + partialResult

		}
	}

	stringResult = strings.TrimSuffix(stringResult, ",")
	var messages []tgbotapi.MessageConfig

	{
		msg := tgbotapi.NewMessage(Message.Chat.ID, "Craft postazioni generato. Ora puoi inoltrarlo a *CLB*.")
		msg.ReplyToMessageID = Message.MessageID
		msg.ParseMode = "Markdown"

		messages = append(messages, msg)
	}
	{
		msg := tgbotapi.NewMessage(Message.Chat.ID, stringResult)
		msg.ParseMode = "Markdown"

		messages = append(messages, msg)
	}

	for _, message := range messages {
		if _, err := bot.Send(message); err != nil {
			log.Println(err)
		}
	}
}

// CLB - Made for custom craft
func clbCraftList(Message *tgbotapi.Message) {
	if Message.Document != nil {
		fileName := "C-" + strconv.Itoa(Message.Date)
		fileID := Message.Document.FileID
		fileURL, _ := bot.GetFileDirectURL(fileID)

		if fileURL != "" {
			donwloadStatus := DownloadFileFromURL(fileName, fileURL)
			if donwloadStatus == true {
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
