---
name: post-task-tests
description: >
  After completing any task the user requested — adding an endpoint, a use case, an aggregate,
  a middleware, a VO, or any code change — automatically review and update unit tests.
  Always invoke this skill at the end of every task, without the user having to ask.
  This includes: new handlers, new use case services, new middlewares, new aggregates,
  changed service return types, changed middleware behavior, renamed files, and refactors.
---

# Post-task test review

After completing the main task, always do this pass before reporting the task as done.

## What to check

### 1. Locate affected test files

For each file you created or modified, find its corresponding test file:

- `application/<use_case>/service.go` → `application/<use_case>/service_test.go` + `mother_test.go`
- `server/middlewares/<middleware>.go` → `server/middlewares/<middleware>_test.go`
- `server/middlewares/<subfolder>/<middleware>.go` → `server/middlewares/<subfolder>/<middleware>_test.go`
- `domain/vos/<vo>.go` → `domain/vos/<vo>_test.go` (or `shared/domain/vo/<vo>_test.go`)
- Aggregate root (`aggregate_root.go`) → `<aggregate>_test.go` at the same level
- Event handlers (`_event_handlers/<handler>.go`) → `_event_handlers/<handler>_test.go`

If a test file doesn't exist yet for something you created, create it.

### 2. Cases that always need new or updated tests

| What changed | What to test |
|---|---|
| New HTTP handler (`server/handlers/<pkg>/<handler>.go`) | Table-driven `TestXxxHandler`: valid body → no error, missing/invalid body → error, bus/service error → error, `Success` no-op or correct status, `Default` → 500. |
| New use case service | Happy path + not-found/error paths. Use `mother_test.go` for fixture builders. |
| Changed service return type | Update type assertions in existing `service_test.go`. |
| New middleware | Table-driven test with: missing context → abort, valid context → next, invalid value → abort. |
| Changed middleware behavior | Update or add cases to cover the new behavior. |
| New aggregate constructor | Valid input passes, each invalid field fails with the right error code. |
| New VO | Valid input, each invalid input returns the expected error. |
| New event handler | Mock the command bus, verify the right command is dispatched. |
| Renamed file/symbol | Update any test that references the old name. |

### 3. Run linter and tests to confirm everything passes

```bash
golangci-lint run ./...
go test ./...
```

Run the linter **before** running tests. Fix any linter errors first (unused imports, style violations, line-length, etc.), then confirm all tests pass.

Fix any failures before reporting the task as done. If a pre-existing test breaks because of your change, that's a sign the test was testing the wrong thing — update it to match the new intended behavior.

> Note: the project uses golangci-lint **v2** (`version: "2"` in `.golangci.yml`). If the installed binary is v1 it will fail with a config version error — install v2 with:
> ```bash
> go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
> ```

## Test conventions in this project

- **Package**: always `<package>_test` (external test package).
- **Parallelism**: `t.Parallel()` on every test and subtest.
- **Table-driven**: use `[]struct{...}` for cases that share the same setup/assert shape.
- **Fixtures**: test data builders go in `mother_test.go` in the same package, using `FromPrimitives` (never `New<Aggregate>`, which fires domain events).
- **Mocks**: use the generated mockery mock (`<aggregate>/mocks/MockRepository`). Use `repo.EXPECT()` for expectations, `mock.Anything` for context args.
- **Assertions**: `require.NoError` before inspecting results; `assert.*` for individual fields.
- **Service constructor helper**: define a `newService(t *testing.T)` at the top of `service_test.go` that wires up the mock repo and returns the service + mock.
- **Middleware tests**: use `gin.New()` + `httptest.NewRecorder()`. Pre-seed context values with a setup handler placed before the middleware under test.

## What NOT to test here

- Infrastructure (PostgreSQL repositories) — those are integration tests.
- End-to-end HTTP flows — those live in `test/e2e/`.
- The Gin router wiring itself — route providers have no unit tests.
- `db_seed` or CLI commands.

## When there is nothing to test

If the change is purely config, route wiring, constant definitions, or a trivial rename with no logic, state so explicitly: "No unit tests needed — change is [reason]."
