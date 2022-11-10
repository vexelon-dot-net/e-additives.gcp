package rs

import (
	"net/http"
)

func (api *RestApi) handleCategories() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		id, err := getKeyParam(r, slashCategories)
		if err != nil {
			w.writeError(err)
			return
		}

		// TODO
		locales, err := api.provider.Locales.FetchAll()
		if err != nil {
			w.writeError(err)
			return
		}
		loc := *locales[1]

		if id > 0 {
			cat, err := api.provider.Additives.Categories.FetchOne(id, loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(cat)
			}
		} else {
			// TODO: decorate add urls for each item
			categories, err := api.provider.Additives.Categories.FetchAll(loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(categories)
			}
		}
	}
}
