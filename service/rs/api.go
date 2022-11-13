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

	paramJSONP  = "callback"
	paramLocale = "locale"
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
	loc, ok := api.locales[code]
	if ok {
		return *loc
	}
	return api.defaultLocale
}

func (api *RestApi) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashIndex, w, r)
		routes := map[string]string{
			slashIndex:                     "Fetches this list of junctions",
			slashLocales:                   "Fetches list of locales",
			slashCategories:                "Fetches list of additive categories",
			slashCategories + "/:category": "Fetches a single additive category by id",
			slashAdditives:                 "Fetches list of additives",
			slashAdditives + "/:code":      "Fetches a single additive by code",
		}
		h.writeJson(routes)
	}
}
