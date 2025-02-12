package controllers

import (
	"avito-shop/models"
	"avito-shop/services"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var jwtKey = []byte("my_secret_key")

type AuthController struct {
	UserService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{
		UserService: userService,
	}
}

func (c *AuthController) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var authReq models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := c.UserService.CreateUser(authReq.Username, authReq.Password); err != nil {
		if isAuthenticated, _ := c.UserService.AuthenticateUser(authReq.Username, authReq.Password); isAuthenticated {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": authReq.Username,
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			})

			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(models.AuthResponse{Token: tokenString})
			if err != nil {
				return
			}
			return
		}
		if err.Error() == "user already exists" {
			http.Error(w, "User already exists", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": authReq.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(models.AuthResponse{Token: tokenString})
	if err != nil {
		return
	}
}
