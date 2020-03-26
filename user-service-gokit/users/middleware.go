package users

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const tokenDuration = 5 * 60                      // In Seconds
const tokenPassword = "Y[6=fJMxa1PRxnW4hxEg@5.Pu" // CHANGE ME FOR PRODUCTION!!!!!

// CustomClaims data type.
type CustomClaims struct {
	UserID             string `json:"UserID"`
	FamilyID           string `json:"FamilyID"`
	TTL                int    `json:"TTL"`
	jwt.StandardClaims `json:"Claims"`
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func listBackendMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ifaces, err := net.Interfaces()
		if err == nil {
			// handle err
			for _, i := range ifaces {
				addrs, err := i.Addrs()
				if err != nil {
					break
				}

				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if !ip.IsLoopback() {
						w.Header().Add("Backend-Server", ip.String())
						break
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateMiddleware is how we ensure that only authorized persons are able to access endpoints.
func validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusUnauthorized)
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
			if serverIP != claims.Audience {

				log.Println("Mismatch token")
				// w.WriteHeader(http.StatusUnauthorized)
				// return
			}
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}
