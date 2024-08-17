package service

type AccountStore interface {
	Find(email string) (*Account, error)
	Save(account *Account) error
}
