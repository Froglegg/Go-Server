package routes

import (
	"todo/ent"

	"github.com/go-chi/jwtauth/v5"
)

type Handler struct {
	Client    *ent.Client
	TokenAuth *jwtauth.JWTAuth
	User      *ent.User
}
