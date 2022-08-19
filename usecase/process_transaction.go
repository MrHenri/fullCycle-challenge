package usecase

import (
	"github.com/MrHenri/fullCycle-challenge/domain"
	"github.com/MrHenri/fullCycle-challenge/dto"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) *UseCaseTransaction {
	return &UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	fromAccount, err := u.hidratateAccount(transactionDto.From)
	if err != nil {
		return domain.Transaction{}, err
	}

	toAccount, err := u.hidratateAccount(transactionDto.To)
	if err != nil {
		return domain.Transaction{}, err
	}

	transaction := u.hidrateTransaction(transactionDto)

	err = u.TransactionRepository.SaveTransaction(*fromAccount, *toAccount, *transaction)

	if err != nil {
		return domain.Transaction{}, err
	}

	return *transaction, nil

}

func (u UseCaseTransaction) hidrateTransaction(transactionDto dto.Transaction) *domain.Transaction {
	t := domain.NewTransaction()
	t.Amount = transactionDto.Amount
	t.From = transactionDto.From
	t.To = transactionDto.To

	return t
}

func (u UseCaseTransaction) hidratateAccount(accNumber string) (*domain.Account, error) {
	tempAccount, err := u.TransactionRepository.GetAccount(accNumber)
	if err != nil {
		return &domain.Account{}, err
	}
	account := domain.NewAccount()
	account.ID = tempAccount.ID
	account.Number = tempAccount.Number
	account.Balance = tempAccount.Balance

	return account, nil
}

func (u UseCaseTransaction) CreateAccount(account *domain.Account) error {
	err := u.TransactionRepository.CreateAccount(*account)

	if err != nil {
		return err
	}

	return nil
}
