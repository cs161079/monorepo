FROM golang:1.23 AS builder
#Install Git (To fetch Go Modules)
RUN apt-get update && apt-get install -y git

# Set Go environment variables
# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=direct
    
WORKDIR /app

RUN mkdir common \
    webApplication \
    webApplication/config \
    webApplication/controllers \
    webApplication/db \
    db \
    db/migrations \
    var 

COPY common/ common/
# COPY api/ api/
COPY webApplication/config webApplication/config
COPY webApplication/controllers webApplication/controllers
COPY webApplication/keycloak webApplication/keycloak
COPY webApplication/db/migrations/release/ db/migrations/
COPY webApplication/various various/

#COPY go.sum go.sum
COPY go.mod go.mod
COPY auth auth
COPY webApplication/main.go .
COPY webApplication/.env .

RUN go mod download
RUN go mod tidy

RUN go build -o bin/executable_go

EXPOSE 8081

ENV SERVER_PORT=8081

CMD ["./bin/executable_go"]
