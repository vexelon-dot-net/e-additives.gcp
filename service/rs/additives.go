package rs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
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
				a.Url = getUrl(r, slashAdditives, a.Code)
				w.writeJson(a)
			}
		} else {
			var additives []*db.AdditiveMeta

			category := strings.TrimSpace(r.URL.Query().Get("category"))
			if len(category) > 0 {
				catId, err := strconv.Atoi(category)
				if err != nil {
					w.writeError(err)
					return
				}
				additives, err = api.provider.Additives.FetchAllByCategory(catId, loc)

			} else {
				additives, err = api.provider.Additives.FetchAll(loc)
			}

			if err != nil {
				w.writeError(err)
			} else {
				for _, a := range additives {
					a.Url = getUrl(r, slashAdditives, a.Code)
				}
				w.writeJson(additives)
			}
		}
	}
}
