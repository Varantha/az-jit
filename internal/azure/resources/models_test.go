package azure

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/stretchr/testify/assert"
)

func Test_NewEligibleRoleAssignments(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name      string
		inputFile string
		activeMap map[string]any
		want      *[]EligibleRoleAssignment
	}{
		{
			name:      "returns correct EligibleRoleAssignments struct",
			inputFile: "testdata/eligibleschedule_single.json",
			activeMap: map[string]any{},
			want: &[]EligibleRoleAssignment{{
				AssignmentID:      "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57",
				RoleDefinitionID:  "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
				RoleName:          "Owner",
				ResourceScope:     "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				ResourceScopeName: "Tenant Root Group",
				ResourceScopeType: "managementgroup",
				Active:            false,
			}},
		},
		{
			name:      "returns correct EligibleRoleAssignments struct when active",
			inputFile: "testdata/eligibleschedule_single.json",
			activeMap: map[string]any{
				"/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57": struct{}{},
			},
			want: &[]EligibleRoleAssignment{{
				AssignmentID:      "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57",
				RoleDefinitionID:  "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
				RoleName:          "Owner",
				ResourceScope:     "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				ResourceScopeName: "Tenant Root Group",
				ResourceScopeType: "managementgroup",
				Active:            true,
			}},
		},
		{
			name:      "returns correct EligibleRoleAssignments struct with multiples",
			inputFile: "testdata/eligibleschedule.json",
			activeMap: map[string]any{
				"/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57": struct{}{},
			},
			want: &[]EligibleRoleAssignment{
				{
					AssignmentID:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.KeyVault/vaults/kv-example/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/68165d74-d817-4777-9250-c97f4645d009",
					RoleDefinitionID:  "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b86a8fe4-44ce-4948-aee5-eccb2c155cd7",
					RoleName:          "Key Vault Secrets Officer",
					ResourceScope:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.KeyVault/vaults/kv-example",
					ResourceScopeName: "kv-example",
					ResourceScopeType: "Microsoft.KeyVault/vaults",
					Active:            false,
				},
				{
					AssignmentID:      "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/99777ae5-fa9e-91c4-4cf2-70ba3a432365",
					RoleDefinitionID:  "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b86a8fe4-44ce-4948-aee5-eccb2c155cd7",
					RoleName:          "Key Vault Secrets Officer",
					ResourceScope:     "/subscriptions/00000000-0000-0000-0000-000000000000",
					ResourceScopeName: "Subscription - Test",
					ResourceScopeType: "subscription",
					Active:            false,
				},
				{
					AssignmentID:      "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57",
					RoleDefinitionID:  "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
					RoleName:          "Owner",
					ResourceScope:     "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					ResourceScopeName: "Tenant Root Group",
					ResourceScopeType: "managementgroup",
					Active:            true,
				}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleEligibilityData, err := os.ReadFile(tt.inputFile)
			if err != nil {
				assert.NoError(t, err, "failed to read test data file: %v")
			}
			var response struct {
				Value []*armauthorization.RoleEligibilityScheduleInstance `json:"value"`
			}
			json.Unmarshal(roleEligibilityData, &response)
			got := NewEligibleRoleCollection(ctx, response.Value, tt.activeMap)
			assert.Equal(t, tt.want, got)
		})
	}
}
