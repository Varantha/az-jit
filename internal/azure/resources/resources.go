package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/varantha/az-jit/internal/clients"
)

func ListEligibleAssignments(client *clients.OmniClient) {

}

func GetRoleDefinitionName(ctx context.Context, roleDefinitionID string, client *clients.OmniClient) (string, error) {
	resp, err := client.RoleDefinitionsClient.GetByID(ctx, roleDefinitionID, nil)
	if err != nil {
		return "(unknown)", fmt.Errorf("failed to resolve role definition %s: %w", roleDefinitionID, err)
	}

	if resp.Properties != nil && resp.Properties.RoleName != nil {
		return *resp.Properties.RoleName, nil
	}

	return "(unknown)", nil
}

func GetActiveLinkedRoleAssignmentIDs(ctx context.Context, client *clients.OmniClient) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var linkedRoleID string
	filter := "asTarget()"
	pager := client.RoleAssignmentClient.NewListForScopePager("/", &armauthorization.RoleAssignmentScheduleInstancesClientListForScopeOptions{
		Filter: &filter,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("error fetching next page: %v", err)
		}

		for _, inst := range page.Value {
			props := inst.Properties
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

type pagerReturn struct {
	NextLink *string
	Value    []*any
}

type pager interface {
	More() bool
	NextPage(ctx context.Context) (pagerReturn, error)
}

func ExpandPage(p pager, ctx context.Context) ([]any, error) {
	results := []any{}
	for p.More() {
		page, err := p.NextPage(ctx)
		if err != nil {
			log.Fatalf("error fetching next page: %v", err)
			return nil, err
		}
		for _, val := range page.Value {
			results = append(results, val)
		}
	}
	return results, nil
}
