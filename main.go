package  main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgconn"
	"github.com/shikhar0507/requestJSON"
	"log"
	rand2 "math/rand"
	"net/http"
	"strings"
	"time"

	//"net/url"
)
var db *pgxpool.Pool
type SuccesRes struct {
	Status int `json:"status"`
	Url string `json:"url"`
}
type stop struct {
	error
}

func main() {

	dbpool, err := pgxpool.Connect(context.Background(),"postgres://pujacapital:Dhruv.2017@localhost:5432/url_short")
	if err != nil {
		log.Fatal(err)
	}
	db = dbpool
	defer db.Close()
	http.HandleFunc("/",handleHome)
	http.HandleFunc("/shorten",handleShortner)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080",nil))

}
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w,r,"index.html")
		return
	}
	if r.URL.Path == "/favicon.ico" {
		return
	}
	fmt.Println(r.URL.Path)
	id := strings.Split(r.URL.Path,"/")[1]
	fmt.Println(id)
	var queryId string
	var originalUrl string
	err := db.QueryRow(context.Background(),"select * from urls where id=$1",id).Scan(&queryId,&originalUrl)
	if err != nil {
		fmt.Println(err)
		http.Error(w,"something went wrong",http.StatusInternalServerError)
		return
	}
	http.Redirect(w,r,originalUrl,http.StatusPermanentRedirect)

}

func handleShortner(w http.ResponseWriter, r *http.Request) {

	type ReqURL struct {
		 Url string
	}

	var reqURL ReqURL
	result := requestJSON.Decode(w,r,&reqURL)
	if result.Status != 200 {
		sendJSONResponse(w, result.Status, result)
		return
	}


	id,err := setId(reqURL.Url)

	if err != nil {
		if err.Error() == "failed to assign a unique value" {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		http.Error(w,"Something went wrong",http.StatusInternalServerError)
	}
	fmt.Println("used id",id)
	succ := SuccesRes{Status: 200,Url: id}
	sendJSONResponse(w,200,succ)

}

func setId(reqURL string) (string,error) {
	value := createId()
	//value := "RsWxP"
	mainErr := retry(3,1000, func() error {
		_,err := db.Exec(context.Background(),"insert into urls values($1,$2)",value,reqURL)
		if err == nil {
			return nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err,&pgErr) {
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
		return "",mainErr
	}
	return value,nil


}


func retry(count int,sleep time.Duration, f func() error) error {
	err := f()
	if err != nil {
		if s,ok := err.(stop); ok {
			return s.error
		}
		count--
		if count > 0 {
			time.Sleep(sleep)
			return retry(count,1*sleep,f)
		}
		return err
	}
	return nil
}

func createId() string{
	  letterString := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	  result := ""

	  for i := 0; i < 6;i++ {
	  	randStr := letterString[rand2.Intn(len(letterString))]
	  	result  = result+string(randStr)
	  }
	  return result
}

func sendJSONResponse(w http.ResponseWriter,status int, s interface{})  {
	w.WriteHeader(status)
	j, err:= json.Marshal(s)
	if err != nil {
		http.Error(w,"Something went wrong",http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w,string(j))

}