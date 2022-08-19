package domain

import uuid "github.com/satori/go.uuid"

type Account struct {
	ID      string
	Number  string
	Balance float64
}

func NewAccount() *Account {
	a := &Account{}
	a.ID = uuid.NewV4().String()
	a.Balance = 300.0
	return a
}
