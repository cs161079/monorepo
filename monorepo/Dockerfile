FROM golang:1.20

# Set Go environment variables
# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=direct
    
WORKDIR /app

RUN mkdir common api

COPY common/ common/
COPY api/ api/

COPY go.sum go.sum
COPY go.mod go.mod
COPY main_api.go main.go
COPY .env .

RUN go mod download

RUN go build -o bin/executable_go

EXPOSE 8081

ENV SERVER_PORT=8081

CMD ["./bin/executable_go"]