# Multi-stage build for the query-server binary. Mirrors
# command-server.Dockerfile; differs only in the build target.

ARG GO_VERSION=1.26.2

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o /out/query-server ./cmd/query-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/query-server /usr/local/bin/query-server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/query-server"]
