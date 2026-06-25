# Build
FROM golang:1.26 AS build

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download
COPY . .

# "static linux amd64" binary
# Flags for linker - https://pkg.go.dev/cmd/link
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o /build/goblncr \
    /build/cmd/goblncr.go

# Run
FROM gcr.io/distroless/static-debian13
WORKDIR /app

COPY --from=build /build/goblncr /app/goblncr
ENTRYPOINT ["/app/goblncr"]

