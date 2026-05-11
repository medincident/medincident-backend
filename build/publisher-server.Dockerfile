# Multi-stage build for the publisher-server binary.
# The first stage compiles a static Linux binary with the same Go
# toolchain version as go.mod. The second stage is a distroless base
# that runs as a non-root user with no shell.

ARG GO_VERSION=1.26.2

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o /out/publisher-server ./cmd/publisher-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/publisher-server /usr/local/bin/publisher-server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/publisher-server"]
