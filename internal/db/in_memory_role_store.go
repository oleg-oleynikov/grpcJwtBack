package db

import (
	"fmt"
	"grpcJwt/internal/service"
)

type InMemoryRoleStore struct {
	service.RoleStore

	roles map[string]*service.Role
}

func NewInMemoryRoleStore() *InMemoryRoleStore {
	return &InMemoryRoleStore{
		roles: make(map[string]*service.Role),
	}
}

func (store *InMemoryRoleStore) Save(role *service.Role) error {
	r, exist := store.roles[role.Title]
	if exist {
		return fmt.Errorf("role with Title %s already exist", role.Title)
	}

	store.roles[role.Title] = r
	return nil
}

func (store *InMemoryRoleStore) FindByTitle(titleRole string) (*service.Role, error) {
	r, exist := store.roles[titleRole]
	if !exist {
		return nil, fmt.Errorf("role with title %s doesnt exist", titleRole)
	}

	return r, nil
}
