# Pkg

[![Test & Lint](https://github.com/katallaxie/pkg/actions/workflows/main.yml/badge.svg)](https://github.com/katallaxie/pkg/actions/workflows/main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/katallaxie/pkg.svg)](https://pkg.go.dev/github.com/katallaxie/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/katallaxie/pkg)](https://goreportcard.com/report/github.com/katallaxie/pkg)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`pkg` is a collection of Go packages that make the life of a Go developers at ZEISS easier.

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/katallaxie/pkg?quickstart=1)

## Installation

```bash
go get github.com/katallaxie/pkg
```

Go has a pretty good standard library, but there are some things that are missing. This collection of small packages is meant to fill those gaps.

## Casting values

There is the typical case of pointers you need to deference. For example in the strict interfaces generated by [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

```go

Somestruct {
  Field: cast.Ptr(req.Field),
}

```

Or the other way around, you have a pointer and you want to get the value or a default value.

```go

SomeStruct {
  Field: cast.Value(req.Field),
}

```

Return a default value of a `nil` pointer.

```go
type Foo struct {}
cast.Zero(Foo) // &Foo{}
```

## Converting values

There is also always the case to covert a value to a specific other value.

```go
// String converts a value to a string.
b := true
s := cast.String(b)

fmt.Println(s) // "true"
```

There are functions to convert `int`, `string` and `bool` values.

## Operators 

There is the implementation of various operators.

```go
// Is the ternary operator.
utilx.IfElse(cond, 100, 0)
```

## Databases

There are also more complex tools like the `Database` interface which enables to easliy implement database wrappers.

```go
// Database provides methods for transactional operations.
type Database[R, W any] interface {
	// ReadTx starts a read only transaction.
	ReadTx(context.Context, func(context.Context, R) error) error
	// ReadWriteTx starts a read write transaction.
	ReadWriteTx(context.Context, func(context.Context, W) error) error

	Migrator
	io.Closer
}
```

Or a simple interface to implement servers. This takes for all `signal` and `context` handling.

```go

s, _ := server.WithContext(ctx)
s.Listen(&srv{}, true)

serverErr := &server.ServerError{}
if err := s.Wait(); errors.As(err, &serverErr) {
  log.Print(err)
	os.Exit(1)
}
```

## FGA with OpenFGA

There is also a package to work with the OpenFGA API.

```go
// Store is an interface that provides methods for transactional operations on the authz database.
type Store[Tx any] interface {
	// Allowed checks if the user is allowed to perform the operation on the object.
	Allowed(context.Context, User, Object, Relation) (bool, error)
	// WriteTx starts a read write transaction.
	WriteTx(context.Context, func(context.Context, Tx) error) error
}

// StoreTx is an interface that provides methods for transactional operations on the authz database.
type StoreTx interface {
	// WriteTuple writes a tuple to the authz database.
	WriteTuple(context.Context, User, Object, Relation) error
	// DeleteTuple deletes a tuple from the authz database.
	DeleteTuple(context.Context, User, Object, Relation) error
}
```

This can be used with the package.

```go
authzStore, err := authx.NewStore(fgaClient, authz.NewWriteTx())
if err != nil {
  return err
}
```

## License

[MIT](/LICENSE)
