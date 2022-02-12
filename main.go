package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
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
	LongUrl             string   `json:"longUrl"`
	Expiration          string   `json:"expiration"`
	Tag                 string   `json:"tag"`
	Password            string   `json:"psswd"`
	NotFoundUrl         string   `json:"not_found_url"`
	AndroidDeepLink     string   `json:"android_deep_link"`
	IosDeepLink         string   `json:"ios_deep_link"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Campaign            Campaign `json:"campaign"`
	HTTPStatus          int      `json:"http_status"`
	PlayStoreLink       string   `json:"play_store_link"`
	Qrcode              bool     `json:"qr_code"`
	Country_block       []string `json:"country_block"`
	Mobile_url          string   `json:"mobile_url"`
	Destkop_url         string   `json:"desktop_url"`
	has_country_block   bool
	has_device_redirect bool
	id                  string
	username            string
	redirect_url_main   string
}

type LinkAddResponse struct {
	ShortUrl string `json:"shortUrl"`
}

func (l *LinkAdd) Value() (driver.Value, error) {
	return json.Marshal(l)
}

type Campaigns struct {
}
type stop struct {
	error
}

type Stat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Link struct {
	Browser     []Stat      `json:"top_browser"`
	Os          []Stat      `json:"top_os"`
	Device_type []Stat      `json:"top_device"`
	Referrer    []Stat      `json:"top_referrer"`
	Countries   []Stat      `json:"top_countries"`
	LongUrl     *string     `json:"long_url"`
	Qrcode      pgtype.Bool `json:"qr_code"`
	Name        *string     `json:"link_name"`
	Description *string     `json:"link_description"`
	Tag         *string     `json:"link_tag"`
	Createdon   string      `json:"created_on"`
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

	if err != nil && r.Method != http.MethodOptions {
		// no session
		fmt.Println(err)
		if err == pgx.ErrNoRows {
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
	fmt.Println("request", r.Method)
	allowed := false
	switch r.Method {
	case http.MethodOptions:
		utils.HandleCors(w, r, []string{http.MethodPost})

	case http.MethodPost:
		allowed = true
	default:
		fmt.Println("not allowed", r.Method)
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
	err := db.QueryRow(context.Background(), "select password,url from urls where id=$1", id).Scan(&rowPsswd, &url)
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
		handleRedirectionForDevice(r, w, link)
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
		var parseErr *time.ParseError
		if errors.As(insErr, &parseErr) {
			utils.SendResponse(w, http.StatusBadRequest, utils.SendErrorToClient("Expiration time is incorrect"))
			return
		}
		fmt.Println(insErr)
		utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		return
	}

	fmt.Println(shortId)
	utils.SendResponse(w, http.StatusOK, &LinkAddResponse{ShortUrl: "http://localhost:8080/" + shortId})
}

/**
GetAll gets the summary logs of every shortened url that user has created
**/
func (link Links) GetAll(w http.ResponseWriter, r *http.Request, sesstion auth.Session) {

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
		query = query + " ORDER BY create_time DESC"
	case "time_oldest":
		query = query + " ORDER BY create_time ASC"
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
	linkDetail := Link{}
	var created_on pgtype.Timestamptz

	checkQuery := `SELECT url,qr_code,TO_CHAR(create_timestamp,'YYYY-MM-DD"T"HH24:MI:SS"Z"')::timestamptz,link_name,link_tag,link_description FROM urls where id=$1 AND username=$2`

	err := db.QueryRow(context.Background(), checkQuery, id, username).Scan(&linkDetail.LongUrl, &linkDetail.Qrcode, &created_on, &linkDetail.Name, &linkDetail.Tag, &linkDetail.Description)
	switch err {

	case nil:
		linkDetail.Createdon = created_on.Time.UTC().String()
		query := `with t as (select browser,b_c,os,o_c,device_type,o_d,referrer,ref_count,country,country_count from (select browser,id,count(browser) as b_c, row_number() over(order by count(browser)  desc) rank from logs  where id=$1 group by browser,id)t1
		FULL OUTER JOIN
		(select  os,id,count(os) as o_c, row_number() over(order by count(os) desc) rank from logs  where id=$1 group by os,id)t2 on t1.rank=t2.rank FULL OUTER JOIN
		(select  device_type,id,count(device_type) as o_d, row_number() over(order by count(device_type) desc) rank from logs where id=$1 group by device_type,id)t3 on t2.rank=t3.rank
FULL OUTER JOIN
(select referrer,id,count(referrer) as ref_count, row_number() over(order by count(referrer) desc) rank from logs where id=$1 group by referrer,id)t4 on t3.rank = t4.rank
FULL OUTER JOIN
(select country,id,count(country) as country_count,row_number() OVER(order by count(country) desc)rank from logs where id=$1 group by country,id)t5 on t4.rank=t5.rank
)SELECT * FROM t`

		rows, err := db.Query(context.Background(), query, id)
		switch err {

		case nil:

			browserArr := make([]Stat, 0)
			osArr := make([]Stat, 0)
			deviceArr := make([]Stat, 0)
			refArr := make([]Stat, 0)
			countryArr := make([]Stat, 0)
			defer rows.Close()
			for rows.Next() {
				var browser, os, device_type, country pgtype.Varchar
				var browserCount, osCount, deviceCount, refCount, countryCount pgtype.Int8
				var ref pgtype.Text
				scanErr := rows.Scan(&browser, &browserCount, &os, &osCount, &device_type, &deviceCount, &ref, &refCount, &country, &countryCount)

				if scanErr != nil {
					fmt.Println(scanErr)
					return
				}

				if browser.Status != pgtype.Null {
					fmt.Println("browser", browser.String, browserCount)
					browserArr = append(browserArr, Stat{
						Name:  browser.String,
						Value: int(browserCount.Int),
					})

				}

				if os.Status != pgtype.Null {
					osArr = append(osArr, Stat{
						Name:  os.String,
						Value: int(osCount.Int),
					})

				}

				if device_type.Status != pgtype.Null {
					deviceArr = append(deviceArr, Stat{
						Name:  device_type.String,
						Value: int(deviceCount.Int),
					})
				}
				if ref.Status != pgtype.Null {
					refArr = append(refArr, Stat{
						Name:  ref.String,
						Value: int(refCount.Int),
					})
				}
				if country.Status != pgtype.Null {

					countryArr = append(countryArr, Stat{
						Name:  country.String,
						Value: int(countryCount.Int),
					})
				}

			}
			rows.Close()
			linkDetail.Browser = browserArr
			linkDetail.Os = osArr
			linkDetail.Device_type = deviceArr
			linkDetail.Referrer = refArr
			linkDetail.Countries = countryArr
			utils.SendResponse(w, http.StatusOK, linkDetail)
		default:
			fmt.Println(err)

			utils.SendResponse(w, http.StatusInternalServerError, utils.SendErrorToClient("Try again later"))
		}
	case pgx.ErrNoRows:
		utils.SendResponse(w, http.StatusNotFound, utils.SendErrorToClient("Not found"))
	default:
		fmt.Println(err)
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
		expTime, _ := time.Parse("2006-01-02 15:04", link.Expiration)
		isAfter := link.Expiration != "" && time.Now().After(expTime)
		link.redirect_url_main = link.LongUrl
		if isAfter {
			link.redirect_url_main = link.NotFoundUrl
		}

		if link.has_country_block {
			country, err := getCountry(r)
			if err != nil {
				return
			}
			var country_count *int
			dbErr := db.QueryRow(context.Background(), "SELECT count(country) from country_block where id=$1 and country_code=$2", country.Registered_country.Names.En, country.Registered_country.Iso_code).Scan(&country_count)
			if dbErr != nil {
				return
			}
			if *country_count > 1 {
				return
			}
			redirect(w, r, link)
			return
		}

		if link.Password != "" {
			temp := template.Must(template.ParseFiles("views/authentication.html"))
			su := map[string]string{
				"Id": link.id,
			}
			temp.Execute(w, su)
			break
		}
		handleRedirectionForDevice(r, w, link)
		break
	case pgx.ErrNoRows:
		http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
		break
	default:
		fmt.Fprintf(w, "Error redirecting")
	}
}

func handleRedirectionForDevice(r *http.Request, w http.ResponseWriter, link LinkAdd) {
	deviceName := getMobileType(r.UserAgent())
	if deviceName == "android" {
		if link.AndroidDeepLink != "" {
			link.redirect_url_main = link.AndroidDeepLink
			redirect(w, r, link)
		}
		if link.PlayStoreLink != "" {
			link.redirect_url_main = link.PlayStoreLink
			redirect(w, r, link)
			return
		}
		return
	}

	if link.IosDeepLink != "" && deviceName == "ios" {
		link.redirect_url_main = link.IosDeepLink
		redirect(w, r, link)
		return
	}

	if link.Mobile_url != "" && isMobile(r.UserAgent()) {
		link.redirect_url_main = link.Mobile_url
		redirect(w, r, link)
		return
	}
	if link.Destkop_url != "" {
		link.redirect_url_main = link.Destkop_url
		redirect(w, r, link)
	}
	redirect(w, r, link)

}

func getMobileType(useragent string) string {
	isAndroid, _ := regexp.MatchString("/android/i", useragent)
	isIos, _ := regexp.MatchString("/iPad|iPhone|iPod/i", useragent)
	if isAndroid {
		return "android"
	}
	if isIos {
		return "ios"
	}
	return "unknown"
}

func isMobile(useragent string) bool {

	isMobile, err := regexp.MatchString(`/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino/i.test(ua)||/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i`, useragent[0:4])
	if err != nil {
		return false
	}
	return isMobile
}

func redirect(w http.ResponseWriter, r *http.Request, link LinkAdd) {
	w.Header().Set("Cache-Control", "no-cache,max-age=0")
	http.Redirect(w, r, link.redirect_url_main, link.HTTPStatus)
	if link.username != "" {
		go updateLogs(r, link)
	}
}

func updateLogs(r *http.Request, result LinkAdd) {

	browser_info := uasurfer.Parse(r.UserAgent())
	u, err := url.Parse(result.LongUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	query := u.Query()
	campaign, medium, source := query.Get("camapgin"), query.Get("meidum"), query.Get("source")
	var insertedId *int
	err = db.QueryRow(context.Background(), "insert into logs(id,campaign,source,medium,os,browser,device_type,referrer)values($1,$2,$3,$4,$5,$6,$7,$8) returning serial_id", result.id, campaign, source, medium, browser_info.OS.Platform.StringTrimPrefix(), browser_info.Browser.Name.StringTrimPrefix(), browser_info.DeviceType.StringTrimPrefix(), r.Referer()).Scan(&insertedId)
	if err != nil {
		var pgError *pgconn.PgError
		errors.As(err, &pgError)
		fmt.Printf("%s at %d col %s %s %s", pgError.Message, pgError.Line, pgError.ColumnName, pgError.Detail, pgError.Hint)
		return
	}

	go insertCountryToLogs(r, *insertedId)
	fmt.Println("written to logs")
}

type Country struct {
	Registered_country struct {
		Iso_code string `json:"iso_code"`
		Names    struct {
			En string `json:"en"`
			Es string `json:"es"`
		} `json:"names"`
	} `json:"registered_country"`
}

func getIp(r *http.Request) string {

	parsedForward := net.ParseIP(strings.Split(r.Header.Get("X-FORWARDED-FOR"), ",")[0])
	parsedRemote := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
	fmt.Println(parsedForward, parsedRemote)
	if parsedForward != nil {
		return parsedForward.String()
	}
	if parsedRemote != nil {
		fmt.Println("got remote")
		return parsedRemote.String()
	}
	return net.ParseIP(r.Header.Get("X-REAL-IP")).String()

}

func getCountry(r *http.Request) (Country, error) {
	var c Country
	ip := getIp(r)
	if ip == "<nil>" {
		return c, errors.New("ip-not-found")
	}

	fmt.Println("fetching country from ip")
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://geolite.info/geoip/v2.1/country/%s?pretty", ip), nil)
	if err != nil {
		fmt.Println(err)
		return c, err
	}
	req.SetBasicAuth("670080", "0vj6DccZISiy45fu")
	resultResp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return c, err
	}
	if resultResp.StatusCode >= 400 && resultResp.StatusCode <= 500 {
		return c, errors.New("Bad request")
	}
	defer resultResp.Body.Close()
	data, err := ioutil.ReadAll(resultResp.Body)
	if err != nil {
		fmt.Println(err)
		return c, err
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Println(err)
		return c, err
	}
	fmt.Println(c.Registered_country)
	return c, nil
}

func insertCountryToLogs(r *http.Request, serial int) {
	country, err := getCountry(r)
	if err != nil {
		return
	}
	_, dbErr := db.Exec(context.Background(), "UPDATE logs SET country=$1 WHERE serial_id=$2", country.Registered_country.Names.En, serial)
	fmt.Println(dbErr)
}

func getRedirectUrl(path string) (LinkAdd, error) {
	var id, username, psswd pgtype.Varchar
	var longUrl, notFoundUrl, android, ios, play_store, mobile_url, desktop_url pgtype.Text
	var exp_time pgtype.Timestamptz
	var country_block, device_redirect pgtype.Bool
	var http_status *int

	err := db.QueryRow(context.Background(), "select id,url,username,password,not_found_url,android_deep_link,ios_deep_link,play_store_link,country_block,device_type_redirect,http_status,mobile_url,desktop_url,others_url,expiration from urls left join expiration on urls.id=expiration.id where urls.id=$1", path).Scan(&id, &longUrl, &username, &psswd, &notFoundUrl, &android, &ios, &play_store, &country_block, &device_redirect, &mobile_url, &desktop_url, &exp_time)
	if err != nil {
		return LinkAdd{}, err
	}
	return LinkAdd{id: id.String, LongUrl: longUrl.String, Password: psswd.String, NotFoundUrl: notFoundUrl.String, AndroidDeepLink: android.String, IosDeepLink: ios.String, PlayStoreLink: play_store.String, HTTPStatus: *http_status, Mobile_url: mobile_url.String, Destkop_url: desktop_url.String, has_country_block: country_block.Bool, has_device_redirect: device_redirect.Bool, username: username.String}, nil

}

func setId(r *http.Request, reqBody LinkAdd, session auth.Session) (string, error) {
	/*
		var timeT time.Time
		if reqBody.Expiration != "" {
			timeT, err := time.Parse("2006-01-02T15:04", reqBody.Expiration)
			if err != nil {
				return "", err
			}
		}
	*/
	psswdHash := reqBody.Password
	if psswdHash != "" {
		hash, err := auth.GeneratePsswdHash(reqBody.Password)
		if err != nil {
			return "", err
		}
		psswdHash = string(hash)
	}

	var uniqId string

	err := db.QueryRow(context.Background(), "SELECT insertLongUrl($1,$2,$3,$4)", session.Username, psswdHash, reqBody).Scan(&uniqId)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return uniqId, nil

}
func onNotify(c *pgconn.PgConn, n *pgconn.Notice) {

	fmt.Println("Message:", *n)
}
