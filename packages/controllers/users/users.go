package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/AriJaya07/go-rest-api/packages/config"
	"github.com/AriJaya07/go-rest-api/packages/controllers/auth"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	"github.com/AriJaya07/go-rest-api/packages/models/types"
	"github.com/AriJaya07/go-rest-api/packages/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store     store.Store
	jwtSecret []byte
}

func NewUserService(s store.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *types.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	hashedPW, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error creating user"})
		return
	}
	payload.Password = hashedPW

	u, err := s.store.CreateUser(payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error creating user"})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error creating session"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	// 1. Finf user in db by email
	var input types.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. compare password with hashed password
	var user types.User
	err := s.store.QueryRow("SELECT id, email, password FROM users WHERE email = ?", input.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Create JWY and set it in a cookie
	token, err := auth.CreateJWT(s.jwtSecret, user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	// 4. Return JWT in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func validateUserPayload(user *types.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
