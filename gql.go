package main

import (
	"github.com/graphql-go/graphql"
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	"log"
)

// var schema graphql.Schema

var schema graphql.Schema = makeSchema()

func makeSchema() graphql.Schema {
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
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}
