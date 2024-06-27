package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/AriJaya07/go-rest-api/packages/utils"

	"github.com/AriJaya07/go-rest-api/packages/config"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	"github.com/AriJaya07/go-rest-api/packages/models/types"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the request (Auth header)
		tokenString := GetTokenFromRequest(r)
		// validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("failed to authentication token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("failed to authentication token")
			permissionDenied(w)
			return
		}
		// get the userId from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Println("failed to get user")
			permissionDenied(w)
			return
		}

		// call the handler func and continue to the endpoint
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"email":     email,
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := config.Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
