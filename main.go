package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"strings"
	"time"
	"url-shortner/auth"
	"url-shortner/utils"

	"github.com/avct/uasurfer"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
)

var db *pgxpool.Pool

type SuccesRes struct {
	Status int    `json:"status"`
	Url    string `json:"url"`
}
type stop struct {
	error
}
type CampaignStrc struct {
	Source   string
	Medium   string
	Campaign string
	ReqURL
}
type ReqURL struct {
	Url string
}

func main() {

	dbpool, err := pgxpool.Connect(context.Background(), "postgres://xanadu:xanadu@localhost:5432/tracker")
	if err != nil {
		log.Fatal(err)
	}
	db = dbpool
	defer db.Close()

	// favicon
	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		return
	})
	// http.Handle("/public/build/", http.StripPrefix("/public/build/", http.FileServer(http.Dir("public/build"))))

	//redirect
	http.HandleFunc("/", handleRedirect)

	//auth
	http.HandleFunc("/signup-user", func(rw http.ResponseWriter, r *http.Request) {
		auth.Signup(rw, r, db)
	})
	http.HandleFunc("/login-user", func(rw http.ResponseWriter, r *http.Request) {
		auth.Signin(rw, r, db)
	})
	http.HandleFunc("/logout", func(rw http.ResponseWriter, r *http.Request) {
		auth.Logout(rw, r, db)
	})

	http.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
		auth.CheckAuth(rw, r, db)
	})

	http.HandleFunc("/campaign", handleCampaign)
	// url-shortner api
	http.HandleFunc("/shorten", handleShortner)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request for", r.URL.String())
	id := strings.Split(r.URL.Path, "/")[1]
	originalUrl, err := getRedirectUrl(id)
	switch err {
	case nil:
		http.Redirect(w, r, originalUrl, http.StatusPermanentRedirect)
	case pgx.ErrNoRows:
		fmt.Println("serving", r.URL.String())
		fs := http.FileServer(http.Dir("public/build/"))
		fs.ServeHTTP(w, r)
		// http.Redirect(w, r, "/public/build/", http.StatusTemporaryRedirect)
	default:
		fmt.Println("query error", err)
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusInternalServerError, resp)

	}

}

func getRedirectUrl(id string) (string, error) {
	var queryId string
	var originalUrl string
	err := db.QueryRow(context.Background(), "select id from urls where id=$1", id).Scan(&queryId, &originalUrl)
	if err != nil {
		return "", err

	}
	return originalUrl, nil
}

func handleCampaign(w http.ResponseWriter, r *http.Request) {
	optns := utils.HandleCors(w, r, http.MethodPost)
	if optns == true {
		return
	}

	var bod CampaignStrc
	result := requestJSON.Decode(w, r, &bod)
	if result.Status != 200 {
		utils.SendResponse(w, result.Status, result)
		return
	}
	if bod.Campaign == "" {
		resp := utils.Response{Status: http.StatusBadRequest, Message: "Campaign name cannot be empty"}
		utils.SendResponse(w, http.StatusBadRequest, resp)
		return
	}
	_, err := setId(r)

	if err != nil {
		fmt.Println(err)
		if err.Error() == "failed to assign a unique value" {
			resp := utils.Response{Status: http.StatusInternalServerError, Message: err.Error()}
			utils.SendResponse(w, http.StatusInternalServerError, resp)
			return
		}
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
		return
	}
	resp := utils.Response{Status: http.StatusOK, Message: "campaign created"}
	utils.SendResponse(w, http.StatusOK, resp)
}

func handleShortner(w http.ResponseWriter, r *http.Request) {

	optns := utils.HandleCors(w, r, http.MethodPost)
	if optns == true {
		return
	}

	var reqURL ReqURL
	result := requestJSON.Decode(w, r, &reqURL)
	fmt.Println(result)
	if result.Status != 200 {

		utils.SendResponse(w, result.Status, result)
		return
	}

	id, err := setId(r)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "failed to assign a unique value" {
			resp := utils.Response{Status: http.StatusInternalServerError, Message: err.Error()}
			utils.SendResponse(w, http.StatusInternalServerError, resp)
			return
		}
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusInternalServerError, resp)
	}
	fmt.Println("used id", id)

	succ := SuccesRes{Status: 200, Url: "http://localhost:8080/" + id}
	utils.SendResponse(w, 200, succ)

}

func parseQuery(r *http.Request) CampaignStrc {
	query := r.URL.Query()
	campignParams := CampaignStrc{Campaign: query.Get("campaign"), Source: query.Get("source"), Medium: query.Get("medium")}
	return campignParams

}

func setId(r *http.Request) (string, error) {
	value := createId()
	//value := "RsWxP"
	mainErr := retry(100, 1000, func() error {
		parsedQuery := parseQuery(r)
		ua := uasurfer.Parse(r.UserAgent())
		username, _, err := auth.GetSession(r, db)
		_, err = db.Exec(context.Background(), "insert into urls values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", value, r.URL.String(), parsedQuery.Campaign, parsedQuery.Source, parsedQuery.Medium, ua.Browser.Name, ua.Browser.Version, ua.OS.Platform, ua.OS.Name, ua.OS.Version, ua.DeviceType.String(), username)

		if err == nil {
			return nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				fmt.Println("creating a new id")
				value = createId()
				return err
			}
		}
		return err
	})

	if mainErr != nil {
		return "", mainErr
	}
	return value, nil

}

func retry(count int, sleep time.Duration, f func() error) error {
	err := f()
	if err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}
		count--
		if count > 0 {
			time.Sleep(sleep)
			return retry(count, 1*sleep, f)
		}
		return err
	}
	return nil
}

func createId() string {
	letterString := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := ""

	for i := 0; i < 6; i++ {
		randStr := letterString[rand2.Intn(len(letterString))]
		result = result + string(randStr)
	}
	return result
}
