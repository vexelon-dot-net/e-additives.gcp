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
	slashAdditives  = slashIndex + "/additives"
)

type RestApi struct {
	provider      *db.DBProvider
	defaultLocale db.Locale
}

func AttachRestApi(router *http.ServeMux, provider *db.DBProvider) error {
	fmt.Printf("Attaching http API at %s ...\n", slashIndex)

	loc, err := provider.Locales.FetchOne("en")
	if err != nil {
		return fmt.Errorf("Error fetching default 'en' locale: %w", err)
	}

	api := &RestApi{provider, *loc}
	router.HandleFunc(slashIndex, api.handleIndex())
	router.HandleFunc(slashIndex+"/", api.handleIndex())
	router.HandleFunc(slashLocales, api.handleLocales())
	router.HandleFunc(slashLocales+"/", api.handleLocales())
	router.HandleFunc(slashCategories, api.handleCategories())
	router.HandleFunc(slashCategories+"/", api.handleCategories())
	router.HandleFunc(slashAdditives, api.handleAdditives())
	router.HandleFunc(slashAdditives+"/", api.handleAdditives())

	return nil
}

func (api *RestApi) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &MyResponseWriter{w}
		routes := map[string]string{
			slashIndex:                     "Fetches this list of junctions",
			slashLocales:                   "Fetches list of locales",
			slashCategories:                "Fetches list of additive categories",
			slashCategories + "/:category": "Fetches a single additive category by id",
			slashAdditives:                 "Fetches list of additives",
			slashAdditives + "/:code":      "Fetches a single additive by code",
		}
		writer.writeJson(routes)
	}
}
