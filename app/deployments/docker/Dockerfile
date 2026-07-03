# Build
FROM golang:1.26 AS build

WORKDIR /build

COPY go.mod go.sum /build/
RUN go mod download
COPY . .

# "no-c-binding static linux amd64" binary
# Flags for linker - https://pkg.go.dev/cmd/link
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o /build/goblncr \
    /build/cmd/goblncr/goblncr.go

# Run
FROM debian:13-slim
WORKDIR /app

COPY --from=build /build/goblncr /app/goblncr

EXPOSE 8080

ENTRYPOINT ["/app/goblncr"]
