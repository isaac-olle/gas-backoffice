package add_allowed_grant_type_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/require"
)

type grantTypeRepository struct { client *clients.Client; findErr, saveErr error; saved int }
func (r *grantTypeRepository) Save(context.Context, *clients.Client) error { r.saved++; return r.saveErr }
func (r *grantTypeRepository) SearchById(context.Context, ids.ClientId) (*clients.Client, error) { return r.client,r.findErr }
func (r *grantTypeRepository) Matching(context.Context, criteria.Criteria) ([]*clients.Client,error) { return nil,nil }
func (r *grantTypeRepository) SoftDelete(context.Context, ids.ClientId) error { return nil }
func (r *grantTypeRepository) HardDelete(context.Context, ids.ClientId) error { return nil }
type grantTypeBus struct { events int }
func (b *grantTypeBus) Publish(_ context.Context, e []event.Event) { b.events+=len(e) }

func TestAdd(t *testing.T) {
	clientID, subscriptionID, adminID := ids.NewClientId(), ids.NewSubscriptionId(), ids.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminID.String(),adminID.String(),vo.SystemAdministratorRole)
	require.NoError(t,err)
	regular, err := clients.NewClient(clientID,"example",clients.NewOwnership(subscriptionID),nil,nil,false,false)
	require.NoError(t,err)
	regular.PullDomainEvents()
	system, err := clients.NewClient(clientID,"system",clients.NewOwnership(subscriptionID),nil,nil,false,true)
	require.NoError(t,err)
	system.PullDomainEvents()
	unauthorized, err := vo.NewUserAuthorFromString(ids.NewUserId().String(),ids.NewUserId().String(),vo.SubscriptionUserRole)
	require.NoError(t,err)
	tests:=[]struct{name,grantType string; client *clients.Client; findErr,saveErr error; author *vo.UserAuthor; wantCode app_errors.ErrorCode; wantEvents int}{
		{name:"adds a valid grant type",grantType:"client_credentials",client:regular,author:author,wantEvents:1},
		{name:"rejects unsupported grant types",grantType:"implicit",client:regular,author:author,wantCode:clients.ErrInvalidOwnership},
		{name:"returns not found for missing clients",grantType:"client_credentials",author:author,wantCode:clients.ErrClientNotFound},
		{name:"wraps repository search errors",grantType:"client_credentials",findErr:errors.New("query failed"),author:author,wantCode:clients.ErrClientNotFound},
		{name:"rejects authors without subscription access",grantType:"client_credentials",client:regular,author:unauthorized,wantCode:clients.ErrInvalidOwnership},
		{name:"rejects system clients",grantType:"client_credentials",client:system,author:author,wantCode:clients.ErrSystemResource},
		{name:"wraps repository save errors",grantType:"client_credentials",client:regular,author:author,saveErr:errors.New("write failed"),wantCode:clients.ErrClientNotSaved},
		{name:"maps concurrent modification errors",grantType:"client_credentials",client:regular,author:author,saveErr:persistance.ErrConcurrentModification,wantCode:clients.ErrClientConcurrentModification},
	}
	for _,tt:=range tests {t.Run(tt.name,func(t *testing.T){
		repo,bus:=&grantTypeRepository{client:tt.client,findErr:tt.findErr,saveErr:tt.saveErr},&grantTypeBus{}
		cmd:=add_allowed_grant_type.NewCommand(clientID,tt.grantType,tt.author)
		err:=add_allowed_grant_type.NewService(domain_services.NewClientFinder(repo),repo,bus,sharedDomainServices.NewSubscriptionAuthorizer()).Add(context.Background(),cmd)
		if tt.wantCode!="" { var appErr *app_errors.AppError; require.ErrorAs(t,err,&appErr); require.Equal(t,tt.wantCode,appErr.Code) } else { require.NoError(t,err) }
		require.Equal(t,tt.wantEvents,bus.events)
		if tt.wantCode=="" {require.Equal(t,1,repo.saved)}
	})}
}
