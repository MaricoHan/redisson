package handlers

import "context"

type AccountValidator struct {
	next IAccount
}

// CreateAccount Create one or more accounts
// return creation result
func (h AccountValidator) CreateAccount(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.CreateAccount(ctx, nil)
}

// Accounts return account list
func (h AccountValidator) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.CreateAccount(ctx, nil)
}
