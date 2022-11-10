package service

import (
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/db"
	"github.com/vexelon-dot-net/e-additives.gcp/service/rs"
)

type ServerContext struct {
	config   *config.Config
	router   *http.ServeMux
	provider *db.DBProvider
	// workerPool *WokerPool
}

func NewServer(config *config.Config) *ServerContext {
	return &ServerContext{
		config,
		http.NewServeMux(),
		nil,
	}
}

func (sc *ServerContext) Run() (err error) {
	sc.provider, err = db.NewProvider(sc.config.DatabasePath)
	if err != nil {
		return err
	}

	_ = rs.NewRestApi(sc.router, sc.provider)

	fmt.Printf("Serving at %s:%d ...\n", sc.config.ListenAddress,
		sc.config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		sc.config.ListenAddress, sc.config.ListenPort), sc.router); err != nil {
		return err
	}

	return nil
}
