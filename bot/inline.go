package bot

import (
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

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
