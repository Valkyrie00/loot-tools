package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Valkyrie00/loot-tools/loot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/joho/godotenv/autoload" // .env
)

var (
	bot               *tgbotapi.BotAPI
	craftableItems    loot.Items
	craftingItemsList loot.CraftingMapType

	botMode   string
	adminID   int
	botAPIKey string
)

func init() {
	adminID, _ = strconv.Atoi(os.Getenv("ID_ADMIN"))
	botMode = os.Getenv("BOT_MODE") // Private or public
	botAPIKey = os.Getenv("TELEGRAM_APIKEY")

	var botErr error
	bot, botErr = tgbotapi.NewBotAPI(botAPIKey)
	bot.Debug = true

	if botErr != nil {
		log.Panic(botErr)
	}
	log.Println(fmt.Sprintf("Bot connected: %s", bot.Self.UserName))

	// Load craftable items
	// craftableItems = loot.GetCraftableItems()

	// Load crafting map
	craftingItemsList = loot.GetCraftingMap(craftableItems)
}

//Handler - Updates Handler
func Handler() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, chanErr := bot.GetUpdatesChan(u)
	if chanErr != nil {
		log.Panicln(chanErr)
	}

	for update := range updates {
		// Check bot access mode
		if botMode == "private" {
			if update.Message.From.ID != adminID || update.InlineQuery.From.ID != adminID {
				continue
			}
		}

		// Message
		if update.Message != nil {
			message(update.Message)
		}

		// Inline query
		if update.InlineQuery != nil && update.InlineQuery.Query != "" {
			inline(update.InlineQuery)
		}

		// CallbackQuery
		if update.CallbackQuery != nil {
			callback(update.CallbackQuery)
		}
	}
}
