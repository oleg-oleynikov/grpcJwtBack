package db

import (
	"fmt"
	"grpcJwt/internal/service"
)

type InMemorySessionStore struct {
	service.SessionStore

	sessions map[string]*service.Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*service.Session),
	}
}

func (store *InMemorySessionStore) Save(session *service.Session) error {
	// _, ok := store.sessions[session.AccountId]
	// if ok {
	// 	return fmt.Errorf("session with accountId %s already exist", session.AccountId)
	// }

	store.sessions[session.AccountId] = session.Clone()

	return nil
}

func (store *InMemorySessionStore) FindById(accountId string) (*service.Session, error) {
	session, ok := store.sessions[accountId]
	if !ok {
		return nil, fmt.Errorf("session with id %s doesnt exist", accountId)
	}

	return session, nil
}
