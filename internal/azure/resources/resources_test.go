package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/varantha/az-jit/internal/clients"
)

func Test_listEligibleAssignments(t *testing.T) {
	defer gock.Off()

	route := "/providers/Microsoft.Authorization/roleEligibilityScheduleInstances"
	response, err := os.ReadFile("testdata/eligibleschedule_single.json")
	if err != nil {
		t.Fatalf("failed to read test data file: %v", err)
	}

	gock.New("https://management.azure.com").
		Get(route).
		Reply(200).
		BodyString(string(response))

	hc := &http.Client{}
	gock.InterceptClient(hc)

	armOptions := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: hc,
		},
	}

	OCOptions := clients.OmniClientOptions{
		ARM: armOptions,
	}

	cred := &clients.CredentialMock{}
	client, err := clients.NewOmniClient(cred, OCOptions)

	eligibleRoles, err := listEligibleAssignments(context.Background(), client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("Returns expected count", func(t *testing.T) {
		expectedCount := 1
		assert.Equal(t, expectedCount, len(eligibleRoles))
	})
	t.Run("Returns correct object", func(t *testing.T) {
		roleEligibilityData, err := os.ReadFile("testdata/eligibleschedule_single.json")
		if err != nil {
			assert.NoError(t, err, "failed to read test data file: %v")
		}
		var armResponse struct {
			Value []*armauthorization.RoleEligibilityScheduleInstance `json:"value"`
		}
		json.Unmarshal(roleEligibilityData, &armResponse)
		assert.Equal(t, armResponse.Value, eligibleRoles)
	})
}

func Test_getActiveLinkedRoleAssignmentIDs(t *testing.T) {
	defer gock.Off()

	route := "/providers/Microsoft.Authorization/roleAssignmentScheduleInstances"
	response, err := os.ReadFile("testdata/activeassignments_single.json")
	if err != nil {
		t.Fatalf("failed to read test data file: %v", err)
	}

	gock.New("https://management.azure.com").
		Get(route).
		Reply(200).
		BodyString(string(response))

	hc := &http.Client{}
	gock.InterceptClient(hc)

	armOptions := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: hc,
		},
	}

	OCOptions := clients.OmniClientOptions{
		ARM: armOptions,
	}

	cred := &clients.CredentialMock{}
	client, err := clients.NewOmniClient(cred, OCOptions)

	activeIDs, err := getActiveLinkedRoleAssignmentIDs(context.Background(), client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("Returns expected count", func(t *testing.T) {
		expectedCount := 1
		assert.Equal(t, expectedCount, len(activeIDs))
	})
	t.Run("Returns correct object", func(t *testing.T) {
		expectedMap := map[string]any{
			"/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilitySchedules/46760108-a1a0-4f53-a808-a32fbce153c7": struct{}{},
		}
		assert.Equal(t, expectedMap, activeIDs)
	})
}
