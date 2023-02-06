package authservice

import (
	"arshsuri96/AUTHENTICATION_MICRO/authservice/jwt"
	"net/http"
)

func tokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Token"]; !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token missing"))
			return
		}
		token := r.Header["token"][0]
		check, err := jwt.ValidateToken(token, "Secure_Random_String")

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Token Validation missing"))
			return
		}

		if !check {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Invalid"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Authorized Token"))
	})
}
