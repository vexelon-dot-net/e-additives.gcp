package service

import (
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/db"
)

type ServerContext struct {
	config *config.Config
	router *http.ServeMux
	// workerPool *WokerPool
}

func NewServer(config *config.Config) *ServerContext {
	return &ServerContext{
		config,
		http.NewServeMux(),
	}
}

func (sc *ServerContext) Run() (err error) {
	if err = db.Open(sc.config.DatabasePath); err != nil {
		return err
	}

	attachApi(sc)

	fmt.Printf("Serving at %s:%d ...\n", sc.config.ListenAddress,
		sc.config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		sc.config.ListenAddress, sc.config.ListenPort), sc.router); err != nil {
		return err
	}

	return nil
}
