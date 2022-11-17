package rs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

const (
	slashIndex           = "/api"
	slashLocales         = slashIndex + "/locales"
	slashCategories      = slashIndex + "/categories"
	slashAdditives       = slashIndex + "/additives"
	slashAdditivesSearch = slashIndex + "/additives/search"

	paramJSONP    = "callback"
	paramLocale   = "locale"
	paramCategory = "category"
	paramKeyword  = "q"
)

type RestApi struct {
	provider      *db.DBProvider
	apiKeys       map[string]*db.ApiKey
	locales       map[string]*db.Locale
	defaultLocale db.Locale
}

func AttachRestApi(router *http.ServeMux, provider *db.DBProvider) error {
	fmt.Printf("Attaching http API at %s ...\n", slashIndex)

	// prefetch API keys to ease on db queries
	apiKeysMap := make(map[string]*db.ApiKey)
	apiKeys, err := provider.ApiKeys.FetchAll()
	if err != nil {
		return fmt.Errorf("Error prefetching api keys: %w", err)
	}
	for _, apiKey := range apiKeys {
		fmt.Printf("\t|-> Cached api key '%s...' for '%s'\n", apiKey.Key[:7], apiKey.Host)
		apiKeysMap[apiKey.Key] = apiKey
	}

	// prefetch localesMap to ease on db queries
	localesMap := make(map[string]*db.Locale)
	locales, err := provider.Locales.FetchAll()
	if err != nil {
		return fmt.Errorf("Error prefetching locales: %w", err)
	}
	for _, loc := range locales {
		fmt.Printf("\t|-> Cached locale '%s'\n", loc.Code)
		localesMap[loc.Code] = loc
	}

	api := &RestApi{provider, apiKeysMap, localesMap, *localesMap["en"]}
	router.HandleFunc(slashIndex, api.handleIndex())
	router.HandleFunc(slashIndex+"/", api.handleIndex())
	router.HandleFunc(slashLocales, api.handleLocales())
	router.HandleFunc(slashLocales+"/", api.handleLocales())
	router.HandleFunc(slashCategories, api.handleCategories())
	router.HandleFunc(slashCategories+"/", api.handleCategories())
	router.HandleFunc(slashAdditives, api.handleAdditives())
	router.HandleFunc(slashAdditives+"/", api.handleAdditives())
	router.HandleFunc(slashAdditivesSearch, api.handleAdditivesSearch())
	router.HandleFunc(slashAdditivesSearch+"/", api.handleAdditivesSearch())

	return nil
}

func (api *RestApi) getLocale(code string) db.Locale {
	loc, ok := api.locales[code]
	if ok {
		return *loc
	}
	return api.defaultLocale
}

func (api *RestApi) isValidApiKey(key string, host string) bool {
	apiKey, ok := api.apiKeys[key]
	fmt.Printf("APIK: %s \n", key)
	if ok {
		return apiKey.Host == "*" || strings.EqualFold(apiKey.Host, host)
	}
	return false
}

func (api *RestApi) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashIndex, w, r)
		routes := map[string]string{
			slashIndex:                     "Fetches this list of junctions",
			slashLocales:                   "Fetches list of locales",
			slashCategories:                "Fetches list of additive categories",
			slashCategories + "/:category": "Fetches a single additive category",
			slashAdditives:                 "Fetches list of additives",
			slashAdditives + "/:code":      "Fetches a single additive by code",
			slashAdditivesSearch:           "Search additives by keyword",
		}
		h.writeJson(routes)
	}
}
