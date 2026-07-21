---
name: aggregate-root
description: >
  Generates or edits the aggregate root file (e.g. user.go, client.go) for a bounded context aggregate.
  Use this whenever the user asks to create or modify an aggregate root, add a mutation method,
  or asks "com faig el New/FromPrimitives?", "afegeix un mètode al agregat", "crea l'agregat Xxx".
---

# Aggregate Root

The aggregate root file lives at `internal/bounded_contexts/<bc>/<aggregate>/<aggregate_singular>.go`.

## Rules

- **No getters**. Fields are private. The rest of the system works with `XxxPrimitives` (read via `ToPrimitives()` in the service layer) or typed IDs passed in from outside.
- **`New<Xxx>`** — real creation constructor. Validates raw inputs, accumulates errors with `app_errors.NewEmptyErrorArray()`, fires a domain event via `a.Record(...)`.
- **`FromPrimitives`** — reconstruction from DB. Validates and builds from `XxxPrimitives`. Never calls `Record` — no events.
- **`CheckOwnership(o Ownership) bool`** — compares `a.ownership == o` (or `.Equals(o)` if Ownership has that method).
- **Mutation methods** (e.g. `ChangeName`, `Delete`, `AddClientId`) record a domain event after mutating state.
- **`IsSystemResource() bool`** — only if the aggregate has an `isSystem` field.

## File structure

```go
package <aggregate>

import (
    aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
    // ... other imports
)

var <aggregate>AggregateName aggregate2.Name = aggregate2.<AggregateName>

type <Xxx> struct {
    *aggregate2.Root
    id        ids.<Xxx>Id
    ownership Ownership
    // other fields — all private, no getters
}

func New<Xxx>(/* raw string params */, ownership Ownership) (*<Xxx>, error) {
    errArr := app_errors.NewEmptyErrorArray()

    // validate each raw param, accumulate errors
    xxxId, err := ids.New<Xxx>IdFromString(rawId)
    if err != nil { errArr.AddError(err) }

    // ...

    if errArr.HasErrors() {
        return nil, errArr
    }

    a := &<Xxx>{
        Root:      aggregate2.NewRoot(<aggregate>AggregateName),
        id:        xxxId,
        ownership: ownership,
        // ...
    }
    a.Record(<bc>events.New<Xxx>Created(a.id.String(), /* relevant fields */))
    return a, nil
}

func FromPrimitives(p <Xxx>Primitives) (*<Xxx>, error) {
    errArr := app_errors.NewEmptyErrorArray()

    id, err := ids.New<Xxx>IdFromString(p.Id)
    if err != nil { errArr.AddError(err) }

    ownership, err := NewOwnership(p.AdminId /* or relevant owner IDs */)
    if err != nil { errArr.AddError(err) }

    // reconstruct all other fields from primitives...

    if errArr.HasErrors() {
        return nil, errArr
    }

    return &<Xxx>{
        Root:      aggregate2.NewRoot(<aggregate>AggregateName),
        id:        id,
        ownership: ownership,
        // ...
    }, nil
    // NOTE: no Record() call here
}

func (a *<Xxx>) CheckOwnership(o Ownership) bool { return a.ownership == o }

// Mutation example — always records an event:
func (a *<Xxx>) ChangeName(newName vos.<Xxx>Name) {
    a.name = newName
    a.Record(<bc>events.New<Xxx>NameChanged(a.id.String(), a.name.String()))
}

func (a *<Xxx>) Delete() {
    a.Record(<bc>events.New<Xxx>Deleted(a.id.String(), a.ownership.AdminId.String()))
}
```

## Key conventions

- The `New<Xxx>` constructor receives **raw string params** for own fields and **already-typed VO IDs** for ownership (e.g. `adminId ids.AdminId`). See `NewUser(userId ids.UserId, name string, ownership Ownership)`.
- `FromPrimitives` receives an `XxxPrimitives` struct (all strings/primitives) and re-parses everything.
- `SetVo` projections (child IDs like `clientIds`, `protectedResourceIds`) are reconstructed in `FromPrimitives` by iterating the string slice and calling `NewXxxIdFromString`.
- The aggregate name constant must be registered in `internal/shared/types/` and re-exported in `internal/shared/domain/aggregate/name.go`.
- Event types live in `internal/shared/events/<bc>/`.
