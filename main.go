package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shadowshot-x/micro-product-go/authservice"
)

func main() {
	mainRouter := mux.NewRouter()
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authservice.SignupHandler)

	authRouter.HandleFunc("/signin", authservice.SignInHandler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mainRouter,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error Booting the Server")
	}
}
