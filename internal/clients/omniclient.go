package clients

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type OmniClient struct {
	cliCredential         *azidentity.AzureCLICredential
	RoleEligibilityClient *armauthorization.RoleEligibilityScheduleInstancesClient
	RoleDefinitionsClient *armauthorization.RoleDefinitionsClient
	RoleAssignmentClient  *armauthorization.RoleAssignmentScheduleInstancesClient
}

func NewOmniClient() (*OmniClient, error) {
	cliCredential, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, err
	}

	instanceClient, err := armauthorization.NewRoleEligibilityScheduleInstancesClient(cliCredential, nil)
	if err != nil {
		log.Fatalf("failed to create instance client: %v", err)
	}

	roleDefClient, err := armauthorization.NewRoleDefinitionsClient(cliCredential, nil)
	if err != nil {
		log.Fatalf("failed to create role definitions client: %v", err)
	}

	assignmentClient, err := armauthorization.NewRoleAssignmentScheduleInstancesClient(cliCredential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment client: %w", err)
	}

	return &OmniClient{
		cliCredential:         cliCredential,
		RoleEligibilityClient: instanceClient,
		RoleDefinitionsClient: roleDefClient,
		RoleAssignmentClient:  assignmentClient,
	}, nil
}
