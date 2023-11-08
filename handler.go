package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"legitt-backend/helper"
	"legitt-backend/model"
	"net/http"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

type checkSessionRequest struct {
	Email string `json:"email"`
}

// TODO fix code structure
func handleGeneratePublicContractWithGptApi(data chatToContractRequest) (any, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	res, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   2000,
			Temperature: 0.6,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: data.Title,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	content := res.Choices[0].Message.Content

	return content, nil
}

func HandleCheckSession(r *http.Request) error {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return err
	}

	var payload checkSessionRequest
	err = helper.DecodeReqJsonToMap(r.Body, &payload)

	if err != nil {
		return err
	}

	res, err := model.GetRedisValue(cookie.Value)

	if err != nil {
		fmt.Print("Session not found")
		return err
	}

	var redisSessionData model.RedisSession
	json.Unmarshal([]byte(res), &redisSessionData)

	if redisSessionData.Email == payload.Email &&
		helper.CompareUnixTimeStamp(redisSessionData.Exp, time.Now()) == 1 {
		return nil
	} else {
		return errors.New("Session Invalid")
	}
}

func InvalidateSession(sid string) error {
	model.DeleteRedisValue(sid)

	timeToUpdate := time.Now().UTC()
	whereClause := make(map[string]interface{})
	whereClause["session_id"] = sid

	model.UpdateRowsInDb("sessions", []string{"session_end_time"}, []interface{}{timeToUpdate}, whereClause)

	return nil
}

func HandleRegisterApp(email string, domain string, permissions []string) error {
	secretKey, err := helper.GenerateRandomString(32)
	if err != nil {
		return err
	}

	columns := []string{"secret_key", "owner_email", "domain"}
	values := [][]interface{}{{secretKey, email, domain}}
	_, err = model.InsertIntoDb("api_key_info", columns, values)

	if err != nil {
		return err
	}

	columns = []string{"secret_key", "permitted_end_point"}
	values = [][]interface{}{}

	for _, permission := range permissions {
		values = append(values, []interface{}{secretKey, permission})
	}

	_, err = model.InsertIntoDb("api_key_permissions", columns, values)
	return err
}
