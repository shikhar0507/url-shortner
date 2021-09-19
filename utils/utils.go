package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func init() {
	fmt.Println("queue initialized")
}

type ClientError struct {
	Message string
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccesRes struct {
	Status int    `json:"status"`
	Url    string `json:"url"`
}

func SendErrorToClient(message string) ClientError {
	return ClientError{Message: message}
}

func HandleCors(w http.ResponseWriter, r *http.Request, methods []string) bool {
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ","))
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)

		return true
	}
	return false
}

func SendResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	j, err := json.Marshal(body)
	fmt.Println(string(j))
	if err != nil {
		errResp := Response{Message: "Error", Status: http.StatusInternalServerError}
		finalMsg, err := json.Marshal(errResp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(finalMsg))
		return
	}
	fmt.Fprintf(w, string(j))

}
