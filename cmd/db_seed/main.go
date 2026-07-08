package main

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

// Seed inicialitza la base de dades amb les dades mínimes necessàries per arrancar el sistema.
// Només es pot executar una vegada — falla si el seed ja ha estat aplicat.
//
// Crea:
//   - Un admin "Gas-Seed-Admin" amb la seva clau d'assertion RSA
//   - Un user "Gas-User" associat a l'admin
//   - Un protected resource "Gas-Backoffice" amb 3 scopes (users:create, users:read, users:read_all)
//   - Un key pair RSA per signar tokens del backoffice
//
// Imprimeix els outputs per stdout com a KEY=VALUE.
// Amb --env-file escriu també els outputs al fitxer .env indicat (ús local).
//
// go run ./cmd/seed
// go run ./cmd/db_seed --env-file config/.env
func main() {
	envFile := flag.String("env-file", "", "path to a .env file to load and update (optional, local use)")
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
	userID := uuid.New().String()
	prID := uuid.New().String()
	keyPairID := uuid.New().String()
	scopeCreateID := uuid.New().String()
	scopeReadID := uuid.New().String()
	scopeReadAllID := uuid.New().String()

	adminPrivKey, adminPubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		fatalf("generate admin assertion key pair: %v", err)
	}
	adminPrivKeyPEM, err := marshalPrivateKey(adminPrivKey)
	if err != nil {
		fatalf("marshal admin private key: %v", err)
	}

	backofficePrivKey, backofficePubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		fatalf("generate backoffice key pair: %v", err)
	}
	backofficePrivKeyPEM, err := marshalPrivateKey(backofficePrivKey)
	if err != nil {
		fatalf("marshal backoffice private key: %v", err)
	}

	var existing string
	err = pool.QueryRow(ctx, `SELECT id FROM admins WHERE name = $1 LIMIT 1`, "Gas-Seed-Admin").Scan(&existing)
	if err == nil {
		fatalf("seed already applied (admin %s exists)", existing)
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

	exec_(`INSERT INTO admins (id, name, is_system) VALUES ($1, $2, TRUE)`,
		adminID, "Gas-Seed-Admin")

	exec_(`INSERT INTO admin_assertion_keys (id, admin_id, algorithm, public_key)
		      VALUES ($1, $2, $3, $4)`,
		adminAssertionKeyID, adminID, "RS256", adminPubKeyPEM)

	exec_(`INSERT INTO users (id, name, admin_id, protected_resource_ids, is_system)
		      VALUES ($1, $2, $3, $4, TRUE)`,
		userID, "Gas-User", adminID, []string{prID})

	exec_(`INSERT INTO protected_resources (id, name, user_id, admin_id, scope_ids, is_system)
		      VALUES ($1, $2, $3, $4, $5, TRUE)`,
		prID, "Gas-Backoffice", userID, adminID,
		[]string{scopeCreateID, scopeReadID, scopeReadAllID})

	exec_(`INSERT INTO key_pairs (id, protected_resource_id, algorithm, public_key, private_key)
		      VALUES ($1, $2, $3, $4, $5)`,
		keyPairID, prID, "RS256", backofficePubKeyPEM, backofficePrivKeyPEM)

	exec_(`INSERT INTO scopes (id, name, protected_resource_id, user_id, admin_id, is_system)
		      VALUES ($1, $2, $3, $4, $5, TRUE)`,
		scopeCreateID, "users:create", prID, userID, adminID)

	exec_(`INSERT INTO scopes (id, name, protected_resource_id, user_id, admin_id, is_system)
		      VALUES ($1, $2, $3, $4, $5, TRUE)`,
		scopeReadID, "users:read", prID, userID, adminID)

	exec_(`INSERT INTO scopes (id, name, protected_resource_id, user_id, admin_id, is_system)
		      VALUES ($1, $2, $3, $4, $5, TRUE)`,
		scopeReadAllID, "users:read_all", prID, userID, adminID)

	if err = tx.Commit(ctx); err != nil {
		fatalf("commit transaction: %v", err)
	}

	adminPrivKeyB64 := base64.StdEncoding.EncodeToString([]byte(adminPrivKeyPEM))

	outputs := map[string]string{
		"BACKOFFICE_PR_ID":            prID,
		"SEED_ADMIN_ID":               adminID,
		"SEED_USER_ID":                userID,
		"ADMIN_ASSERTION_KEY_ID":      adminAssertionKeyID,
		"ADMIN_ASSERTION_PRIVATE_KEY": adminPrivKeyB64,
	}

	for k, v := range outputs {
		fmt.Printf("%s=%s\n", k, v)
	}

	if *envFile != "" {
		if err = writeEnvFile(*envFile, outputs); err != nil {
			fatalf("write env file: %v", err)
		}
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  adminID,
		"sub":  adminID,
		"aud":  "gas",
		"iat":  now.Unix(),
		"exp":  now.Add(99999 * time.Hour).Unix(),
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
