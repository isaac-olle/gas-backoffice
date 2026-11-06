# golang-migrate

Migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

The schema migrations are kept as a clean-install baseline. Recreate the
database before applying this migration set to an existing installation; the
removed incremental migrations are not an upgrade path for an existing schema.

## Format

`{timestamp}_{description}.up.sql`

## Run

```sh
migrate -path deploy/golang-migrate -database "$POSTGRES_DSN" up
```
