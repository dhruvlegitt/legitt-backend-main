package auth

import (
	"legitt-backend/model"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandleGenerateSession(w http.ResponseWriter, email string) error {
	expirationDuration := time.Hour * 10 * 24
	expirationTimestamp := time.Now().Add(expirationDuration).Unix()

	sessionId := uuid.New().String()

	redisSession := model.RedisSession{Email: email, Location: "UAE", Exp: expirationTimestamp}

	err := model.SetRedisValue(sessionId, redisSession, expirationDuration)

	if err != nil {
		return err
	}

	columns := []string{"session_id", "email", "exp"}
	values := [][]interface{}{{sessionId, email, expirationTimestamp}}

	_, err = model.InsertIntoDb("sessions", columns, values)

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    sessionId,
		MaxAge:   int(expirationDuration),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	return nil
}

func HandleCheckPassword(w http.ResponseWriter, email string, password string) error {
	var hashedPassword string

	whereFields := make(map[string]interface{})
	whereFields["email"] = email

	row := model.GetRowFromDb("users", whereFields, "password")
	err := row.Scan(&hashedPassword)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}
