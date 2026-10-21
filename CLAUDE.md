# gas — Claude Code Guidelines

## Arquitectura general

Projecte Go seguint DDD minimalista, sense Temporal. Tres bounded contexts principals:
- **backoffice** — admins, tenants, protected_resources, endpoints, scopes, clients
- **security** — client_assertion_keys, admin_assertion_keys, key_pairs, client_secrets
- **protocol** — session (OAuth2/PKCE)

Flows cross-BC a `internal/flows/` (encara no implementat a `dev`; vegeu secció "Flows" més avall). Servidor HTTP amb Gin a `internal/server/`.

Estructura del servidor:
- `server/route_providers/` — un fitxer per "user" (`admins.go`, `users.go`, `clients.go`), tots al mateix package `route_providers`. Defineixen quins endpoints existeixen i apunten als handlers.
- `server/handlers/admins|users|clients/` — carpetes separades per "user", cada una amb un `handler.go` que conté la lògica HTTP dels endpoints.

---

## Convencions de codi

### JSON (respostes HTTP)
- Tots els camps de les respostes JSON van en `snake_case`, mai PascalCase. Cada `XxxPrimitives` (i qualsevol altre struct que es retorni directament via `c.JSON(...)`) ha de portar `json:"snake_case"` explícit a tots els camps exportats — `encoding/json` sense tags serialitza amb el nom Go tal qual (PascalCase), que no és el que volem exposar.
- Excepció: structs que representen un format extern/estàndard ja definit per una spec (resposta de token OAuth2 RFC 6749, JWKS RFC 7517) es construeixen amb `gin.H{...}` a mà amb els noms exactes de la spec, no amb un struct `Primitives` — no calen tags perquè no hi ha struct.
- Els `Primitives` que són purament interns (mai es passen a `c.JSON`, p. ex. `auth.SessionPrimitives`, `key_pairs.KeyPairPrimitives`, `protocol/session` primitives de persistència) no necessiten tags `json` per aquesta raó, encara que sí que els caldria si mai s'exposessin per algun endpoint.

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
- Projeccions de fills amb `*vo.SetVo[ids.XxxId]` **només quan l'agregat les muta síncronament ell mateix** (ex: `Scope.allowedClientIds`, mutat per `AddAllowedClient`/`RemoveAllowedClient`). **No duplicar relacions als dos costats** sense un benefici real ni cachejar al pare una llista deduïble del fill (ex: ja no hi ha `User.clientIds`/`protectedResourceIds`) — consulta-ho en viu amb `Repository.Matching(ctx, criteria.Criteria{...})` (patró ja usat a `clients/application/list`, `protected_resources/application/list`, `scopes/domain_services/duplicate_name_checker.go`).
  - **Excepcions conscients: `Client.allowedScopes` i `ProtectedResource.scopeIds`.** Totes dues són `*vo.SetVo[ids.ScopeId]` que existeixen de nou com a **cache de lectura desnormalitzada i mantinguda per events**; mai es muten directament des d'un command, només des dels seus event handlers (mètodes d'agregat `AddXxx`/`RemoveXxx` que no registren cap event: són una reacció, no una acció de domini primària).
  - **`Client.allowedScopes`** (columna `backoffice.clients.allowed_scopes`). Existeix de nou com a **cache de lectura desnormalitzada i mantinguda per events**: els event handlers `add_allowed_scope_on_scope_client_added.go` / `remove_allowed_scope_on_scope_client_removed.go` de `clients/application/_event_handlers/` la sincronitzen a partir de `ScopeClientAdded`/`ScopeClientRemoved`. `Scope.allowedClientIds` continua sent l'**única font de veritat d'escriptura** — sempre es muta via `scope.AddAllowedClient`/`RemoveAllowedClient`, mai `Client.allowedScopes` directament des d'un command — però la lectura de "quins scopes té aquest client" ara surt directament del `Client`, sense query/endpoint addicional. És el mateix tipus d'excepció conscient que l'aplanament d'`AdminId` (decisió 2026-10-12 més avall): aquí sí que hi ha un benefici real (estalvia un endpoint de llistat dedicat / un `JOIN` per a un flux del frontend), a diferència de l'antic anti-patró contra el qual advertia aquesta regla.
  - **`ProtectedResource.scopeIds`** (columna `backoffice.protected_resources.scope_ids`). Relació estructural pare→fill: cada `Scope` pertany a un sol PR via `Scope.ownership.ProtectedResourceId` (fixat a la creació, immutable), que és l'**única font de veritat**. Els event handlers `add_scope_id_on_scope_created.go` / `remove_scope_id_on_scope_deleted.go` de `protected_resources/application/_event_handlers/` la sincronitzen a partir de `ScopeCreated`/`ScopeDeleted`; la creació/esborrat del `Scope` és l'única manera que canviï aquesta llista. Benefici: el `GET` d'un PR mostra els seus scopes sense `JOIN` ni endpoint de llistat dedicat.
    - ⚠️ Forat conegut (preexistent, no resolt aquí): `scopes/application/soft_delete` no crida `scope.Delete()` ni publica cap event, així que un scope esborrat amb soft delete **continua apareixent** a `ProtectedResource.scopeIds`. Només `hard_delete` emet `ScopeDeleted`.

