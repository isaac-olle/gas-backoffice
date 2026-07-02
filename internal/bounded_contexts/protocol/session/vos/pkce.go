package vos

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type CodeChallengeMethod string

const (
	CodeChallengePlain CodeChallengeMethod = "plain"
	CodeChallengeS256  CodeChallengeMethod = "S256"
)

func NewCodeChallengeMethod(method string) (CodeChallengeMethod, error) {
	m := CodeChallengeMethod(method)
	if m != CodeChallengePlain && m != CodeChallengeS256 {
		return "", errors.New("invalid code_challenge_method")
	}
	return m, nil
}

func (c CodeChallengeMethod) String() string { return string(c) }

type Pkce struct {
	CodeChallenge string
	method        CodeChallengeMethod
}

func NewPkce(challenge string, methodString string) (*Pkce, error) {
	if challenge == "" && methodString == "" {
		return nil, nil
	}
	if challenge == "" {
		return nil, errors.New("code_challenge is required when code_challenge_method is provided")
	}
	if methodString == "" {
		return nil, errors.New("code_challenge_method is required when code_challenge is provided")
	}
	method := CodeChallengeMethod(methodString)
	if method != CodeChallengePlain && method != CodeChallengeS256 {
		return nil, errors.New("invalid code_challenge_method: must be 'plain' or 'S256'")
	}
	return &Pkce{CodeChallenge: challenge, method: method}, nil
}

func (p *Pkce) Challenge() string           { return p.CodeChallenge }
func (p *Pkce) Method() CodeChallengeMethod { return p.method }

func (p *Pkce) Verify(codeVerifier string) bool {
	switch p.method {
	case CodeChallengePlain:
		return codeVerifier == p.CodeChallenge
	case CodeChallengeS256:
		h := sha256.Sum256([]byte(codeVerifier))
		computed := base64.RawURLEncoding.EncodeToString(h[:])
		return computed == p.CodeChallenge
	}
	return false
}
