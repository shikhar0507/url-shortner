package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
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

	"url-shortner/props"

	"github.com/avct/uasurfer"
	"github.com/gocolly/colly"
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
type Campaign struct {
	Name    string `json:"name"`
	Medium  string `json:"medium"`
	Source  string `json:"source"`
	Term    string `json:"term"`
	Content string `json:"content"`
	Id      string `json:"id"`
}
type country_redirect struct {
	id           string
	Country_code string `json:"country_code"`
	Country_url  string `json:"country_url"`
}
type LinkAdd struct {
	LongUrl          string             `json:"longUrl"`
	Expiration       string             `json:"expiration"`
	Tag              string             `json:"tag"`
	Password         string             `json:"psswd"`
	NotFoundUrl      string             `json:"not_found_url"`
	AndroidDeepLink  string             `json:"android_deep_link"`
	IosDeepLink      string             `json:"ios_deep_link"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Campaign         Campaign           `json:"campaign"`
	HTTPStatus       int                `json:"http_status"`
	PlayStoreLink    string             `json:"play_store_link"`
	Qrcode           bool               `json:"qr_code"`
	Country_block    []string           `json:"country_block"`
	Country_redirect []country_redirect `json:"country_redirect"`
	Mobile_url       string             `json:"mobile_url"`
	Destkop_url      string             `json:"desktop_url"`
	Others_url       string             `json:"others_url"`
}

type LinkAddResponse struct {
	ShortUrl string `json:"shortUrl"`
}

func (l *LinkAdd) Value() (driver.Value, error) {
	return json.Marshal(l)
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
	Browser     []map[string]interface{} `json:"top_browser"`
	Os          []map[string]interface{} `json:"top_os"`
	Device_type []map[string]interface{} `json:"top_device"`
	LongUrl     *string                  `json:"long_url"`
	Qrcode      bool                     `json:"qr_code"`
	Name        *string                  `json:"link_name"`
	Description *string                  `json:"link_description"`
	Tag         *string                  `json:"link_tag"`
	Createdon   string                   `json:"created_on"`
}

type GetLinks struct {
	Result      []map[string]interface{} `json:"result"`
	Topbrowser  string                   `json:"top_browser"`
	Topos       string                   `json:"top_os"`
	Topdevice   string                   `json:"top_device"`
	TotalClicks int                      `json:"total_clicks"`
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
	if err != nil {
		// no session
		if err == pgx.ErrNoRows {
			fmt.Println("shortening")
			shorten(w, r, session, link)
			return
		}
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
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

	if head == "opengraph" {
		parseRequest := colly.NewCollector(colly.AllowURLRevisit())
		websiteRequest(w, r, parseRequest)
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
func websiteRequest(w http.ResponseWriter, r *http.Request, parseRequest *colly.Collector) {
	fmt.Println("request")
	allowed := false
	switch r.Method {
	case "options":
		return
	case "POST":
		allowed = true
	default:
		utils.SendResponse(w, http.StatusMethodNotAllowed, utils.SendErrorToClient("Method not allowed"))
	}
	if !allowed {
		return
	}
	var request props.OpenGraphReq
	result := requestJSON.Decode(w, r, &request)

	if result.Status != 200 {
		utils.SendResponse(w, result.Status, utils.SendErrorToClient(result.Message))
		return
	}

	var basic props.Basic
	imageHash := make(map[string]props.Image)
	videoHash := make(map[string]props.Video)
	audioHash := make(map[string]props.Audio)
	type types struct {
		currentImage props.Image
		currentVideo props.Video
		currentMusic props.Music
		currentAudio props.Audio
	}
	ty := types{}

	parseRequest.OnHTML("meta", func(e *colly.HTMLElement) {
		prop := e.Attr("property")
		value := e.Attr("content")
		switch prop {
		case "og:title":
			basic.Title = value
		case "og:type":
			basic.Type = value
		case "og:image":
			currImage, ok := imageHash[value]
			img := props.Image{Url: value}
			if ok {
				ty.currentImage = currImage
			} else {
				imageHash[value] = img
				ty.currentImage = img
			}
		case "og:image:width":
			ty.currentImage.Width = value
			imageHash[ty.currentImage.Url] = ty.currentImage
		case "og:image:height":
			ty.currentImage.Height = value
			imageHash[ty.currentImage.Url] = ty.currentImage
		case "og:image:alt":
			ty.currentImage.Alt = value
			imageHash[ty.currentImage.Url] = ty.currentImage
		case "og:image:type":
			ty.currentImage.Type = value
			imageHash[ty.currentImage.Url] = ty.currentImage
		case "og:image:secure_url":
			ty.currentImage.Secure_url = value
			imageHash[ty.currentImage.Url] = ty.currentImage
		case "og:url":
			basic.Url = value
		case "og:description":
			basic.Description = value
		case "og:site_name":
			basic.SiteName = value
		case "og:determiner":
			basic.Determiner = value
		case "og:video":
			vid := props.Video{Url: value}
			currVideo, ok := videoHash[value]
			if ok {
				ty.currentVideo = currVideo
			} else {
				videoHash[value] = vid
				ty.currentVideo = vid
			}
		case "og:video:secure_url":
			ty.currentVideo.Secure_url = value
			videoHash[ty.currentVideo.Url] = ty.currentVideo
		case "og:video:type":
			ty.currentVideo.Type = value
			videoHash[ty.currentVideo.Url] = ty.currentVideo
		case "og:video:width":
			ty.currentVideo.Width = value
			videoHash[ty.currentVideo.Url] = ty.currentVideo
		case "og:video:height":
			ty.currentVideo.Height = value
			videoHash[ty.currentVideo.Url] = ty.currentVideo
		case "og:audio":
			aud := props.Audio{Url: value}
			currAudio, ok := audioHash[value]
			if ok {
				ty.currentAudio = currAudio
			} else {
				audioHash[value] = aud
				ty.currentAudio = aud
			}
		case "og:audio:secure_url":
			ty.currentAudio.Secure_url = value
			audioHash[ty.currentAudio.Url] = ty.currentAudio
		case "og:audio:type":
			ty.currentAudio.Type = value
			audioHash[ty.currentAudio.Url] = ty.currentAudio
		case "og:locale":
			if basic.Locale != "" {
				basic.Locales = append(basic.Locales, value)
			} else {
				basic.Locale = value
			}
		}

	})
	parseRequest.OnError(func(r *colly.Response, e error) {
		fmt.Println(e, string(r.Body))
	})

	parseRequest.OnRequest(func(collyReq *colly.Request) {

		fmt.Println("visiting", collyReq.URL.String())
	})
	parseRequest.OnResponse(func(collyResp *colly.Response) {
		fmt.Println("parsing response", collyResp.StatusCode)
	})
	parseRequest.OnScraped(func(re *colly.Response) {
		fmt.Println("parsed", re.StatusCode)
		for _, v := range imageHash {
			basic.Image = append(basic.Image, v)
		}
		for _, v := range videoHash {
			basic.Video = append(basic.Video, v)
		}
		for _, v := range audioHash {
			basic.Audio = append(basic.Audio, v)
		}
		utils.SendResponse(w, http.StatusOK, basic)
	})
	parseRequest.Visit(request.Url)

}

func shorten(w http.ResponseWriter, r *http.Request, session auth.Session, link Links) {
	switch r.Method {
	case http.MethodOptions:
		utils.HandleCors(w, r, []string{http.MethodGet})
	case http.MethodPost:
		link.Add(w, r, session)
	default:
		utils.SendResponse(w, http.StatusMethodNotAllowed, utils.SendErrorToClient("Method not supported"))

	}
	return
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
	var username pgtype.Varchar

	err := db.QueryRow(context.Background(), "select password,url,username from urls where id=$1", id).Scan(&rowPsswd, &url, &username)
	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("The requested link is not found or is expired"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(rowPsswd), []byte(psswd))
	if err != nil {
		fmt.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Incorrect password"))
		return
	}

	link, err := getRedirectUrl(id)
	switch err {
	case nil:
		fmt.Println("username status", username.Status, username.String)
		if username.Status != pgtype.Null {
			go updateLogs(r, link)
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	case pgx.ErrNoRows:
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("The requested link is not found or is expired"))
	default:
		fmt.Println(err)
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
		fmt.Println(insErr)
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}
	if reqBody.Expiration != "" {
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
	fmt.Println(shortId)
	utils.SendResponse(w, http.StatusOK, &LinkAddResponse{ShortUrl: "http://localhost:8080/" + shortId})
}

func addExpiration(reqBody LinkAdd, shortId string) error {

	timeT, err := time.Parse("2006-01-02T15:04", reqBody.Expiration)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(timeT.UTC())
	expUrl := reqBody.NotFoundUrl
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
	fmt.Println("fetching for", sesstion.Username)

	query := `with t as (
select urls.id,urls.link_tag,urls.link_name,urls.username,urls.url,t2.device_type,t1.browser,t0.os,t4.total_clicks,urls.create_timestamp,sum(case when t4.total_clicks is not null then t4.total_clicks else 0 end) from urls
left join (select id,count(*) as total_clicks from logs group by id order by total_clicks desc)t4 on t4.id = urls.id
left join (select t.* from (select logs.id,logs.device_type, rank() over(partition by id order by count(device_type) desc) from logs group by logs.id,logs.device_type) t where rank=1)t2 on urls.id = t2.id
left join (select t.* from (select logs.id,logs.browser, rank() over(partition by id order by count(browser) desc) from logs group by logs.id,logs.browser) t where rank=1) t1 on t2.id = t1.id
left join  (select t.* from (select logs.id,logs.os, rank() over(partition by id order by count(os) desc) from logs group by logs.id,logs.os) t where rank=1)t0 on t1.id = t0.id group by urls.id,urls.username,urls.url,t2.device_type,t1.browser,t0.os,t4.total_clicks
order by sum desc)
select distinct t.id,t.link_tag,t.link_name,t.url,t.browser,t.os,t.device_type,coalesce(t.total_clicks,0) as total_clicks,t.create_timestamp as create_time from t where t.username=$1`

	sort := r.URL.Query().Get("sort")
	switch sort {
	case "time_latest":
		query = query + " ORDER BY create_time ASC"
	case "time_oldest":
		query = query + " ORDER BY create_time DESC"
	case "clicks_desc":
		query = query + " ORDER BY total_clicks DESC"
	case "clicks_asc":
		query = query + " ORDER BY total_clicks ASC"
	default:
		query = query + " ORDER BY total_clicks DESC"
	}

	rows, err := db.Query(context.Background(), query, sesstion.Username)
	if err != nil {
		fmt.Println(err)
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

		var id, device_type, browser, os, url, name, tag *string
		var total_clicks int
		var create_time pgtype.Timestamptz
		hash := make(map[string]interface{}, 0)
		scanErr := rows.Scan(&id, &tag, &name, &url, &browser, &os, &device_type, &total_clicks, &create_time)
		if scanErr != nil {
			fmt.Println("Scan err", scanErr)
		}

		hash["id"] = id
		hash["url"] = url
		hash["browser"] = browser
		hash["os"] = os
		hash["device_type"] = device_type
		hash["total_clicks"] = total_clicks
		hash["tag"] = tag
		hash["name"] = name
		hash["timestamp"] = create_time.Time.String()
		list = append(list, hash)
	}

	var topBrowser pgtype.Varchar
	var topOs pgtype.Varchar
	var topDevice pgtype.Varchar
	var totalClicks *int

	metaErr := db.QueryRow(context.Background(), `with t as (SELECT urls.id,logs.browser,logs.os,logs.device_type from urls INNER JOIN logs ON urls.id=logs.id where username=$1)
SELECT browser,os,device_type,total_clicks.total_count FROM (SELECT browser,count(browser) as c FROM t GROUP BY browser ORDER BY c DESC LIMIT 1) browser, (SELECT os,count(os) as c FROM t GROUP BY os ORDER BY c DESC LIMIT 1)os,(SELECT device_type,count(device_type) as c FROM t GROUP BY device_type ORDER BY c DESC LIMIT 1)device_type,(SELECT count(*) as total_count FROM t ORDER BY total_count DESC)total_clicks`, sesstion.Username).Scan(&topBrowser, &topOs, &topDevice, &totalClicks)
	if metaErr != nil {
		fmt.Println(metaErr)
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}

	links := GetLinks{Result: list, Topbrowser: topBrowser.String, Topos: topOs.String, Topdevice: topDevice.String, TotalClicks: *totalClicks}
	utils.SendResponse(w, http.StatusOK, links)
}

// Get the summary detail of  short url for a given id
func (link Links) Get(w http.ResponseWriter, r *http.Request, id, username string) {

	query := `with t as (select browser,b_c,os,o_c,device_type,o_d,t3.rank,t1.id from (select distinct browser,id,count(browser) as b_c, row_number() over(order by count(browser)  desc) rank from logs where id='sg' group by browser,id)t1
FULL OUTER JOIN
(select  os,id,count(os) as o_c, row_number() over(order by count(os) desc) rank from logs where id='sg' group by os,id)t2 on t1.rank=t2.rank FULL OUTER JOIN
(select  device_type,id,count(device_type) as o_d, row_number() over(order by count(device_type) desc) rank from logs where id='sg' group by device_type,id)t3 on t2.rank=t3.rank)
SELECT urls.url,qr_code,create_timestamp,link_name,link_tag,link_description,browser,b_c,os,o_c,device_type,o_d FROM URLS left join t on urls.id=t.id where urls.id=$1 and urls.username=$2`

	rows, err := db.Query(context.Background(), query, id, username)

	switch err {
	case nil:
		linkDetail := Link{}
		for rows.Next() {
			var browser, os, device_type pgtype.Varchar
			var browserCount, osCount, deviceCount *int
			var created_on pgtype.Timestamptz
			scanErr := rows.Scan(&linkDetail.LongUrl, &linkDetail.Qrcode, &created_on, &linkDetail.Name, &linkDetail.Tag, &linkDetail.Description, &browser, &browserCount, &os, &osCount, &device_type, &deviceCount)
			if scanErr != nil {
				fmt.Println(scanErr)
				return
			}
			linkDetail.Createdon = created_on.Time.String()
			linkDetail.Browser = append(linkDetail.Browser, map[string]interface{}{
				"name":  browser,
				"value": browserCount,
			})
			linkDetail.Os = append(linkDetail.Os, map[string]interface{}{
				"name":  os,
				"value": osCount,
			})
			linkDetail.Device_type = append(linkDetail.Device_type, map[string]interface{}{
				"name":  device_type,
				"value": deviceCount,
			})

		}
		utils.SendResponse(w, http.StatusOK, &linkDetail)
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
		fmt.Println("og time", link.data.Expiration)
		expTime, _ := time.Parse("2006-01-02 15:04", link.data.Expiration)
		fmt.Println("exp time", expTime.String())
		isAfter := link.data.Expiration != "" && time.Now().After(expTime)
		redirectUrl := link.data.LongUrl
		fmt.Println("password", link.data.Password)
		if link.data.Password != "" {
			temp := template.Must(template.ParseFiles("views/authentication.html"))
			su := map[string]string{
				"Id": link.id,
			}
			temp.Execute(w, su)
			break
		}
		if isAfter {
			redirectUrl = link.data.NotFoundUrl
		}
		fmt.Println("check if website exist")

		conn, err := net.DialTimeout("tcp", "instagram.fdel1-5.fna.fbcdn.net:https", time.Second*5)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, link.data.NotFoundUrl, 301)
			return
		}
		defer conn.Close()
		http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
		if link.username != "" {
			go updateLogs(r, link)
		}
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
	fmt.Println(browser_info.DeviceType.StringTrimPrefix())
	_, err = db.Exec(context.Background(), "insert into logs(id,campaign,source,medium,os,browser,device_type,created_on,ip)values($1,$2,$3,$4,$5,$6,$7,$8,$9)", result.id, campaign, source, medium, browser_info.OS.Platform.StringTrimPrefix(), browser_info.Browser.Name.StringTrimPrefix(), browser_info.DeviceType.StringTrimPrefix(), time.Now(), ip)
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
	var id, username, psswd pgtype.Varchar
	var longUrl, notFoundUrl pgtype.Text
	var exp_time pgtype.Timestamptz
	err := db.QueryRow(context.Background(), "select urls.id,urls.url,urls.username,urls.password,urls.not_found_url,expiration from urls left join expiration on urls.id=expiration.id where urls.id=$1", path).Scan(&id, &longUrl, &username, &psswd, &notFoundUrl, &exp_time)

	fmt.Println(id)

	if err != nil {
		return su, err
	}
	var exp string
	if exp_time.Status != pgtype.Null {
		exp = exp_time.Time.String()
	}

	link := LinkAdd{LongUrl: longUrl.String, Password: psswd.String, NotFoundUrl: notFoundUrl.String, Expiration: exp}
	return storedUrl{id: id.String, username: username.String, data: link}, nil
}

func setId(r *http.Request, reqBody LinkAdd, session auth.Session) (string, error) {
	psswdHash := reqBody.Password

	if psswdHash != "" {
		fmt.Println("psswd", psswdHash)
		hash, err := auth.GeneratePsswdHash(reqBody.Password)
		if err != nil {
			return "", err
		}
		psswdHash = string(hash)
	}

	var uniqId string

	err := db.QueryRow(context.Background(), "SELECT insertLongUrl($1,$2,$3)", session.Username, psswdHash, reqBody).Scan(&uniqId)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return uniqId, nil

}
func onNotify(c *pgconn.PgConn, n *pgconn.Notice) {

	fmt.Println("Message:", *n)
}
