package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

const (
	API_INDEX   = "/api"
	API_LOCALES = API_INDEX + "/locales"
)

type Api struct {
	ctx *ServerContext
}

func attachApi(serverCtx *ServerContext) {
	fmt.Printf("Attaching API junctions at %s ...\n", API_INDEX)

	api := &Api{serverCtx}
	api.ctx.router.HandleFunc(API_INDEX, api.handleIndex())
	api.ctx.router.HandleFunc(API_INDEX+"/", api.handleIndex())
	// api.ctx.router.HandleFunc(API_DEVICES+"/search", api.handleDeviceSearch())
	// api.ctx.router.HandleFunc(API_DEVICES+"/search/", api.handleDeviceSearch())
	api.ctx.router.HandleFunc(API_LOCALES, api.handleLocales())
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := map[string]string{
			API_INDEX:   "Shows this",
			API_LOCALES: "Fetches a list of all locales",
			// API_DEVICES + "/:product":        "Fetches a single Apple device by product name",
			// API_DEVICES + "/search?key=:key": "Fetches a list of devices given a key parameter",
			// API_UPDATES + "/:product":        "Fetches device updates by product name",
		}
		resp, _ := json.Marshal(routes)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (api *Api) handleLocales() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
		} else {
			resp, _ := json.Marshal(locales)
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
		}
		// }
	}
}
