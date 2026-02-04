package azure

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/varantha/az-jit/internal/clients"
)

func Test_ListEligibleAssignments(t *testing.T) {

}

func Test_GetRoleDefinitionName(t *testing.T) {

}

func Test_GetActiveLinkedRoleAssignmentIDs(t *testing.T) {
	defer gock.Off()

	route := "/providers/Microsoft.Authorization/roleAssignmentScheduleInstances"
	response, err := os.ReadFile("testdata/activeassignments.json")
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

	activeIDs, err := GetActiveLinkedRoleAssignmentIDs(context.Background(), client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedCount := 1
	assert.Equal(t, expectedCount, len(activeIDs))
}
