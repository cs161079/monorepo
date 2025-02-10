FROM golang:1.22

# Set Go environment variables
# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct

# Install necessary dependencies for Alpine (e.g., git for Go modules)
# RUN apk add --no-cache git
# RUN apk add --no-cache curl

# RUN curl -v https://filippo.io/edwards25519?go-get=1

WORKDIR /app

RUN mkdir common cronjob cronjob/config cronjob/dao

COPY common/ common/
COPY cronjob/config cronjob/config
COPY cronjob/dao cronjob/dao

# Δοκιμή να μην αντιγράψω το sum. Μπορεί να παραχθεί από την εντολή go mod tidy
COPY go.sum go.sum
COPY go.mod go.mod
COPY cronjob/main.go .
COPY cronjob/.env .

# RUN go mod tidy
RUN go mod download

RUN go build -o bin/executable_go

CMD ["./bin/executable_go"]


