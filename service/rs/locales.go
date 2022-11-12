package rs

import (
	"net/http"
)

func (api *RestApi) handleLocales() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		code := getKeyParam(r, slashLocales)

		if len(code) > 0 {
			loc, err := api.provider.Locales.FetchOne(code)
			if err != nil {
				w.writeError(err)
			} else {
				loc.Url = getUrl(r, slashLocales, loc.Code)
				w.writeJson(loc)
			}
		} else {
			locales, err := api.provider.Locales.FetchAll()
			if err != nil {
				w.writeError(err)
			} else {
				for _, loc := range locales {
					loc.Url = getUrl(r, slashLocales, loc.Code)
				}
				w.writeJson(locales)
			}
		}
	}
}
