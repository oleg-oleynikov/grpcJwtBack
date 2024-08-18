package db

import (
	"fmt"
	"grpcJwt/internal/service"
	"sync"
)

type InMemoryAccountStore struct {
	service.AccountStore
	sync.RWMutex

	accounts map[string]*service.Account
}

func NewInMemoryAccountStore() *InMemoryAccountStore {
	return &InMemoryAccountStore{
		accounts: make(map[string]*service.Account),
	}
}

func (store *InMemoryAccountStore) Find(email string) (*service.Account, error) {
	store.RLock()
	defer store.RUnlock()

	account, ok := store.accounts[email]
	if !ok {
		return nil, fmt.Errorf("store doesnt contains account with email %s", email)
	}

	return account, nil
}

func (store *InMemoryAccountStore) Save(account *service.Account) error {
	store.Lock()
	defer store.Unlock()

	_, ok := store.accounts[account.Email]
	if ok {
		return fmt.Errorf("store already exist account with email %s", account.Email)
	}

	store.accounts[account.Email] = account.Clone()
	return nil
}
