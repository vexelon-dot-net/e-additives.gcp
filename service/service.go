package service

import (
	"fmt"
	"net/http"

	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/db"
	"github.com/vexelon-dot-net/e-additives.gcp/service/rs"
	"github.com/vexelon-dot-net/e-additives.gcp/service/www"
)

type Service struct {
	config   config.Config
	router   *http.ServeMux
	provider *db.DBProvider
}

func New(config config.Config) Service {
	return Service{
		config,
		http.NewServeMux(),
		nil,
	}
}

func (sc *Service) Run() (err error) {
	if sc.provider, err = db.NewProvider(sc.config.DatabasePath); err != nil {
		return err
	}

	if err = rs.AttachRestApi(sc.router, sc.provider); err != nil {
		return err
	}
	www.AttachWWW(sc.router, sc.config.IsDevMode)

	fmt.Printf("Serving at %s:%d ...\n", sc.config.ListenAddress,
		sc.config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		sc.config.ListenAddress, sc.config.ListenPort), sc.router); err != nil {
		return err
	}

	return nil
}
