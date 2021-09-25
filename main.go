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
	"text/template"
	"time"
	"url-shortner/auth"
	"url-shortner/utils"

	"github.com/avct/uasurfer"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
	"golang.org/x/crypto/bcrypt"
)

var db *pgxpool.Pool

type Api struct {
	Links     Links
	Campaigns Campaigns
}

type Links struct {
}
type Expiration struct {
	Time          string `json:time`
	ExpirationUrl string `json:expirationUrl`
	id            string `json: id`
}
type LinkAdd struct {
	LongUrl         string     `json: longUrl`
	Expiration      Expiration `json:expiration`
	Tag             string     `json:tag`
	Password        string     `json:psswd`
	NotFoundUrl     string     `json: 404_page`
	AndroidDeepLink string     `json: android_deep_link`
	IosDeepLink     string     `json: ios_deep_link`
}

type storedUrl struct {
	username string
	data     LinkAdd
	id       string
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

var count int = 0

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
	case "link-auth":
		validateLinkPassword(w, r)
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
	fmt.Println("links")
	unauthorized := true

	if err != nil {
		fmt.Println("links err", err)
		if err != pgx.ErrNoRows {
			utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Unauthorized"))
			return
		}
		if len(head) == 0 && r.Method == http.MethodPost {
			unauthorized = false
		}
	} else {
		unauthorized = false
	}
	if unauthorized {
		utils.SendResponse(w, http.StatusUnauthorized, utils.SendErrorToClient("Unauthorized"))
		return
	}
	// request to /links
	if len(head) == 0 {
		switch r.Method {
		case http.MethodOptions:
			utils.HandleCors(w, r, []string{http.MethodGet, http.MethodPost})
		case http.MethodPost:
			link.Add(w, r, session)
		case http.MethodGet:
			link.GetAll(w, r, session)
		default:
			utils.SendResponse(w, http.StatusMethodNotAllowed, utils.SendErrorToClient("Method not supported"))

		}
		return
	}

	// request to /link/<id>
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
		utils.SendResponse(w, http.StatusMethodNotAllowed, utils.SendErrorToClient("Method not supported"))
	}

}

func validateLinkPassword(w http.ResponseWriter, r *http.Request) {
	valid := false
	switch r.Method {
	case http.MethodOptions:
		utils.HandleCors(w, r, []string{http.MethodPost})
		break
	case http.MethodPost:
		valid = true
	default:
		utils.SendResponse(w, http.StatusMethodNotAllowed, utils.SendErrorToClient("Method not supported"))
	}

	if valid == false {
		return
	}
	if err := r.ParseForm(); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Parsing error"))
		return
	}
	psswd := r.Form.Get("psswd")
	id := r.Form.Get("id")
	if psswd == "" {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Incorrect password"))
		return
	}
	if id == "" {
		utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Link id is missing"))
		return
	}
	var rowPsswd, url string

	err := db.QueryRow(context.Background(), "select password,url from urls where id=$1", id).Scan(&rowPsswd, &url)
	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("The requested link is not found or is expired"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(rowPsswd), []byte(psswd))
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}

	link, err := getRedirectUrl(id)
	switch err {
	case nil:
		go updateLogs(r, link)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	case pgx.ErrNoRows:
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("The requested link is not found or is expired"))
	default:
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
	}
}

/**
  Add a new shortened url
 **/
func (link Links) Add(w http.ResponseWriter, r *http.Request, session auth.Session) {

	fmt.Println("add link")
	var reqBody LinkAdd
	result := requestJSON.Decode(w, r, &reqBody)
	if result.Status != 200 {
		utils.SendResponse(w, http.StatusInternalServerError, result.Message)
		return
	}

	shortId, insErr := setId(r, reqBody, session)
	if insErr != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}
	if reqBody.Expiration.Time != "" {
		err := addExpiration(reqBody, shortId)
		if err != nil {
			var parseErr *time.ParseError
			var dbErr *pgconn.PgError
			if errors.As(err, &parseErr) {
				utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Expiration time is incorrect"))
				return
			}
			if errors.As(err, &dbErr) {
				utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
				return
			}
			return
		}
	}
	utils.SendResponse(w, http.StatusOK, &LinkAdd{LongUrl: "http://localhost:8080/" + shortId, Password: reqBody.Password})
}

func addExpiration(reqBody LinkAdd, shortId string) error {
	timeT, err := time.Parse("Mon Jan 02 2006 15:04:05 MST-0700", reqBody.Expiration.Time)
	if err != nil {
		return err
	}
	expUrl := reqBody.Expiration.ExpirationUrl
	if expUrl == "" {
		expUrl = "http://localhost:8080"
	}
	_, dbErr := db.Exec(context.Background(), "insert into expiration values($1,$2,$3)", shortId, timeT.UTC(), expUrl)
	if dbErr != nil {
		return dbErr
	}
	return nil
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
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}

	list := make([]map[string]interface{}, 0)

	if rows.Err() != nil {
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
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
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("Not found"))
	default:
		log.Fatal(err)
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
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
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("Not found"))
		return
	}

	utils.SendResponse(w, http.StatusOK, utils.SendErrorToClient("Link removed"))
}

