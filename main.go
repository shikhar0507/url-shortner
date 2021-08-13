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
type Expiration struct {
	Time          time.Time `json:time`
	ExpirationUrl string    `json: url`
	id            string    `json: id`
}
type LinkAdd struct {
	LongUrl         string     `json: longUrl`
	Expiration      Expiration `json:expiration`
	Tag             string     `json:tag`
	Password        string     `json:psswd`
	NotFoundUrl     string     `json: 404_page`
	AndroidDeepLink string     `json: android_deep_link`
	IosDeepLink     string     `json: ios_deep_link`
	//	CountryBlock    []string          `json: country_block`
	//IpRedirection   map[string]string `json:ip_redirection`
}

type Campaigns struct {
}
type stop struct {
	error
}

type Link struct {
	Browser          string `json:browser`
	Os               string `json:os`
	Device_type      string `json:device_type`
	Total_clicks     int    `json:total_clicks`
	BrowserCount     int    `json:browser_count`
	OsCount          int    `json:os_count`
	DeviceType_Count int    `json:device_type_count`
}

type GetLinks struct {
	Result []map[string]interface{} `json:result`
}

type storedUrl struct {
	id              string
	url             string
	username        string
	Expiration      Expiration
	Tag             string
	Password        string
	NotFoundUrl     string
	AndroidDeepLink string
	IosDeepLink     string
	CountryBlock    []string
	IpRedirection   map[string]string
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
	case "index.html":
		http.ServeFile(w, r, "index.html")
	default:
		handleRedirect(w, r, head)
	}
}

func (link Links) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string
	head, r.URL.Path = getSegment(r)
	session, err := auth.GetSession(r, db)
	fmt.Println(head)
	if err != nil {
		if err == pgx.ErrNoRows {
			if r.Method == http.MethodPost && len(head) == 0 {
				link.Add(w, r, session)
				return
			}
			utils.SendResponse(w, http.StatusUnauthorized, "not authorized")
			return
		}
		utils.SendResponse(w, http.StatusInternalServerError, "Try again later")
	}

	if len(head) != 0 {
		switch r.Method {
		case http.MethodOptions:
			utils.HandleCors(w, r, []string{http.MethodDelete, http.MethodPut, http.MethodGet})
		case http.MethodDelete:
			link.Delete(w, r, head, session.Username)
		case http.MethodPut:
			link.Update(w, r, head, session.Username)
		case http.MethodGet:
			link.Get(w, r, head, session.Username)
		default:
			utils.SendResponse(w, http.StatusMethodNotAllowed, "Wrong method")
		}

		return
	}

	switch r.Method {
	case http.MethodOptions:
		utils.HandleCors(w, r, []string{http.MethodGet, http.MethodPost})
	case http.MethodPost:
		link.Add(w, r, session)
	case http.MethodGet:
		link.GetAll(w, r, session)
	default:
		fmt.Fprintf(w, "Option not supported")
	}
}

/**
  Add a new shortened url
 **/
func (link Links) Add(w http.ResponseWriter, r *http.Request, session auth.Session) {

	var reqBody LinkAdd
	result := requestJSON.Decode(w, r, &reqBody)
	if result.Status != 200 {
		utils.SendResponse(w, http.StatusInternalServerError, result.Message)
		return
	}

	shortId, insErr := setId(r, reqBody, session)
	if insErr != nil {
		fmt.Println(insErr)
		utils.SendResponse(w, http.StatusInternalServerError, "Try again later")
		return
	}

	utils.SendResponse(w, http.StatusOK, &LinkAdd{LongUrl: "http://localhost:8080/" + shortId})
}

/**
GetAll gets the summary logs of every shortened url that user has created
**/
func (link Links) GetAll(w http.ResponseWriter, r *http.Request, sesstion auth.Session) {
	query := `with t as (
select urls.id,urls.username,urls.url,t2.device_type,t1.browser,t0.os,t4.total_clicks,sum(case when t4.total_clicks is not null then t4.total_clicks else 0 end) from urls
left join (select id,count(*) as total_clicks from logs group by id order by total_clicks desc)t4 on t4.id = urls.id
left join (select t.* from (select logs.id,logs.device_type, rank() over(partition by id order by count(device_type) desc) from logs group by logs.id,logs.device_type) t where rank=1)t2 on urls.id = t2.id
left join (select t.* from (select logs.id,logs.browser, rank() over(partition by id order by count(browser) desc) from logs group by logs.id,logs.browser) t where rank=1) t1 on t2.id = t1.id
left join  (select t.* from (select logs.id,logs.os, rank() over(partition by id order by count(os) desc) from logs group by logs.id,logs.os) t where rank=1)t0 on t1.id = t0.id group by urls.id,urls.username,urls.url,t2.device_type,t1.browser,t0.os,t4.total_clicks
order by sum desc)
select t.id,t.url,t.browser,t.os,t.device_type,coalesce(t.total_clicks,0) from t where t.username=$1 limit 10`

	rows, err := db.Query(context.Background(), query, sesstion.Username)
	if err != nil {
		fmt.Println("links fetch", err)
		utils.SendResponse(w, http.StatusInternalServerError, utils.Response{Status: http.StatusInternalServerError, Message: "Try again later"})
		return
	}

	list := make([]map[string]interface{}, 0)

	if rows.Err() != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "try again later")
		return
	}

	defer rows.Close()
	for rows.Next() {

		var id, device_type, browser, os, url *string
		var total_clicks int
		hash := make(map[string]interface{}, 0)
		scanErr := rows.Scan(&id, &url, &browser, &os, &device_type, &total_clicks)
		if scanErr != nil {
			fmt.Println("Scan err", scanErr)
		}

		hash["id"] = id
		hash["url"] = url
		hash["browser"] = browser
		hash["os"] = os
		hash["device_type"] = device_type
		hash["total_clicks"] = total_clicks
		list = append(list, hash)
	}

	links := GetLinks{Result: list}
	utils.SendResponse(w, http.StatusOK, links)
}

