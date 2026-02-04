package clients

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type OmniClient struct {
	cliCredential         azcore.TokenCredential
	RoleEligibilityClient *armauthorization.RoleEligibilityScheduleInstancesClient
	RoleAssignmentClient  *armauthorization.RoleAssignmentScheduleInstancesClient
}

type OmniClientOptions struct {
	ARM *arm.ClientOptions
}

func NewOmniClient(token azcore.TokenCredential, opt OmniClientOptions) (*OmniClient, error) {
	instanceClient, err := armauthorization.NewRoleEligibilityScheduleInstancesClient(token, opt.ARM)
	if err != nil {
		log.Fatalf("failed to create instance client: %v", err)
	}

	assignmentClient, err := armauthorization.NewRoleAssignmentScheduleInstancesClient(token, opt.ARM)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment client: %w", err)
	}

	return &OmniClient{
		cliCredential:         token,
		RoleEligibilityClient: instanceClient,
		RoleAssignmentClient:  assignmentClient,
	}, nil
}