- `Admin` (bc `backoffice/admins`) i `Tenant` (bc `backoffice/tenants`) comparteixen la taula base `backoffice.users`. `Tenant` afegeix `tenant_profiles`; ambdós tenen `subscription_ids TEXT[]` com a projecció de lectura sincronitzada pels events de membership. `admin_id` continua sent un `ids.UserId` per a la relació Admin/Tenant, però no forma part de l'ownership dels recursos.

### Value Objects
- Un fitxer per VO dins `vos/`.
- El VO valida al constructor i retorna `app_errors.BadRequestError(ErrInvalidXxx)` si falla.
- Exposa `.String()` com a mètode principal.

### Use cases
- Cada use case en el seu propi package: `application/create/`, `application/get/`, etc.
- `command.go` / `query.go` conté: tipus constant, struct Command/Query, NewXxx, Type(), CommandHandler/QueryHandler, NewXxxHandler, Handle().
- `service.go` conté la lògica de negoci.
- Actualment `clients`, `protected_resources`, `scopes`, `users` i `identity_providers` tenen tots un use case `create` (a `dev` ja no és cert que "la majoria no en tingui" — revisar aquest punt si torna a quedar desfasat).
- `Client`, `ProtectedResource`, `Scope`, `ClientAssertionKey` i `ClientSecret` tenen ownership per `SubscriptionId`. `Scope` conserva també `ProtectedResourceId` com a relació estructural.
- Els tokens del backoffice porten `subscription_ids`; el middleware els posa a `authors.UserAuthor`. Els serveis comproven que la subscripció del recurs hi sigui. Els tokens del BC `protocol` no inclouen aquest claim.
- **On va cada comprovació d'ownership (servei vs. handler):** si el check es resol amb l'agregat que el servei ja carrega del seu propi mòdul (ex: `scope.EnsureOwnership(...)` després d'un `scopeFinder.Find`), es queda al **servei**. Si el check necessita consultar un agregat d'un **altre mòdul** (un altre `Finder`/`Repository`/domini), va al **handler**, com una query separada abans d'emetre el command — mai important el domini d'un altre mòdul dins del servei. Prova ràpida: si per fer el check el `service.go` hauria d'importar el paquet de domini d'un altre mòdul, el check no toca allà — puja'l al handler. El handler pot fer "dues coses" (dispatch query de check + dispatch command) sense que sigui un smell: la seva responsabilitat és precisament orquestrar, i és l'únic lloc que travessa mòduls de manera natural. Un cop existeixi `check_ownership`, el handler hi cridarà en lloc de reciclar una query `get` només per aprofitar-ne l'`EnsureOwnership`.
- L'ownership de cada agregat es gestiona en dos fitxers al mateix nivell que l'aggregate root:
  - `ownership.go`: `Ownership` (VO intern), `newOwnership(typed ids...)` (intern), `NewOwnership(raw strings...) (Ownership, error)` (públic, parseig intern). Per als recursos subscrits, `EnsureOwnership` compara l'`SubscriptionId` del recurs amb el conjunt de subscripcions del caller.
  - `primitives.go`: `OwnershipPrimitives` (struct públic amb strings), `func (o Ownership) ToPrimitives() OwnershipPrimitives`.
- La Postgres guarda els camps d'ownership flattened al mateix nivell que `id`, `name`, etc.

### Domain services (Finders)
- Els domain services de tipus Finder (`domain_services/xxx_finder.go`) encapsulen la lògica de cerca + validació de nil: reben el repositori, fan el `FindById`, comproven si és nil i retornen el error tipat corresponent.
- Els use cases **injecten el domain service Finder**, no el repositori directament. Així el service no repeteix la comprovació de nil ni els errors de "not found".
- Exemple: `clientFinder *domain_services.ClientFinder` al service, i `client, err := s.clientFinder.Find(ctx, id)` — si no existeix ja retorna `NotFoundError` sense codi addicional al service.

### Event handlers
- A `application/_event_handlers/`, un fitxer per handler.
- El handler té un mètode `Handle(ctx, TypedEvent) error`.
- Es registra al bootstrap amb `event.NewInMemorySubscriber[EventType](handler.Handle)`.
- Si un event handler ha de mutar un altre agregat com a reacció, ho fa pel cicle de domini estàndard: **Finder → mètode de mutació de l'agregat → `Repository.Save()`**, confiant en el check OCC de `version` del `Save` (mapeja `persistance.ErrConcurrentModification` a `app_errors.ConflictError(ErrXxxConcurrentModification)`). **Mai** un `UPDATE`/`array_append` SQL atòmic cru que salti el `Save()`. Això substitueix un patró anterior erroni (`AppendAssertionKeyId`/`AppendSecretId` al repositori de `clients`) que s'arriscava a *lost updates* en silenci perquè no participava gens en l'OCC (no incrementava `version` i el `Save` normal ni tan sols actualitzava aquelles columnes). No hi ha reintent: si hi ha conflicte, el handler retorna error i l'`InMemoryEventBus` el registra (limitació acceptada).

### Flows
- **Encara no implementat a `dev`** — no existeix `internal/flows/`. Disseny previst: a `internal/flows/create_xxx/`, un fitxer principal amb `Flow`, `Input`, `Output` i `Execute`, orquestrant múltiples BCs via `command.Bus`.

---

## Decisions de disseny rellevants

### 2026-05-14 — KrakenD com a API Gateway i authorization — ⚠️ SUPERADA, no vigent
> Aquesta decisió ha quedat desfasada; no la citis com a disseny actual. Es manté aquí només com a referència històrica fins que es documenti el disseny vigent.

KrakenD actua com a gateway davant del servidor Gin. El flux d'autorització és:
1. El client envia una request a KrakenD amb les seves credencials.
2. KrakenD crida al nostre endpoint `POST /token` (al route provider de `clients`) per validar i obtenir els claims.
3. KrakenD injecta els claims com a headers HTTP a les requests downstream (ex: `X-Client-Id`, `X-User-Id`, `X-Admin-Id`).
4. El backend confia cegament en aquests headers (KrakenD ja ha fet la validació).

El endpoint `POST /oidc/token` s'implementarà quan s'abordi el BC `protocol` (route provider de `clients`). Quan estigui llest, cal afegir a cada endpoint protegit del `krakend.json` el bloc `extra_config["auth/validator"]` amb la URL del JWKS i el mapeig de claims a headers. Els handlers llegiran la identitat dels headers injectats per KrakenD via un middleware Gin que la posi al `context.Context`.

### 2026-05-14 — PublicKey VO a assertion keys
`ClientAssertionKey` i `AdminAssertionKey` tenen un camp `publicKey vos.PublicKey`. El VO valida que sigui un bloc PEM vàlid (`pem.Decode`). La clau pública arriba des de fora (el client la genera i la passa al flow `create_client`).

**Motivació:** Les assertion keys serveixen per verificar JWTs signats pel client/admin; cal guardar la clau pública per fer la verificació sense dependre del key pair intern.

### 2026-10-08 — Concurrència optimista (OCC) només on hi ha un cicle carrega-muta-desa real
Per evitar *lost updates* (un `Save()` que sobreescriu amb una foto en memòria desactualitzada el que un altre escriptor ja havia guardat), els agregats `Client`, `Scope`, `ProtectedResource` i `User` (BC `backoffice`) tenen un camp `version int`:
- Columna `version INTEGER NOT NULL DEFAULT 1` a la taula corresponent.
- El `Save()` del repositori fa `INSERT ... ON CONFLICT (id) DO UPDATE SET ..., version = EXCLUDED.version + 1 WHERE <table>.version = EXCLUDED.version`, via `PostgresRepository[T].SaveWithResult` (mètode genèric, afegit sense tocar `Save`/`Delete` existents).
- Si `RowsAffected() == 0`, el repositori torna `persistance.ErrConcurrentModification` (sentinel cru, `internal/shared/persistance/errors.go`); el servei el mapeja a `app_errors.ConflictError(ErrXxxConcurrentModification)`.

**No s'aplica a tot arreu — decisió conscient, no un oblit:**
- **BC `security`** (`client_assertion_keys`, `key_pairs`, `client_secrets`): cap dels tres té cap cas d'ús que carregui un agregat existent i el torni a desar — només `create` (INSERT net, mai hi ha conflicte) i `hard_delete` (mai crida `Save`). L'`ON CONFLICT DO UPDATE` que ja hi ha als seus repositoris no es dispara amb cap flux actual. Si mai s'hi afegeix un cas d'ús que faci "carrega → muta → desa" (p. ex. rotar una clau), revisa aquest punt i aplica el mateix patró.
- **BC `protocol`, agregat `session`**: exclòs expressament. Una sessió OAuth2/PKCE no té cap flux previst amb múltiples escriptors concurrents sobre el mateix recurs (no hi ha processos en background que la toquin, i no és un recurs molt compartit); afegir-hi OCC seria protecció sense cap carrera real a evitar.

### 2026-10-08 — Joins SQL directes entre BCs: acceptats, decisió conscient
`protocol` llegeix taules de `backoffice` directament (`client_repository.go` fa `JOIN backoffice.scopes`/`backoffice.protected_resources`; `_shared/infrastructure/postgres_identity_provider_repository.go` fa `SELECT` contra `backoffice.identity_providers` i en desxifra el `client_secret` amb el mateix `Encrypter` que `backoffice`). `jwks` (BC `protocol`) fa un `JOIN` entre `security.key_pairs` i `backoffice.protected_resources`.

**Això és intencionat, no un oblit d'arquitectura:** el projecte és un monòlit amb una sola BD compartida; separar-ho en crides via command/query bus o en vistes SQL per evitar aquest acoblament tindria un cost de manteniment i disseny que no es justifica donat que no hi ha cap frontera de desplegament real entre BCs (mateix binari, mateix procés, mateixa BD). Si mai es desplega `protocol` o `security` per separat, **revisa aquest punt** — deixaria de ser vàlid.

### 2026-10-12 — `AdminId` aplanat a l'`Ownership` de `clients`, `protected_resources` i `scopes` — ⚠️ SUPERADA
> El model d'ownership per `AdminId` d'aquesta secció queda superat per la decisió del 2026-10-18 sobre `SubscriptionId`. Es conserva com a historial; no descriu el codi vigent.
Model de tenancy: un admin és només una fila de `backoffice.users` amb `role="admin"`, i cada `User` tenant té `adminId` apuntant al `UserId` del seu admin. Aquest `adminId` és la frontera de tenancy: un admin només ha de veure/gestionar recursos seus o dels tenants que li pertanyen, mai els d'un altre admin. Fins ara els `Ownership` de `clients`, `protected_resources` i `scopes` només portaven `UserId` (+ `ProtectedResourceId` a `scopes`) i es comparaven per igualtat exacta amb el `UserId` del caller, de manera que un admin no podia tocar cap recurs d'un tenant seu.

Ara els tres `Ownership` porten un camp `AdminId ids.UserId` (`clients.NewOwnership(adminId, userId)`, `protected_resources.NewOwnership(adminId, userId)` / `NewOwnershipFromStrings(adminId, userId)`, `scopes.NewOwnershipFromVos(prId, adminId, userId)` / `NewOwnership(prId, adminId, userId)`), i les taules `backoffice.clients`, `backoffice.protected_resources` i `backoffice.scopes` tenen la columna `admin_id TEXT NOT NULL` al costat de `user_id`. L'`AdminId` del caller arriba com a claim `admin_id` de l'access token (`UserProjection.AdminId` → `Session.adminId` → `BusTokenGenerator`), el middleware d'autenticació el posa a `authors.UserAuthor.AdminId`, i els handlers el llegeixen de l'author (`vo.AuthorFromContext[*authors.UserAuthor]`) — no hi ha cap context key propi per a l'`AdminId`.

**No es crea cap tipus `AdminId`:** un admin és un `User`, així que el seu id és un `ids.UserId` — mateix tipus, diferent slot semàntic (al VO, a l'author i als claims).

**Semàntica d'`EnsureOwnership`:** el caller (`o`) passa si és l'usuari propietari (`ownership.UserId == o.UserId`, cas self-service) **o** si és l'admin propietari (`ownership.AdminId == o.AdminId`), però aquesta segona branca només s'aplica quan el caller actua com a admin, és a dir quan `o.AdminId == o.UserId` (un admin és el seu propi admin; ho garanteixen `create_admin` i el seed). Sense aquesta condició, com que el claim `admin_id` d'un tenant és el seu admin, qualsevol tenant podria accedir als recursos de tots els tenants germans del mateix admin. A `scopes`, a més, sempre es requereix `ownership.ProtectedResourceId == o.ProtectedResourceId` (és l'id del pare estructural, no part de la identitat). Els `list` de `clients` i `protected_resources` filtren per `admin_id` quan el caller és admin i per `user_id` quan és tenant.

**És una excepció conscient a la regla de "no cachejar al fill una relació deduïble"**, i no el mateix anti-patró que els antics `ProtectedResource.scopeIds` i `Client.allowedScopes` originals (abans de reintroduir-los com a caches mantingudes per events — vegeu "Agregats"):
- Aquells caches no aportaven res: `get`/`list` havien de fer `JOIN` igualment, i només afegien una segona font de veritat a mantenir sincronitzada.
- Aquí, aplanar `admin_id` estalvia un `JOIN` contra `backoffice.users` a cada `list`/`get`/check d'ownership, i a més l'`AdminId` ha d'anar dins l'access token com a claim: les peticions downstream que només verifiquen el JWT no poden resoldre'l amb un `JOIN` en viu a cada crida.
- L'`AdminId` d'un recurs és immutable (es fixa a la creació i cap cas d'ús el muta), així que no hi ha cap sincronització a mantenir.

**Fora d'abast:** `users/ownership.go` (l'ownership propi de l'agregat `User`) **no canvia** — `User.adminId` ja existia com a camp de domini i `User.EnsureTenantAccess` ja implementa l'equivalent admin-vs-tenant per a aquest agregat. Aquesta decisió només tanca el mateix forat per a `clients`, `protected_resources` i `scopes` (i, per l'addenda de sota, `client_assertion_keys`/`client_secrets` de `security`). No toca el BC `protocol`.

**Addenda (2026-10-13; superada també el 2026-10-18) — aplanament antic a `security`:** la implementació vigent usa `SubscriptionId`; aquesta descripció queda aquí només com a historial.

### 2026-10-14 — Separació d'`Admin` i `Tenant` en dos agregats (taula base compartida)
Fins ara `backoffice/users` era un sol agregat `User` amb un camp `role` (`admin`/`tenant`) i un `adminId` que, per als admins, apuntava a ells mateixos. La separació lògica entre Admin (gestor: veu/gestiona els seus tenants, controla l'estat del sistema) i Tenant (consumidor: només usa el sistema per ell mateix) ja era prou forta, i mantenir-los junts hauria portat a columnes condicionals a mesura que cadascun necessités atributs propis.

Ara hi ha dos bounded contexts Go separats: `backoffice/admins` (nou, agregat `Admin`) i `backoffice/tenants` (renombrat de l'antic `backoffice/users`, agregat `Tenant`). Comparteixen la taula base `backoffice.users` (`id, role, name, email, is_system, version, created_at, updated_at, deleted_at`) com a "class-table inheritance": `Tenant` afegeix una taula d'extensió 1:1 `backoffice.tenant_profiles (user_id PK/FK, admin_id)` per al seu `adminId` (abans un camp propi de `User`, ara ve d'un `JOIN`); `Admin` no necessita cap extensió pròpia. No s'ha fet servir herència nativa de Postgres (`INHERITS`) — es descarta perquè no propaga `PRIMARY KEY`/`UNIQUE`/`FOREIGN KEY` del pare als fills i és un mecanisme pensat per partitioning, no per aquest patró.

Com que `admin_id` encara és un `ids.UserId` per a la relació entre Admin i Tenant, no s'ha creat un tipus `AdminId` separat. Aquesta decisió de tipus continua vigent per `authors.UserAuthor.AdminId`; els ownership dels recursos, però, han passat a `SubscriptionId`.

### 2026-10-18 — `Subscription` com a propietari dels recursos
`backoffice/subscriptions` és el bounded context que crea subscripcions i gestiona membres. `Subscription.memberIds` és la font de veritat; `backoffice.users.subscription_ids` és una cache de lectura que els handlers de `SubscriptionMemberAdded`/`SubscriptionMemberRemoved` actualitzen amb el cicle normal Finder → mutate → Save.

`Client`, `ProtectedResource`, `Scope`, `ClientAssertionKey` i `ClientSecret` guarden `subscription_id`. `Scope` manté a més `protected_resource_id` com a relació estructural. La pertinença al recurs es valida comprovant que el seu `subscription_id` és al conjunt `subscription_ids` del caller. El rol continua sent global al `User` i els permisos de rol s'apliquen a les rutes/use cases.

Només el token de backoffice inclou `subscription_ids`; el BC `protocol` queda intacte. La revocació de membership s'aplica quan el token es refresca o expira. Encara queda pendent definir com una petició de creació tria subscripció quan el caller en té més d'una; fins que es tanqui aquesta regla, la creació de clients i protected resources només accepta un caller amb una única subscripció accessible. El seed crea una subscripció `Gas` i n'afegeix l'Admin com a membre.

La migració pressuposa recrear la BDD de desenvolupament abans d'aplicar-la; no s'hi afegeix migració retroactiva de les files antigues. La facturació i el tracking d'ús per subscripció queden fora d'aquest canvi.

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
