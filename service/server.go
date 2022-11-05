package service

import (
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

type ServerContext struct {
	router *http.ServeMux
	// workerPool *WokerPool
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	if err = db.InitEadDb(config.DatabasePath); err != nil {
		return err
	}

	attachApi(ctx)

	fmt.Printf("Serving at %s:%d ...\n", config.ListenAddress,
		config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
