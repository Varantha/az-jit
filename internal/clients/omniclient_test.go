package clients

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
)

type credentialMock struct{}

func (c credentialMock) GetToken(ctx context.Context, req policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token: "fake-token",
	}, nil
}

func Test_NewOmniClient(t *testing.T) {
	t.Run("returns no error", func(t *testing.T) {
		cred := &credentialMock{}
		_, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
	})

	t.Run("initializes RoleEligibilityClient", func(t *testing.T) {
		cred := &credentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleEligibilityClient)
	})

	t.Run("initializes RoleDefinitionsClient", func(t *testing.T) {
		cred := &credentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleDefinitionsClient)
	})

	t.Run("initializes RoleAssignmentClient", func(t *testing.T) {
		cred := &credentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleAssignmentClient)
	})
}
