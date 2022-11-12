package rs

import (
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

const (
	slashIndex      = "/api"
	slashLocales    = slashIndex + "/locales"
	slashCategories = slashIndex + "/categories"
)

type RestApi struct {
	provider *db.DBProvider
}

func NewRestApi(router *http.ServeMux, provider *db.DBProvider) *RestApi {
	fmt.Printf("Attaching http API at %s ...\n", slashIndex)

	api := &RestApi{provider}
	router.HandleFunc(slashIndex, api.handleIndex())
	router.HandleFunc(slashIndex+"/", api.handleIndex())
	router.HandleFunc(slashLocales, api.handleLocales())
	router.HandleFunc(slashLocales+"/", api.handleLocales())
	router.HandleFunc(slashCategories, api.handleCategories())
	router.HandleFunc(slashCategories+"/", api.handleCategories())

	return api
}

func (api *RestApi) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &MyResponseWriter{w}
		routes := map[string]string{
			slashIndex:               "Fetches this list of junctions",
			slashLocales:             "Fetches list of locales",
			slashCategories:          "Fetches list of additive categories",
			slashCategories + "/:id": "Fetches a single additive category by id",
			// API_DEVICES + "/search?key=:key": "Fetches a list of devices given a key parameter",
		}
		writer.writeJson(routes)
	}
}