package rs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

func (api *RestApi) handleAdditives() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashAdditives, w, r)

		code := h.pathParam()
		if len(code) > 0 {
			a, err := api.provider.Additives.FetchOne(code, h.locale())
			if err != nil {
				h.writeError(err)
			} else {
				a.Url = h.relUrl(a.Code)
				h.writeJson(a)
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
					h.writeError(err)
					return
				}
				additives, err = api.provider.Additives.FetchAllByCategory(catId, h.locale())

			} else {
				additives, err = api.provider.Additives.FetchAll(h.locale())
			}

			if err != nil {
				h.writeError(err)
			} else {
				for _, a := range additives {
					a.Url = h.relUrl(a.Code)
				}
				h.writeJson(additives)
			}
		}
	}
}
