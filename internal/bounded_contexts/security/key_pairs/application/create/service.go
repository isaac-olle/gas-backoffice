package create

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs"
)

type Service struct {
	repository key_pairs.Repository
	eventBus   event.Bus
}

func NewService(repository key_pairs.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	privateKeyPEM, publicKeyPEM, err := generateKeyPair(cmd.algorithm)
	if err != nil {
		return app_errors.InternalError(key_pairs.ErrKeyPairNotCreated).
			WithDescription("failed to generate key pair").
			WithCause(err)
	}

	keyPair, err := key_pairs.NewKeyPair(cmd.id, cmd.protectedResourceId, cmd.algorithm, privateKeyPEM, publicKeyPEM)
	if err != nil {
		return app_errors.BadRequestError(key_pairs.ErrKeyPairNotCreated).
			WithDescription("failed to create key pair").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, keyPair); err != nil {
		return app_errors.InternalError(key_pairs.ErrKeyPairNotSaved).
			WithDescription("failed to save key pair").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, keyPair.PullDomainEvents())
	return nil
}

func generateKeyPair(algorithm string) (privateKeyPEM string, publicKeyPEM string, err error) {
	switch algorithm {
	case "RS256", "RS384", "RS512":
		return generateRSAKeyPair()
	case "ES256", "ES384", "ES512":
		return generateECKeyPair(algorithm)
	default:
		return "", "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func generateRSAKeyPair() (string, string, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(privBytes), hex.EncodeToString(pubBytes), nil
}

func generateECKeyPair(algorithm string) (string, string, error) {
	var curve elliptic.Curve
	switch algorithm {
	case "ES256":
		curve = elliptic.P256()
	case "ES384":
		curve = elliptic.P384()
	case "ES512":
		curve = elliptic.P521()
	}

	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(privBytes), hex.EncodeToString(pubBytes), nil
}
