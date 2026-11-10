package authenticate_client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func TestService_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		methods     func() []authenticate_client.Method
		wantErr     bool
		wantErrCode app_errors.ErrorCode
		wantResult  string
	}{
		{
			name: "delegates to the first method that can process",
			methods: func() []authenticate_client.Method {
				skipped := &stubMethod{canProcess: false}
				matched := &stubMethod{canProcess: true, result: authenticate_client.Result{ClientID: "client-b"}}
				return []authenticate_client.Method{skipped, matched}
			},
			wantResult: "client-b",
		},
		{
			name: "propagates the matched method's error",
			methods: func() []authenticate_client.Method {
				return []authenticate_client.Method{&stubMethod{canProcess: true, err: assert.AnError}}
			},
			wantErr: true,
		},
		{
			name: "returns ErrNoAuthenticationMethod when no method matches",
			methods: func() []authenticate_client.Method {
				return []authenticate_client.Method{&stubMethod{canProcess: false}}
			},
			wantErr:     true,
			wantErrCode: authenticate_client.ErrNoAuthenticationMethod,
		},
		{
			name: "returns ErrNoAuthenticationMethod with zero methods registered",
			methods: func() []authenticate_client.Method {
				return nil
			},
			wantErr:     true,
			wantErrCode: authenticate_client.ErrNoAuthenticationMethod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := authenticate_client.NewService(tt.methods()...)

			result, err := svc.Authenticate(context.Background(), authenticate_client.Credentials{})

			if tt.wantErr {
				if tt.wantErrCode != "" {
					var appErr *app_errors.AppError
					require.ErrorAs(t, err, &appErr)
					assert.Equal(t, tt.wantErrCode, appErr.Code)
					return
				}
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResult, result.ClientID)
		})
	}
}
