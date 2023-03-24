package urn

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	// Separator is the separator used to separate the segments of a URN.
	Seperator = ":"
)

// ErrorInvalid is returned when parsing an URN with an invalid format.
var ErrorInvalid = errors.New("invalid URN format")

var validate = validator.New()

// URN represents a unique, uniform identifier for a resource
type URN struct {
	// Namespace is the namespace segment of the URN.
	Namespace string `validate:"required"`
	// Partition is the partition segment of the URN.
	Partition string `validate:"max=256"`
	// Service is the service segment of the URN.
	Service string `validate:"max=256"`
	// Region is the region segment of the URN.
	Region string `validate:"max=256"`
	// Identifier is the identifier segment of the URN.
	Identifier string `validate:"max=64"`
	// Resource is the resource segment of the URN.
	Resource string `validate:"required,max=256"`
}

// String returns the string representation of the URN.
func (u *URN) String() string {
	return strings.Join([]string{u.Namespace, u.Partition, u.Service, u.Region, u.Identifier, u.Resource}, Seperator)
}

// Match returns true if the URN matches the given URN.
func (u *URN) Match(urn *URN) bool {
	return u.String() == urn.String()
}

// New takes a namespace, partition, service, region, identifier and resource and returns a URN.
func New(namespace, partition, service, region, identifier, resource string) (*URN, error) {
	urn := &URN{
		Namespace:  namespace,
		Partition:  partition,
		Service:    service,
		Region:     region,
		Identifier: identifier,
		Resource:   resource,
	}

	validate = validator.New()

	if err := validate.Struct(urn); err != nil {
		return nil, err
	}

	return urn, nil
}

// Parse takes a string and parses it to a URN.
func Parse(s string) (*URN, error) {
	segments := strings.SplitN(s, Seperator, 6)
	if len(segments) < 5 {
		return nil, ErrorInvalid
	}

	urn := &URN{
		Namespace:  segments[0],
		Partition:  segments[1],
		Service:    segments[2],
		Region:     segments[3],
		Identifier: segments[4],
		Resource:   segments[5],
	}

	validate = validator.New()

	if err := validate.Struct(urn); err != nil {
		return nil, err
	}

	return urn, nil
}

// ProtoMessage returns the proto.ResourceURN representation of the URN.
func (u *URN) ProtoMessage() *ResourceURN {
	return &ResourceURN{
		Canonical:  u.String(),
		Namespace:  u.Namespace,
		Partition:  u.Partition,
		Service:    u.Service,
		Region:     u.Region,
		Identifier: u.Identifier,
		Resource:   u.Resource,
	}
}

// FromProto returns the URN representation from a proto.ResourceURN representation.
func FromProto(r *ResourceURN) (*URN, error) {
	return New(r.Namespace, r.Partition, r.Service, r.Region, r.Identifier, r.Resource)
}
