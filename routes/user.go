package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	auth "todo/auth"
	"todo/ent"
	"todo/ent/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "userID"

func (handler *Handler) MakeToken(user ent.User) (string, error) {
	_, tokenString, err := handler.TokenAuth.Encode(map[string]interface{}{"username": user.Name, "email": user.Email, "userID": user.ID})
	if err != nil {
		log.Println("Failed to create token: ", err)
	}
	return tokenString, err
}

func UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Get the userID from the jwt claims, set as int in context
		_, claims, _ := jwtauth.FromContext(ctx)
		userID, ok := claims[string(userIDKey)]
		if !ok {
			fmt.Println("no userID in claims")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		userID = int(userID.(float64))
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx.Value(userIDKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the user by name
	user, err := handler.Client.User.Query().Where(user.Email(loginDetails.Email)).Only(ctx)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Validate password
	if !auth.CheckPasswordHash(loginDetails.Password, user.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := handler.MakeToken(*user)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User logged in successfully",
	})
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1, // Delete the cookie.
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt",
		Value: "",
	})
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User logged out successfully",
	})
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var registrationDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registrationDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create the user
	user, err := handler.Client.User.Create().
		SetEmail(registrationDetails.Email).
		SetPassword(string(hashedPassword)). // Use the hashed password
		SetName(registrationDetails.Name).
		SetAge(registrationDetails.Age).
		Save(ctx)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := handler.MakeToken(*user)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		// Uncomment below for HTTPS:
		// Secure: true,
		Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value: token,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
	})
}

func (handler *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := handler.Client.User.Query().All(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching users: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (handler *Handler) QueryUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name := chi.URLParam(r, "name")
	user, err := handler.Client.User.Query().Where(user.Name(name)).Only(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("no user or multiple users found: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("user returned: ", user)
	json.NewEncoder(w).Encode(user)
}
