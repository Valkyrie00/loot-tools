package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func callback(CallbackQuery *tgbotapi.CallbackQuery) {
	data := CallbackQuery.Data

	// Hack for lcbNeededListShop
	if strings.Contains(data, "lcbNeededListShop") {
		data = "lcbNeededListShop"
	}

	switch data {
	case "lootPlatformCraftParser":
		// From loot plataform generate craft string for clb
		lootPlatformCraftParser(CallbackQuery.Message.ReplyToMessage)

	case "lootPlatformShopParser":
		// From loot plataform generate shop string for clb
		lootPlatformShopParser(CallbackQuery.Message.ReplyToMessage)

	case "lcbNeededListEmpty":
		// Clear needed list map (CLB)
		lcbNeededListEmpty(CallbackQuery.Message)

	case "lcbNeededListCalculate":
		lcbNeededListCalculate(CallbackQuery.Message)

	case "lcbNeededListShop":
		callbackSplitted := strings.Split(CallbackQuery.Data, "-")
		lcbNeededListShop(callbackSplitted[1], CallbackQuery.Message)

	case "deleteMessage":
		deleteConfig := tgbotapi.DeleteMessageConfig{
			ChatID:    CallbackQuery.Message.Chat.ID,
			MessageID: CallbackQuery.Message.MessageID,
		}

		if _, err := bot.Send(deleteConfig); err != nil {
			log.Println(err)
		}
	}

	returnCallBack := tgbotapi.NewCallback(CallbackQuery.ID, "")
	if _, err := bot.AnswerCallbackQuery(returnCallBack); err != nil {
		log.Println(err.Error())
	}
}
