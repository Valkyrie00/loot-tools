package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Valkyrie00/loot-tools/loot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot               *tgbotapi.BotAPI
	craftableItems    loot.Items
	craftingItemsList loot.ItemsCraftingMapType

	botMode   string
	adminID   int
	botApiKey string
)

func init() {
	// errEnv := godotenv.Load()
	// if errEnv != nil {
	// 	log.Println("Error loading .env file")
	// }

	adminID, _ = strconv.Atoi(os.Getenv("ID_ADMIN"))
	botMode = os.Getenv("BOT_MODE") // Private or public
	botApiKey = os.Getenv("TELEGRAM_APIKEY")

	var botErr error
	bot, botErr = tgbotapi.NewBotAPI(botApiKey)
	bot.Debug = true

	if botErr != nil {
		log.Panic(botErr)
	}
	log.Println(fmt.Sprintf("Bot connected: %s", bot.Self.UserName))

	// Load craftable items
	craftableItems = loot.GetCraftableItems()

	// Load crafting map
	craftingItemsList = loot.GetCraftingMap(craftableItems)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, chanErr := bot.GetUpdatesChan(u)
	if chanErr != nil {
		log.Panicln(chanErr)
	}

	for update := range updates {
		// Message
		if update.Message != nil {
			if botMode == "private" {
				if update.Message.From.ID != adminID {
					continue
				}
			}

			message(update.Message)
		}

		// Inline query
		if update.InlineQuery != nil && update.InlineQuery.Query != "" {
			if botMode == "private" {
				if update.InlineQuery.From.ID != adminID {
					continue
				}
			}

			inline(update.InlineQuery)
		}
	}
}

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

// Logic Structure (Inline Process)
// -- Telegram
// ---- Search
// ---- Craft
// ------ Base Craft
// ------ Custom Crat (CLB)

func inline(InlineQuery *tgbotapi.InlineQuery) {
	var inlineResults []interface{}

	switch ParseInlineType(InlineQuery.Query) {
	case 1:
		inlineResults = inlineSearch(InlineQuery)
	case 2:
		inlineResults = inlineBaseCraft(InlineQuery)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       inlineResults,
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
	}
}

func inlineSearch(InlineQuery *tgbotapi.InlineQuery) []interface{} {
	var resultsForInlineQuery []interface{}

	// Search item by Query
	itemLists := SearchItem(InlineQuery.Query)
	if len(itemLists) > 0 {
		for _, item := range itemLists {

			// Give first crafting item
			firstCraftingItem := craftingItemsList[item.ID][0]
			itemInterface := tgbotapi.NewInlineQueryResultArticle(string(item.ID), item.Name, "Crea "+firstCraftingItem)
			itemInterface.Description = "Rarity: " + item.Rarity + " - Crafts: " + strconv.Itoa(len(craftingItemsList[item.ID])) + " 🔨 "

			if len(craftingItemsList[item.ID]) > 1 {
				replyText := "Next"
				replyInputMessage := item.Name + " 🔨 B-" + strconv.Itoa(item.ID) + ":1"
				itemInterface.ReplyMarkup = SetterCraftInlineKeyboard(replyText, replyInputMessage)
			}

			resultsForInlineQuery = append(resultsForInlineQuery, itemInterface)
		}
	}

	return resultsForInlineQuery
}

func inlineBaseCraft(InlineQuery *tgbotapi.InlineQuery) []interface{} {
	var resultsForInlineQuery []interface{}

	craftingType, itemID, craftingIndex := ParseInlineCraftType(InlineQuery.Query)

	switch craftingType {
	case "B":
		//craftingIndex is already incremented
		baseItem := craftableItems.FindItemByID(itemID)

		if craftingIndex <= len(craftingItemsList[itemID]) {
			nextItem := craftingItemsList[itemID][craftingIndex]

			itemInterface := tgbotapi.NewInlineQueryResultArticle(string(baseItem.ID), nextItem, "Crea "+nextItem)
			itemInterface.Description = "Need for " + baseItem.Name + " ( " + strconv.Itoa(craftingIndex) + " / " + strconv.Itoa(len(craftingItemsList[itemID])-1) + " ) "

			// Controllo se se non è l'ultimo
			if len(craftingItemsList[itemID]) > craftingIndex+1 {
				replyText := "Next " + " ( " + strconv.Itoa(craftingIndex+1) + " / " + strconv.Itoa(len(craftingItemsList[itemID])-1) + " ) "
				replyInputMessage := baseItem.Name + " 🔨 B-" + strconv.Itoa(baseItem.ID) + ":" + strconv.Itoa(craftingIndex+1)
				itemInterface.ReplyMarkup = SetterCraftInlineKeyboard(replyText, replyInputMessage)
			}

			resultsForInlineQuery = append(resultsForInlineQuery, itemInterface)
		}

	case "C":
		fileName := "C-" + strconv.Itoa(itemID)
		craftingItemsList := GetLinesFromFile("storage/clb/" + fileName + ".txt")
		if craftingIndex <= len(craftingItemsList) {
			nextItem := craftingItemsList[craftingIndex]

			itemInterface := tgbotapi.NewInlineQueryResultArticle(fileName, nextItem, nextItem)
			itemInterface.Description = "Custom craft (" + strconv.Itoa(craftingIndex) + " / " + strconv.Itoa(len(craftingItemsList)-1) + " ) "

			if len(craftingItemsList) > craftingIndex+1 {
				replyText := "Next " + " ( " + strconv.Itoa(craftingIndex+1) + " / " + strconv.Itoa(len(craftingItemsList)-1) + " ) "
				replyInputMessage := "Custom Craft 🔨 " + fileName + ":" + strconv.Itoa(craftingIndex+1)
				itemInterface.ReplyMarkup = SetterCraftInlineKeyboard(replyText, replyInputMessage)
			}

			resultsForInlineQuery = append(resultsForInlineQuery, itemInterface)
		}
	}

	return resultsForInlineQuery
}
