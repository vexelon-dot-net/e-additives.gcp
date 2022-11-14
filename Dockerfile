# see https://github.com/GoogleCloudPlatform/cloud-build-samples/tree/main/golang-sample
FROM golang:1.16 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -tags "fts5" -o eadditives

FROM alpine

COPY --from=builder /app/eadditives /app/eadditives

CMD ["/app/eadditives"]