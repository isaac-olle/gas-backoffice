package seed_handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	idpCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/create"
	prCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeAddAllowedClientCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	scopeCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	subscriptionCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/create"
	createAdminCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create_user"
	clientSecretCreateCmd "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedSeed "github.com/isaac-olle/gas-backoffice/internal/shared/seed"
)

type seedRequest struct {
	Admin             seedAdmin              `json:"admin"              binding:"required"`
	Subscription      seedSubscription       `json:"subscription"       binding:"required"`
	ProtectedResource seedProtectedResource  `json:"protected_resource" binding:"required"`
	Scopes            []seedScope            `json:"scopes"              binding:"required"`
	Client            seedClient             `json:"client"             binding:"required"`
	IdentityProviders []seedIdentityProvider `json:"identity_providers"`
}

type seedAdmin struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type seedSubscription struct {
	Name string `json:"name" binding:"required"`
}

type seedProtectedResource struct {
	Id   string `json:"id"   binding:"required"`
	Name string `json:"name" binding:"required"`
}

type seedScope struct {
	Id   string `json:"id"   binding:"required"`
	Name string `json:"name" binding:"required"`
}

// seedClient is the backoffice frontend: a confidential client that uses
// authorization_code with PKCE and authenticates with a client_secret.
type seedClient struct {
	Id           string   `json:"id"            binding:"required"`
	Name         string   `json:"name"          binding:"required"`
	RedirectUris []string `json:"redirect_uris" binding:"required,min=1"`
	Secret       string   `json:"secret"        binding:"required"`
}

type seedIdentityProvider struct {
	Id                    string   `json:"id"                     binding:"required"`
	ReferenceName         string   `json:"reference_name"         binding:"required"`
	Name                  string   `json:"name"                   binding:"required"`
	Issuer                string   `json:"issuer"                 binding:"required"`
	ClientId              string   `json:"client_id"              binding:"required"`
	ClientSecret          string   `json:"client_secret"`
	AuthorizationEndpoint string   `json:"authorization_endpoint" binding:"required"`
	TokenEndpoint         string   `json:"token_endpoint"         binding:"required"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
	JwksUri               string   `json:"jwks_uri"               binding:"required"`
	Scopes                []string `json:"scopes"`
	BackofficeRedirectUri string   `json:"backoffice_redirect_uri"`
}

type seedResponse struct {
	AdminId string `json:"admin_id"`
}

type Handler struct {
	repo     *sharedSeed.Repository
	executor command.CommandExecutor
}

func NewHandler(repo *sharedSeed.Repository, executor command.CommandExecutor) *Handler {
	return &Handler{repo: repo, executor: executor}
}

func (h *Handler) Execute(c *gin.Context) {
	var req seedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isSeeded, err := h.repo.IsSeeded(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check seed state"})
		return
	}
	if isSeeded {
		c.JSON(http.StatusConflict, gin.H{"error": "system already seeded"})
		return
	}

	adminId := ids.NewUserId()
	subscriptionId := ids.NewSubscriptionId()
	commands, appErr := h.buildCommands(adminId, subscriptionId, req)
	if appErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	// The seed actor is a system administrator. System administrators are
	// global and therefore do not receive a subscription membership.
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build seed author"})
		return
	}
	ctx := vo.WithAuthor(c.Request.Context(), author)
	ctx = ids.WithSubscriptionId(ctx, subscriptionId)

	if err = h.executor.MultipleExecution(ctx, commands...); err != nil {
		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "seed failed"})
		return
	}

	if err = h.repo.MarkSeeded(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark seed as done"})
		return
	}

	c.JSON(http.StatusCreated, seedResponse{AdminId: adminId.String()})
}

func (h *Handler) buildCommands(adminId ids.UserId, subscriptionId ids.SubscriptionId, req seedRequest) ([]command.Command, error) {
	prId, err := ids.NewProtectedResourceIdFromString(req.ProtectedResource.Id)
	if err != nil {
		return nil, err
	}

	clientId, err := ids.NewClientIdFromString(req.Client.Id)
	if err != nil {
		return nil, err
	}

	prOwnership := prCreate.NewOwnership(subscriptionId)
	clientOwnership := clientCreate.NewOwnership(subscriptionId)

	clientSecretId := ids.NewClientSecretId()

	cmds := []command.Command{
		createAdminCmd.NewCommand(
			adminId,
			adminId,
			req.Admin.Name,
			vo.SystemAdministratorRole,
			true,
			req.Admin.Email,
		),
		subscriptionCreateCmd.NewCommand(subscriptionId, req.Subscription.Name, adminId, true),
		prCreateCmd.NewCommand(
			prId,
			req.ProtectedResource.Name,
			prOwnership,
			true,
		),
		clientCreateCmd.NewCommand(
			clientId,
			req.Client.Name,
			clientOwnership,
			[]string{"authorization_code"},
			req.Client.RedirectUris,
			false,
			true,
		),
		clientSecretCreateCmd.NewCommand(
			clientSecretId,
			clientId,
			req.Client.Secret,
		),
	}

	for _, scope := range req.Scopes {
		scopeId, err := ids.NewScopeIdFromString(scope.Id)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, scopeCreateCmd.NewCommand(
			scopeId,
			scope.Name,
			scopes.NewOwnershipFromVos(prId, subscriptionId),
			true,
		))

		// The seed client needs these scopes allowed to be able to log in via
		// the OIDC flow (it requests "openid email" from itself as a resource).
		cmds = append(cmds, scopeAddAllowedClientCmd.NewCommand(
			scopeId.String(),
			clientId.String(),
		))
	}

	for _, idp := range req.IdentityProviders {
		idpId, err := ids.NewIdentityProviderIdFromString(idp.Id)
		if err != nil {
			return nil, err
		}

		scopes := idp.Scopes
		if len(scopes) == 0 {
			scopes = []string{"openid", "email", "profile"}
		}

		cmds = append(cmds, idpCreateCmd.NewCommand(
			idpId,
			idp.ReferenceName,
			idp.Name,
			idp.Issuer,
			idp.ClientId,
			idp.ClientSecret,
			idp.AuthorizationEndpoint,
			idp.TokenEndpoint,
			idp.UserinfoEndpoint,
			idp.JwksUri,
			scopes,
			idp.BackofficeRedirectUri,
		))
	}

	return cmds, nil
}
