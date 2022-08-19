package domain

import uuid "github.com/satori/go.uuid"

type TransactionRepository interface {
	SaveTransaction(fromAccount Account, toAccount Account, transaction Transaction) error
	GetAccount(accountNumber string) (Account, error)
	CreateAccount(account Account) error
}

type Transaction struct {
	ID     string
	From   string
	To     string
	Amount float64
}

func NewTransaction() *Transaction {
	t := &Transaction{}
	t.ID = uuid.NewV4().String()
	return t
}
