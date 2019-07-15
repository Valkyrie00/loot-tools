package loot

// ItemResponse - Fenix item response struct
type ItemResponse struct {
	Code int
	Res  []Item `json:"res"`
}

// Item - Item struct
type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Rarity      string  `json:"rarity"`
	Description string  `json:"description"`
	Value       int     `json:"value"`
	MaxValue    int     `json:"max_value"`
	Estimate    int     `json:"estimate"`
	Spread      int     `json:"spread"`
	SpreadTot   float32 `json:"spread_tot"`
	Craftable   int     `json:"craftable"`
	Reborn      int     `json:"reborn"`
	Power       int     `json:"power"`
	PowerArmor  int     `json:"power_armor"`
	PowerShield int     `json:"power_shield"`
	DragonPower int     `json:"dragon_power"`
	Critical    int     `json:"critical"`
	Category    int     `json:"category"`
	Cons        int     `json:"cons"`
	AllowSell   int     `json:"allow_sell"`
	RarityName  string  `json:"rarity_name"`
	CraftPnt    int     `json:"craft_pnt"`
	ConsVal     float32 `json:"cons_val"`
}

// Items - Item list
type Items []Item

// FindItemByID - Search item By name
// func (items Items) FindItemByID(id int) Item {
// 	for _, item := range items {
// 		if item.ID == id {
// 			return item
// 		}
// 	}

// 	return Item{}
// }

// CraftResponse - struct
type CraftResponse struct {
	Code int
	Item int
	Res  []Craft `json:"res"`
}

// Craft - struct
type Craft struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	Craftable int    `json:"craftable"`
}

// CraftingMapType - Crafting Map Type
type CraftingMapType map[int][]string
