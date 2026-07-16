package vo

const (
	// users
	ScopeUsersRead    = "users:read"
	ScopeUsersReadAll = "users:read_all"
	ScopeUsersCreate  = "users:create"
	ScopeUsersRename  = "users:rename"
	ScopeUsersDelete  = "users:delete"

	// protected_resources
	ScopeProtectedResourcesCreate = "protected_resources:create"
	ScopeProtectedResourcesRead   = "protected_resources:read"
	ScopeProtectedResourcesRename = "protected_resources:rename"
	ScopeProtectedResourcesDelete = "protected_resources:delete"

	// scopes
	ScopeScopesCreate              = "scopes:create"
	ScopeScopesRead                = "scopes:read"
	ScopeScopesDelete              = "scopes:delete"
	ScopeScopesAddAllowedClient    = "scopes:add_allowed_client"
	ScopeScopesRemoveAllowedClient = "scopes:remove_allowed_client"

	// clients
	ScopeClientsCreate         = "clients:create"
	ScopeClientsRead           = "clients:read"
	ScopeClientsDelete         = "clients:delete"
	ScopeClientsAddFlow        = "clients:add_flow"
	ScopeClientsRemoveFlow     = "clients:remove_flow"
	ScopeClientsAddRedirectUri = "clients:add_redirect_uri"
)
