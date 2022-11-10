package rs

import (
	"net/http"
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
			cat, err := api.provider.Locales.FetchOne(id)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(cat)
			}
		} else {
			locales, err := api.provider.Locales.FetchAll()
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(locales)
			}
		}
	}
}
