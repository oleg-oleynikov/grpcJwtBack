package service

import (
	"fmt"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id             string
	Email          string
	HashedPassword string
	Role           string
	Age            int
}

func NewAccount(email string, password string, role string, age int) (*Account, error) {
	uuid, _ := uuid.NewV7()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed hash password, %v", err)
	}

	account := &Account{
		Id:             uuid.String(),
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           role,
		Age:            age,
	}

	return account, nil
}

func (account *Account) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password))
	return err == nil
}

func (account *Account) Clone() *Account {
	return &Account{
		Email:          account.Email,
		HashedPassword: account.HashedPassword,
		Role:           account.Role,
		Age:            account.Age,
	}
}
