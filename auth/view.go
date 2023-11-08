package auth

import (
	"database/sql"
	"fmt"
	"io"
	"legitt-backend/helper"
	"legitt-backend/model"
	"net/http"

	"github.com/rs/zerolog/log"
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

// func GeneratePublicContractWithGptApiHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}

// 	var data chatToContractRequest
// 	helper.DecodeReqJsonToMap(r.Body, &data)

// 	res, err := handleGeneratePublicContractWithGptApi(data)

// 	if err != nil {
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 	}

// 	fmt.Fprint(w, res)
// }

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

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	var data userCredentials
// 	helper.DecodeReqJsonToMap(r.Body, &data)

// 	userExists, err := CheckUserExists(data.Email)
// 	if err != nil {
// 		http.Error(w, "", http.StatusInternalServerError)
// 		return
// 	}
// 	if userExists {
// 		http.Error(w, "User already exists try to login", http.StatusConflict)
// 		return
// 	}
// 	if data.Email == "" || data.Password == "" {
// 		http.Error(w, "Bad Credentials", http.StatusBadRequest)
// 		return
// 	}

// 	err = HandleRegister(w, data)

// 	if err != nil {
// 		fmt.Print("error", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	helper.SendJsonResponse(w, "ok")
// }

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

// func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	err := HandleCheckSession(r)

// 	if err != nil {
// 		http.Error(w, "session invalid", http.StatusUnauthorized)
// 		return
// 	}

// 	helper.SendJsonResponse(w, "session valid")
// }

// func LogoutHandler(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("sid")

// 	if err == nil {
// 		HandleLogout(w, *cookie)
// 	}

// 	helper.SendJsonResponse(w, "ok")
// }

// Generates a secret key with appropriate permissions for a domain
// type RegisterAppPayload struct {
// 	OwnerEmail  string   `json:"ownerEmail"`
// 	Domain      string   `json:"domain"`
// 	Permissions []string `json:"permissions"`
// }Domainrror(w, "", http.StatusBadRequest)
// 		return
// 	}

// 	err = HandleRegisterApp(payload.OwnerEmail, payload.Domain, payload.Permissions)
// 	if err != nil {
// 		http.Error(w, "", http.StatusInternalServerError)
// 	}

// 	helper.SendJsonResponse(w, "apps registered")
// }

// func CadAnalyzeHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload
// }

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// IsGoogleSignIn  string `json:"isGoogleSignIn"`
	// RedirectSite    string `json:"redirectSite"`
	CaptchaResponse string `json:"captcha"`
}

func LoginView(w http.ResponseWriter, r *http.Request) {
	var data loginPayload
	err := helper.DecodeReqJsonToMap(r.Body, &data)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = LoginHandler(w, data)
	if err != nil {
		print(err)
		return
	}

	helper.SendJsonResponse(w, "ok")
}

func RegisterView(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	err := helper.DecodeReqJsonToMap(r.Body, &payload)

	if err != nil || payload.Captcha == "" || payload.Email == "" || payload.Password == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = HandleRegister(w, payload)

	if err != nil {
		fmt.Print("error", err.Error())
		return
	}

	helper.SendJsonResponse(w, "ok")
}
