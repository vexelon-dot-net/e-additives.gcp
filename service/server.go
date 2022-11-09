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

func NewServer() *ServerContext {
	return &ServerContext{http.NewServeMux()}
}

func (sc *ServerContext) Start() (err error) {
	if err = db.Open(config.DatabasePath); err != nil {
		return err
	}

	attachApi(sc)

	fmt.Printf("Serving at %s:%d ...\n", config.ListenAddress,
		config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), sc.router); err != nil {
		return err
	}

	return nil
}
