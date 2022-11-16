E-additives on Google Cloud
=============================

A web platform that provides detailed information about [food additives](http://en.wikipedia.org/wiki/Food_additive).

This project has the following two components:

  * An HTTP API web service that provides API access to the data.
  * A responsive web app that consumes the API and presents the data.

# Usage

API docs are WIP, but [the old docs](https://github.com/vexelon-dot-net/e-additives.server/blob/master/docs/API.md) are very much compatible.

# Development

To init the `www` web app submodule run:

    git submodule init
    git submodule update

To build and deploy locally at `http://localhost:8080/api` run:

    DEVMODE=1 PORT=8080 DB_PATH=./data/eadditives.sqlite3 go run main.go

To run with full-text search support use the following tags:

    go run -tags "eadfts,fts5"  main.go 

To build and deploy on GCP run:

    gcloud app deploy

# License

  * e-additives server under [MIT License](LICENSE)
  * e-additives data under [CC BY-SA 4.0](data/LICENSE)