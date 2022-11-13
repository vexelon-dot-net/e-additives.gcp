package rs

import (
	"fmt"
	"net/http"
	"strings"

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
	locales       map[string]*db.Locale
	defaultLocale db.Locale
}

func AttachRestApi(router *http.ServeMux, provider *db.DBProvider) error {
	fmt.Printf("Attaching http API at %s ...\n", slashIndex)

	locales := make(map[string]*db.Locale)
	fetched, err := provider.Locales.FetchAll()
	if err != nil {
		return fmt.Errorf("Error fetching default 'en' locale: %w", err)
	}
	for _, loc := range fetched {
		locales[loc.Code] = loc
	}

	api := &RestApi{provider, locales, *locales["en"]}
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

func (api *RestApi) getLocale(code string) db.Locale {
	//code := r.URL.Query().Get("locale")
	if len(code) > 0 {
		loc, ok := api.locales[strings.TrimSpace(code)]
		if ok {
			return *loc
		}
	}
	return api.defaultLocale
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
