package main

import (
	"github.com/graphql-go/graphql"
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	"log"
)

// var schema graphql.Schema

var schema graphql.Schema = makeSchema()

type loginResponse struct {
	Token string
}

func makeSchema() graphql.Schema {
	login := graphql.NewObject(graphql.ObjectConfig{
		Name: "Login",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "fake-token", nil
				},
			},
		},
	})
	subObj := graphql.NewObject(graphql.ObjectConfig{
		Name: "subobj",
		Fields: graphql.Fields{
			"foo": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "bar", nil
				},
			},
		},
	})
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "obj",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "James", nil
				},
			},
			"surname": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Kozianski", nil
				},
			},
			"subobj": &graphql.Field{
				Type: subObj,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return 1, nil
				},
			},
		},
	})
	// Schema
	fields := graphql.Fields{
		"obj": &graphql.Field{
			Type: obj,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Println("I'm resolving obj!")
				for i, f := range p.Info.FieldASTs {
					fmt.Printf("%d: %v %v\n", i, f.GetKind(), f.Name.Value)
					for j, s := range f.SelectionSet.Selections {
						nss := 0
						if ss := s.GetSelectionSet(); ss != nil {
							nss = len(ss.Selections)
						}
						if f, ok := s.(*ast.Field); ok {
							fmt.Printf("  %d: %s (%d subs)\n", j, f.Name.Value, nss)
						}
					}
				}
				return &struct{}{}, nil
			},
		},
		"bar": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// fmt.Printf("%#v\n", p)
				for i, f := range p.Info.FieldASTs {
					fmt.Printf("%d: %v %v\n", i, f.GetKind(), f.Name.Value)
				}
				return "baz", nil
			},
		},
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	mutationFields := graphql.Fields{
		"login": &graphql.Field{
			Type: login,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Printf("args: %v\n", p.Args)
				return &loginResponse{"fake-token"}, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}
