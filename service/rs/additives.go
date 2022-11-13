package rs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

func (api *RestApi) handleAdditives() http.HandlerFunc {
	return func(_w http.ResponseWriter, _r *http.Request) {
		r := &MyRequest{_r, slashAdditives}
		w := &MyResponseWriter{_w}

		code := r.pathParam()
		if len(code) > 0 {
			a, err := api.provider.Additives.FetchOne(code, r.locale(api))
			if err != nil {
				w.writeError(err)
			} else {
				a.Url = r.relUrl(a.Code)
				w.writeJson(a)
			}
		} else {
			var (
				additives []*db.AdditiveMeta
				err       error
			)

			category := strings.TrimSpace(r.URL.Query().Get("category"))
			if len(category) > 0 {
				catId, err := strconv.Atoi(category)
				if err != nil {
					w.writeError(err)
					return
				}
				additives, err = api.provider.Additives.FetchAllByCategory(catId, r.locale(api))

			} else {
				additives, err = api.provider.Additives.FetchAll(r.locale(api))
			}

			if err != nil {
				w.writeError(err)
			} else {
				for _, a := range additives {
					a.Url = r.relUrl(a.Code)
				}
				w.writeJson(additives)
			}
		}
	}
}
