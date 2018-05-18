package scalar_test

import (
	"fmt"
	"strings"
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/graph-gophers/graphql-go/scalar"
)

type MyScalar struct {
	Name string
}

func (m *MyScalar) UnmarshalGraphQL(input interface{}) error {
	if str, ok := input.(string); ok {
		m.Name = str
		return nil
	}
	return fmt.Errorf("%s is not a string", input)
}

func (m *MyScalar) String() string {
	return m.Name
}

var _ scalar.Custom = &MyScalar{}

type MyNamedReturn struct {
	Result string
}

func (m *MyNamedReturn) ImplementsGraphQLType(name string) bool {
	return name == "MyReturn"
}

func (m *MyNamedReturn) UnmarshalGraphQL(input interface{}) error {
	if str, ok := input.(string); ok {
		m.Result = str
		return nil
	}
	return fmt.Errorf("%s is not a string", input)
}

func (m *MyNamedReturn) String() string {
	return m.Result
}

var _ scalar.NamedCustom = &MyNamedReturn{}

type customScalarResolver struct{}

func (c *customScalarResolver) ID() graphql.ID {
	return graphql.ID("test")
}

func (c *customScalarResolver) Exists(args struct{ SearchID graphql.ID }) bool {
	return graphql.ID("success") == args.SearchID
}

func (c *customScalarResolver) ToUpper(args struct{ Input MyScalar }) string {
	return strings.ToUpper(args.Input.Name)
}

func (c *customScalarResolver) ToLower(args struct{ Input string }) *MyNamedReturn {
	return &MyNamedReturn{Result: strings.ToLower(args.Input)}
}

func TestCustomScalarWithReturnString(t *testing.T) {
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(`
				schema {
					query: Query
				}
				scalar MyScalar
				type Query {
					id: ID!
					exists(searchID: ID!): Boolean!
					toUpper(input: MyScalar!): String!
				}
			`, &customScalarResolver{}),
			Query: `
				{
					id
					exists(searchID: "success")
					toUpper(input: "phr\"ase")
				}
			`,
			ExpectedResult: `
				{
					"id": "test",
					"exists": true,
					"toUpper": "PHR\"ASE"
				}
			`,
		},
	})
}

func TestStringWithReturnCustomScalar(t *testing.T) {
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(`
				schema {
					query: Query
				}
				scalar MyReturn
				type Query {
					id: ID!
					exists(searchID: ID!): Boolean!
					toLower(input: String!): MyReturn!
				}
			`, &customScalarResolver{}),
			Query: `
				{
					id
					exists(searchID: "success")
					toLower(input: "PHRASE")
				}
			`,
			ExpectedResult: `
				{
					"id": "test",
					"exists": true,
					"toLower": "phrase"
				}
			`,
		},
	})
}

func TestCustomScalarError(t *testing.T) {
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(`
				schema {
					query: Query
				}
				scalar MyScalar
				type Query {
					id: ID!
					exists(searchID: ID!): Boolean!
					toUpper(input: MyScalar!): String!
				}
			`, &customScalarResolver{}),
			Query: `
				{
					id
					exists(searchID: "success")
					toUpper(input: 123)
				}
			`,
			ExpectedResult: `
				{
					"id": "test",
					"exists": true
				}
			`,
			ExpectedErrors: []*errors.QueryError{
				errors.Errorf("%s is not a string", int32(123)),
			},
		},
	})
}
