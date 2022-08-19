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
	fromAccount := domain.NewAccount()
	toAccount := domain.NewAccount()

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

		fromAccountTemp, err := repo.GetAccount(transaction.From)
		if err != nil {
			log.Fatal(err)
		}

		fromAccount.ID = fromAccountTemp.ID
		fromAccount.Number = fromAccountTemp.Number
		fromAccount.Balance = fromAccountTemp.Balance

		toAccountTemp, err := repo.GetAccount(transaction.To)
		if err != nil {
			log.Fatal(err)
		}

		toAccount.ID = toAccountTemp.ID
		toAccount.Number = toAccountTemp.Number
		toAccount.Balance = toAccountTemp.Balance

		err = repo.SaveTransaction(*fromAccount, *toAccount, *transaction)

		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
		}

		fromAccountTemp, err = repo.GetAccount(transaction.From)
		if err != nil {
			log.Fatal(err)
		}
		fromAccount.Balance = fromAccountTemp.Balance

		toAccountTemp, err = repo.GetAccount(transaction.To)
		if err != nil {
			log.Fatal(err)
		}
		toAccount.Balance = toAccountTemp.Balance
	} else {
		msg := "Only POST Method"
		http.Error(w, msg, http.StatusBadRequest)
	}

	transactionSavedResponse(w, *fromAccount, *toAccount)
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

func transactionSavedResponse(w http.ResponseWriter, fromAccount domain.Account, toAccount domain.Account) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	resp["from Account Number"] = fromAccount.Number
	resp["from Account Amount"] = strconv.FormatFloat(fromAccount.Balance, 'f', -1, 64)
	resp["to Account Number"] = toAccount.Number
	resp["to Account Amount"] = strconv.FormatFloat(toAccount.Balance, 'f', -1, 64)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
