package rs

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

type MyRequest struct {
	*http.Request
	junction string
}

type MyResponseWriter struct {
	http.ResponseWriter
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

func (r *MyRequest) pathParam() string {
	return strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, r.junction), "/")
}

func (r *MyRequest) idParam() (int, error) {
	key := r.pathParam()
	if len(key) > 0 {
		id, err := strconv.Atoi(key)
		if err != nil {
			return 0, fmt.Errorf("Error parsing id '%s': %w", key, err)
		}
		return id, nil
	}

	return 0, nil
}

func (r *MyRequest) relUrl(id string) string {
	return fmt.Sprintf("%s%s/%s", r.Referer(), r.junction, id)
}

func (r *MyRequest) locale(api *RestApi) db.Locale {
	return api.getLocale(r.URL.Query().Get("locale"))
}
