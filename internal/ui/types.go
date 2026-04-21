package ui

import "time"

// Domain is the source directory for a role.
type Domain string

const (
	DomainEntra  Domain = "entra"  // Microsoft Graph: directoryRoles / roleDefinitions
	DomainAzure  Domain = "azure"  // Azure RBAC: Microsoft.Authorization/roleDefinitions
	DomainGroups Domain = "groups" // Microsoft Graph: PIM for Groups
)

// Role is one eligible (or active) role assignment.
// Populate one of these per row returned from the three APIs:
//   - Entra: GET /roleManagement/directory/roleEligibilityScheduleInstances + roleAssignmentScheduleInstances
//   - Azure: GET {scope}/providers/Microsoft.Authorization/roleEligibilityScheduleInstances + roleAssignmentScheduleInstances
//   - Groups: GET /identityGovernance/privilegedAccess/group/eligibilityScheduleInstances + assignmentScheduleInstances
type Role struct {
	ID     string // stable identifier for selection tracking (the instance ID is fine)
	Name   string // "Global Administrator", "Key Vault Secrets Officer", etc.
	Domain Domain

	// Scope is split into a type and a value so we can render two columns.
	// Entra:  ScopeType="Directory" or "AdministrativeUnit", Scope=tenant display name or AU name
	// Azure:  ScopeType="Microsoft.KeyVault/vaults" | "subscription" | "managementgroup" | "resourcegroup"
	//         Scope=resource name (e.g. "kv-example")
	// Groups: ScopeType="Group", Scope=group display name
	ScopeType string
	Scope     string

	// Active is true when this role is currently elevated (came from
	// *AssignmentScheduleInstances rather than *EligibilityScheduleInstances).
	// Active roles are shown in amber and cannot be re-selected.
	Active   bool
	TimeLeft time.Duration // remaining on the active elevation; zero when !Active

	// MaxDuration is the policy cap for activation (from the role's
	// activation policy). Bulk defaults get clamped to this per row.
	MaxDuration time.Duration

	// Requirements, per the role's activation policy (PIM policy rules).
	// RequireTicket: ticketing system info required.
	// RequireReason: justification required (nearly always true).
	RequireReason bool
	RequireTicket bool
}

// Activation is a role selected for activation with its per-row config.
// Empty string / zero time means "inherit the bulk default".
// An explicit empty value (after the user cleared an inherited field) is
// not modelled here for simplicity — if you want that distinction, switch
// Reason/Ticket to *string and Duration to *time.Duration.
type Activation struct {
	Role     *Role
	Reason   string        // empty = inherit default
	Ticket   string        // empty = inherit default
	Duration time.Duration // zero = inherit default (and clamp to Role.MaxDuration)
}
