package loot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// CallFenixWs - Call Fenix WebService
func CallFenixWs(url string) []byte {
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
