e-additives on Google Cloud
=============================

A web platform that provides detailed information about [food additives](http://en.wikipedia.org/wiki/Food_additive).

That two components that make this project are:

  * An HTTP API web service that provides API access to the data.
  * A responsive web app that uses the API to present the data.

# Usage

See the [API documentation](https://github.com/vexelon-dot-net/e-additives.server/blob/master/docs/API.md).

# Development

To build and deploy locally at `http://localhost:8080/api` run:

    PORT=8080 DB_PATH=./data/eadditives.sqlite3 go run main.go

To build and deploy on GCP run:

    gcloud app deploy

# License

  * e-additives server under [MIT License](LICENSE)
  * e-additives data under [CC BY-SA 4.0](data/LICENSE)