package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewOmniClient(t *testing.T) {
	t.Run("returns no error", func(t *testing.T) {
		cred := &CredentialMock{}
		_, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
	})

	t.Run("initializes RoleEligibilityClient", func(t *testing.T) {
		cred := &CredentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleEligibilityClient)
	})

	t.Run("initializes RoleDefinitionsClient", func(t *testing.T) {
		cred := &CredentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleDefinitionsClient)
	})

	t.Run("initializes RoleAssignmentClient", func(t *testing.T) {
		cred := &CredentialMock{}
		c, err := NewOmniClient(cred, OmniClientOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, c.RoleAssignmentClient)
	})
}
