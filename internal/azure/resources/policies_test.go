package azure

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/stretchr/testify/assert"
)

func Test_extractEndUserActivationPolicy(t *testing.T) {
	roleData, err := os.ReadFile("testdata/rolemanagementpolicies.json")
	if err != nil {
		t.Fatalf("failed to read test data file: %v", err)
	}

	valueTests := []struct {
		name        string
		expectError bool
		wantMaxDur  time.Duration
		wantReason  bool
		wantTicket  bool
	}{
		{name: "keyvault", expectError: false, wantMaxDur: 7 * time.Hour, wantReason: true, wantTicket: false},
		{name: "subscription", expectError: false, wantMaxDur: 7 * time.Hour, wantReason: true, wantTicket: false},
		{name: "mgmt-group", expectError: false, wantMaxDur: 7 * time.Hour, wantReason: true, wantTicket: true},
		{name: "empty-rules", expectError: true},
		{name: "no-enduser-rules", expectError: true},
		{name: "nil-maxduration", expectError: true},
	}

	for _, tt := range valueTests {
		t.Run(tt.name, func(t *testing.T) {
			var response struct {
				Value []armauthorization.RoleManagementPolicy `json:"value"`
			}
			json.Unmarshal(roleData, &response)

			policyFound := false
			for _, rolePolicy := range response.Value {
				if *rolePolicy.Name == tt.name {
					policyFound = true
					policy, err := extractEndUserActivationPolicy(rolePolicy.Properties.Rules)
					if tt.expectError {
						assert.Error(t, err)
						break
					}
					assert.NoError(t, err)
					assert.Equal(t, tt.wantMaxDur, policy.MaxDuration)
					assert.Equal(t, tt.wantReason, policy.RequireReason)
					assert.Equal(t, tt.wantTicket, policy.RequireTicket)
					break
				}
			}
			if !policyFound {
				t.Fatalf("policy %s not found in test data", tt.name)
			}
		})
	}
}
