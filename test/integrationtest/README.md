# Integration test database

PostgreSQL-backed integration tests use `TEST_POSTGRES_DSN`. They skip when it
is unset. Point it at a disposable database with the migrations in
`deploy/golang-migrate` already applied; tests create uniquely identified rows
and clean them up when they finish.

For example, after creating an empty test database and applying the migrations:

```sh
TEST_POSTGRES_DSN='postgres://user:password@localhost:5432/gas_test?sslmode=disable' go test ./internal/bounded_contexts/...
```
