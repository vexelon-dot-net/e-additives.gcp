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

		if err := h.verifyAuth(); err != nil {
			h.writeError(err)
			return
		}

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

			category := strings.TrimSpace(h.qvCache.Get(paramCategory))
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

func (api *RestApi) handleAdditivesSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashAdditives, w, r)

		if err := h.verifyAuth(); err != nil {
			h.writeError(err)
			return
		}

		keyword := strings.TrimSpace(h.qvCache.Get(paramKeyword))
		if len(keyword) > 0 {
			additives, err := api.provider.Additives.Search(keyword, h.locale())
			if err != nil {
				h.writeError(err)
				return
			}

			for _, a := range additives {
				a.Url = h.relUrl(a.Code)
			}

			h.writeJson(additives)
		} else {
			h.writeJson(make([]interface{}, 0))
		}
	}
}
