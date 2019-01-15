package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Response - struct
type Response struct {
	Code int
	Res  Items `json:"res"`
}

// Items - struct
type Items struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rarity      string `json:"rarity"`
	RarityName  string `json:"rarity_name"`
	Value       int    `json:"value"`
	Estimate    int    `json:"estimate"`
	Craftable   bool   `json:"craftable"`
	Reborn      int    `json:"reborn"`
	Power       int    `json:"power"`
	PowerArmor  int    `json:"power_armor"`
	PowerShield int    `json:"power_shield"`
	DragonPower int    `json:"dragon_power"`
	Critical    int    `json:"critical"`
	CraftPnt    int    `json:"craft_pnt"`
	ConsVal     string `json:"cons_val"`
}

// ResponseCrafts - struct
type ResponseCrafts struct {
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

func main() {
	// Messaggio contrabbandiere
	message := "Benvenuto casteponters! Puoi creare oggetti per il Contrabbandiere ed egli provvederà a valutarli e ricompensarti adeguatamente, purtroppo però è disponibile solamente di giorno. Quando lascia la piazza, aggiorna la sua fornitura e quando torna ti propone affari diversi.Falce dello Stregone (L) al prezzo di 198.579 §	Accetti l'incarico di questo oggetto? Se l'offerta che ti propone non ti sembra valida, puoi cambiarla. Hai ancora a disposizione 5 offerte per oggi."

	// Estraggo items
	messageItems := strings.Split(strings.Split(message, ".")[2], "(")[0]
	log.Println(messageItems)

	// Recupero items da api
	responseItemsData := callFenix("http://fenixweb.net:3300/api/v2/V71L06foMz3ajlp811224/items/Falce%20dello%20Stregone")

	var items Response
	json.Unmarshal(responseItemsData, &items)
	fmt.Println(items)

	// Ricerco cosa ha bisogno per essere craftato
	responseBaseCraft := callFenix(fmt.Sprintf("http://fenixweb.net:3300/api/v2/V71L06foMz3ajlp811224/crafts/%v/needed", items.Res.ID))

	var crafts ResponseCrafts
	json.Unmarshal(responseBaseCraft, &crafts)
	fmt.Println(crafts)

	var listForCraft []string
	for _, craft := range crafts.Res {
		if craft.Craftable == 1 {
			listForCraft = append(listForCraft, craft.Name)
		}
	}

	// TOFIX
	log.Println(listForCraft)
}

func callFenix(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}
