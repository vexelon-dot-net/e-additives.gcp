package rs

import (
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

func (api *RestApi) handleLocales() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}

		id, err := getKeyParam(r, slashLocales)
		if err != nil {
			w.writeError(err)
			return
		}

		if id > 0 {
			cat, err := db.FetchOneLocale(id)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(cat)
			}
		} else {
			locales, err := db.FetchAllLocales()
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(locales)
			}
		}
	}
}