type UpdateMessage struct {
	Message string
	oldUrl  string
	newUrl  string
}

// Update the long url for a given id
func (link Links) Update(w http.ResponseWriter, r *http.Request, id string, username string) {
	var reqBody LinkAdd
	result := requestJSON.Decode(w, r, &reqBody)
	if result.Status != 200 {
		utils.SendResponse(w, http.StatusInternalServerError, result.Message)
		return
	}
	updateMessage := &UpdateMessage{}
	err := db.QueryRow(context.Background(), "BEGIN; SELECT url FROM urls WHERE id=$1 AND username=$2 FOR UPDATE; UPDATE urls SET url=$3 where id=$4 AND username=$5;COMMIT;", id, username, reqBody.LongUrl, id, username).Scan(updateMessage.oldUrl)

	if err != nil {
		if err == pgx.ErrNoRows {
			utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("Link not found"))
			return
		}
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}
	updateMessage.newUrl = reqBody.LongUrl
	updateMessage.Message = "Link updated"
	utils.SendResponse(w, http.StatusOK, updateMessage)
}

func getSegment(r *http.Request) (string, string) {
	split := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	head := split[0]
	rest := "/" + strings.Join(split[1:], "/")
	return head, rest
}

func main() {

	fmt.Println("starting server")

	poolConfig, err := pgxpool.ParseConfig("postgres://xanadu:xanadu@localhost:5432/tracker")
	if err != nil {
		log.Fatal(err)
	}
	poolConfig.ConnConfig.OnNotice = onNotify
	poolConfig.AfterConnect = func(c1 context.Context, c2 *pgx.Conn) error {

		c2.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			OID:   pgtype.UUIDOID,
			Name:  "uuid",
		})
		return nil
	}
	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	api := Api{}
	err = http.ListenAndServe(":8080", api)
	log.Fatal(err)

}

func handleRedirect(w http.ResponseWriter, r *http.Request, id string) {
	link, err := getRedirectUrl(id)
	switch err {
	case nil:
		expTime, _ := time.Parse("Mon Jan 02 2006 15:04:05 MST-0700", link.data.Expiration.Time)
		isAfter := link.data.Expiration.Time != "" && time.Now().After(expTime)
		redirectUrl := link.data.LongUrl
		if link.data.Password != "" {
			temp := template.Must(template.ParseFiles("views/authentication.html"))
			su := map[string]string{
				"Id": link.id,
			}
			temp.Execute(w, su)
			break
		}
		if isAfter {
			redirectUrl = link.data.Expiration.ExpirationUrl
		}

		http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
		go updateLogs(r, link)
		break
	case pgx.ErrNoRows:
		http.Redirect(w, r, "http://localhost:3000", http.StatusPermanentRedirect)
		break
	default:
		fmt.Fprintf(w, "Error redirecting")
	}
}

func updateLogs(r *http.Request, result storedUrl) {

	browser_info := uasurfer.Parse(r.UserAgent())
	u, err := url.Parse(result.data.LongUrl)
	if err != nil {
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
		var pgError *pgconn.PgError
		errors.As(err, &pgError)
		fmt.Printf("%s at %d col %s %s %s", pgError.Message, pgError.Line, pgError.ColumnName, pgError.Detail, pgError.Hint)
		return
	}
	fmt.Println("written to logs")
}

func getRedirectUrl(path string) (storedUrl, error) {
	var su storedUrl
	var id, longUrl, username, psswd, notFoundUrl, exp_url *string
	var exp_time *time.Time
	err := db.QueryRow(context.Background(), "select urls.id,urls.url,urls.username,urls.password,urls.not_found_url,expiration,expired_url from urls left join expiration on urls.id=expiration.id where urls.id=$1", path).Scan(&id, &longUrl, &username, &psswd, &notFoundUrl, &exp_time, &exp_url)

	if err != nil {
		return su, err

	}
	exp := Expiration{}
	if exp_time != nil {
		exp.Time = exp_time.String()
	}
	if exp_url != nil {
		exp.ExpirationUrl = *exp_url
	}
	link := LinkAdd{LongUrl: *longUrl, Password: *psswd, NotFoundUrl: *notFoundUrl, Expiration: exp}
	return storedUrl{id: *id, username: *username, data: link}, nil

}

func setId(r *http.Request, reqBody LinkAdd, session auth.Session) (string, error) {
	psswdHash := reqBody.Password
	if psswdHash != "" {
		hash, err := auth.GeneratePsswdHash(reqBody.Password)
		if err != nil {
			return "", err
		}
		psswdHash = string(hash)
	}

	var uniqId string
	err := db.QueryRow(context.Background(), "SELECT insertLongUrl($1)", reqBody.LongUrl).Scan(&uniqId)
	if err != nil {
		return "", err
	}

	return uniqId, nil

}
func onNotify(c *pgconn.PgConn, n *pgconn.Notice) {

	fmt.Println("Message:", *n)
}
