# Multi-stage build for the gateway-server binary. Mirrors
# command-server.Dockerfile and query-server.Dockerfile; differs
# only in the build target.

ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/gateway-server ./cmd/gateway-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/gateway-server /usr/local/bin/gateway-server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/gateway-server"]
