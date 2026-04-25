# Multi-stage build for the migrate image.
# Stage 1 compiles dbmate from the pinned go.mod tool directive.
# Stage 2 is a minimal Alpine image with a shell (required by dbmate)
# that bundles the binary and the SQL migration files.
#
# Runtime: set DATABASE_URL and run the container once before services start.

ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/dbmate github.com/amacneil/dbmate/v2

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/dbmate /usr/local/bin/dbmate
COPY db/migrations/ /migrations/
ENTRYPOINT ["dbmate", "--migrations-dir", "/migrations", "up"]
