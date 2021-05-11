package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"url-shortner/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
	"golang.org/x/crypto/bcrypt"
)

type AuthBody struct {
	Username string
	Psswd    string
}
type loggedIn struct {
	Authenticated bool
}

func Signup(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, http.MethodPost)
	if optns == true {
		return
	}
	var signupBody AuthBody
	//var result requestDecoder.Result

	result := requestJSON.Decode(w, r, &signupBody)
	if result.Status != 200 {
		utils.SendResponse(w, result.Status, result)
		return
	}

	if signupBody.Username == "" {
		resp := utils.Response{Status: http.StatusBadRequest, Message: "Username cannot be empty"}
		utils.SendResponse(w, http.StatusBadRequest, resp)
		return
	}

	if signupBody.Psswd == "" {
		resp := utils.Response{Status: http.StatusBadRequest, Message: "Password cannot be empty"}
		utils.SendResponse(w, http.StatusBadRequest, resp)
		return
	}
	fmt.Println("check for  user")
	if userExists(signupBody.Username, signupBody.Psswd, db) {
		resp := utils.Response{Message: "Account-already-exist", Status: http.StatusConflict}
		utils.SendResponse(w, http.StatusConflict, resp)
		return
	}

	err := createUser(signupBody.Username, signupBody.Psswd, db)
	if err != nil {
		resp := utils.Response{Status: http.StatusInternalServerError, Message: "Error creating user"}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
		return
	}
	successResponse, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		resp := utils.Response{Status: http.StatusInternalServerError, Message: "Error creating response"}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
		return
	}
	fmt.Println("account created", successResponse)
	utils.SendResponse(w, http.StatusOK, successResponse)

}
func Signin(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, http.MethodPost)
	if optns == true {
		return
	}
	var signinBody AuthBody
	result := requestJSON.Decode(w, r, &signinBody)
	if result.Status != 200 {
		utils.SendResponse(w, result.Status, result)
		return
	}

	if signinBody.Username == "" {
		resp := utils.Response{Status: http.StatusBadRequest, Message: "Username cannot be empty"}
		utils.SendResponse(w, http.StatusBadRequest, resp)
		return
	}
	if signinBody.Psswd == "" {
		resp := utils.Response{Status: http.StatusBadRequest, Message: "Password cannot be empty"}
		utils.SendResponse(w, http.StatusBadRequest, resp)
		return
	}
	if userExists(signinBody.Username, signinBody.Psswd, db) == false {
		fmt.Println("User does not exist")
		resp := utils.Response{Status: http.StatusNotFound, Message: "Account does not exist"}
		utils.SendResponse(w, http.StatusNotFound, resp)
		return
	}
	uuid, err := generateUUID()
	if err != nil {
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusNotFound, resp)
		return
	}
	_, err = db.Exec(context.Background(), "insert into sessions values($1,$2)", signinBody.Username, uuid)
	if err != nil {
		fmt.Println(err)
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusNotFound, resp)
		return
	}
	//fmt.Print("Failed to create user")
	cookie := http.Cookie{Name: "sessionId", Value: uuid, HttpOnly: true, Expires: time.Now().AddDate(2022, 11, 22), Path: "/", Secure: true, SameSite: http.SameSiteNoneMode}

	http.SetCookie(w, &cookie)
	resp := utils.Response{Status: http.StatusOK, Message: "Logged in"}
	utils.SendResponse(w, http.StatusOK, resp)

}

func CheckAuth(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, http.MethodGet)
	if optns == true {
		return
	}
	fmt.Println(r.Method)
	_, uid, err := GetSession(r, db)
	switch err {
	case nil:
		fmt.Println("found user", uid)
		resp := loggedIn{Authenticated: true}
		utils.SendResponse(w, http.StatusOK, resp)
	case pgx.ErrNoRows:
		resp := loggedIn{Authenticated: false}
		utils.SendResponse(w, http.StatusOK, resp)
	default:
		resp := loggedIn{Authenticated: false}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
	}
}

func Logout(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, http.MethodDelete)
	if optns == true {
		return
	}

	uname, uid, err := GetSession(r, db)
	fmt.Println("logout")
	if err != nil {
		fmt.Println(err)
		resp := utils.Response{Status: http.StatusInternalServerError, Message: "Error logging out"}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
		return
	}
	_, err = db.Exec(context.Background(), "delete from sessions where username=$1 AND sessionid=$2", uname, uid)
	if err != nil {
		resp := utils.Response{Status: http.StatusInternalServerError, Message: "Error logging out"}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
		return
	}
	coo, err := r.Cookie("sessionId")
	coo.Path = "/"
	coo.Secure = true
	coo.HttpOnly = true
	coo.Expires = time.Now().AddDate(-1, -1, -1)
	coo.Value = ""
	http.SetCookie(w, coo)
	type Logout struct {
		Message string
		Status  int
	}
	var logout Logout
	success, err := json.Marshal(&logout)
	utils.SendResponse(w, http.StatusOK, success)
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

func GetSession(r *http.Request, db *pgxpool.Pool) (string, string, error) {
	sessionId := getSessionCookie(r)
	fmt.Println("sid", sessionId)
	var sessionResult string
	var sessionUsername string
	err := db.QueryRow(context.Background(), "select * from sessions where sessionid=$1", sessionId).Scan(&sessionUsername, &sessionResult)
	if err != nil {
		return "", "", err
	}
	return sessionUsername, sessionResult, nil
}

func getSessionCookie(r *http.Request) string {
	coo, err := r.Cookie("sessionId")
	if err != nil {
		return ""
	}
	return coo.Value
}
