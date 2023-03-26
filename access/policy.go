package access

import "github.com/katallaxie/pkg/urn"

// Policy is a set of rules that define how a user can access a resource.
type Policy struct {
	// ID is the unique identifier of the policy.
	ID string `json:"id"`
	// Name is the name of the policy.
	Name string `json:"name"`
	// Description is the description of the policy.
	Description string `json:"description"`
	// Rules is the list of rules that define how a user can access a resource.
	Rules []Rule `json:"rules"`
}

// Rule is a set of conditions that define how a user can access a resource.
type Rule struct {
	// Resources is the list of resources that the rule applies to.
	Resources []string `json:"resources"`
	// Actions is the list of actions that the rule applies to.
	Actions []string `json:"actions"`
	// Effect is the effect of the rule, it can be allow or deny.
	Effect Effect `json:"effect"`
	// Conditions is the list of conditions that the rule applies to.
	Conditions []Condition `json:"conditions"`
}

// Condition is a set of key-value pairs that define how a user can access a resource.
type Condition struct {
	// Key is the key of the condition.
	Key string `json:"key"`
	// Value is the value of the condition.
	Value string `json:"value"`
	// Operator is the operator of the condition.
	Operator string `json:"operator"`
}

// Effect is the effect of the rule, it can be allow or deny.
type Effect string

// Allow effect.
const Allow Effect = "allow"

// Deny effect.
const Deny Effect = "deny"

// Action is the action that the rule applies to.
type Action string

// Allowed returns true if the effect is allow.
func (p *Policy) Allowed(res *urn.URN, action Action) (bool, error) {
	for _, r := range p.Rules {
		for _, a := range r.Actions {
			if Action(a) == action {
				for _, rr := range r.Resources {
					u, err := urn.Parse(rr)
					if err != nil {
						return false, err
					}

					if !res.Match(u) {
						continue
					}

					// only allow if effect is explicitly set to allow
					if r.Effect == Allow {
						return true, nil
					}

					return false, nil
				}
			}
		}
	}

	return false, nil
}
