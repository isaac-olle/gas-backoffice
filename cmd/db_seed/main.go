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
// go run ./cmd/db_seed --env-file .env.local
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
	userID := uuid.New().String()
	prID := uuid.New().String()
	keyPairID := uuid.New().String()
	// users
	scopeUsersCreateID := uuid.New().String()
	scopeUsersReadID := uuid.New().String()
	scopeUsersReadAllID := uuid.New().String()
	scopeUsersDeleteID := uuid.New().String()
	scopeUsersRenameID := uuid.New().String()
	// protected_resources
	scopePRCreateID := uuid.New().String()
	scopePRReadID := uuid.New().String()
	scopePRRenameID := uuid.New().String()
	scopePRDeleteID := uuid.New().String()
	// scopes
	scopeScopeCreateID := uuid.New().String()
	scopeScopeReadID := uuid.New().String()
	scopeScopeDeleteID := uuid.New().String()
	scopeScopeAddAllowedClientID := uuid.New().String()
	scopeScopeRemoveAllowedClientID := uuid.New().String()
	// clients
	scopeClientCreateID := uuid.New().String()
	scopeClientReadID := uuid.New().String()
	scopeClientDeleteID := uuid.New().String()
	scopeClientAddFlowID := uuid.New().String()
	scopeClientRemoveFlowID := uuid.New().String()
	scopeClientAddRedirectUriID := uuid.New().String()

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

	alreadyApplied := false

	// Recover existing seed IDs if already applied
	var existingAdminID string
	err = pool.QueryRow(ctx, `SELECT id FROM backoffice.admins WHERE name = $1 LIMIT 1`, "Gas-Seed-Admin").Scan(&existingAdminID)
	if err == nil {
		alreadyApplied = true
		adminID = existingAdminID
		fmt.Fprintf(os.Stderr, "⚠ Seed already applied, recovering existing IDs\n")
		_ = pool.QueryRow(ctx, `SELECT id FROM security.admin_assertion_keys WHERE admin_id = $1 LIMIT 1`, adminID).Scan(&adminAssertionKeyID)
		_ = pool.QueryRow(ctx, `SELECT id FROM backoffice.users WHERE admin_id = $1 AND name = $2 LIMIT 1`, adminID, "Gas-User").Scan(&userID)
		_ = pool.QueryRow(ctx, `SELECT id FROM backoffice.protected_resources WHERE name = $1 LIMIT 1`, "Gas-Backoffice").Scan(&prID)
		_ = pool.QueryRow(ctx, `SELECT id FROM security.key_pairs WHERE protected_resource_id = $1 LIMIT 1`, prID).Scan(&keyPairID)

		existingPrivKeyB64 := os.Getenv("ADMIN_ASSERTION_PRIVATE_KEY")
		if existingPrivKeyB64 != "" {
			existingPrivKeyPEM, decErr := base64.StdEncoding.DecodeString(existingPrivKeyB64)
			if decErr == nil {
				recovered, parseErr := parsePrivateKey(string(existingPrivKeyPEM))
				if parseErr == nil {
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

	exec_(`INSERT INTO backoffice.admins (id, name, is_system) VALUES ($1, $2, TRUE)
		ON CONFLICT (id) DO NOTHING`,
		adminID, "Gas-Seed-Admin")

	exec_(`INSERT INTO security.admin_assertion_keys (id, admin_id, algorithm, public_key, author_id, author_role)
		      VALUES ($1, $2, $3, $4, '', '') ON CONFLICT (id) DO NOTHING`,
		adminAssertionKeyID, adminID, "RS256", adminPubKeyPEM)

	exec_(`INSERT INTO backoffice.users (id, name, admin_id, protected_resource_ids, is_system, author_id, author_role)
		      VALUES ($1, $2, $3, $4, TRUE, '', '') ON CONFLICT (id) DO NOTHING`,
		userID, "Gas-User", adminID, []string{prID})

	seedScopes := []struct {
		id   string
		name string
	}{
		{scopeUsersCreateID, "users:create"},
		{scopeUsersReadID, "users:read"},
		{scopeUsersReadAllID, "users:read_all"},
		{scopeUsersDeleteID, "users:delete"},
		{scopeUsersRenameID, "users:rename"},
		{scopePRCreateID, "protected_resources:create"},
		{scopePRReadID, "protected_resources:read"},
		{scopePRRenameID, "protected_resources:rename"},
		{scopePRDeleteID, "protected_resources:delete"},
		{scopeScopeCreateID, "scopes:create"},
		{scopeScopeReadID, "scopes:read"},
		{scopeScopeDeleteID, "scopes:delete"},
		{scopeScopeAddAllowedClientID, "scopes:add_allowed_client"},
		{scopeScopeRemoveAllowedClientID, "scopes:remove_allowed_client"},
		{scopeClientCreateID, "clients:create"},
		{scopeClientReadID, "clients:read"},
		{scopeClientDeleteID, "clients:delete"},
		{scopeClientAddFlowID, "clients:add_flow"},
		{scopeClientRemoveFlowID, "clients:remove_flow"},
		{scopeClientAddRedirectUriID, "clients:add_redirect_uri"},
	}

	// Recover existing scope IDs by name to keep them stable across re-runs
	for i, s := range seedScopes {
		var existingID string
		if err2 := pool.QueryRow(ctx, `SELECT id FROM backoffice.scopes WHERE name = $1 AND protected_resource_id = $2 LIMIT 1`, s.name, prID).Scan(&existingID); err2 == nil {
			seedScopes[i].id = existingID
		}
	}

	allScopeIDs := make([]string, len(seedScopes))
	for i, s := range seedScopes {
		allScopeIDs[i] = s.id
	}

	exec_(`INSERT INTO backoffice.protected_resources (id, name, user_id, admin_id, scope_ids, is_system, author_id, author_role)
		      VALUES ($1, $2, $3, $4, $5, TRUE, '', '')
		ON CONFLICT (id) DO UPDATE SET scope_ids = EXCLUDED.scope_ids`,
		prID, "Gas-Backoffice", userID, adminID, allScopeIDs)

	exec_(`INSERT INTO security.key_pairs (id, protected_resource_id, algorithm, public_key, private_key, author_id, author_role)
		      VALUES ($1, $2, $3, $4, $5, '', '') ON CONFLICT (id) DO NOTHING`,
		keyPairID, prID, "RS256", backofficePubKeyPEM, backofficePrivKeyPEM)

	for _, s := range seedScopes {
		exec_(`INSERT INTO backoffice.scopes (id, name, protected_resource_id, user_id, admin_id, is_system, author_id, author_role)
			      VALUES ($1, $2, $3, $4, $5, TRUE, '', '')
			ON CONFLICT (protected_resource_id, name) WHERE deleted_at IS NULL DO NOTHING`,
			s.id, s.name, prID, userID, adminID)
	}

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
		if err = writeEnvFile(*envFile, map[string]string{"BACKOFFICE_PR_ID": prID}); err != nil {
			fatalf("write env file: %v", err)
		}
	}

	targetEnv := *adminTokenEnvFile
	if targetEnv == "" {
		targetEnv = *envFile
	}
	if targetEnv != "" && !alreadyApplied {
		if err = writeEnvFile(targetEnv, map[string]string{
			"SEED_ADMIN_ID":               adminID,
			"ADMIN_ASSERTION_KEY_ID":      adminAssertionKeyID,
			"ADMIN_ASSERTION_PRIVATE_KEY": adminPrivKeyB64,
		}); err != nil {
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
