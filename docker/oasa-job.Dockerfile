FROM golang:1.20

# Set Go environment variables
# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=direct
    
WORKDIR /app


RUN mkdir common

COPY ../common/ common/

COPY ../go.sum .
COPY ../go.mod .
COPY ../cronjob/syncService.go .
COPY ../cronjob/uVersionsRepository.go .
COPY ../cronjob/main.go .
COPY .env .

RUN go mod download

RUN go build -o bin/oasa-job

ENV SERVER_PORT=8081

CMD ["./bin/oasa-job"]


