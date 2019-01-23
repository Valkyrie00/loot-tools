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
				lootPlatformForwardHandler(Message)
			}
		}
	}
}

// Loot platform message handler
func lootPlatformForwardHandler(Message *tgbotapi.Message) {
	craftInputMessage := "lootPlatformCraftParser"
	negozioInputMessage := "lootPlatformShopParser"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "🔨 Craft",
				CallbackData: &craftInputMessage,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "💰 Negozio",
				CallbackData: &negozioInputMessage,
			},
		),
	)
	msg := tgbotapi.NewMessage(Message.Chat.ID, "Come vuoi convertire la lista? 👻")
	msg.ReplyToMessageID = Message.MessageID
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

// Loot - Platflorm shop parser
func lootPlatformShopParser(Message *tgbotapi.Message) {
	lines := strings.Split(Message.Text, "\n")

	var itemIndex int
	var results []string
	stringResult := "/negozio "

	itemIndex = 0
	for i, line := range lines {
		if strings.Contains(line, "✅") {
			if itemIndex >= 10 {
				stringResult = strings.TrimSuffix(stringResult, ",")
				results = append(results, stringResult)
				itemIndex = 0
				stringResult = "/negozio "
			}

			itemName := GetStringInBetween(line, "> ", " (")
			todoItems := strings.Split(GetStringInBetween(line, ") ", " ✅"), "/")
			partialResult := itemName + "::" + todoItems[1] + ","

			stringResult = stringResult + partialResult
			itemIndex++
		}

		if len(lines)-1 == i {
			log.Println(len(lines), i)
			stringResult = strings.TrimSuffix(stringResult, ",")
			results = append(results, stringResult)
		}
	}

	if len(results) >= 1 {
		deleteInputMessage := "deleteMessage"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:         "🗑 Cancella",
					CallbackData: &deleteInputMessage,
				},
			),
		)

		for _, result := range results {
			msg := tgbotapi.NewMessage(Message.Chat.ID, result)
			msg.ReplyMarkup = keyboard
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
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

	deleteInputMessage := "deleteMessage"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "🗑 Cancella",
				CallbackData: &deleteInputMessage,
			},
		),
	)

	msg := tgbotapi.NewMessage(Message.Chat.ID, stringResult)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
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
