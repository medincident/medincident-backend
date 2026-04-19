# Multi-stage build for the query-server binary. Mirrors
# command-server.Dockerfile; differs only in the build target and the
# default exposed gRPC port.

ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/query-server ./cmd/query-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/query-server /usr/local/bin/query-server
USER nonroot:nonroot
EXPOSE 9091
ENTRYPOINT ["/usr/local/bin/query-server"]
