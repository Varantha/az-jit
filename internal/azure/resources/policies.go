package azure

import (
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

// EndUserActivationPolicy holds the self-activation settings extracted from an
// EndUser assignment policy.
type EndUserActivationPolicy struct {
	MaxDuration   time.Duration
	RequireReason bool
	RequireTicket bool
}

func extractEndUserActivationPolicy(
	rules []armauthorization.RoleManagementPolicyRuleClassification,
) (EndUserActivationPolicy, error) {
	policy := EndUserActivationPolicy{
		MaxDuration:   0,
		RequireReason: false,
		RequireTicket: false,
	}
	for _, rule := range rules {
		base := rule.GetRoleManagementPolicyRule()

		if base.Target == nil ||
			base.Target.Caller == nil || *base.Target.Caller != "EndUser" ||
			base.Target.Level == nil || *base.Target.Level != "Assignment" {
			continue
		}

		switch r := rule.(type) {
		case *armauthorization.RoleManagementPolicyExpirationRule:
			if r.MaximumDuration == nil {
				break
			}
			// MaxDuration: r.MaximumDuration is *string, e.g. "PT8H"
			expDuration, err := parseISO8601Duration(*r.MaximumDuration)
			if err != nil {
				return EndUserActivationPolicy{}, fmt.Errorf("parse max duration: %w", err)
			}
			policy.MaxDuration = expDuration

		case *armauthorization.RoleManagementPolicyEnablementRule:
			for _, e := range r.EnabledRules {
				switch *e {
				case armauthorization.EnablementRulesJustification:
					policy.RequireReason = true
				case armauthorization.EnablementRulesTicketing:
					policy.RequireTicket = true
				}
			}
		}
	}
	// No duration means no EndUser assignment expiration rule was found (or it
	// had a nil/zero maximum), so the policy is unusable for self-activation —
	// treat that as an error rather than returning a silently-degraded struct.
	if policy.MaxDuration == 0 {
		return EndUserActivationPolicy{}, errors.New("unable to determine max duration of rule")
	}
	return policy, nil
}
