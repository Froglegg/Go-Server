package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"todo/ent"
	routes "todo/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
)

var (
	PG_USER     = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	PG_DB       = os.Getenv("PG_DB")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	connectionString := fmt.Sprintf("host=localhost port=5432 user=%s dbname=%s password=%s sslmode=disable", PG_USER, PG_DB, PG_PASSWORD)
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to pg: %v", err)
	}
	defer client.Close()

	r := chi.NewRouter()

	// middleware stack
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://foo.com", "http://localhost:3000"}, // client origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // allow for cookies (jwt in our case)
		MaxAge:           300,  // Preflight request cache duration
	}))

	// auth & handler
	tokenAuth = jwtauth.New("HS256", []byte(JWT_SECRET), nil)
	handler := &routes.Handler{Client: client, TokenAuth: tokenAuth}

	// Public routes
	r.Group(func(r chi.Router) {
		//  prevent brute force attacks
		r.Use(httprate.LimitByIP(10, 1*time.Minute))
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/register", handler.Register)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(routes.UserContextMiddleware)
		r.Get("/users/{name}", handler.QueryUser)
		r.Get("/users", handler.GetAllUsers)
		r.Post("/todos", handler.CreateTodo)
		r.Get("/todos", handler.GetTodos)
		r.Post("/todos/{id}/complete", handler.MarkTodoComplete)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

// test function for github actions
func Sum(x, y int) int {
	return x + y
}
