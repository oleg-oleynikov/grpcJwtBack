package service

type SessionStore interface {
	Save(session *Session) error
	FindById(accountId string) (*Session, error)
}
