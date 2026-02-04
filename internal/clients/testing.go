package clients

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type CredentialMock struct{}

func (c CredentialMock) GetToken(ctx context.Context, req policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token: "fake-token",
	}, nil
}
