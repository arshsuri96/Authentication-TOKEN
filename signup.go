package main

import (
	"net/http"

	"github.com/shadowshot-x/micro-product-go/authservice/data"
)

func SignupHandler(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("passwordhash missing"))
	}
	if _, ok := r.Header["Fullname"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Fullname Missing"))
	}
	if _, ok := r.Header["Username"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username Missing"))
	}

	check := data.AddUserObject(r.Header["Email"][0], r.Header["Passwordhash"][0], r.Header["Fullname"][0],
		r.Header["Username"][0], 0)

	if !check {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Email or Username already exists"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User created"))

}
