# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pre-compress static CSS/JS so the embedded gzip handler can serve them.
RUN go generate ./...

ARG VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X 'jaronjones/ace-of-base/internal/version.Version=${VERSION}'" \
    -o /out/ace-of-base .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/ace-of-base /app/ace-of-base

EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/ace-of-base"]
