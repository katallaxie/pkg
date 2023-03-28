package access

import (
	"strings"

	"github.com/katallaxie/pkg/urn"
)

const (
	// ActionAccessService is the action that the rule applies to.
	AccessService urn.Match = "access"
)

// IsRole returns true if the resource is a role.
func IsRole(u *urn.URN) bool {
	return u.Service == AccessService && strings.HasPrefix(u.Resource.String(), "roles")
}

// IsUser returns true if the resource is a user.
func IsUser(u *urn.URN) bool {
	return u.Service == AccessService && strings.HasPrefix(u.Resource.String(), "users")
}

// IsGroup returns true if the resource is a group.
func IsGroup(u *urn.URN) bool {
	return u.Service == AccessService && strings.HasPrefix(u.Resource.String(), "groups")
}

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

// Matcher is a function that returns true if the URN matches.
type Matcher func(l *urn.URN, r *urn.URN) bool

// IdentityBasedMatcher is a matcher that matches the URN based on the identity.
//
//nolint:gocyclo
var IdentityBasedMatcher = func(l *urn.URN, r *urn.URN) bool {
	return (l.Namespace == r.Namespace || (l.Namespace == urn.Wildcard && r.Namespace == urn.Wildcard) || (l.Namespace == urn.Empty && r.Namespace == urn.Empty) || r.Namespace == urn.Wildcard || r.Namespace == urn.Empty) &&
		(l.Partition == r.Partition || (l.Partition == urn.Wildcard && r.Partition == urn.Wildcard) || (l.Partition == urn.Empty && r.Partition == urn.Empty) || r.Partition == urn.Wildcard || r.Partition == urn.Empty) &&
		(l.Service == r.Service || (l.Service == urn.Wildcard && r.Service == urn.Wildcard) || (l.Service == urn.Empty && r.Service == urn.Empty) || r.Service == urn.Wildcard || r.Service == urn.Empty) &&
		(l.Region == r.Region || (l.Region == urn.Wildcard && r.Region == urn.Wildcard) || (l.Region == urn.Empty && r.Region == urn.Empty) || r.Region == urn.Wildcard || r.Region == urn.Empty) &&
		(l.Identifier == r.Identifier || (l.Identifier == urn.Wildcard && r.Identifier == urn.Wildcard) || (l.Identifier == urn.Empty && r.Identifier == urn.Empty) || r.Identifier == urn.Wildcard || r.Identifier == urn.Empty) &&
		(l.Resource == r.Resource || (l.Resource == urn.Wildcard && r.Resource == urn.Wildcard) || (l.Resource == urn.Empty && r.Resource == urn.Empty) || r.Resource == urn.Wildcard || r.Resource == urn.Empty)
}

// NewEvaluator returns a new evaluator.
func NewEvaluator(m Matcher, a Accessor) *Evaluator {
	return &Evaluator{m, a}
}

// Accessor is the interface to allow or deny access.
type Accessor interface {
	// Allow returns true if the user is allowed to perform the action on the resource.
	Allow(res *urn.URN, action Action, policies ...Policy) (bool, error)
}

// Evaluator evaluates the policies and returns true if the user is allowed to perform the action on the resource.
type Evaluator struct {
	Matcher
	Accessor
}

// Allow returns true if the user is allowed to perform the action on the resource.
func (d *Evaluator) Allow(res *urn.URN, action Action, policies ...Policy) (bool, error) {
	var allow bool // default to deny

	for _, p := range policies {
		for _, r := range p.Rules {
			for _, a := range r.Actions {
				if a == action {
					for _, rr := range r.Resources {
						u, err := urn.Parse(rr.String())
						if err != nil {
							return false, err
						}

						if !d.Matcher(u, res) {
							continue
						}

						allow = r.Effect == Allow
					}
				}
			}
		}
	}

	return allow, nil
}

// Rule is a set of conditions that define how a user can access a resource.
type Rule struct {
	// ID is the unique identifier of the rule.
	ID string `json:"id"`
	// Resources is the list of resources that the rule applies to.
	Resources []Resource `json:"resources"`
	// Actions is the list of actions that the rule applies to.
	Actions []Action `json:"actions"`
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

// String ...
func (a Action) String() string {
	return string(a)
}

// Resource is the resource that the rule applies to.
type Resource string

// String returns the string representation of the resource.
func (r Resource) String() string {
	return string(r)
}

// URN returns the URN representation of the resource.
func (r Resource) URN() (*urn.URN, error) {
	return urn.Parse(r.String())
}