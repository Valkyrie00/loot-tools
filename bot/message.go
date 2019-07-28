package bot

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func message(Message *tgbotapi.Message) {
	if Message.ForwardFrom != nil {
		// Forward from Craft Loot Boot
		if Message.ForwardFrom.ID == 280391978 {
			// File with craft list
			clbCraftFile(Message)

			// Message with needed items
			if strings.Contains(Message.Text, "Lista oggetti necessari per") || strings.Contains(Message.Text, "Per eseguire i craft spenderai") {
				ClbNeededList(Message)
			}
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
		if strings.Contains(line, "🚫") {
			if itemIndex >= 10 {
				stringResult = strings.TrimSuffix(stringResult, ",")
				results = append(results, stringResult)
				itemIndex = 0
				stringResult = "/negozio "
			}

			itemName := GetStringInBetween(line, "> ", " (")
			todoItems := strings.Split(GetStringInBetween(line, ") ", " 🚫"), "/")
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

// CLB - Craft
func clbCraftFile(Message *tgbotapi.Message) {
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

// *********************
// Needed items
// *********************

var clbNeededMap = make(map[int64][]string)

func ClbNeededList(message *tgbotapi.Message) {
	clbNeededMap[int64(message.From.ID)] = append(clbNeededMap[int64(message.From.ID)], message.Text)

	// Send inline command
	emptyInputMessage := "lcbNeededListEmpty"
	calculateInputMessage := "lcbNeededListCalculate"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "🗑 Svuota",
				CallbackData: &emptyInputMessage,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "🤯 Calcola",
				CallbackData: &calculateInputMessage,
			},
		),
	)

	replyMsg := "Lista oggetti necessari caricata correttamente, puoi continuare a mandare altri liste e ultimare scegliendo calcola. 👻"
	replyMsg += fmt.Sprintf("\n- %v x liste caricate.", len(clbNeededMap[int64(message.From.ID)]))

	msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func lcbNeededListEmpty(message *tgbotapi.Message) {
	if _, ok := clbNeededMap[message.Chat.ID]; ok {
		clbNeededMap[message.Chat.ID] = nil
	}

	replyMsg := "✅ Lista oggetti necessari svuotata correttamente."
	msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func lcbNeededListCalculate(message *tgbotapi.Message) {
	var moneyNeeded int
	var itemsNeeded = make(map[string]int)
	var higherQuantity int

	if _, ok := clbNeededMap[message.Chat.ID]; ok {
		for _, message := range clbNeededMap[message.Chat.ID] {
			lines := strings.Split(message, "\n")
			for _, line := range lines {

				// Item
				if strings.Contains(line, "> ") {
					itemQuantity, _ := strconv.Atoi(GetStringInBetween(line, "> ", " di"))
					itemName := GetStringInBetween(line, "di", ")") + ")"

					// It's for generate inline kayboard
					if higherQuantity < itemQuantity {
						higherQuantity = itemQuantity
					}

					// Add or summ
					if _, ok := itemsNeeded[itemName]; !ok {
						itemsNeeded[itemName] = itemQuantity
						continue
					}

					itemsNeeded[itemName] = itemsNeeded[itemName] + itemQuantity
				}

				// Money
				if strings.Contains(line, "spenderai") {
					strNeededQauntity := GetStringInBetween(line, "spenderai: ", "§")
					strNeededQauntity = strings.ReplaceAll(strNeededQauntity, "'", "")
					strNeededQauntity = strings.ReplaceAll(strNeededQauntity, ".", "")

					intNeededQuantity, _ := strconv.Atoi(strNeededQauntity)
					moneyNeeded = moneyNeeded + intNeededQuantity
				}
			}
		}
	}

	replyMsg := "Hai consumato:\n"
	for items, quantity := range itemsNeeded {
		replyMsg += fmt.Sprintf("> %v di %v\n", quantity, items)
	}
	replyMsg += fmt.Sprintf("\nPer un totale di: %v§", moneyNeeded)
	replyMsg += "Qui di seguito hai la possibilità di creare delle stringhe per generare negozi.\n Indica per quante persone vuoi dividere la spesa. 🛒"

	// Add 20 inline buttons
	var tempButtons []tgbotapi.InlineKeyboardButton
	var inlineKeyboardRows [][]tgbotapi.InlineKeyboardButton

	// Filter for list under 20 quantity
	minNumberOfButtons := 20
	if higherQuantity < minNumberOfButtons {
		minNumberOfButtons = higherQuantity
	}

	for index := 1; index <= minNumberOfButtons; index++ {
		shopInputMessage := fmt.Sprintf("lcbNeededListShop-%v", index)
		button := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf("👤 %v", index),
			CallbackData: &shopInputMessage,
		}

		tempButtons = append(tempButtons, button)
		if index%5 == 0 || index == 1 {
			row := tgbotapi.NewInlineKeyboardRow(tempButtons...)
			inlineKeyboardRows = append(inlineKeyboardRows, row)
			tempButtons = nil
		}
	}

	// Append delete button
	deleteInputMessage := "deleteMessage"
	deleteButton := tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{
		Text:         "🗑 Cancella",
		CallbackData: &deleteInputMessage,
	})
	inlineKeyboardRows = append(inlineKeyboardRows, deleteButton)

	msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboardRows,
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func lcbNeededListShop(divider string, message *tgbotapi.Message) {
	lines := strings.Split(message.Text, "\n")
	divisor, _ := strconv.Atoi(divider)

	var itemIndex int
	var results []string
	stringResult := "/negozio "
	canCreateShop := false

	itemIndex = 0
	for i, line := range lines {
		if strings.Contains(line, "> ") {
			if itemIndex >= 10 {
				stringResult = strings.TrimSuffix(stringResult, ",")
				results = append(results, stringResult)
				itemIndex = 0
				stringResult = "/negozio "
			}

			itemName := GetStringInBetween(line, "di ", " (")
			itemQuantity, _ := strconv.Atoi(GetStringInBetween(line, "> ", " di"))

			if itemQuantity < divisor {
				itemQuantity = 1
				if itemQuantity < 5 {
					continue
				}
			} else {
				fItemQuantity := float64(itemQuantity) / float64(divisor)
				itemQuantity = int(math.Ceil(fItemQuantity))
			}

			if itemQuantity > 0 {
				canCreateShop = true
			}

			partialResult := fmt.Sprintf("%s::%v,", itemName, itemQuantity)
			stringResult = stringResult + partialResult
			itemIndex++
		}

		if len(lines)-1 == i {
			stringResult = strings.TrimSuffix(stringResult, ",")
			results = append(results, stringResult)
		}
	}

	if canCreateShop {
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
			msg := tgbotapi.NewMessage(message.Chat.ID, result)
			msg.ReplyMarkup = keyboard
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Non è stato possibile creare il negozio.")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}

}
