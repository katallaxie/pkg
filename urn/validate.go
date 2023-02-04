package urn

import (
	"errors"
)

var (
	// ErrMissingCanonical signals that the canonical representation of the ResourceURN is missing
	ErrMissingCanonical = errors.New("urn: the URN is missing the canonical representation")
	// ErrInvalidCanonical  signals that the canonical representation  is invalid
	ErrInvalidCanonical = errors.New("urn: the canonical representation is invalid")
	// ErrMissingNamespace signals that the namespace of the ResourceURN is missing
	ErrMissingNamespace = errors.New("urn: the URN is missing a namespace")
	// ErrMissingCollection signals that the collection of the ResourceURN is missing
	ErrMissingCollection = errors.New("urn: the URN is missing a collection")
	// ErrMissingIdentifier signals that the identifier of the ResourceURN is missing
	ErrMissingIdentifier = errors.New("urn: the URN is missing a identifier")
	// ErrMissingResource signals that the resource of the ResourceURN is missing
	ErrMissingResource = errors.New("urn: the URN is missing a identifier")
)

// Validate returns a non-nil error in case the given URN
// is not fully valid. Meaning it does not transport all information
// in the message. This can be the case in transitioning to use
// the URN.proto and not yet fully publish everything.
//
//	urn := req.GetUrn()
//	if err := urn.Validate(); err != nil{
//		return err
//	}
func (r *ResourceURN) Validate(validateFuncs ...ValidateFunc) error {
	canonical := r.GetCanonical()
	if canonical == "" {
		return ErrMissingCanonical
	}

	_, err := Parse(r.GetCanonical(), validateFuncs...)
	if err != nil {
		return ErrInvalidCanonical
	}

	namespace := r.GetNamespace()
	if namespace == "" {
		return ErrMissingNamespace
	}

	collection := r.GetCollection()
	if collection == "" {
		return ErrMissingCollection
	}

	identifier := r.GetIdentifier()
	if identifier == "" {
		return ErrMissingIdentifier
	}

	resource := r.GetResource()
	if resource == "" {
		return ErrMissingResource
	}

	return nil
}
