package loot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	lootToken      string
	cacheCraftsMap map[int]CraftResponse
)

func init() {
	lootToken = "V71L06foMz3ajlp811224"

	cacheCraftsMap = make(map[int]CraftResponse)
}

// GetCraftableItems - Get from remote all craftable items
func GetCraftableItems() []Item {
	url := fmt.Sprintf("http://fenixweb.net:3300/api/v2/%v/items", lootToken)

	response := CallFenixWs(url)
	var responseData ItemResponse
	json.Unmarshal(response, &responseData)

	var items []Item
	for _, item := range responseData.Res {
		if item.Craftable == 1 {
			items = append(items, item)
		}
	}

	return items
}

// GetCraftingMap -
func GetCraftingMap(itemsList []Item) ItemsCraftingMapType {
	itemsCraftingMap := make(ItemsCraftingMapType)

	if _, err := os.Stat("assets/craftingMaps.json"); err == nil {
		// Exist
		craftingMapsFile, err := os.Open("assets/craftingMaps.json")
		if err != nil {
			fmt.Println(err)
		}

		log.Println("craftingMaps.json - Successfully Opened!")

		byteValue, _ := ioutil.ReadAll(craftingMapsFile)
		json.Unmarshal([]byte(byteValue), &itemsCraftingMap)

	} else if os.IsNotExist(err) {
		log.Println("craftingMaps.json - Not found!")

		// Not Exist
		for i, item := range itemsList {
			log.Println(len(itemsList)-1, i, "Getting crafting list: "+item.Name)

			var itemID int
			itemID = item.ID

			var needsToCraft []string
			needsToCraft = []string{item.Name}

			mapCraft([]int{itemID}, &needsToCraft)

			itemsCraftingMap[itemID] = needsToCraft

			//Convert to JSON - ONLY FOR MOCK OR FOR SAVE JSON
			// jsonMap, err := json.Marshal(itemsCraftingMap)
			// if err != nil {
			// 	log.Panicln(err)
			// }
			// log.Println(string(jsonMap))
		}

	}

	return itemsCraftingMap
}

// listitems è la lista di item da ricercare, toCraft è la lista di tutte le cose da craftare
func mapCraft(listItems []int, toCrafts *[]string) {
	for index := 0; index < 1; index++ {
		item := listItems[index]

		var crafts CraftResponse
		if cacheCraftsMap[item].Item != "" {
			crafts = cacheCraftsMap[item]
		} else {
			url := fmt.Sprintf("http://fenixweb.net:3300/api/v2/%v/crafts/%v/needed", lootToken, item)
			responseBaseCraft := CallFenixWs(url)
			json.Unmarshal(responseBaseCraft, &crafts)

			// Add response in cache
			cacheCraftsMap[item] = crafts
		}

		for _, craft := range crafts.Res {
			if craft.Craftable == 1 {
				listItems = append([]int{craft.ID}, listItems...)
				*toCrafts = append([]string{craft.Name}, *toCrafts...)
				mapCraft(listItems, toCrafts)
			}
		}
	}
}
