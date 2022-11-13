package rs

import (
	"net/http"
)

func (api *RestApi) handleLocales() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashLocales, w, r)

		code := h.pathParam()
		if len(code) > 0 {
			loc, err := api.provider.Locales.FetchOne(code)
			if err != nil {
				h.writeError(err)
			} else {
				loc.Url = h.relUrl(loc.Code)
				h.writeJson(loc)
			}
		} else {
			locales, err := api.provider.Locales.FetchAll()
			if err != nil {
				h.writeError(err)
			} else {
				for _, loc := range locales {
					loc.Url = h.relUrl(loc.Code)
				}
				h.writeJson(locales)
			}
		}
	}
}
