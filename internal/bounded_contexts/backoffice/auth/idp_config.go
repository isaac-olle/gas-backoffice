package auth

// IdpConfig és el read model de la configuració d'un IdP,
// propietat del mòdul identity_providers (dins del BC backoffice). El mòdul auth
// només en consumeix aquesta vista de lectura per saber que l'IdP existeix
// i està actiu: on redirigir a l'IdP ho resol l'authorization server del gas.
type IdpConfig struct {
	IdpReferenceName string
}
