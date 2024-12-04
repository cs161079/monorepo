FROM golang:1.20

# Set Go environment variables
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct \
    GOCACHE=/go/cache \
    GOMODCACHE=/go/pkg/mod
    
WORKDIR /app

RUN mkdir common api

COPY common/* common/.
COPY api/* api/.

COPY go.mod .
COPY main_api.go main.go
COPY .env .

RUN go mod tidy

RUN go mod download

RUN go build -o bin/executable_go

EXPOSE 8081

ENV SERVER_PORT=8081

CMD ["./bin/executable_go"]
