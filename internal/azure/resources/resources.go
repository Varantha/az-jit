package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/varantha/az-jit/internal/clients"
)

func listEligibleAssignments(ctx context.Context, client *clients.OmniClient) ([]*armauthorization.RoleEligibilityScheduleInstance, error) {
	var allEligibleRoles []*armauthorization.RoleEligibilityScheduleInstance
	pager := client.RoleEligibilityClient.NewListForScopePager("/", nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("error fetching next page: %v", err)
			return nil, err
		}

		allEligibleRoles = append(allEligibleRoles, page.Value...)
	}

	return allEligibleRoles, nil
}

func getActiveLinkedRoleAssignmentIDs(ctx context.Context, client *clients.OmniClient) (map[string]any, error) {
	output := make(map[string]any)
	var linkedRoleID string
	filter := "asTarget()"
	pager := client.RoleAssignmentClient.NewListForScopePager("/", &armauthorization.RoleAssignmentScheduleInstancesClientListForScopeOptions{
		Filter: &filter,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("error fetching next page: %v", err)
			return nil, err
		}

		for _, roleAssignment := range page.Value {
			props := roleAssignment.Properties
			if props == nil || props.RoleDefinitionID == nil || props.Scope == nil {
				continue
			}

			if props.LinkedRoleEligibilityScheduleID != nil {
				linkedRoleID = *props.LinkedRoleEligibilityScheduleID
				output[linkedRoleID] = struct{}{}
			}

		}
	}
	return output, nil
}

func GetEligibleRoles(ctx context.Context, client *clients.OmniClient) ([]EligibleRoleAssignment, error) {
	eligibleRoles, err := listEligibleAssignments(ctx, client)
	if err != nil {
		return nil, err
	}

	activeLinkedRoleIDs, err := getActiveLinkedRoleAssignmentIDs(ctx, client)
	if err != nil {
		return nil, err
	}

	return *NewEligibleRoleCollection(ctx, eligibleRoles, activeLinkedRoleIDs), nil
}
