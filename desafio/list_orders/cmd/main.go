package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"orders-app/internal/order"
	"orders-app/proto/orderpb"
)

func main() {
	// Conexão com Postgres
	connStr := "postgres://postgres:postgres@db:5432/orders?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &order.Repository{DB: db}

	// -------------------
	// REST
	// -------------------
	go func() {
		fmt.Println("REST running on :8081")
		if err := http.ListenAndServe(":8081", order.RESTHandler(repo)); err != nil {
			log.Fatal(err)
		}
	}()

	// -------------------
	// gRPC
	// -------------------
	go func() {
		fmt.Println("gRPC running on :50051")
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}
		s := grpc.NewServer()
		orderpb.RegisterOrderServiceServer(s, &order.GRPCServer{Repo: repo})

		// Habilita Reflection para grpcurl e outras ferramentas
		reflection.Register(s)

		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// -------------------
	// GraphQL
	// -------------------
	schema, err := order.GraphQLSchema(repo)
	if err != nil {
		log.Fatal(err)
	}
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
	fmt.Println("GraphQL running on :8082/graphql")
	if err := http.ListenAndServe(":8082", h); err != nil {
		log.Fatal(err)
	}
}
