package rs

import (
	"net/http"
	"strconv"
)

func (api *RestApi) handleCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := newHandlerContext(api, slashCategories, w, r)

		id, err := h.idParam()
		if err != nil {
			h.writeError(err)
			return
		}

		if id > 0 {
			cat, err := api.provider.Additives.Categories.FetchOne(id, h.locale())
			if err != nil {
				h.writeError(err)
			} else {
				cat.Url = h.relUrl(strconv.Itoa(cat.Category))
				h.writeJson(cat)
			}
		} else {
			categories, err := api.provider.Additives.Categories.FetchAll(h.locale())
			if err != nil {
				h.writeError(err)
			} else {
				for _, cat := range categories {
					cat.Url = h.relUrl(strconv.Itoa(cat.Category))
				}
				h.writeJson(categories)
			}
		}
	}
}
