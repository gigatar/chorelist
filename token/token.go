package token

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const tokenDuration = 5 * 60                      // In Seconds
const tokenPassword = "Y[6=fJMxa1PRxnW4hxEg@5.Pu" // CHANGE ME FOR PRODUCTION!!!!!

// JWTToken data type.
type JWTToken struct {
	Token  string `json:"token"`
	Expire int64  `json:"expires"`
}

// CustomClaims data type.
type CustomClaims struct {
	UserID   string             `json:"UserID"`
	FamilyID string             `json:"FamilyID"`
	TTL      int                `json:"TTL"`
	Standard jwt.StandardClaims `json:"Claims"`
}

// Valid is implementing the Valid interface from JWT to prevent an error
func (c CustomClaims) Valid() error {
	return nil
}

// CreateJWT Creates and returns a new JWT.
func (v *JWTToken) CreateJWT(remoteHost, userID, familyID string) (string, error) {
	// Get IP from Remote Addr:
	serverIP, _, err := net.SplitHostPort(remoteHost)
	if err != nil {
		return "", err
	}

	// Expire
	expire := time.Now().Add(time.Duration(time.Minute) * tokenDuration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		userID,
		familyID,
		tokenDuration * 60,
		jwt.StandardClaims{
			ExpiresAt: expire,
			Audience:  serverIP,
		},
	})

	tokenString, err := token.SignedString([]byte(tokenPassword))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUser returns the user from the JWT Token.
func (v JWTToken) GetUser(inputToken string) (string, error) {
	bearerToken := strings.Split(inputToken, " ")
	if len(bearerToken) != 2 {
		return "", errors.New("Invalid token")
	}

	token, err := jwt.ParseWithClaims(bearerToken[1], CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenPassword), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return "", errors.New("Invalid token")
		} else if strings.Contains(err.Error(), "signature is invalid") {

			return "", errors.New("Invalid Token")
		}
		return "", err
	}

	if claims, ok := token.Claims.(CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", nil
}

// // RefreshToken regenerates our token.
// func (v *JWTToken) RefreshToken(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	userID, err := v.GetUser(r.Header.Get("authorization"))
// 	if err != nil {
// 		if strings.Contains(err.Error(), "Invalid token") {
// 			w.WriteHeader(http.StatusUnauthorized)
// 		}
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	token, err := v.CreateJWT(r.RemoteAddr, userID)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Authorization", token)
// 	w.WriteHeader(http.StatusNoContent)

// }

// ValidateMiddleware is how we ensure that only authorized persons are able to access endpoints.
func (v *JWTToken) ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenPassword), nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				w.WriteHeader(http.StatusUnauthorized)
				log.Println("Expired Token")
				return
			} else if strings.Contains(err.Error(), "signature is invalid") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		// Get IP from Remote Addr:
		serverIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			if serverIP != claims.Standard.Audience {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}
