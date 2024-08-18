package service

import "time"

type Session struct {
	RefreshToken string
	AccountId    string
	ExpiresRt    int64
}

func NewSession(refreshToken string, accountId string, exp time.Duration) *Session {
	return &Session{
		RefreshToken: refreshToken,
		AccountId:    accountId,
		ExpiresRt:    time.Now().Add(exp).Unix(),
	}
}

func (s *Session) Clone() *Session {
	return &Session{
		RefreshToken: s.RefreshToken,
		AccountId:    s.AccountId,
		ExpiresRt:    s.ExpiresRt,
	}
}
