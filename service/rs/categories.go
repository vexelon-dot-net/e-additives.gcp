package rs

import (
	"net/http"
	"strconv"
)

func (api *RestApi) handleCategories() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		id, err := getIdParam(r, slashCategories)
		if err != nil {
			w.writeError(err)
			return
		}

		if id > 0 {
			cat, err := api.provider.Additives.Categories.FetchOne(id, api.defaultLocale)
			if err != nil {
				w.writeError(err)
			} else {
				cat.Url = getUrl(r, slashCategories, strconv.Itoa(cat.Category))
				w.writeJson(cat)
			}
		} else {
			categories, err := api.provider.Additives.Categories.FetchAll(api.defaultLocale)
			if err != nil {
				w.writeError(err)
			} else {
				for _, cat := range categories {
					cat.Url = getUrl(r, slashCategories, strconv.Itoa(cat.Category))
				}
				w.writeJson(categories)
			}
		}
	}
}
