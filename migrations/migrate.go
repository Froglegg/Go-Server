package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"todo/ent"

	_ "github.com/lib/pq"
)

var (
	PG_USER     = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	PG_DB       = os.Getenv("PG_DB")
)

func main() {
	connectionString := fmt.Sprintf("host=localhost port=5432 user=%s dbname=%s password=%s sslmode=disable", PG_USER, PG_DB, PG_PASSWORD)
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to pg: %v", err)
	}
	defer client.Close()

	// Run auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
