// generate_admin_token generates a signed RS256 JWT assertion for the admin token endpoint.
//
// Usage:
//
//	go run ./scripts/generate_admin_token
//
// Reads SEED_ADMIN_ID, ADMIN_ASSERTION_KEY_ID and ADMIN_ASSERTION_PRIVATE_KEY from .env.local (or --env flag).
// The output JWT can then be used to obtain an admin access token:
//
//	curl -X POST http://localhost/v1/backoffice/admins/token \
//	  -d "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer" \
//	  -d "assertion=<JWT>"
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("env", ".env.local", "path to the .env file")
	ttl := flag.Duration("ttl", 15*time.Minute, "token TTL (e.g. 5m, 1h)")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatalf("load %s: %v", *envFile, err)
	}

	adminID := os.Getenv("SEED_ADMIN_ID")
	kid := os.Getenv("ADMIN_ASSERTION_KEY_ID")
	keyB64 := os.Getenv("ADMIN_ASSERTION_PRIVATE_KEY")

	if adminID == "" || kid == "" || keyB64 == "" {
		log.Fatal("SEED_ADMIN_ID, ADMIN_ASSERTION_KEY_ID and ADMIN_ASSERTION_PRIVATE_KEY must be set in the env file")
	}

	keyPEM, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		log.Fatalf("decode private key: %v", err)
	}

	privKey, err := parsePrivateKey(keyPEM)
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  adminID,
		"sub":  adminID,
		"aud":  "gas",
		"iat":  now.Unix(),
		"exp":  now.Add(*ttl).Unix(),
		"jti":  uuid.New().String(),
		"role": "admin",
	})
	token.Header["kid"] = kid

	signed, err := token.SignedString(privKey)
	if err != nil {
		log.Fatalf("sign token: %v", err)
	}

	fmt.Println(signed)
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA")
	}
	return rsaKey, nil
}
