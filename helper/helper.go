package helper

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

type HttpResponse struct {
	Data interface{} `json:"data"`
}

func DecodeReqJsonToMap(rawBody io.ReadCloser, data interface{}) error {
	decoder := json.NewDecoder(rawBody)
	return decoder.Decode(&data)
}

func PrintHeaders(headerMap http.Header) {
	for k, v := range headerMap {
		fmt.Printf("%s : ", k)
		for _, val := range v {
			fmt.Println(val)
		}
	}
}

func SendJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := HttpResponse{Data: data}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func CompareUnixTimeStamp(unixStartTime int64, timeToCheck time.Time) int {
	timeFromUnix := time.Unix(unixStartTime, 0)
	if timeFromUnix.Before(timeToCheck) {
		return -1
	} else if timeFromUnix.After(timeToCheck) {
		return 1
	} else {
		return 0
	}
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
