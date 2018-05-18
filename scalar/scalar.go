package scalar

import "fmt"

// Custom defines an interface capable of deserializing a custom scalar.
//
// This interface assumes that the name of the type implementing and the name of the
// custom scalar match.
//
// The spec defines that all scalars should be representable as strings, therefore
// a custom scalar Unmarshaler must also satisfy the fmt.Stringer interface.
// See http://facebook.github.io/graphql/draft/#sec-Scalars
type Custom interface {
	fmt.Stringer
	// UnmarshalGraphQL takes an interface and populates the underlying instance accordingly.
	//
	// An error is returned if something failes during this process.
	UnmarshalGraphQL(input interface{}) error
}

// NamedCustom extends Custom with a convinience method that allows implementations
// to have different naming than the corresponding GraphQL Scalar.
type NamedCustom interface {
	Custom
	// ImplementsGraphQLType returns true if the implementation supports a GraphQL type
	// of the provided name.
	ImplementsGraphQLType(name string) bool
}
