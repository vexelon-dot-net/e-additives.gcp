package rs

import (
	"bytes"
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
	w := h.w
	resp, _ := json.Marshal(data)
	if h.qvCache.Has(paramJSONP) {
		var buf bytes.Buffer
		buf.WriteString(h.qvCache.Get(paramJSONP))
		buf.WriteString("(")
		buf.Write(resp)
		buf.WriteString(")")

		w.Header().Set("Content-Type", "application/javascript")
		w.Write(buf.Bytes())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (h *HandlerContext) pathParam() string {
	urlPath := h.r.URL.Path
	idx := strings.Index(urlPath, h.path)
	if idx > -1 {
		return strings.TrimPrefix(urlPath[idx+len(h.path):], "/")
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
