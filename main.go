package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/service"
)

const (
	HEART = "\u2764"
)

var (
	ctx = context.Background()
)

//var db *sql.DB

// func InitDb(path string) (err error) {
// 	db, err := sql.Open("sqlite3", path)
// 	if err != nil {
// 		return fmt.Errorf("Error opening database at %s : %w", path, err)
// 	}

// 	if err = db.Ping(); err != nil {
// 		return fmt.Errorf("Failed database ping : %w", err)
// 	}

// 	return nil
// }

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Querying additives ...")

	db, err := sql.Open("sqlite3", "./data/eadditives.sqlite3")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	rows, err := db.Query(`
	SELECT c.id, p.name, p.description, p.last_update,
    (SELECT COUNT(id) FROM ead_Additive as a WHERE a.category_id=c.id) as additives
	FROM ead_AdditiveCategory as c
	LEFT JOIN ead_AdditiveCategoryProps as p ON p.category_id = c.id
	WHERE p.locale_id = $1`, "3")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var desc string
		var lu string
		var count int
		err := rows.Scan(&id, &name, &desc, &lu, &count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// device.LastCheckedOnParsed, err = time.Parse(DATE_TIME_LAYOUT, device.LastCheckedOn)
		// if err != nil {
		// 	return nil, fmt.Errorf("Error parsing last update check date time '%s' : %w", device.LastCheckedOn, err)
		// }

		fmt.Fprintf(w, "%d %s %s %s", id, name, desc, lu)
	}
}

func main() {
	fmt.Printf("%s e-additives API service v%s %s\n\n", HEART, config.VERSION, HEART)

	if err := config.ParseConfig(); err != nil {
		log.Fatalf("Config error: %v\n", err)
	}

	if err := service.Start(); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
