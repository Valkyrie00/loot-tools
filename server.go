package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func server() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8011"
		log.Println("$PORT must be set")
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Pong")
	})

	http.ListenAndServe(":"+port, nil)
}
