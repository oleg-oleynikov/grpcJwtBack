package service

import "time"

type Session struct {
	RefreshToken string
	AccountId    string
	ExpiresAt    int64
}

func NewSession(refreshToken string, accountId string, exp time.Duration) *Session {
	return &Session{
		RefreshToken: refreshToken,
		AccountId:    accountId,
		ExpiresAt:    time.Now().Add(exp).Unix(),
	}
}

func (s *Session) Clone() *Session {
	return &Session{
		RefreshToken: s.RefreshToken,
		AccountId:    s.AccountId,
		ExpiresAt:    s.ExpiresAt,
	}
}
