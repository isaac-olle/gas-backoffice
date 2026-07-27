package main

// Seed inicialitza la base de dades amb les credencials inicials irrepetibles del sistema.
// Genera secrets (claus RSA) i els persiste — es pot tornar a córrer per recuperar-los.
//
// Crea:
//   - Un admin "Gas Seed Admin" amb la seva clau d'assertion RSA
//
// Imprimeix els outputs per stdout com a KEY=VALUE.
// Amb --env-file escriu també els outputs al fitxer .env indicat (ús local).
//
// Prerequisit per al sistema complet: córrer app_setup després del seed.
//
// go run ./cmd/db_seed
// go run ./cmd/db_seed --env-file .env.local

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/isaac-olle/gas-backoffice/config"
)

func main() {
	envFile := flag.String("env-file", "", "path to a .env file to load and update (optional, local use)")
	adminTokenEnvFile := flag.String("admin-token-env-file", "", "path to the .env file to write admin token vars (defaults to --env-file if set)")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil && !os.IsNotExist(err) {
			fatalf("load env file %s: %v", *envFile, err)
		}
	}

	ctx := context.Background()
	dsn := config.EnvOr("POSTGRES_DSN", "postgres://gas_admin:gas_admin@localhost:5432/gas?sslmode=disable")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fatalf("connect to postgres: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		fatalf("ping postgres: %v", err)
	}

	adminID := uuid.New().String()
	adminAssertionKeyID := uuid.New().String()

	adminPrivKey, adminPubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		fatalf("generate admin assertion key pair: %v", err)
	}
	adminPrivKeyPEM, err := marshalPrivateKey(adminPrivKey)
	if err != nil {
		fatalf("marshal admin private key: %v", err)
	}

	// Recover existing IDs if seed was already applied
	alreadyApplied := false
	var existingAdminID string
	if err = pool.QueryRow(ctx,
		`SELECT id FROM backoffice.users WHERE name = $1 AND is_system = TRUE AND role = 'admin' LIMIT 1`,
		"Gas Seed Admin",
	).Scan(&existingAdminID); err == nil {
		alreadyApplied = true
		adminID = existingAdminID
		fmt.Fprintf(os.Stderr, "⚠ Seed already applied, recovering existing IDs\n")
		_ = pool.QueryRow(ctx,
			`SELECT id FROM security.user_assertion_keys WHERE user_id = $1 LIMIT 1`,
			adminID,
		).Scan(&adminAssertionKeyID)

		existingPrivKeyB64 := os.Getenv("ADMIN_ASSERTION_PRIVATE_KEY")
		if existingPrivKeyB64 != "" {
			existingPrivKeyPEM, decErr := base64.StdEncoding.DecodeString(existingPrivKeyB64)
			if decErr == nil {
				if recovered, parseErr := parsePrivateKey(string(existingPrivKeyPEM)); parseErr == nil {
					adminPrivKey = recovered
					adminPrivKeyPEM = string(existingPrivKeyPEM)
				}
			}
		}
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		fatalf("begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	exec_ := func(query string, args ...any) {
		if _, err = tx.Exec(ctx, query, args...); err != nil {
			fatalf("exec query: %v\nquery: %s", err, query)
		}
	}

	exec_(`INSERT INTO security.user_assertion_keys (id, user_id, algorithm, public_key, author_id, author_role)
		VALUES ($1, $2, $3, $4, $2, 'admin')
		ON CONFLICT (id) DO UPDATE SET author_id = EXCLUDED.author_id, author_role = EXCLUDED.author_role`,
		adminAssertionKeyID, adminID, "RS256", adminPubKeyPEM)

	if err = tx.Commit(ctx); err != nil {
		fatalf("commit transaction: %v", err)
	}

	backofficePRName := config.EnvOr("GAS_BACKOFFICE_PR_NAME", "Gas Backoffice")
	if err = runSetup(ctx, pool, adminID, backofficePRName); err != nil {
		fatalf("run setup: %v", err)
	}

	adminPrivKeyB64 := base64.StdEncoding.EncodeToString([]byte(adminPrivKeyPEM))

	outputs := map[string]string{
		"SEED_ADMIN_ID":               adminID,
		"ADMIN_ASSERTION_KEY_ID":      adminAssertionKeyID,
		"ADMIN_ASSERTION_PRIVATE_KEY": adminPrivKeyB64,
	}
	for k, v := range outputs {
		fmt.Printf("%s=%s\n", k, v)
	}

	targetEnv := *adminTokenEnvFile
	if targetEnv == "" {
		targetEnv = *envFile
	}
	if targetEnv != "" && !alreadyApplied {
		if err = writeEnvFile(targetEnv, outputs); err != nil {
			fatalf("write admin token vars: %v", err)
		}
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  adminID,
		"sub":  adminID,
		"aud":  "gas",
		"iat":  now.Unix(),
		"exp":  now.Add(15 * time.Minute).Unix(),
		"jti":  uuid.New().String(),
		"role": "admin",
	})
	token.Header["kid"] = adminAssertionKeyID
	assertion, err := token.SignedString(adminPrivKey)
	if err != nil {
		fatalf("sign admin assertion: %v", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Seed completed\n\n")
	fmt.Fprintf(os.Stderr, "Admin assertion token:\n%s\n", assertion)
}

func writeEnvFile(file string, vars map[string]string) error {
	existing, _ := godotenv.Read(file)
	if existing == nil {
		existing = map[string]string{}
	}
	for k, v := range vars {
		existing[k] = v
	}
	return godotenv.Write(existing, file)
}

func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaKey, nil
}

func generateRSAKeyPair() (*rsa.PrivateKey, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, "", fmt.Errorf("marshal public key: %w", err)
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return priv, pubPEM, nil
}

func marshalPrivateKey(priv *rsa.PrivateKey) (string, error) {
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("marshal private key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "seed: "+format+"\n", args...)
	os.Exit(1)
}
