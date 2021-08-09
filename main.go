package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
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

type Api struct {
	Links     Links
	Campaigns Campaigns
}

type Links struct {
}

type LinkAdd struct {
	LongUrl string `json: longUrl`
}

type Campaigns struct {
}
type stop struct {
	error
}

func (api Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string
	head, r.URL.Path = getSegment(r)
	fmt.Println("api head", head)
	switch head {
	case "login-user":
		auth.Signin(w, r, db)
	case "signup-user":
		auth.Signup(w, r, db)
	case "auth":
		auth.CheckAuth(w, r, db)
	case "logout":
		auth.Logout(w, r, db)
	case "favicon.ico":
		return
	case "links":
		fmt.Println("got links")
		api.Links.ServeHTTP(w, r)
	case "campaigns":
		api.Campaigns.ServeHTTP(w, r)
	default:
		handleRedirect(w, r, head)

	}
}

func (link Links) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string
	head, r.URL.Path = getSegment(r)
	if len(head) != 0 {
		switch r.Method {
		case http.MethodOptions:
			utils.HandleCors(w, r, []string{http.MethodDelete, http.MethodPut, http.MethodGet})
		case http.MethodDelete:
			link.Delete(w, r, head)
		case http.MethodPut:
			link.Update(w, r, head)
		case http.MethodGet:
			link.Get(w, r, head)
		default:
			utils.SendResponse(w, http.StatusMethodNotAllowed, "Wrong method")
		}

		return
	}

	switch r.Method {
	case http.MethodOptions:
		utils.HandleCors(w, r, []string{http.MethodGet, http.MethodPost})
	case http.MethodPost:
		link.Add(w, r)
	case http.MethodGet:
		link.GetAll(w, r)
	default:
		fmt.Fprintf(w, "Option not supported")
	}
}

/**
  Add a new shortened url
 **/
func (link Links) Add(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetSession(r, db)

	var reqBody LinkAdd
	result := requestJSON.Decode(w, r, &reqBody)
	if result.Status != 200 {
		fmt.Println(result.Message)
		utils.SendResponse(w, http.StatusInternalServerError, result.Message)
		return
	}
	fmt.Println("error", err)
	if err == nil || err == pgx.ErrNoRows {
		shortId, insErr := setId(r, reqBody.LongUrl, session)
		if insErr != nil {
			fmt.Println(insErr)
			utils.SendResponse(w, http.StatusInternalServerError, "Try again later")
			return
		}

		utils.SendResponse(w, http.StatusOK, &LinkAdd{LongUrl: "http://localhost:8080/" + shortId})
		return
	}
	utils.SendResponse(w, http.StatusInternalServerError, "Try again later")
}

/**
GetAll gets the summary logs of every shortened url that user has created
**/
func (link Links) GetAll(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, fmt.Sprintf("getAll url"))
}

// Get the summary detail of  short url for a given id
func (link Links) Get(w http.ResponseWriter, r *http.Request, id string) {

}

// Delete the short url for a given id
func (link Links) Delete(w http.ResponseWriter, r *http.Request, id string) {

}

// Update the short url for a given id
func (link Links) Update(w http.ResponseWriter, r *http.Request, id string) {

}

func (campaign Campaigns) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func getSegment(r *http.Request) (string, string) {
	split := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	head := split[0]
	rest := "/" + strings.Join(split[1:], "/")
	return head, rest
}

func main() {

	fmt.Println("starting server")
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://xanadu:xanadu@localhost:5432/tracker")
	if err != nil {
		log.Fatal(err)
	}
	db = dbpool
	defer db.Close()
	api := Api{}
	err = http.ListenAndServe(":8080", api)
	log.Fatal(err)
}

func handleRedirect(w http.ResponseWriter, r *http.Request, id string) {
	result, err := getRedirectUrl(id)
	switch err {
	case nil:
		http.Redirect(w, r, result.url, http.StatusPermanentRedirect)
		sessionR, serr := auth.GetSession(r, db)
		if serr != nil || sessionR.SessionId == "" {
			break
		}
		if sessionR.SessionId != "" {
			updateLogs(r, result)
		}
		break
	case pgx.ErrNoRows:
		fs := http.FileServer(http.Dir("public/build/"))
		fs.ServeHTTP(w, r)
		break
	default:
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

	parsedIP := net.ParseIP(r.Header.Get("ip"))
	var ip string
	if parsedIP == nil {
		ip = "0.0.0.0"
	} else {
		ip = parsedIP.String()
	}

	_, err = db.Exec(context.Background(), "insert into logs(id,campaign,source,medium,os,browser,device_type,created_on,ip)values($1,$2,$3,$4,$5,$6,$7,$8,$9)", result.id, campaign, source, medium, browser_info.OS.Platform.String(), browser_info.Browser.Name.String(), browser_info.DeviceType.String(), time.Now(), ip)
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

func setId(r *http.Request, largeUrl string, session auth.Session) (string, error) {
	value := createId()
	mainErr := retry(100, 1000, func() error {

		_, err := db.Exec(context.Background(), "insert into urls values($1,$2,$3)", value, largeUrl, session.Username)
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
		randStr := letterString[rand.Intn(len(letterString))]
		result = result + string(randStr)
	}
	return result
}
