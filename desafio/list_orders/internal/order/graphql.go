package order

import (
	"github.com/graphql-go/graphql"
)

func GraphQLSchema(repo *Repository) (graphql.Schema, error) {
	orderType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},
			"customer": &graphql.Field{Type: graphql.String},
			"amount":   &graphql.Field{Type: graphql.Float},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return repo.List()
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})
}
