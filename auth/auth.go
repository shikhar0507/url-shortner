package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"url-shortner/utils"

	"github.com/gofrs/uuid"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
	"golang.org/x/crypto/bcrypt"
)

type AuthBody struct {
	Username string `json:username`
	Psswd    string `json:psswd`
}
type loggedIn struct {
	Authenticated bool
}
type Session struct {
	Username  string
	SessionId string
}

func Signup(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, []string{"http.MethodPost"})
	if optns {
		return
	}
	var signupBody AuthBody
	result := requestJSON.Decode(w, r, &signupBody)
	if result.Status != 200 {
		utils.SendResponse(w, result.Status, result)
		return
	}

	if signupBody.Username == "" {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Username cannot be empty"))
		return
	}

	if signupBody.Psswd == "" {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Password cannot be empty"))
		return
	}

	if userExists(signupBody.Username, signupBody.Psswd, db) {
		utils.SendResponse(w, http.StatusConflict, utils.SendErrorToClient("Account already exist"))
		return
	}

	err := createUser(signupBody.Username, signupBody.Psswd, db)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Error creating user"))
		return
	}
	successResponse, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Error creating response"))
		return
	}
	utils.SendResponse(w, http.StatusOK, successResponse)

}
func Signin(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, []string{"http.MethodPost"})
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
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Username cannot be empty"))
		return
	}
	if signinBody.Psswd == "" {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Password cannot be empty"))
		return
	}
	if userExists(signinBody.Username, signinBody.Psswd, db) == false {
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("Account does not exist"))
		return
	}

	var sessionId pgtypeuuid.UUID
	err := db.QueryRow(context.Background(), "insert into sessions(username) values($1) RETURNING sessionId", signinBody.Username).Scan(&sessionId)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}
	cookie := http.Cookie{Name: "sessionId", Value: sessionId.UUID.String(), HttpOnly: true, Expires: time.Now().AddDate(2022, 11, 22), Path: "/", Secure: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookie)
	resp := utils.Response{Status: http.StatusOK, Message: "Logged in"}
	utils.SendResponse(w, http.StatusOK, resp)

}

func CheckAuth(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	optns := utils.HandleCors(w, r, []string{"http.MethodGet"})
	if optns == true {
		return
	}
	_, err := GetSession(r, db)
	switch err {
	case nil:
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
	optns := utils.HandleCors(w, r, []string{http.MethodDelete})
	if optns == true {
		return
	}

	session, err := GetSession(r, db)

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Error logging out"))
		return
	}
	_, err = db.Exec(context.Background(), "delete from sessions where username=$1 AND sessionid=$2", session.Username, session.SessionId)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Error logging out"))
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
	success, err := json.Marshal(logout)
	utils.SendResponse(w, http.StatusOK, success)
}

func userExists(username string, psswd string, db *pgxpool.Pool) bool {
	var savedUsername string
	var savedHash string
	err := db.QueryRow(context.Background(), "select * from auth where username=$1", username).Scan(&savedUsername, &savedHash)
	if err != nil {
		return false
	}
	er := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(psswd))
	if er != nil {
		return false
	}
	return true

}

func createUser(username string, psswd string, db *pgxpool.Pool) error {
	hash, err := GeneratePsswdHash(psswd)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), "insert into auth values($1,$2)", username, hash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func GeneratePsswdHash(psswd string) ([]byte, error) {

	bytePsswd := []byte(psswd)
	hash, err := bcrypt.GenerateFromPassword(bytePsswd, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil

}

/**
  Returns username,sessionId and error
*/
func GetSession(r *http.Request, db *pgxpool.Pool) (Session, error) {
	sessionCookie := getSessionCookie(r)
	var sessionId pgtypeuuid.UUID
	var username string
	session := Session{SessionId: "", Username: ""}
	if sessionCookie == "" {
		return session, pgx.ErrNoRows
	}
	sessionUUID, err := uuid.FromString(sessionCookie)
	if err != nil {
		return session, err
	}

	err = db.QueryRow(context.Background(), "select * from sessions where sessionid=$1", sessionUUID).Scan(&username, &sessionId)
	session.SessionId = sessionId.UUID.String()
	session.Username = username
	if err != nil {
		return session, err
	}
	return session, nil

}

func getSessionCookie(r *http.Request) string {
	coo, err := r.Cookie("sessionId")
	if err != nil {
		return ""
	}
	return coo.Value
}
