package graphql_test

import (
	"fmt"
	"strings"
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/gqltesting"
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

func TestCustomScalar(t *testing.T) {
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
					toUpper(input: "phrase")
				}
			`,
			ExpectedResult: `
				{
					"id": "test",
					"exists": true,
					"toUpper": "PHRASE"
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
