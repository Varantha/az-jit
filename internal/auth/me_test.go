package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockClaims = &userClaims{
	Name:     "Test User",
	Email:    "test@test.com",
	TenantID: "11111111-1111-1111-1111-111111111111",
	ObjectID: "22222222-2222-2222-2222-222222222222",
}

var errClaims = jwt.MapClaims{
	"oid": 123, // Invalid: should be string, not int
}

func createToken(claims jwt.Claims) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString([]byte("fake-signing-key"))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func TestGetIdentity(t *testing.T) {
	tests := []struct {
		name    string
		cred    azcore.TokenCredential
		want    AuthedUser
		wantErr bool
	}{
		{
			name: "valid token",
			cred: &MockCred{Tok: createToken(mockClaims)},
			want: *convertClaims(mockClaims),
		},
		{
			name:    "invalid token",
			cred:    &MockCred{Tok: createToken(errClaims)},
			wantErr: true,
		},
		{
			name:    "token error",
			cred:    &MockCred{Err: errors.New("blah")},
			wantErr: true,
		},
		{
			name:    "az cli not logged in",
			cred:    &MockCred{Err: context.Canceled},
			wantErr: true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIdentity(ctx, tt.cred)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, *got)
		})
	}
}