// Get the summary detail of  short url for a given id
func (link Links) Get(w http.ResponseWriter, r *http.Request, id, username string) {
	var browser, os, device_type *string
	var browserC, osC, device_typeC, total_clicks *int
	query := `select urls.username,urls.id,urls.url, browser,browser_count,os,os_count,device_type,device_count,t0.total_clicks from urls
 left join 
 (select id,count(*) as total_clicks from logs group by id)t0 on urls.id = t0.id full outer join  (select id,browser,count(*) as browser_count from logs group by id,browser order by browser_count desc) t1 on t0.id=t1.id
left join
(select id,os,count(*) as os_count from logs group by id,os order by os_count desc) t2 on t1.id = t2.id
left join
(select id,device_type,count(*) as device_count from logs  group by id,device_type order by device_count desc)t3 on t2.id = t3.id where t0.id=$1 and urls.username=$2`

	err := db.QueryRow(context.Background(), query, id, username).Scan(&browser, &browserC, &os, &osC, &device_type, &device_typeC, &total_clicks)
	switch err {
	case nil:
		utils.SendResponse(w, http.StatusOK, Link{Browser: *browser, BrowserCount: *browserC, Os: *os, OsCount: *osC, Device_type: *device_type, DeviceType_Count: *device_typeC, Total_clicks: *total_clicks})
	case pgx.ErrNoRows:
		utils.SendResponse(w, http.StatusNotFound, "Not found")
	default:
		log.Fatal(err)
		utils.SendResponse(w, http.StatusInternalServerError, "try again later")
	}

}

// Delete the short url for a given id
func (link Links) Delete(w http.ResponseWriter, r *http.Request, id, username string) {
	tag, err := db.Exec(context.Background(), "delete from urls where id=$1 and username=$2", id, username)

	if err != nil {
		fmt.Println(err)
		return
	}
	if tag.RowsAffected() == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Not found")
		return
	}
	fmt.Println("rows affected", tag.RowsAffected())
	utils.SendResponse(w, http.StatusOK, "Link removed")
}

// Update the long url for a given id
func (link Links) Update(w http.ResponseWriter, r *http.Request, id string, username string) {
	var reqBody LinkAdd
	result := requestJSON.Decode(w, r, &reqBody)
	if result.Status != 200 {
		utils.SendResponse(w, http.StatusInternalServerError, result.Message)
		return
	}

	tag, err := db.Exec(context.Background(), "update urls set url=$1 where id=$2 and username=$3", reqBody.LongUrl, id, username)

	if err != nil {
		fmt.Println(err)
		return
	}
	if tag.RowsAffected() == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Not found")
		return
	}
	fmt.Println("rows affected", tag.RowsAffected())
	utils.SendResponse(w, http.StatusOK, "Link updated")
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
		go updateLogs(r, result)
		break
	case pgx.ErrNoRows:
		//	fs := http.FileServer(http.Dir("public/build/"))
		//fs.ServeHTTP(w, r)
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

func getRedirectUrl(path string) (LinkAdd, error) {
	var su storedUrl
	err := db.QueryRow(context.Background(), "select * from urls where id=$1", path).Scan(&su.id, &su.url, &su.username)
	if err != nil {
		return su, err

	}
	return su, nil
}

func setId(r *http.Request, reqBody LinkAdd, session auth.Session) (string, error) {
	value := createId()
	if reqBody.Expiration.Time > 0 && len(reqBody.Expiration.ExpirationUrl) == 0 {
		return "", errors.New("Expiration Url not found")
	}

	mainErr := retry(100, 1000, func() error {

		_, err := db.Exec(context.Background(), "insert into urls values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", value, reqBody.LongUrl, session.Username, reqBody.Tag, reqBody.Password, reqBody.NotFoundUrl, reqBody.AndroidDeepLink, reqBody.IosDeepLink)
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
