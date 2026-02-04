package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type EligibleRoleAssignment struct {
	AssignmentID      string
	RoleDefinitionID  string
	RoleName          string
	ResourceScope     string
	ResourceScopeName string
	ResourceScopeType string
	Active            bool
}

func NewEligibleRoleCollection(ctx context.Context, EligibleRoles []*armauthorization.RoleEligibilityScheduleInstance, activeLinkedRoleIDs map[string]any) *[]EligibleRoleAssignment {
	output := []EligibleRoleAssignment{}
	if len(EligibleRoles) == 0 {
		return &output
	}

	for _, v := range EligibleRoles {
		roleName := v.Properties.ExpandedProperties.RoleDefinition.DisplayName
		_, isActive := activeLinkedRoleIDs[*v.ID]
		output = append(output, EligibleRoleAssignment{
			AssignmentID:      *v.ID,
			RoleDefinitionID:  *v.Properties.RoleDefinitionID,
			RoleName:          *roleName,
			ResourceScope:     *v.Properties.Scope,
			ResourceScopeName: *v.Properties.ExpandedProperties.Scope.DisplayName,
			ResourceScopeType: *v.Properties.ExpandedProperties.Scope.Type,
			Active:            isActive,
		})
	}

	return &output
}
