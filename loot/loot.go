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
	ForceMap               string
)

func init() {
	ForceMap = os.Getenv("FORCE_MAP")
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
	var itemsCraftingMap CraftingMapType

	if ForceMap == "true" {
		itemsCraftingMap = createMapFile(items)
	} else {
		itemsCraftingMap = loadMapFile(items)
	}

	return itemsCraftingMap
}

func loadMapFile(items map[int]Item) CraftingMapType {
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
		itemsCraftingMap = createMapFile(items)
	}

	return itemsCraftingMap
}

func createMapFile(items map[int]Item) CraftingMapType {
	itemsCraftingMap := make(CraftingMapType)

	// Not Exist
	for i, item := range items {
		// Debugger limit
		// if item.ID != 229 {
		// 	continue
		// }

		log.Println(len(items), i, "Getting crafting list: "+item.Name)

		// Get all needed IDs and list
		var neededIDS []int
		makeNeeded(item.ID, &neededIDS)
		neededIDS = append(neededIDS, item.ID)

		// Minimize - return string:quantity
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

	return itemsCraftingMap
}

func reordersCrafting(minifiedCrafting []MinyCrafting, items *map[int]Item) []string {
	var reorderedCrafting []string

	// Sort crafting
	for _, miny := range minifiedCrafting {
		item := (*items)[miny.ItemID]

		// Stupid split ( Max 3 qty for element in a row )
		if miny.Quantity == 1 {
			reorderedCrafting = append(reorderedCrafting, fmt.Sprintf("%s", item.Name))
		} else if miny.Quantity > 1 && miny.Quantity < 3 {
			reorderedCrafting = append(reorderedCrafting, fmt.Sprintf("%s,%v", item.Name, miny.Quantity))
		} else if miny.Quantity >= 3 {
			x := 0
			for index := 0; index < miny.Quantity; index++ {
				x++

				if x == 3 {
					reorderedCrafting = append(reorderedCrafting, fmt.Sprintf("%s,%v", item.Name, x))
					x = 0
					continue
				}
			}
			if x > 0 {
				reorderedCrafting = append(reorderedCrafting, fmt.Sprintf("%s,%v", item.Name, x))
			}
		}
	}

	return reorderedCrafting
}

func minimizeCrafting(itemsList []int) []MinyCrafting {
	var minifieds []MinyCrafting

	for _, item := range itemsList {
		exists := false

		for kMiny, vMiny := range minifieds {
			if vMiny.ItemID == item {
				minifieds[kMiny].Quantity++
				exists = true
			}
		}

		if exists == false {
			minifieds = append(minifieds, MinyCrafting{ItemID: item, Quantity: 1})
		}
	}

	return minifieds
}

func makeNeeded(itemID int, neededIDS *[]int) {
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
			makeNeeded(craft.ID, neededIDS)
		}
	}
}
