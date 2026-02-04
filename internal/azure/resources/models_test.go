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
		activeMap map[string]interface{}
		want      *EligibleRoleAssignments
	}{
		{
			name:      "returns correct EligibleRoleAssignments struct",
			inputFile: "testdata/eligibleschedulesingle.json",
			activeMap: map[string]interface{}{},
			want: &EligibleRoleAssignments{
				AssignmentID:     "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57",
				RoleDefinitionID: "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
				RoleName:         "Owner",
				ResourceScope:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				Active:           false,
			},
		},
		{
			name:      "returns correct EligibleRoleAssignments struct when active",
			inputFile: "testdata/eligibleschedulesingle.json",
			activeMap: map[string]interface{}{
				"/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57": struct{}{},
			},
			want: &EligibleRoleAssignments{
				AssignmentID:     "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/627f4610-5e8a-4f61-a773-dca39e03ff57",
				RoleDefinitionID: "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
				RoleName:         "Owner",
				ResourceScope:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				Active:           true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleEligibilityData, err := os.ReadFile(tests[0].inputFile)
			if err != nil {
				assert.NoError(t, err, "failed to read test data file: %v")
			}
			var response struct {
				Value []*armauthorization.RoleEligibilityScheduleInstance `json:"value"`
			}
			json.Unmarshal(roleEligibilityData, &response)
			got := NewEligibleRoleAssignments(ctx, response.Value[0], tt.activeMap)
			assert.Equal(t, tt.want, got)
		})
	}
}
