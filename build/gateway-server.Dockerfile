# Multi-stage build for the gateway-server binary. Mirrors
# command-server.Dockerfile and query-server.Dockerfile; differs
# only in the build target.

ARG GO_VERSION=1.26.3

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o /out/gateway-server ./cmd/gateway-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/gateway-server /usr/local/bin/gateway-server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/gateway-server"]
