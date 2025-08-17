package proto

import "time"

type User struct {
	Id          string
	UserName    string
	UserPwdHash string
	FamilyName  string
	GivenName   string
	BirthDate   string
	Email       string
}

type Token struct {
	TokenID      string
	UserID       string
	Scope        string // copied from AuthCode
	ExpiresAt    time.Time
	RevocationID string // optional, for batch revocation
}

type AuthCode struct {
	Code        string
	UserID      string // user authenticated
	Scope       string // requested scopes
	ExpiresAt   time.Time
	ClientID    string // which client requested it
	RedirectURI string
}
