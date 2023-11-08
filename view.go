package main

import (
	"database/sql"
	"fmt"
	"io"
	"legitt-backend/helper"
	"legitt-backend/model"
	"net/http"
)

type userCredentials struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	isGoogleSignIn string `json:isGoogleSignIn`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to my website")
}

type chatToContractRequest struct {
	Title string
}

func GeneratePublicContractWithGptApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var data chatToContractRequest
	helper.DecodeReqJsonToMap(r.Body, &data)

	res, err := handleGeneratePublicContractWithGptApi(data)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	fmt.Fprint(w, res)
}

func CheckUserExists(email string) (bool, error) {
	whereFields := map[string]interface{}{"email": email}
	rowCursor := model.GetRowFromDb("users", whereFields, "id")
	var id int
	err := rowCursor.Scan(&id)

	if err == nil {
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func RegisterView(w http.ResponseWriter, r *http.Request) {
	var data userCredentials
	helper.DecodeReqJsonToMap(r.Body, &data)

	userExists, err := CheckUserExists(data.Email)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if userExists {
		http.Error(w, "User already exists try to login", http.StatusConflict)
		return
	}
	if data.Email == "" || data.Password == "" {
		http.Error(w, "Bad Credentials", http.StatusBadRequest)
		return
	}

	// err = HandleRegister(w, data)

	if err != nil {
		fmt.Print("error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.SendJsonResponse(w, "ok")
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	err := HandleCheckSession(r)

	if err != nil {
		http.Error(w, "session invalid", http.StatusUnauthorized)
		return
	}

	helper.SendJsonResponse(w, "session valid")
}

// Generates a secret key with appropriate permissions for a domain
type RegisterAppPayload struct {
	OwnerEmail  string   `json:"ownerEmail"`
	Domain      string   `json:"domain"`
	Permissions []string `json:"permissions"`
}

func RegisterAppHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterAppPayload

	err := helper.DecodeReqJsonToMap(r.Body, &payload)
	fmt.Print(err)
	fmt.Print(payload)
	if err != nil || payload.Domain == "" || len(payload.Permissions) == 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = HandleRegisterApp(payload.OwnerEmail, payload.Domain, payload.Permissions)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	helper.SendJsonResponse(w, "apps registered")
}
