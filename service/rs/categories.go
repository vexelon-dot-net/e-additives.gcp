package rs

import (
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
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
		locales, err := db.FetchAllLocales()
		if err != nil {
			w.writeError(err)
			return
		}
		loc := *locales[1]

		if id > 0 {
			cat, err := db.FetchOneCategory(id, loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(cat)
			}
		} else {
			// TODO: decorate add urls for each item
			categories, err := db.FetchAllCategories(loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(categories)
			}
		}
	}
}
