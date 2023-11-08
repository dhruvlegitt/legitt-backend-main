package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"legitt-backend/model"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type CaptchaPayload struct {
	secret          string
	captchaResponse string
}

type CaptchaResponse struct {
	success bool
}

func LoginHandler(w http.ResponseWriter, data loginPayload) error {
	fmt.Print("here")
	captchaVerifyPayload := CaptchaPayload{secret: captchaSecretKey, captchaResponse: data.CaptchaResponse}
	captchaVerifyPayloadString, err := json.Marshal(captchaVerifyPayload)

	res, err := http.Post(captchaVerifyUrl, "application/json", bytes.NewBuffer(captchaVerifyPayloadString))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var captchaRes CaptchaResponse
	err = json.NewDecoder(res.Body).Decode(&captchaRes)

	fmt.Print("Captcha", captchaRes.success)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return err
	} else if captchaRes.success == false {
		http.Error(w, "Invalid Captcha", http.StatusBadRequest)
		return errors.New("Invalid Captcha")
	}

	err = HandleCheckPassword(w, data.Email, data.Password)
	if err != nil {
		http.Error(w, "Invalid username/password", http.StatusBadRequest)
		return err
	}

	err = HandleGenerateSession(w, data.Email)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return err
	}

	return nil
}

func HandleRegister(w http.ResponseWriter, payload RegisterPayload) error {
	userExists, err := CheckUserExists(payload.Email)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return err
	}
	if userExists {
		http.Error(w, "User already exists try to login", http.StatusConflict)
		return errors.New("User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return err
	}

	columns := []string{"email", "password"}
	values := [][]interface{}{{payload.Email, string(hashedPassword[:])}}
	_, err = model.InsertIntoDb("users", columns, values)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return err
	}

	return nil
}
