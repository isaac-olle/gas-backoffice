---
name: add-endpoint
description: >
  Adds a new HTTP endpoint to the gas-backoffice project. Use whenever the user asks to add
  a new route, endpoint, or HTTP operation — e.g. "afegeix un endpoint per X", "vull un GET de Y",
  "fes la ruta per Z". Covers: handler, route registration, Postman collection, and db_seed scopes.
---

# add-endpoint — Checklist

When adding a new endpoint, always do ALL of these steps in order.

## 1. Handler

Create `internal/server/handlers/<aggregate>/my_action.go`:

```go
package <aggregate>

type MyActionHandler struct { ... }

func NewMyActionHandler(...) *MyActionHandler { ... }

func (h *MyActionHandler) Execute(c *gin.Context) (any, error) { ... }
func (h *MyActionHandler) Success(c *gin.Context, result any) error { ... }
func (h *MyActionHandler) Default(c *gin.Context) { ... }
```

## 2. Route provider

Register in the correct `internal/server/route_providers/backoffice_*.go`.

Always apply middlewares in this order:
1. `path_params.XxxIdParam()` — at group level (`.Use(...)`)
2. `ownership.Xxx(p.queryBus)` — per route, on mutating and read ops
3. `system_resource.Xxx(p.queryBus)` — only on DELETE and PATCH/PUT (mutations), never on GET or POST to sub-resources like allowed-clients

## 3. Postman collection

Add the request to `deploy/gas.postman_collection.json` under the correct role folder (Admin/Tenant/Client/Protocol) and aggregate subfolder.

Rules:
- Path params use `:param` syntax (NOT `{{param}}`), declared in the `variable` array of the url object
- Auth header: `"Bearer {{access_token}}"` for protected endpoints
- Add a test script asserting the expected status code
- If the endpoint returns a new resource ID, store it: `pm.environment.set('flow_xxx_id', body.id)`

Example request structure:
```json
{
  "name": "My Action",
  "request": {
    "method": "POST",
    "header": [
      {"key": "Authorization", "value": "Bearer {{access_token}}"},
      {"key": "Content-Type", "value": "application/json"}
    ],
    "url": {
      "raw": "{{base_url}}/v1/backoffice/admins/users/:user_id/resource/:resource_id/action",
      "host": ["{{base_url}}"],
      "path": ["v1", "backoffice", "admins", "users", ":user_id", "resource", ":resource_id", "action"],
      "variable": [
        {"key": "user_id", "value": "", "description": ""},
        {"key": "resource_id", "value": "", "description": ""}
      ]
    },
    "body": {"mode": "raw", "raw": "{}", "options": {"raw": {"language": "json"}}}
  },
  "event": [{"listen": "test", "script": {"type": "text/javascript", "exec": [
    "pm.test(\"Status 200/201\", () => pm.expect(pm.response.code).to.be.oneOf([200, 201]));"
  ]}}]
}
```

## 4. db_seed — add the scope

In `cmd/db_seed/main.go`, add:

1. A new scope ID variable (grouped with its aggregate):
```go
scopeMyResourceActionID := uuid.New().String()
```

2. Add it to `allScopeIDs` slice.

3. Add it to `seedScopes`:
```go
{scopeMyResourceActionID, "my_resource:action"},
```

Scope naming convention: `resource:action` (snake_case, colon separator).
Examples: `users:delete`, `scopes:add_allowed_client`, `clients:add_redirect_uri`.
