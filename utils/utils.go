package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleCors(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Methods", method)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func SendResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
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
