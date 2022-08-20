package main

import (
	"log"
	"net/http"

	"github.com/MrHenri/fullCycle-challenge/services"
)

func main() {
	log.Println("Server started on: http://localhost:8000")
	http.HandleFunc("/bank-accounts", services.PostAccount)
	http.HandleFunc("/bank-accounts/transfer", services.PostTransaction)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
