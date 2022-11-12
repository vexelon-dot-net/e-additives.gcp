package www

import (
	"fmt"
	"net/http"
)

func AttachWWW(router *http.ServeMux, isDevMode bool) {
	if isDevMode {
		fmt.Println("Attaching dev mode web app ...")
		router.Handle("/", http.FileServer(http.Dir("www")))
	} else {
		fmt.Println("Redirect / to /api ...")
		router.Handle("/", http.RedirectHandler("/api", http.StatusMovedPermanently))
	}
}
