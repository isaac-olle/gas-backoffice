package session

import "time"

const AccessTokenTTL = time.Hour

type TokenResponse struct {
	AccessToken string
	IDToken     string
	TokenType   string
	ExpiresIn   int
	Scope       string
}
