package session

type TokenClaims struct {
	Scope               string
	ClientID            string
	UserID              string
	ProtectedResourceId string
}

type Session interface {
	TokenClaims() TokenClaims
}

type AccessTokenClaims struct {
	Issuer string
	Scope  string
}

type IDTokenClaims struct {
	Issuer   string
	Subject  string
	Audience string
}
