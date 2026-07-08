# gas — Claude Code Guidelines

## Arquitectura general

Projecte Go seguint DDD minimalista, sense Temporal. Tres bounded contexts principals:
- **backoffice** — admins, users, protected_resources, endpoints, scopes, clients
- **security** — client_assertion_keys, admin_assertion_keys, key_pairs, client_secrets
- **protocol** — session (OAuth2/PKCE)

Flows cross-BC a `internal/flows/`. Servidor HTTP amb Gin a `internal/server/`.

Estructura del servidor:
- `server/route_providers/` — un fitxer per "user" (`admins.go`, `users.go`, `clients.go`), tots al mateix package `route_providers`. Defineixen quins endpoints existeixen i apunten als handlers.
- `server/handlers/admins|users|clients/` — carpetes separades per "user", cada una amb un `handler.go` que conté la lògica HTTP dels endpoints.

---

## Convencions de codi

### Errors
- Usa sempre `app_errors.ErrorCode = "STRING_CONSTANT"` (no `errors.New(...)`).
- Els constructors d'agregat retornen errors via `app_errors.NewEmptyErrorArray()` — acumula tots els errors de validació i retorna al final.
- Mai retornis un `ErrorCode` directament com a `error`; embolcalla'l: `app_errors.BadRequestError(ErrXxx).WithDescription("...")`.
- **Infraestructura usa `fmt.Errorf`**, mai `app_errors`. Els errors d'infraestructura són crus; és la capa d'aplicació (service) qui els embolcalla amb `app_errors.InternalError(ErrXxx).WithCause(err)`.

### IDs
- Tots els IDs són structs tipats: `type AdminId struct{ vo.Uuid }`, mai type aliases.
- Constructors: `NewXxxId()` (genera UUID nou) i `NewXxxIdFromString(raw string) (XxxId, error)`.

