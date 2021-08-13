package api

import (
	"fmt"
	"net/http"
	"strings"
	"url-shortner/auth"
	"url-shortner/utils"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shikhar0507/requestJSON"
)

func init() {

}

type ReqURL struct {
	Url string `json: url`
}

func getResourceId(path string, sep string) string {
	linkid := strings.TrimPrefix(path, sep)
	if linkid == path {
		return ""
	}
	return linkid

}

func HandleCampaigns(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {

}

// / GET , / POST
// /id/ GET ,DELETE , PUT

func HandleLinks(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {
	username, _, err := auth.GetSession(r, db)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			resp := utils.Response{Status: http.StatusUnauthorized, Message: "Unauthorized user"}
			utils.SendResponse(w, http.StatusUnauthorized, resp)
			break
		default:
			resp := utils.Response{Status: http.StatusInternalServerError, Message: "Try again later"}
			utils.SendResponse(w, http.StatusInternalServerError, resp)
		}
		return
	}

	resId := getResourceId(r.URL.Path, "/api/v1/links/")
	fmt.Println("resid", resId)
	if resId == "" {
		switch r.Method {
		case http.MethodGet:
			getLinks(w, r, db, username)
			break
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		getLink(w, r, db, resId)
		break
	case http.MethodDelete:
		break
	case http.MethodPut:
		break
	}

}

func handleShortner(w http.ResponseWriter, r *http.Request, db *pgxpool.Pool) {

	optns := utils.HandleCors(w, r, []string{"http.MethodPost"})
	if optns == true {
		return
	}

	var reqURL ReqURL

	result := requestJSON.Decode(w, r, &reqURL)

	if result.Status != 200 {

		utils.SendResponse(w, result.Status, result)
		return
	}

	id, err := setId(r, reqURL.Url, db)
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

	succ := utils.SuccesRes{Status: 200, Url: "http://localhost:8080/" + id}
	utils.SendResponse(w, 200, succ)

}
