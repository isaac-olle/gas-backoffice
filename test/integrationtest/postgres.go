// Package integrationtest contains helpers shared by integration tests.
package integrationtest

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres returns a connection to a database migrated with deploy/golang-migrate.
// Integration tests are skipped unless TEST_POSTGRES_DSN points to a disposable test database.
func Postgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set TEST_POSTGRES_DSN to a migrated, disposable PostgreSQL test database")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	if err = pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Fatalf("connect to PostgreSQL test database: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// BackofficeSubscription creates and cleans up a user and subscription for repository tests.
func BackofficeSubscription(t *testing.T, pool *pgxpool.Pool) (ownerID, subscriptionID string) {
	t.Helper()
	ownerID = uuid.NewString()
	subscriptionID = uuid.NewString()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO backoffice.users (id,name,email) VALUES ($1,$2,$3)`, ownerID, "Integration Owner", "owner-"+ownerID+"@example.test"); err != nil {
		t.Fatalf("insert integration owner: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO backoffice.subscriptions (id,name,owner_id) VALUES ($1,$2,$3)`, subscriptionID, "Integration Subscription", ownerID); err != nil {
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.users WHERE id=$1`, ownerID)
		t.Fatalf("insert integration subscription: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.subscriptions WHERE id=$1`, subscriptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.users WHERE id=$1`, ownerID)
	})
	return ownerID, subscriptionID
}

// ProtectedResource creates a protected resource with its required subscription.
func ProtectedResource(t *testing.T, pool *pgxpool.Pool) (subscriptionID, protectedResourceID string) {
	t.Helper()
	_, subscriptionID = BackofficeSubscription(t, pool)
	protectedResourceID = uuid.NewString()
	if _, err := pool.Exec(context.Background(), `INSERT INTO backoffice.protected_resources (id,name,subscription_id) VALUES ($1,$2,$3)`, protectedResourceID, "Integration Resource", subscriptionID); err != nil {
		t.Fatalf("insert integration protected resource: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM backoffice.protected_resources WHERE id=$1`, protectedResourceID)
	})
	return subscriptionID, protectedResourceID
}
