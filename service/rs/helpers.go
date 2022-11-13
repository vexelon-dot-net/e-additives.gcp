package rs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

type HandlerContext struct {
	api     *RestApi
	path    string
	w       http.ResponseWriter
	r       *http.Request
	qvCache url.Values
}

func newHandlerContext(api *RestApi, path string, w http.ResponseWriter, r *http.Request) *HandlerContext {
	return &HandlerContext{
		api,
		path,
		w,
		r,
		r.URL.Query(),
	}
}

func (h *HandlerContext) writeError(err error) {
	w := h.w
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

func (h *HandlerContext) writeJson(data interface{}) {
	resp, _ := json.Marshal(data)
	h.w.Header().Set("Content-Type", "application/json")
	h.w.Write(resp)
}

func (h *HandlerContext) pathParam() string {
	urlPath := h.r.URL.Path
	idx := strings.Index(urlPath, h.path)
	if idx > -1 {
		return strings.TrimPrefix(urlPath[idx+len(h.path)+1:], "/")
	}
	return ""
}

func (h *HandlerContext) idParam() (int, error) {
	key := h.pathParam()
	if len(key) > 0 {
		id, err := strconv.Atoi(key)
		if err != nil {
			return 0, fmt.Errorf("Error parsing id '%s': %w", key, err)
		}
		return id, nil
	}

	return 0, nil
}

func (h *HandlerContext) relUrl(id string) string {
	return fmt.Sprintf("%s%s/%s", h.r.Referer(), h.path, id)
}

func (h *HandlerContext) locale() db.Locale {
	return h.api.getLocale(h.qvCache.Get(paramLocale))
}
