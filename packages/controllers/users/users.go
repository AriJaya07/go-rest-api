package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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
	store store.Store
}

func NewUserService(s store.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", s.handleGetAllUser).Methods("GET")
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
	r.HandleFunc("/users/edit-profile/{id}", s.handleUserUpdate).Methods("PUT")
	r.HandleFunc("/users/delete/{id}", s.handleUserDelete).Methods("DELETE")
	r.HandleFunc("/users/change-password/{id}", s.handleChangePassword).Methods("PUT")
}

func (s *UserService) handleGetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetAllUsers()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error getting all user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
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

	token, err := createAndSetAuthCookie(u.ID, u.Email, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error creating session"})
		return
	}

	response := types.LoginResponse{
		Email: u.Email,
		Token: token,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
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
	token, err := createAndSetAuthCookie(user.ID, user.Email, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Error user not found"})
		return
	}

	// 4. Return JWT in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "email": user.Email})
}

func (s *UserService) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	// 1. Decode JSON input
	var input types.ChangePassword

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	defer r.Body.Close()

	// 2. Validate input
	if input.CurrentPassword == "" || input.NewPassword == "" {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Fetch user from store using userID (convert userID to string if required by your GetUserByID function)
	user, err := s.store.GetUserByID(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{Error: "Failed to fetch user"})
	}

	// verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{Error: "Invalid current password"})
		return
	}

	// 4. Update password
	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to generate hashed password"})
		return
	}

	user.Password = string(hashNewPassword)
	if err := s.store.UpdateUser(user); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password update successfully"})
}

func (s *UserService) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	// Decode JSON input
	var input types.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	// Fetch user from store
	user, err := s.store.GetUserByID(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	// Update user fields
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	// Update user in store
	if err := s.store.UpdateUser(user); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (s *UserService) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Missing user ID"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	if _, err := s.store.GetUserByID(idStr); err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "User not found"})
			return
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to check user existence"})
		}
	}

	// Delete the user
	if err := s.store.DeleteUser(id); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to delete user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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

func createAndSetAuthCookie(id int64, email string, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, id, email)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	return token, nil
}
