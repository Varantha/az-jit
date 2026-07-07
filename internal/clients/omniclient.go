package clients

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type OmniClient struct {
	CliCredential                       azcore.TokenCredential
	RoleEligibilityClient               *armauthorization.RoleEligibilityScheduleInstancesClient
	RoleAssignmentClient                *armauthorization.RoleAssignmentScheduleInstancesClient
	RoleAssignmentScheduleRequestClient *armauthorization.RoleAssignmentScheduleRequestsClient
}

type OmniClientOptions struct {
	ARM *arm.ClientOptions
}

func NewOmniClient(cred azcore.TokenCredential, opt OmniClientOptions) (*OmniClient, error) {
	instanceClient, err := armauthorization.NewRoleEligibilityScheduleInstancesClient(cred, opt.ARM)
	if err != nil {
		return nil, fmt.Errorf("failed to create eligibility client: %w", err)
	}

	assignmentClient, err := armauthorization.NewRoleAssignmentScheduleInstancesClient(cred, opt.ARM)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment client: %w", err)
	}

	assignmentScheduleRequestClient, err := armauthorization.NewRoleAssignmentScheduleRequestsClient(cred, opt.ARM)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment schedule request client: %w", err)
	}

	return &OmniClient{
		CliCredential:                       cred,
		RoleEligibilityClient:               instanceClient,
		RoleAssignmentClient:                assignmentClient,
		RoleAssignmentScheduleRequestClient: assignmentScheduleRequestClient,
	}, nil
}
