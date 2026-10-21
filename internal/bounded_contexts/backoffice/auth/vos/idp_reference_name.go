package vos

import "fmt"

var allowedIdpReferenceNames = map[string]IdpReferenceName{
	"google": IdpReferenceName("google"),
}

type IdpReferenceName string

func NewFromString(value string) (IdpReferenceName, error) {
	idpReferenceName := IdpReferenceName(value)

	if _, ok := allowedIdpReferenceNames[value]; !ok {
		return "", fmt.Errorf("invalid IdP reference name: %s", value)
	}

	return idpReferenceName, nil
}

func (id IdpReferenceName) String() string {
	return string(id)
}
