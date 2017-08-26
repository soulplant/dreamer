package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/soulplant/dreamer/api"
	"context"
	"google.golang.org/grpc"
	"net"
)

const port = ":1234"
const grpcPort = ":1235"

//go:generate ./gen-protos.sh

func main() {
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

	// Query
	query := `
		{
			hello
			bar
			obj {
				name
				subobj {
					foo
				}
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	apiMux := runtime.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	api.RegisterLoginServiceHandlerFromEndpoint(ctx, apiMux, grpcPort, []grpc.DialOption{grpc.WithInsecure()})
	http.Handle("/graphql", h)
	http.Handle("/api/", http.StripPrefix("/api", apiMux))
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		http.StripPrefix("/files/", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
	})
	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatal("Failed to listen", err)
		}
		rpcServer := grpc.NewServer()
		api.RegisterLoginServiceServer(rpcServer, &loginService{})
		log.Fatal(rpcServer.Serve(lis))
	}()
	fmt.Printf("Listening on %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Failed to listen", err)
	}
}

func OpenTestDb() *gorm.DB {
	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Project{}, &User{}, &Task{}, &Stretch{}, &Category{})
	return db
}

func test(db *gorm.DB) {
	// Read
	var project Project
	if e := db.First(&project, 1000); e.Error != nil {
		fmt.Println("Couldn't find 1000")
	}
	if e := db.First(&project, 1); e.Error != nil {
		fmt.Println("Couldn't find 1")
	}
	db.First(&project, "Name = ?", "Dreamer")

	// Delete - delete project
	// db.Delete(&project)
}
