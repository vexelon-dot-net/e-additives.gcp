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
)

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

func getKeyParam(r *http.Request, junction string) string {
	return strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, junction), "/")
}

func getIdParam(r *http.Request, junction string) (int, error) {
	key := getKeyParam(r, junction)

	if len(key) > 0 {
		id, err := strconv.Atoi(key)
		if err != nil {
			return 0, fmt.Errorf("Error parsing id '%s': %w", key, err)
		}
		return id, nil
	}

	return 0, nil
}
