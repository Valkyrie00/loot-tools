package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Valkyrie00/loot-tools/loot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// ParseInlineType - Return type of inline message 1: Search - 2: Craft
func ParseInlineType(message string) int {
	if strings.Contains(message, "🔨") {
		return 2 // Craft
	}

	return 1 // Search
}

// ParseInlineCraftType - Helper
func ParseInlineCraftType(message string) (string, int, int) {
	smash := strings.Split(message, "🔨")

	if len(smash) > 1 {
		smashed := strings.TrimSpace(smash[1])
		crafting := strings.Split(smashed, ":")

		if len(crafting) > 1 {
			preParsingItem := strings.Split(crafting[0], "-")

			if len(preParsingItem) > 1 {
				// get B-XXX:
				craftingType := preParsingItem[0]
				itemID, _ := strconv.Atoi(preParsingItem[1])

				// get :XXX
				craftingIndex, _ := strconv.Atoi(crafting[1])

				return craftingType, itemID, craftingIndex
			}
		}
	}

	return "", 0, 0
}

//SearchItem - Helper
func SearchItem(s string) []loot.Item {
	var results []loot.Item
	list := craftableItems

	for _, item := range list {
		if strings.Contains(strings.ToUpper(item.Name), strings.ToUpper(s)) {
			results = append(results, item)
		}
	}

	return results
}

//SetterCraftInlineKeyboard - Rerturn craft inline keyboard
func SetterCraftInlineKeyboard(text string, inputMessage string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: text,
				SwitchInlineQueryCurrentChat: &inputMessage,
			},
		),
	)

	return &keyboard
}

//SetterSwitchCLBInlineKeyboard - Rerturn craft inline keyboard
func SetterSwitchCLBInlineKeyboard(text string, inputMessage string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:              text,
				SwitchInlineQuery: &inputMessage,
			},
		),
	)

	return &keyboard
}

//GetEnv - helper
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// DownloadFileFromURL - Download file from URL
func DownloadFileFromURL(filename string, URL string) bool {
	out, errCreate := os.Create("storage/clb/" + filename + ".txt")
	if errCreate != nil {
		log.Println("Err create file for download remote file")
		return false
	}
	defer out.Close()

	resp, errGet := http.Get(URL)
	if errGet != nil {
		log.Println("Err get file from URL")
		return false
	}
	defer resp.Body.Close()

	_, errCopy := io.Copy(out, resp.Body)
	if errCopy != nil {
		return false
	}

	return true
}

//OpenFileByFileName - Open file by filename
func OpenFileByFileName(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}

	return file
}

//ScanFile - Scan openedFile
func ScanFile(file *os.File) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Error
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return lines
}

//GetLineFromFile - Get Line From File
func GetLineFromFile(filename string, indexLine int) string {
	file := OpenFileByFileName(filename)
	lines := ScanFile(file)

	log.Println(lines[indexLine])

	return lines[indexLine]
}

//GetLinesFromFile - Get Line From File
func GetLinesFromFile(filename string) []string {
	file := OpenFileByFileName(filename)
	lines := ScanFile(file)

	return lines
}
