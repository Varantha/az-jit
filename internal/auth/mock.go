package auth

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// MockCred is an azcore.TokenCredential implementation for tests and the dev
// binary. Set Tok to control the access-token string returned, or Err to
// simulate a credential failure (e.g. az CLI not logged in).
type MockCred struct {
	Tok string
	Err error
}

// GetToken satisfies azcore.TokenCredential. If Err is non-nil it is returned
// unwrapped; otherwise Tok is handed back with a one-hour expiry.
func (m MockCred) GetToken(ctx context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if m.Err != nil {
		return azcore.AccessToken{}, m.Err
	}
	return azcore.AccessToken{Token: m.Tok, ExpiresOn: time.Now().Add(time.Hour)}, nil
}