### Agregats
- Cada agregat té: `aggregate_root.go`, `primitives.go`, `contracts.go`, `error_messages.go`, `infrastructure/postgresql.go`, i carpeta `vos/`.
- Cada agregat té dos constructors: `NewXxx(rawParams..., ownershipId ids.XxxId)` per a creació real (llença events de domini; rep camps propis en raw i IDs d'ownership ja com a VO tipats) i `FromPrimitives(p XxxPrimitives)` per a reconstrucció des de persistència (mai llença events).
- Sempre publica events després de cada mutació: `s.eventBus.Publish(ctx, aggregate.PullDomainEvents())`.
- Projeccions de fills amb `*vo.SetVo[ids.XxxId]`; relacions germanes (ex: client→scopes) també com a `SetVo` d'IDs, mai objectes embedats.

### Value Objects
- Un fitxer per VO dins `vos/`.
- El VO valida al constructor i retorna `app_errors.BadRequestError(ErrInvalidXxx)` si falla.
- Exposa `.String()` com a mètode principal.

### Use cases
- Cada use case en el seu propi package: `application/check_ownership/`, etc.
- `command.go` / `query.go` conté: tipus constant, struct Command/Query, NewXxx, Type(), CommandHandler/QueryHandler, NewXxxHandler, Handle().
- `service.go` conté la lògica de negoci.
- **No hi ha use case `create`** per a la majoria d'agregats del BC backoffice — es creen manualment a la BD. Excepció: `clients` té use cases `create` i `add_redirect_uri`.
- `Client` té `adminId` i `userId` com a camps d'ownership (tipats com a VO); `allowedScopes` com a `SetVo[ScopeId]`; `redirectUris` com a `SetVo[RedirectUri]`.
- **`check_ownership`** és una Query (no muta estat) que retorna `(any, error)` — en cas d'èxit retorna `OwnershipPrimitives`; en cas de fallada un `app_errors.ForbiddenError`. La cadena d'ownership és: Scope/Endpoint → ProtectedResource → User → Admin; Client → Admin + User.
- L'ownership de cada agregat es gestiona en dos fitxers al mateix nivell que l'aggregate root:
  - `ownership.go`: `Ownership` (VO intern, camps d'ID privats), `newOwnership(typed ids...)` (intern), `NewOwnership(raw strings...) (Ownership, error)` (públic, parseig intern). L'agregat té `CheckOwnership(Ownership) bool` que compara `a.ownership == o` directament.
  - `primitives.go`: `OwnershipPrimitives` (struct públic amb strings), `func (o Ownership) ToPrimitives() OwnershipPrimitives`. Els serveis `check_ownership` creen l'`Ownership` via `NewOwnership`, la passen a `aggregate.CheckOwnership(ownership)`, i retornen `ownership.ToPrimitives()` — l'agregat no exposa cap mètode de get ownership.
- La Postgres guarda els camps d'ownership flattened al mateix nivell que `id`, `name`, etc.

### Domain services (Finders)
- Els domain services de tipus Finder (`domain_services/xxx_finder.go`) encapsulen la lògica de cerca + validació de nil: reben el repositori, fan el `FindById`, comproven si és nil i retornen el error tipat corresponent.
- Els use cases **injecten el domain service Finder**, no el repositori directament. Així el service no repeteix la comprovació de nil ni els errors de "not found".
- Exemple: `clientFinder *domain_services.ClientFinder` al service, i `client, err := s.clientFinder.Find(ctx, id)` — si no existeix ja retorna `NotFoundError` sense codi addicional al service.

### Event handlers
- A `application/_event_handlers/`, un fitxer per handler.
- El handler té un mètode `Handle(ctx, TypedEvent) error`.
- Es registra al bootstrap amb `event.NewInMemorySubscriber[EventType](handler.Handle)`.

### Flows
- A `internal/flows/create_xxx/`, un fitxer principal amb `Flow`, `Input`, `Output` i `Execute`.
- Orquestren múltiples BCs via `command.Bus`.

---

## Decisions de disseny rellevants

### 2026-05-14 — KrakenD com a API Gateway i authorization
KrakenD actua com a gateway davant del servidor Gin. El flux d'autorització és:
1. El client envia una request a KrakenD amb les seves credencials.
2. KrakenD crida al nostre endpoint `POST /token` (al route provider de `clients`) per validar i obtenir els claims.
3. KrakenD injecta els claims com a headers HTTP a les requests downstream (ex: `X-Client-Id`, `X-User-Id`, `X-Admin-Id`).
4. El backend confia cegament en aquests headers (KrakenD ja ha fet la validació).

El endpoint `POST /v1/token` s'implementarà quan s'abordi el BC `protocol` (route provider de `clients`). Quan estigui llest, cal afegir a cada endpoint protegit del `krakend.json` el bloc `extra_config["auth/validator"]` amb la URL del JWKS i el mapeig de claims a headers. Els handlers llegiran la identitat dels headers injectats per KrakenD via un middleware Gin que la posi al `context.Context`.

### 2026-05-14 — PublicKey VO a assertion keys
`ClientAssertionKey` i `AdminAssertionKey` tenen un camp `publicKey vos.PublicKey`. El VO valida que sigui un bloc PEM vàlid (`pem.Decode`). La clau pública arriba des de fora (el client la genera i la passa al flow `create_client`).

**Motivació:** Les assertion keys serveixen per verificar JWTs signats pel client/admin; cal guardar la clau pública per fer la verificació sense dependre del key pair intern.

---

## Infraestructura

- **Postgres** — repositoris principals per a tots els agregats.
- **Mongo** — classe genèrica disponible a `shared/persistance/infrastructure/mongo.go`, però cap agregat l'usa de moment.
- **InMemoryEventBus** — worker pool de 10 goroutines; `Register(event.Subscriber)` per subscriure.
- **InMemoryCommandBus** — registre per `command.Type`.
- No hi ha Temporal, no hi ha suport BC, no hi ha projeccions MongoDB per agregats.

## Mòdul Go

```
github.com/isaac-olle/gas
```
