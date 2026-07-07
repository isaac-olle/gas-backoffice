# golang-migrate

Migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## Format

`{timestamp}_{description}.up.sql`

## Run

```sh
migrate -path deploy/golang-migrate -database "$POSTGRES_DSN" up
```
