package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type EligibleRoleAssignments struct {
	AssignmentID     string
	RoleDefinitionID string
	RoleName         string
	ResourceScope    string
	Active           bool
}

func NewEligibleRoleAssignments(ctx context.Context, RoleEligibilityScheduleInstance *armauthorization.RoleEligibilityScheduleInstance, activeLinkedRoleIDs map[string]interface{}) *EligibleRoleAssignments {

	roleName := RoleEligibilityScheduleInstance.Properties.ExpandedProperties.RoleDefinition.DisplayName
	_, isActive := activeLinkedRoleIDs[*RoleEligibilityScheduleInstance.ID]

	return &EligibleRoleAssignments{
		AssignmentID:     *RoleEligibilityScheduleInstance.ID,
		RoleDefinitionID: *RoleEligibilityScheduleInstance.Properties.RoleDefinitionID,
		RoleName:         *roleName,
		ResourceScope:    *RoleEligibilityScheduleInstance.Properties.Scope,
		Active:           isActive,
	}
}
