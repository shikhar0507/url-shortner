package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortner/api"
	"url-shortner/auth"
	"url-shortner/utils"

	"github.com/avct/uasurfer"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

type Router struct {
	api Api
}

type Api struct {
	Links link
}

type link struct {
}

func getSegments(r *http.Request) (string, string) {
	split := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	head := split[0]
	rest := "/" + strings.Join(split[1:], "/")
	return head, rest
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = getSegments(r)
	switch head {
	case "":
		handleRedirect(w, r)
	case "api":
		Api.ServeHTTP(w http.ResponseWriter, r *http.Request)
		router.api.ServeHTTP(w, r)
	case "stat":
		fmt.Println("stat")
	}
}

func (api Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("read api")
}

func main() {

	fmt.Println("starting server")
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

	//http.HandleFunc("/campaign", handleCampaign)
	// url-shortner api
	http.HandleFunc("/shorten", func(rw http.ResponseWriter, r *http.Request) {
		//handleShortner
	})

	http.HandleFunc("/api/v1/campaigns", func(rw http.ResponseWriter, r *http.Request) {
		api.HandleCampaigns(rw, r, db)
	})
	http.HandleFunc("/api/v1/links", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("calling links")
		api.HandleLinks(rw, r, db)
	})
	router := Router{}
	err = http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request for", r.URL.String())
	id := strings.Split(r.URL.Path, "/")[1]
	result, err := getRedirectUrl(id)
	fmt.Println(err)
	switch err {
	case nil:
		fmt.Println("redirecting to", result.url)
		http.Redirect(w, r, result.url, http.StatusPermanentRedirect)
		go updateLogs(r, result)
		break
	case pgx.ErrNoRows:
		fmt.Println("serving", r.URL.String())
		fs := http.FileServer(http.Dir("public/build/"))
		fs.ServeHTTP(w, r)
		break
	default:
		fmt.Println("query error", err)
		resp := utils.Response{Status: http.StatusInternalServerError}
		utils.SendResponse(w, http.StatusInternalServerError, resp)

	}

}

func updateLogs(r *http.Request, result storedUrl) {

	browser_info := uasurfer.Parse(r.UserAgent())
	u, err := url.Parse(result.url)
	if err != nil {
		fmt.Println("unable to parse short url long url")
		return
	}
	query := u.Query()
	campaign, medium, source := query.Get("camapgin"), query.Get("meidum"), query.Get("source")

	parsedIP, referer := net.ParseIP(r.Header.Get("ip")), r.Header.Get("referer")
	var ip string
	if parsedIP == nil {
		ip = "0.0.0.0"
	} else {
		ip = parsedIP.String()
	}
	fmt.Println(browser_info.OS.Platform)
	_, err = db.Exec(context.Background(), "insert into logs(url,username,os,browser,device_type,created_on,ip,referer,campaign,medium,source,id)values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", result.url, result.username, browser_info.OS.Name.String(), browser_info.Browser.Name.String(), browser_info.DeviceType.String(), time.Now(), ip, referer, campaign, medium, source, result.id)
	if err != nil {
		fmt.Println(err)
		var pgError *pgconn.PgError
		errors.As(err, &pgError)
		fmt.Printf("%s at %d col %s %s %s", pgError.Message, pgError.Line, pgError.ColumnName, pgError.Detail, pgError.Hint)
		return
	}
	fmt.Println("written to logs")
}

type storedUrl struct {
	id       string
	url      string
	username string
}

func getRedirectUrl(path string) (storedUrl, error) {
	var su storedUrl
	err := db.QueryRow(context.Background(), "select * from urls where id=$1", path).Scan(&su.id, &su.url, &su.username)
	if err != nil {
		return su, err

	}
	return su, nil
}
