package main

import (
	"log"

	"github.com/Valkyrie00/loot-tools/loot"
)

func main() {
	log.Println("LOAD - Bot")
	// bot.Handler()

	loot.SyncItems()
}
