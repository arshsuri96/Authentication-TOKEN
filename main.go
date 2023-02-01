package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type user struct {
	email        string
	username     string
	passwordhash string
	fullname     string
	createDate   string
	role         int
}

func GetUserObject(email string) (user, bool) {
	for _, user := range userList {
		if user.email == email {
			return user, true
		}
	}
	return user{}, false
}

func (u *user) ValidatePassowrdHash(pswdhash string) bool {
	return u.passwordhash == pswdhash
}

func AddUser(email string, username string, passwordhash string, fullname string, role int) {
	newUser := user{
		email:        email,
		passwordhash: passwordhash,
		username:     username,
		fullname:     fullname,
		role:         role,
	}

	for _, ele := range userList {
		if ele.email == email || ele.username == username {
			return false
		}
	}
	userList = append(userList, newUser)
}

func generateToken(header string, payload map[string]string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	payloadstr, error := json.Marshal(payload)
	if error != nil {
		fmt.Println("error generating token")
		return string(payloadstr), error
	}
	payload64 := base64.StdEncoding.EncodeToString(payloadstr)

	message := header64 + "." + payload64

	unsignedstr := header + string(payloadstr)

	h.Write([]byte(unsignedstr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	tokenStr := message + "." + signature
	return tokenStr, nil

}

func ValdiateToken(token string, secret string) (bool, error) {
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 3 {
		return false, nil
	}
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		return false, nil
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		return false, nil
	}
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	if signature != splitToken[2] {
		return false, nil
	}

	return true, nil

}

// SIGNUP HANDLER

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Email"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email Missing"))
		return
	}

	if _, ok := r.Header["PasswordHash"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password hash issue"))
		return
	}

	if _, ok := r.Header["Username"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username missing"))
		return
	}
	if _, ok := r.Header["Fullname"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fullname missing"))
		return
	}
	check := data.AddUserObject(r.Header["Eamil"][0], r.Header["Username"][0], r.Header["Passowrdhas"][0], r.Header["Fullname"][0], 0)

	if !check {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email or username already exists"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Created"))
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Email"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email Missing"))
		return
	}
	if _, ok := r.Header["Passowrdhash"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passowrdhash Missing"))
		return
	}

	valid, err := validateUser(r.Header["Email"][0], r.Header["PasswordHash"][0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not exist"))
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect Password"))
	}

	tokenString, err := getSignedToken()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func getSignedToken() {

}

func validateUser(email string, passwordHash string) (bool, error) {
	usr, exists := data.GetUserObject(email)
	if !exists {
		return false, errors.New("user does not exist")
	}
	passwordCheck := usr.ValidatePassowrdHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}
