package session

type AuthorizationCodeSessionPrimitives struct {
	SessionID     string
	ClientID      string
	RedirectURI   string
	Scope         string
	State         string
	ResponseTypes []string
	Code          string
	Pkce          *PkcePrimitives
	SessionState  string
	Flow          string
	Consumed      bool
}

type PkcePrimitives struct {
	CodeChallenge string
	Method        string
}

type ClientCredentialsSessionPrimitives struct {
	ID       string
	ClientID string
	UserId   string
	AdminID  string
	Scope    string
	IssuedAt string
}
