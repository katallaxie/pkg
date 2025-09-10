package protox

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrUnimplemented is returned when a method is not implemented.
var ErrUnimplemented = errors.New("method not implemented")

// ToProto is an interface for types that can be converted to protobuf messages.
type ToProto interface {
	ToProto() (proto.Message, error)
}

// FromProto is an interface for types that can be populated from protobuf messages.
type FromProto interface {
	FromProto(ptr any) error
}

// UnimplementedToProto is a struct that implements the ToProto interface with a default method.
type UnimplementedToProto struct{}

// ToProto returns nil, indicating that the method is unimplemented.
func (u UnimplementedToProto) ToProto() (proto.Message, error) {
	return nil, ErrUnimplemented
}

// UnimplementedFromProto is a struct that implements the FromProto interface with a default method.
type UnimplementedFromProto struct{}

// FromProto returns ErrUnimplemented, indicating that the method is unimplemented.
func (u UnimplementedFromProto) FromProto(ptr any) error {
	return ErrUnimplemented
}

// ProtoX is an interface that combines ToProto and FromProto.
type ProtoX interface {
	ToProto
	FromProto
}

// UnimplementedProto is a struct that implements the Proto interface with default methods.
type UnimplementedProto struct {
	UnimplementedToProto
	UnimplementedFromProto
}
