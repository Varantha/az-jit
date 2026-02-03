package clients

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type OmniClient struct {
	cliCredential         azcore.TokenCredential
	RoleEligibilityClient *armauthorization.RoleEligibilityScheduleInstancesClient
	RoleDefinitionsClient *armauthorization.RoleDefinitionsClient
	RoleAssignmentClient  *armauthorization.RoleAssignmentScheduleInstancesClient
}

func NewOmniClient(token azcore.TokenCredential) (*OmniClient, error) {
	instanceClient, err := armauthorization.NewRoleEligibilityScheduleInstancesClient(token, nil)
	if err != nil {
		log.Fatalf("failed to create instance client: %v", err)
	}

	roleDefClient, err := armauthorization.NewRoleDefinitionsClient(token, nil)
	if err != nil {
		log.Fatalf("failed to create role definitions client: %v", err)
	}

	assignmentClient, err := armauthorization.NewRoleAssignmentScheduleInstancesClient(token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment client: %w", err)
	}

	return &OmniClient{
		cliCredential:         token,
		RoleEligibilityClient: instanceClient,
		RoleDefinitionsClient: roleDefClient,
		RoleAssignmentClient:  assignmentClient,
	}, nil
}
