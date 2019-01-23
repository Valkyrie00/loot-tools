package bot

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func callback(CallbackQuery *tgbotapi.CallbackQuery) {
	data := CallbackQuery.Data

	switch data {
	case "lootPlatformCraftParser":
		// From loot plataform generate craft string for clb
		lootPlatformCraftParser(CallbackQuery.Message.ReplyToMessage)

	case "lootPlatformShopParser":
		// From loot plataform generate shop string for clb
		lootPlatformShopParser(CallbackQuery.Message.ReplyToMessage)

	case "deleteMessage":
		deleteConfig := tgbotapi.DeleteMessageConfig{
			ChatID:    CallbackQuery.Message.Chat.ID,
			MessageID: CallbackQuery.Message.MessageID,
		}

		if _, err := bot.Send(deleteConfig); err != nil {
			log.Println(err)
		}
	}
}
