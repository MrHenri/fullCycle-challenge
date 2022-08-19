package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MrHenri/fullCycle-challenge/domain"
	"github.com/MrHenri/fullCycle-challenge/dto"
	"github.com/MrHenri/fullCycle-challenge/infrastructure/repository"
	_ "github.com/mattn/go-sqlite3"
)

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := domain.NewTransaction()
	if r.Method == http.MethodPost {
		db := setupDb()
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var transactionDto dto.Transaction

		err := dec.Decode(&transactionDto)

		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
		}

		transaction.From = transactionDto.From
		transaction.To = transactionDto.To
		transaction.Amount = transactionDto.Amount

		repo := repository.NewTransactionRepositoryDb(db)

		fromAccount, err := repo.GetAccount(transaction.From)
		if err != nil {
			log.Fatal(err)
		}

		toAccount, err := repo.GetAccount(transaction.To)
		if err != nil {
			log.Fatal(err)
		}
		err = repo.SaveTransaction(fromAccount, toAccount, *transaction)

		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
		}
	} else {
		msg := "Only POST Method"
		http.Error(w, msg, http.StatusBadRequest)
	}

	transactionSavedResponse(w, *transaction)
}

func PostAccount(w http.ResponseWriter, r *http.Request) {
	account := domain.NewAccount()
	if r.Method == http.MethodPost {
		db := setupDb()
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var accountDto dto.Account

		err := dec.Decode(&accountDto)

		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
		}

		account.Number = accountDto.Number

		repo := repository.NewTransactionRepositoryDb(db)

		err = repo.CreateAccount(*account)

		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
		}
	} else {
		msg := "Only POST Method"
		http.Error(w, msg, http.StatusBadRequest)
	}

	accountCreatedResponse(w, *account)
}

func setupDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./bank.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func accountCreatedResponse(w http.ResponseWriter, account domain.Account) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	resp["account ID"] = account.ID
	resp["account Number"] = account.Number
	resp["balance"] = strconv.FormatFloat(account.Balance, 'f', -1, 64)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func transactionSavedResponse(w http.ResponseWriter, transaction domain.Transaction) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	resp["transaction ID"] = transaction.ID
	resp["transaction From"] = transaction.From
	resp["transaction To"] = transaction.To
	resp["transaction Amount"] = strconv.FormatFloat(transaction.Amount, 'f', -1, 64)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
