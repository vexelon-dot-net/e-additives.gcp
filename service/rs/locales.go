package rs

import (
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

func (api *RestApi) handleLocales() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		// product := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, API_LOCALES), "/")
		// if len(product) > 0 {

		// 	device, err := db.FetchDeviceByProduct(product)
		// 	if errors.Is(err, sql.ErrNoRows) {
		// 		w.WriteHeader(http.StatusNotFound)
		// 		fmt.Fprintf(w, "Error: %v", err)
		// 	} else if err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		fmt.Fprintf(w, "Error: %v", err)
		// 	} else {
		// 		resp, _ := json.Marshal(device)
		// 		w.Header().Set("Content-Type", "application/json")
		// 		w.Write(resp)
		// 	}
		// } else {

		locales, err := db.FetchAllLocales()
		if err != nil {
			w.writeError(err)
		} else {
			w.writeJson(locales)
		}
		// }
	}
}
