//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_presenters"
	commandInfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command/infrastructure"
	eventInfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/infrastructure"
	queryInfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/infrastructure"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
	"github.com/isaac-olle/gas-backoffice/internal/shared/workers"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/isaac-olle/gas-backoffice/cmd/api/bootstrap/modules"
	"github.com/isaac-olle/gas-backoffice/internal/flows"
	protocolHandler "github.com/isaac
	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/flows"
)

const appName = "gas"

var (
	testServer     *httptest.Server
	testAdminID    string
	testAdminToken string
)

type adminSeed struct {
	adminID        string
	assertionKeyID string
	signingKeyID   string
	adminPrivKey   *rsa.PrivateKey
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, pgCfg := mustConnectDB(ctx)
	defer pool.Close()

	seed := mustSeedAdmin(ctx, pool)
	testAdminID = seed.adminID

	gin.SetMode(gin.TestMode)
	testServer = httptest.NewServer(buildEngine(pool, seed.signingKeyID))
	defer testServer.Close()

	testAdminToken = mustAdminToken(seed)

	code := m.Run()

	if p, err := pgxpool.NewWithConfig(context.Background(), pgCfg); err == nil {
		cleanup(context.Background(), p, seed.adminID, seed.signingKeyID)
		p.Close()
	}

	os.Exit(code)
}

func mustConnectDB(ctx context.Context) (*pgxpool.Pool, *pgxpool.Config) {
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = 5432
	}
	cfg, _ := pgxpool.ParseConfig("")
	cfg.ConnConfig.Host = os.Getenv("POSTGRES_HOST")
	cfg.ConnConfig.Port = uint16(port)
	cfg.ConnConfig.User = os.Getenv("POSTGRES_USER")
	cfg.ConnConfig.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.ConnConfig.Database = os.Getenv("POSTGRES_DB")

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		fatalf("connect postgres: %v", err)
	}
	if err = pool.Ping(ctx); err != nil {
		fatalf("ping postgres: %v", err)
	}
	return pool, cfg
}

func mustSeedAdmin(ctx context.Context, pool *pgxpool.Pool) adminSeed {
	adminPrivKey, adminPubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		fatalf("generate admin assertion key pair: %v", err)
	}
	signingPrivKey, signingPubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		fatalf("generate signing key pair: %v", err)
	}
	signingPrivKeyPEM, err := marshalPrivateKey(signingPrivKey)
	if err != nil {
		fatalf("marshal signing private key: %v", err)
	}

	seed := adminSeed{
		adminID:        newID(),
		assertionKeyID: newID(),
		signingKeyID:   newID(),
		adminPrivKey:   adminPrivKey,
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		fatalf("begin tx: %v", err)
	}
	exec := func(q string, args ...any) {
		if _, execErr := tx.Exec(ctx, q, args...); execErr != nil {
			_ = tx.Rollback(ctx)
			fatalf("exec %q: %v", q, execErr)
		}
	}
	exec(`INSERT INTO admins (id, name) VALUES ($1, $2)`, seed.adminID, "e2e-admin")
	exec(`INSERT INTO admin_assertion_keys (id, admin_id, algorithm, public_key) VALUES ($1, $2, $3, $4)`,
		seed.assertionKeyID, seed.adminID, "RS256", adminPubKeyPEM)
	exec(`INSERT INTO key_pairs (id, protected_resource_id, algorithm, public_key, private_key) VALUES ($1, $2, $3, $4, $5)`,
		seed.signingKeyID, newID(), "RS256", signingPubKeyPEM, signingPrivKeyPEM)
	if err = tx.Commit(ctx); err != nil {
		fatalf("commit seed tx: %v", err)
	}

	return seed
}

