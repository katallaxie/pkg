package urn

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrorInvalid is returned when parsing an URN with an
	// invalid format.
	ErrorInvalid = errors.New("invalid URN format")

	segmentFormat = `[a-zA-Z0-9-_~\.]{1,256}`
	urnRegexp     = regexp.MustCompile(`^(` + segmentFormat + `):(` + segmentFormat + `):(` + segmentFormat + `):(` + segmentFormat + `)$`)
)

// ValidateFunc is a function to validate a URN.
type ValidateFunc func(*URN) error

// URN represents a unique, uniform identifier for a resource
type URN struct {
	// namespace is the namespace segment of the URN.
	namespace string
	// collection is the collection segment of the URN.
	collection string
	// identifier is the identifier segment of the URN.
	identifier string
	// resource is the resource segment of the URN.
	resource string
}

// Namespace returns the namespace segment of the URN.
func (u *URN) Namespace() string {
	return u.namespace
}

// Collection returns the collection segment of the URN.
func (u *URN) Collection() string {
	return u.collection
}

// Identifier returns the identifier segment of the URN.
func (u *URN) Identifier() string {
	return u.identifier
}

// Resource returns the resource segment of the URN.
func (u *URN) Resource() string {
	return u.resource
}

// String returns the string representation of the URN.
func (u *URN) String() string {
	return strings.Join([]string{u.namespace, u.collection, u.identifier, u.resource}, ":")
}

// ValidateNamespace is a function to validate that the namespace is not empty.
func ValidateNamespace() ValidateFunc {
	return func(u *URN) error {
		if u.namespace == "" {
			return ErrorInvalid
		}

		return nil
	}
}

// New takes a Collection, an identifier, a resource and an optional list of
// options, and returns a URN. It may return a non-nil error if namespace, coll or identifier
// contain invalid values (e.g. an empty namespace, an unknown
// collection, etc).
func New(namespace, collection, identifier, resource string, validateFunc ...ValidateFunc) (*URN, error) {
	urn := &URN{
		namespace:  namespace,
		collection: collection,
		identifier: identifier,
		resource:   resource,
	}

	for _, fn := range validateFunc {
		if err := fn(urn); err != nil {
			return nil, err
		}
	}

	return urn, nil
}

// Parse takes a string and parses it to a URN.
func Parse(s string, validateFunc ...ValidateFunc) (*URN, error) {
	segments := urnRegexp.FindStringSubmatch(s)
	if len(segments) < 5 { // the first element is the full canonical string
		return nil, ErrorInvalid
	}

	urn := &URN{
		namespace:  segments[1],
		collection: segments[2],
		identifier: segments[3],
		resource:   segments[4],
	}

	return urn, nil
}

// ProtoMessage returns the proto.ResourceURN representation of the URN.
func (u *URN) ProtoMessage() *ResourceURN {
	return &ResourceURN{
		Canonical:  u.String(),
		Namespace:  u.Namespace(),
		Collection: u.Collection(),
		Identifier: u.Identifier(),
		Resource:   u.Resource(),
	}
}

// FromProto returns the URN representation from a proto.ResourceURN representation.
func FromProto(r *ResourceURN, validateFunc ...ValidateFunc) (*URN, error) {
	return &URN{
		namespace:  r.Namespace,
		collection: r.Collection,
		identifier: r.Identifier,
		resource:   r.Resource,
	}, nil
}
