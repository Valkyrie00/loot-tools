package loot

// ItemResponse - Fenix item response struct
type ItemResponse struct {
	Code int
	Res  []Item `json:"res"`
}

// Item - Item struct
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rarity      string `json:"rarity"`
	RarityName  string `json:"rarity_name"`
	Value       int    `json:"value"`
	Estimate    int    `json:"estimate"`
	Craftable   int    `json:"craftable"`
	Reborn      int    `json:"reborn"`
	Power       int    `json:"power"`
	PowerArmor  int    `json:"power_armor"`
	PowerShield int    `json:"power_shield"`
	DragonPower int    `json:"dragon_power"`
	Critical    int    `json:"critical"`
	CraftPnt    int    `json:"craft_pnt"`
	ConsVal     string `json:"cons_val"`
}

// Items - Item list
type Items []Item

// FindItemByID - Search item By name
func (items Items) FindItemByID(id int) Item {
	for _, item := range items {
		if item.ID == id {
			return item
		}
	}

	return Item{}
}

// CraftResponse - struct
type CraftResponse struct {
	Code int
	Item string
	Res  []Craft `json:"res"`
}

// Craft - struct
type Craft struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	Craftable int    `json:"craftable"`
}

// ItemsCraftingMapType - Items Crafting Map Type
type ItemsCraftingMapType map[int][]string
