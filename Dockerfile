FROM golang:1.13

WORKDIR /go/src/github.com/crunchydata/postgres-operator
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a postgres-operator.go

FROM alpine:3.7
COPY --from=0 /go/src/github.com/crunchydata/postgres-operator/postgres-operator /crunchydata/postgres-operator
ENTRYPOINT ["/crunchydata/postgres-operator"]