func mustAdminToken(seed adminSeed) string {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  seed.adminID,
		"sub":  seed.adminID,
		"aud":  appName,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Hour).Unix(),
		"jti":  newID(),
		"role": "admin",
	})
	token.Header["kid"] = seed.assertionKeyID

	assertion, err := token.SignedString(seed.adminPrivKey)
	if err != nil {
		fatalf("sign admin assertion: %v", err)
	}

	grantType := "urn:ietf:params:oauth:grant-type:jwt-bearer"
	body := "grant_type=" + strings.ReplaceAll(grantType, ":", "%3A") + "&assertion=" + assertion
	resp, err := http.Post(testServer.URL+"/v1/backoffice/admins/token", "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		fatalf("call token endpoint: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		fatalf("token endpoint %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fatalf("decode token response: %v", err)
	}
	return result.AccessToken
}

func buildEngine(pool *pgxpool.Pool, signingKeyID string) *gin.Engine {
	workerPool := workers.NewPool(10)
	evtBus := eventInfra.NewInMemoryEventBus(shared_logger.DummyLogger{}, workerPool)
	cmdBus := commandInfra.NewCommandBus()
	qryBus := queryInfra.NewQueryBus()

	modules.RegisterBackoffice(pool, cmdBus, qryBus, evtBus)
	modules.RegisterSecurity(pool, cmdBus, qryBus, evtBus, appName)
	modules.RegisterProtocol(pool, cmdBus, qryBus, evtBus, signingKeyID, appName)

	flowProvider := flows.NewProvider(cmdBus)
	hm := middlewares.NewHandlerMiddleware()
	presenter := error_presenters.NewHttpErrorPresenter()

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(io.Discard))

	v1 := router.Group("/v1")
	v1.Use(middlewares.RequestID(), middlewares.HttpError(presenter))

	route_providers.NewBackoffice(hm, qryBus, cmdBus, flowProvider, signingKeyID, appName).Register(v1)
	route_providers.NewProtocol(protocolHandler.NewHandler(cmdBus, qryBus)).Register(v1)

	return router
}

func cleanup(ctx context.Context, pool *pgxpool.Pool, adminID, signingKeyID string) {
	exec := func(q string, args ...any) {
		if _, err := pool.Exec(ctx, q, args...); err != nil {
			fmt.Fprintf(os.Stderr, "cleanup warning: %v\n", err)
		}
	}
	exec(`DELETE FROM client_assertion_keys WHERE client_id IN (SELECT id FROM clients WHERE admin_id = $1)`, adminID)
	exec(`DELETE FROM client_secrets       WHERE client_id IN (SELECT id FROM clients WHERE admin_id = $1)`, adminID)
	exec(`DELETE FROM scopes               WHERE admin_id = $1`, adminID)
	exec(`DELETE FROM key_pairs            WHERE protected_resource_id IN (SELECT id FROM protected_resources WHERE admin_id = $1)`, adminID)
	exec(`DELETE FROM key_pairs            WHERE id = $1`, signingKeyID)
	exec(`DELETE FROM clients              WHERE admin_id = $1`, adminID)
	exec(`DELETE FROM protected_resources  WHERE admin_id = $1`, adminID)
	exec(`DELETE FROM tenants              WHERE admin_id = $1`, adminID)
	exec(`DELETE FROM admin_assertion_keys WHERE admin_id = $1`, adminID)
	exec(`DELETE FROM admins               WHERE id = $1`, adminID)
}

func do(method, path string, body any, opts ...func(*http.Request)) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("marshal body: %v", err))
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		panic(fmt.Sprintf("new request: %v", err))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, opt := range opts {
		opt(req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("do request: %v", err))
	}
	return resp
}

func withBearer(token string) func(*http.Request) {
	return func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+token) }
}

func withHeader(key, value string) func(*http.Request) {
	return func(r *http.Request) { r.Header.Set(key, value) }
}

func drain(resp *http.Response) {
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
}

func decodeBody(resp *http.Response, target any) {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		panic(fmt.Sprintf("decode response body: %v", err))
	}
}

func newID() string { return uuid.New().String() }

func generateRSAKeyPair() (*rsa.PrivateKey, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, "", err
	}
	return priv, string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})), nil
}

func marshalPrivateKey(priv *rsa.PrivateKey) (string, error) {
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "e2e setup: "+format+"\n", args...)
	os.Exit(1)
}
