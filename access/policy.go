package access

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/katallaxie/pkg/urn"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// DefaultVersion
const DefaultVersion = "2023-03-28"

const (
	defaultAccessService = "access"
)

// Servide is the service name of the access service.
type Service urn.Match

// Services is the list of services.
type Services map[Service]bool

// Add adds a service to the list.
func (s Services) Add(service Service) {
	s[service] = true
}

// DefaultServices is the default list of services.
var DefaultServices = Services{
	defaultAccessService: true,
}

// Region is the region name of the access service.
type Region urn.Match

// Regions is the list of regions.
type Regions map[Region]bool

// Add adds a region to the list.
func (r Regions) Add(region Region) {
	r[region] = true
}

// DefaultRegions is the default list of regions.
var DefaultRegions = Regions{
	"eu-central-1": true,
}

// Partition is the partition name of the access service.
type Partition urn.Match

// Partitions is the list of partitions.
type Partitions map[Partition]bool

// Add adds a partition to the list.
func (p Partitions) Add(partition Partition) {
	p[partition] = true
}

// DefaultPartitions is the default list of partitions.
var DefaultPartitions = Partitions{
	"cloud": true,
}

// ResourceIdentifier is the unique identifier of a resource.
type ResourceIdentifier func(*urn.URN) bool

// RoleResourceIdentifier is the identifier for a role.
var RoleResourceIdentifier = func(u *urn.URN) bool {
	return u.Service == defaultAccessService && strings.HasPrefix(u.Resource.String(), "roles")
}

// UserResourceIdentifier is the identifier for a user.
var UserResourceIdentifier = func(u *urn.URN) bool {
	return u.Service == defaultAccessService && strings.HasPrefix(u.Resource.String(), "users")
}

// GroupResourceIdentifier is the identifier for a group.
var GroupResourceIdentifier = func(u *urn.URN) bool {
	return u.Service == defaultAccessService && strings.HasPrefix(u.Resource.String(), "groups")
}

// Is returns true if the resource matches the identifier.
func Is(u *urn.URN, i ResourceIdentifier) bool {
	return i(u)
}

// Policy is a set of rules that define how a user can access a resource.
type Policy struct {
	// Version is the version of the policy.
	Version string `json:"version" yaml:"version"`
	// ID is the unique identifier of the policy.
	ID string `json:"id" yaml:"id"`
	// Name is the name of the policy.
	Name string `json:"name" yaml:"name"`
	// Description is the description of the policy.
	Description string `json:"description" yaml:"description"`
	// Rules is the list of rules that define how a user can access a resource.
	Rules Rules `json:"rules" yaml:"rules"`
}

// DefaultPolicy returns the default policy.
func DefaultPolicy() *Policy {
	return &Policy{
		Version: DefaultVersion,
		Rules:   []Rule{},
	}
}

// UnmarshalJSON overwrite own policy with values of the given in policy in JSON format
func (p *Policy) UnmarshalJSON(data []byte) error {
	pol := struct {
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		Rules       Rules  `json:"rules"`
		Version     string `json:"version"`
	}{}

	if err := json.Unmarshal(data, &pol); err != nil {
		return errors.WithStack(err)
	}

	p.Rules = pol.Rules
	p.Version = pol.Version

	return nil
}

// UnmarshalYAML overwrite own policy with values of the given policy in YAML format.
func (p *Policy) UnmarshalYAML(data []byte) error {
	pol := struct {
		Description string `yaml:"description"`
		ID          string `yaml:"id"`
		Name        string `yaml:"name"`
		Rules       Rules  `yaml:"rules"`
		Version     string `yaml:"version"`
	}{}

	if err := yaml.Unmarshal(data, &pol); err != nil {
		return errors.WithStack(err)
	}

	p.Rules = pol.Rules
	p.Version = pol.Version

	return nil
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

// UnimplementedAccessor is the default implementation of the Accessor interface.
type UnimplementedAccessor struct{}

// Allow returns true if the user is allowed to perform the action on the resource.
func (u *UnimplementedAccessor) Allow(principal *urn.URN, ressource *urn.URN, action Action) (bool, error) {
	return false, errors.New("not implemented")
}

// Accessor is the interface to allow or deny access.
type Accessor interface {
	// Allow returns true if the user is allowed to perform the action on the resource.
	Allow(ctx context.Context, principal *urn.URN, ressource *urn.URN, action Action) (bool, error)
}

// Policer returns the policy for the given user.
type Policer interface {
	// Policies returns the policy for the given user.
	Policies(ctx context.Context, principal *urn.URN) ([]*Policy, error)
}

// Rule is a set of conditions that define how a user can access a resource.
type Rule struct {
	// ID is the unique identifier of the rule.
	ID string `json:"id" yaml:"id"`
	// Resources is the list of resources that the rule applies to.
	Resources Resources `json:"resources" yaml:"resources"`
	// Actions is the list of actions that the rule applies to.
	Actions Actions `json:"actions" yaml:"actions"`
	// Effect is the effect of the rule, it can be allow or deny.
	Effect Effect `json:"effect" yaml:"effect"`
	// Conditions is the list of conditions that the rule applies to.
	Conditions Conditions `json:"conditions" yaml:"conditions"`
}

// Rules is a list of rules.
type Rules []Rule

// Condition is a set of key-value pairs that define how a user can access a resource.
type Condition struct {
	// Key is the key of the condition.
	Key string `json:"key" yaml:"key"`
	// Value is the value of the condition.
	Value string `json:"value" yaml:"value"`
	// Operator is the operator of the condition.
	Operator string `json:"operator" yaml:"operator"`
}

// Conditions is a list of conditions.
type Conditions []Condition

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

// Actions is a list of actions.
type Actions []Action

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

// Resources is a list of resources.
type Resources []Resource
