package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"todo/ent"
	"todo/ent/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	r := chi.NewRouter()

	// use a permissive allow all handler
	r.Use(cors.AllowAll().Handler)

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		CreateUser(w, r, client)
	})

	r.Get("/users/{name}", func(w http.ResponseWriter, r *http.Request) {
		QueryUser(w, r, client)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		GetAllUsers(w, r, client)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func CreateUser(w http.ResponseWriter, r *http.Request, client *ent.Client) {
	ctx := context.Background()
	u := &ent.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := client.User.Create().SetAge(u.Age).SetName(u.Name).Save(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("user was created: ", user)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, client *ent.Client) {
	ctx := context.Background()
	users, err := client.User.Query().All(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching users: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func QueryUser(w http.ResponseWriter, r *http.Request, client *ent.Client) {
	ctx := context.Background()
	name := chi.URLParam(r, "name")
	user, err := client.User.Query().Where(user.Name(name)).Only(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("no user or multiple users found: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("user returned: ", user)
	json.NewEncoder(w).Encode(user)
}

func Sum(x, y int) int {
	return x + y
}
