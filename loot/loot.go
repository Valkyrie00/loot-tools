package loot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload" //Autoload env
)

var (
	LootToken              string
	CacheResponseCraftsMap map[int]CraftResponse
	RarityList             = []string{"C", "NC", "R", "UR", "L", "E", "UE", "U", "S", "I", "D", "X"}
)

func init() {
	LootToken = os.Getenv("LOOT_APIKEY")
	CacheResponseCraftsMap = make(map[int]CraftResponse)
}

// SyncItems - getter and setter lootbot items
func SyncItems() (map[int]Item, CraftingMapType) {
	// Step 1 - load craftable items
	craftableItems := getCraftableItems()

	mapOfCrafts := mapCrafts(craftableItems)

	return craftableItems, mapOfCrafts
}

// getCraftableItems - Get from remote all craftable items
func getCraftableItems() map[int]Item {
	url := fmt.Sprintf("http://fenixweb.net:3300/api/v2/%s/items", LootToken)

	response := CallFenixWs(url)
	var responseData ItemResponse
	err := json.Unmarshal(response, &responseData)
	if err != nil {
		log.Panicln("Error getting craftable items", err)
	}

	var items = make(map[int]Item)
	for _, item := range responseData.Res {
		if item.Craftable == 1 {
			items[item.ID] = item
		}
	}

	return items
}

func mapCrafts(items map[int]Item) CraftingMapType {
	itemsCraftingMap := make(CraftingMapType)

	if _, err := os.Stat("assets/crafting-map.json"); err == nil {
		// Exist
		craftingMapsFile, err := os.Open("assets/crafting-map.json")
		if err != nil {
			fmt.Println(err)
		}

		log.Println("LOAD - CaftingMap: OK!")

		byteValue, _ := ioutil.ReadAll(craftingMapsFile)
		err = json.Unmarshal([]byte(byteValue), &itemsCraftingMap)
		if err != nil {
			log.Panicln(err)
		}

	} else if os.IsNotExist(err) {

		// Not Exist
		for i, item := range items {
			// Debugger limit
			// if item.ID != 627 {
			// 	continue
			// }

			log.Println(len(items), i, "Getting crafting list: "+item.Name)

			// Get all needed IDs
			var neededIDS []int
			neededIDS = append(neededIDS, item.ID)
			makeNeededIDs(item.ID, &neededIDS)

			// Minimize - return itemID:quantity
			minifiedCrafting := minimizeCrafting(neededIDS)

			// Reorder and get item name
			reorderedCrafiting := reordersCrafting(minifiedCrafting, &items)

			// Associate
			itemsCraftingMap[item.ID] = reorderedCrafiting
		}

		jCrafting, err := json.Marshal(itemsCraftingMap)
		if err != nil {
			log.Panicln(err)
		}

		// Save to file
		err = ioutil.WriteFile("assets/crafting-map.json", jCrafting, 0644)
		if err != nil {
			log.Panicln(err)
		}

	}

	return itemsCraftingMap
}

func reordersCrafting(minifiedCrafting map[int]int, items *map[int]Item) []string {
	var sortedByRarity = make(map[string][]string)

	// Sort crafting
	for key, quantity := range minifiedCrafting {
		item := (*items)[key]

		// Stupid split ( Max 3 qty for element in a row )
		if quantity == 1 {
			sortedByRarity[item.Rarity] = append(sortedByRarity[item.Rarity], fmt.Sprintf("%s", item.Name))
		} else if quantity > 1 && quantity < 3 {
			sortedByRarity[item.Rarity] = append(sortedByRarity[item.Rarity], fmt.Sprintf("%s,%v", item.Name, quantity))
		} else if quantity >= 3 {
			x := 0
			for index := 0; index < quantity; index++ {
				x++

				if x == 3 {
					sortedByRarity[item.Rarity] = append(sortedByRarity[item.Rarity], fmt.Sprintf("%s,%v", item.Name, x))
					x = 0
					continue
				}
			}
			if x > 0 {
				sortedByRarity[item.Rarity] = append(sortedByRarity[item.Rarity], fmt.Sprintf("%s,%v", item.Name, x))
			}
		}
	}

	// Reording by rarity type
	var reorderedCrafting []string
	for _, rarity := range RarityList {
		for _, item := range sortedByRarity[rarity] {
			reorderedCrafting = append(reorderedCrafting, item)
		}
	}

	return reorderedCrafting
}

func minimizeCrafting(itemsIDS []int) map[int]int {
	var minified = make(map[int]int)

	for _, item := range itemsIDS {
		minified[item]++
	}

	return minified
}

func makeNeededIDs(itemID int, neededIDS *[]int) {
	var crafts CraftResponse

	// Get crafting needed
	if _, ok := CacheResponseCraftsMap[itemID]; ok {
		crafts = CacheResponseCraftsMap[itemID]
	} else {
		responseBaseCraft := CallFenixWs(fmt.Sprintf("http://fenixweb.net:3300/api/v2/%v/crafts/%v/needed", LootToken, itemID))
		err := json.Unmarshal(responseBaseCraft, &crafts)
		if err != nil {
			log.Panicln("Error getting needed items", err)
		}

		// Add response in cache
		CacheResponseCraftsMap[itemID] = crafts
	}

	// Recursive crafting loops
	for _, craft := range crafts.Res {
		if craft.Craftable == 1 {
			*neededIDS = append([]int{craft.ID}, *neededIDS...)
			makeNeededIDs(craft.ID, neededIDS)
		}
	}
}
