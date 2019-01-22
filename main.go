package main

import (
	"log"

	"github.com/Valkyrie00/loot-tools/bot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// PingServer
	log.Println("LOAD - Ping Server")
	go server()

	log.Println("LOAD - Bot")
	bot.Handler()
}
