package handlers

import "context"

type IAccount interface {
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	Accounts(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewAccount() IAccount {
	return AccountValidator{next: newAccount()}
}

type account struct {
}

func newAccount() *account {
	return &account{}
}

// CreateAccount Create one or more accounts
// return creation result
func (h account) CreateAccount(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Accounts return account list
func (h account) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
