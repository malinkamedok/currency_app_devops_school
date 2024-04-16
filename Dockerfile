# Step 1: Modules caching
FROM docker.io/golang:1.22-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod tidy

# Step 2: Builder
FROM golang:1.22-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN apk add --no-cache ca-certificates=20240226-r0 && \
    update-ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/main

# Step 3: Final
FROM scratch
EXPOSE 8000
COPY --from=builder /app/internal/config /internal/config
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/app"]
