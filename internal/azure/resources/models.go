package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/varantha/az-jit/internal/clients"
)

type EligibleRoleAssignments struct {
	AssignmentID     string
	RoleDefinitionID string
	RoleName         string
	ResourceScope    string
	Active           bool
}

func NewEligibleRoleAssignments(ctx context.Context, RoleEligibilityScheduleInstance *armauthorization.RoleEligibilityScheduleInstance, client *clients.OmniClient, activeLinkedRoleIDs map[string]interface{}) *EligibleRoleAssignments {

	roleName, err := GetRoleDefinitionName(ctx, *RoleEligibilityScheduleInstance.Properties.RoleDefinitionID, client)
	if err == nil {

	}
	_, isActive := activeLinkedRoleIDs[*RoleEligibilityScheduleInstance.ID]

	return &EligibleRoleAssignments{
		AssignmentID:     *RoleEligibilityScheduleInstance.ID,
		RoleDefinitionID: *RoleEligibilityScheduleInstance.Properties.RoleDefinitionID,
		RoleName:         roleName,
		ResourceScope:    *RoleEligibilityScheduleInstance.Properties.Scope,
		Active:           isActive,
	}
}
