package ui

import "time"

// This file holds placeholder data that matches the design mockups so the
// TUI renders without any network calls. Replace each function with a real
// Azure / Graph loader (see app.go Init and README for endpoints) and then
// drop this file entirely.

func mockUser() User {
	return User{
		Email:  "sam.smith@contoso.com",
		Region: "eu-west",
	}
}

func mockTenant() Tenant {
	return Tenant{
		Name: "Contoso Enterprises",
		ID:   "8a4f3b2e-9c71-4d5a-b0f2-c9d1e4a7b5c8",
	}
}

// mockAllRoles is the combined fixture used to seed the list screen with
// something to render across all three tabs.
func mockAllRoles() []Role {
	out := mockEntraRoles()
	out = append(out, mockAzureRoles()...)
	out = append(out, mockGroupRoles()...)
	return out
}

// mockEntraRoles matches the active tab in the role-list mockup.
// Two rows (Helpdesk, Security) are currently active.
func mockEntraRoles() []Role {
	return []Role{
		{ID: "e1", Name: "Global Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true, RequireTicket: true},
		{ID: "e2", Name: "User Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "e3", Name: "Application Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "e4", Name: "Helpdesk Administrator", Domain: DomainEntra,
			ScopeType: "AdministrativeUnit", Scope: "Finance-AU",
			Active: true, TimeLeft: 1*time.Hour + 42*time.Minute,
			MaxDuration: 4 * time.Hour, RequireReason: true},
		{ID: "e5", Name: "Groups Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "e6", Name: "Authentication Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "e7", Name: "Security Administrator", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			Active: true, TimeLeft: 3*time.Hour + 15*time.Minute,
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "e8", Name: "Security Reader", Domain: DomainEntra,
			ScopeType: "Directory", Scope: "Contoso",
			MaxDuration: 8 * time.Hour, RequireReason: true},
	}
}

// mockAzureRoles seeds the Azure Roles tab with a few representative rows.
func mockAzureRoles() []Role {
	return []Role{
		{ID: "a1", Name: "Contributor", Domain: DomainAzure,
			ScopeType: "subscription", Scope: "Subscription - Test",
			MaxDuration: 8 * time.Hour, RequireReason: true, RequireTicket: true},
		{ID: "a2", Name: "Key Vault Secrets Officer", Domain: DomainAzure,
			ScopeType: "Microsoft.KeyVault/vaults", Scope: "kv-example",
			MaxDuration: 2 * time.Hour, RequireReason: true, RequireTicket: true},
		{ID: "a3", Name: "Reader", Domain: DomainAzure,
			ScopeType: "subscription", Scope: "Subscription - Test",
			MaxDuration: 8 * time.Hour, RequireReason: true},
	}
}

// mockGroupRoles seeds the Groups tab with PIM-for-Groups entries.
func mockGroupRoles() []Role {
	return []Role{
		{ID: "g1", Name: "member", Domain: DomainGroups,
			ScopeType: "Group", Scope: "sg-ops-readers",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "g2", Name: "owner", Domain: DomainGroups,
			ScopeType: "Group", Scope: "sg-ops-admins",
			MaxDuration: 4 * time.Hour, RequireReason: true, RequireTicket: true},
		{ID: "g3", Name: "member", Domain: DomainGroups,
			ScopeType: "Group", Scope: "sg-finance-approvers",
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "g4", Name: "member", Domain: DomainGroups,
			ScopeType: "Group", Scope: "sg-platform-engineers",
			Active: true, TimeLeft: 2*time.Hour + 10*time.Minute,
			MaxDuration: 8 * time.Hour, RequireReason: true},
		{ID: "g5", Name: "owner", Domain: DomainGroups,
			ScopeType: "Group", Scope: "sg-security-leads",
			MaxDuration: 4 * time.Hour, RequireReason: true, RequireTicket: true},
	}
}

// mockActivations matches the config screen: 5 roles chosen across two domains.
// Indexes: 0 Global Admin (inherited), 1 App Admin (reason overridden),
// 2 Key Vault (duration clamped by role max), 3 Contributor (ticket missing),
// 4 Reader (no ticket required).
func mockActivations() []Activation {
	entra := mockEntraRoles()
	globalAdmin := &entra[0]
	appAdmin := &entra[2]

	kv := &Role{
		ID: "a1", Name: "Key Vault Secrets Officer", Domain: DomainAzure,
		ScopeType: "Microsoft.KeyVault/vaults", Scope: "kv-example",
		MaxDuration: 2 * time.Hour, RequireReason: true, RequireTicket: true,
	}
	contributor := &Role{
		ID: "a2", Name: "Contributor", Domain: DomainAzure,
		ScopeType: "subscription", Scope: "Subscription - Test",
		MaxDuration: 8 * time.Hour, RequireReason: true, RequireTicket: true,
	}
	reader := &Role{
		ID: "a3", Name: "Reader", Domain: DomainAzure,
		ScopeType: "subscription", Scope: "Subscription - Test",
		MaxDuration: 8 * time.Hour, RequireReason: true, RequireTicket: false,
	}

	return []Activation{
		{Role: globalAdmin},                                 // all inherited
		{Role: appAdmin, Reason: "Manual fix after deploy"}, // reason override
		{Role: kv},                                          // duration will clamp to 2h
		{Role: contributor, Ticket: " "},                    // explicit "clear" (see note)
		{Role: reader},                                      // ticket not required => "n/a"
	}
}

// Bulk defaults that appear at the top of the config screen.
func mockDefaults() (reason, ticket string, duration time.Duration) {
	return "Testing new setup in Dev Portal", "JIRA-1234", 4 * time.Hour
}

// selectedIDs returns the roleIDs chosen on the list screen.
// On the list mockup: cursor row (Global Admin) and Application Admin are selected.
func selectedIDs() map[string]bool {
	return map[string]bool{"e1": true, "e3": true}
}
