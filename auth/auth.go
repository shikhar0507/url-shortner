package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
	"golang.org/x/crypto/bcrypt"
)

type AuthBody struct {
	Username string
	Psswd    string
}

func Signup(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {

	var signupBody AuthBody
	//var result requestDecoder.Result

	result := requestJSON.Decode(w, r, &signupBody)
	if signupBody.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	if signupBody.Psswd == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}
	fmt.Println("check for  user")
	if userExists(signupBody.Username, signupBody.Psswd, db) {

		result.Message = "Account already  exist"
		result.Status = http.StatusConflict
		r, err := json.Marshal(result)
		w.WriteHeader(http.StatusConflict)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(r))
		return
	}
	if result.Status == 200 {
		// if request is successfully parsed
		err := createUser(signupBody.Username, signupBody.Psswd, db)
		if err != nil {
			http.Error(w, "Error creating success response", http.StatusInternalServerError)
			return
		}
		successResponse, marshalErr := json.Marshal(result)
		if marshalErr != nil {
			http.Error(w, "Error creating success response", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(successResponse))

		return
	}
	fmt.Println("error", result)
	http.Error(w, result.Message, result.Status)
}
func Signin(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	var signinBody AuthBody
	result := requestJSON.Decode(w, r, &signinBody)
	if signinBody.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}
	if signinBody.Psswd == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}
	if result.Status != 200 {
		http.Error(w, result.Message, result.Status)
		return
	}
	if userExists(signinBody.Username, signinBody.Psswd, db) == false {
		fmt.Println("User does not exist")
		result.Message = "Account does not exist"
		result.Status = http.StatusNotFound
		r, err := json.Marshal(result)
		w.WriteHeader(http.StatusNotFound)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(r))
		return
	}
	uuid, err := generateUUID()
	if err != nil {
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(context.Background(), "insert into sessions values($1,$2)", signinBody.Username, uuid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
	//fmt.Print("Failed to create user")
	cookie := http.Cookie{Name: "session", Value: uuid}
	http.SetCookie(w, &cookie)
	result.Message = "Login successfull"
	success, er := json.Marshal(result)
	if er != nil {
		log.Fatal(er)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(success))
}
func Logout(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	uname, uid, err := getSession(r, db)
	fmt.Println("logout")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(context.Background(), "delete from sessions where username=$1 AND sessionid=$2", uname, uid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}
	coo, err := r.Cookie("session")
	coo.MaxAge = -1
	coo.Value = ""
	http.SetCookie(w, coo)
	type Logout struct {
		Message string
		Status  int
	}
	var logout Logout
	success, err := json.Marshal(&logout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(success))
}

func userExists(username string, psswd string, db *pgxpool.Pool) bool {
	var savedUsername string
	var savedHash string
	err := db.QueryRow(context.Background(), "select * from auth where username=$1", username).Scan(&savedUsername, &savedHash)
	if err != nil {
		return false
	}
	fmt.Println(savedUsername, savedHash)

	er := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(psswd))
	if er != nil {
		return false
	}
	fmt.Println("isSame")
	return true

}

func createUser(username string, psswd string, db *pgxpool.Pool) error {
	hash, err := generatePsswdHash(psswd)
	if err != nil {
		return err
	}
	fmt.Println("hash", string(hash))
	_, err = db.Exec(context.Background(), "insert into auth values($1,$2)", username, hash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func generatePsswdHash(psswd string) ([]byte, error) {

	bytePsswd := []byte(psswd)
	hash, err := bcrypt.GenerateFromPassword(bytePsswd, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return hash, nil

}

func getSession(r *http.Request, db *pgxpool.Pool) (string, string, error) {
	sessionId := getSessionCookie(r)
	var sessionResult string
	var sessionUsername string
	err := db.QueryRow(context.Background(), "select * from sessions where sessionid=$1", sessionId).Scan(&sessionUsername, &sessionResult)
	if err != nil {
		return "", "", err
	}
	return sessionUsername, sessionResult, nil
}

func getSessionCookie(r *http.Request) string {
	coo, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	return coo.Value
}
