package rs

import (
	"net/http"
	"strconv"
)

func (api *RestApi) handleCategories() http.HandlerFunc {
	return func(_w http.ResponseWriter, _r *http.Request) {
		r := &MyRequest{_r, slashCategories}
		w := &MyResponseWriter{_w}

		id, err := r.idParam()
		if err != nil {
			w.writeError(err)
			return
		}

		if id > 0 {
			cat, err := api.provider.Additives.Categories.FetchOne(id, r.locale(api))
			if err != nil {
				w.writeError(err)
			} else {
				cat.Url = r.relUrl(strconv.Itoa(cat.Category))
				w.writeJson(cat)
			}
		} else {
			categories, err := api.provider.Additives.Categories.FetchAll(r.locale(api))
			if err != nil {
				w.writeError(err)
			} else {
				for _, cat := range categories {
					cat.Url = r.relUrl(strconv.Itoa(cat.Category))
				}
				w.writeJson(categories)
			}
		}
	}
}
