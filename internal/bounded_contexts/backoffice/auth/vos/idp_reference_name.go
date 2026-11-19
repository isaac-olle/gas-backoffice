package vos

import (
	"fmt"
	"regexp"
)

var idpReferenceNamePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type IdpReferenceName string

func NewFromString(value string) (IdpReferenceName, error) {
	idpReferenceName := IdpReferenceName(value)

	if !idpReferenceNamePattern.MatchString(value) {
		return "", fmt.Errorf("invalid IdP reference name: %s", value)
	}

	return idpReferenceName, nil
}

func (id IdpReferenceName) String() string {
	return string(id)
}
