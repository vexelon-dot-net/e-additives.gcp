package rs

import (
	"net/http"
)

func (api *RestApi) handleAdditives() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		code := getKeyParam(r, slashAdditives)

		// TODO
		locales, err := api.provider.Locales.FetchAll()
		if err != nil {
			w.writeError(err)
			return
		}
		loc := *locales[1]

		if len(code) > 0 {
			a, err := api.provider.Additives.FetchOne(code, loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(a)
			}
		} else {
			// TODO: decorate add urls for each item
			additives, err := api.provider.Additives.FetchAll(loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(additives)
			}
		}
	}
}
