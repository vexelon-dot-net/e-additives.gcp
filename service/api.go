package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

const (
	API_INDEX      = "/api"
	API_LOCALES    = API_INDEX + "/locales"
	API_CATEGORIES = API_INDEX + "/categories"
)

type Api struct {
	ctx *ServerContext
}

type MyResponseWriter struct {
	http.ResponseWriter
}

func attachApi(serverCtx *ServerContext) {
	fmt.Printf("Attaching API junctions at %s ...\n", API_INDEX)

	api := &Api{serverCtx}
	api.ctx.router.HandleFunc(API_INDEX, api.handleIndex())
	api.ctx.router.HandleFunc(API_INDEX+"/", api.handleIndex())
	api.ctx.router.HandleFunc(API_LOCALES, api.handleLocales())
	api.ctx.router.HandleFunc(API_CATEGORIES, api.handleCategories())
	api.ctx.router.HandleFunc(API_CATEGORIES+"/", api.handleCategories())
}

func getIdParam(r *http.Request, junction string) (int, error) {
	parsed := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, junction), "/")
	if len(parsed) > 0 {
		id, err := strconv.Atoi(parsed)
		if err != nil {
			return 0, fmt.Errorf("Error parsing '%s': %w", parsed, err)
		}
		return id, nil
	}
	return 0, nil
}

func (w *MyResponseWriter) writeError(err error) {
	if errors.Is(err, strconv.ErrSyntax) {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: %v\n", err)
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusBadRequest))
	} else if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(w, "Error: %v", err)
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusNotFound))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusInternalServerError))
	}
}

func (w *MyResponseWriter) writeJson(data interface{}) {
	resp, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &MyResponseWriter{w}
		routes := map[string]string{
			API_INDEX:               "Fetches this list of junctions",
			API_LOCALES:             "Fetches list of locales",
			API_CATEGORIES:          "Fetches list of additive categories",
			API_CATEGORIES + "/:id": "Fetches a single additive category by id",
			// API_DEVICES + "/search?key=:key": "Fetches a list of devices given a key parameter",
		}
		writer.writeJson(routes)
	}
}

func (api *Api) handleLocales() http.HandlerFunc {
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

func (api *Api) handleCategories() http.HandlerFunc {
	return func(_w http.ResponseWriter, r *http.Request) {
		w := &MyResponseWriter{_w}
		id, err := getIdParam(r, API_CATEGORIES)
		if err != nil {
			w.writeError(err)
			return
		}

		// TODO
		locales, err := db.FetchAllLocales()
		if err != nil {
			w.writeError(err)
			return
		}
		loc := *locales[1]

		if id > 0 {
			cat, err := db.FetchOneCategory(id, loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(cat)
			}
		} else {
			// TODO: decorate add urls for each item
			categories, err := db.FetchAllCategories(loc)
			if err != nil {
				w.writeError(err)
			} else {
				w.writeJson(categories)
			}
		}
	}
}
