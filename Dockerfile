# First stage: build the application
FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# use mount type=cache packages between rebuilds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    CGO_ENABLED=0 GOOS=linux go build -o rentals-api main.go

# Second stage: create the runtime container
FROM alpine:3

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/rentals-api /app/rentals-api

ENTRYPOINT ["/app/rentals-api"]
