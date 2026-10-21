package session

// UserIdentityProjection is session's own read projection of the resolved
// principal for a session — only the subject claim it needs to build the JWT
// "sub". It is populated from protocol.oidc_user_identities for a federated
// login, or built directly from the client id for a client_credentials grant,
// but session never depends on the user_identity aggregate itself.
type UserIdentityProjection struct {
	Subject string
	Email   string
}

func NewUserIdentityProjection(subject, email string) *UserIdentityProjection {
	return &UserIdentityProjection{Subject: subject, Email: email}
}
