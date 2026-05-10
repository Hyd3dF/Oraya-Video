# syntax=docker/dockerfile:1.7
# Multi-stage build that produces two final images sharing the Go build cache:
#   - oroya-server   : HTTP API
#   - oroya-worker   : background video processor (FFmpeg installed)
# Build a single image by default; pick the binary with the BIN build arg.

ARG BIN=server

FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ARG BIN
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
        -o /out/oroya ./cmd/${BIN}

FROM alpine:3.20
RUN apk add --no-cache ca-certificates ffmpeg tini && \
    addgroup -S oroya && adduser -S -G oroya oroya
WORKDIR /app
COPY --from=builder /out/oroya /app/oroya
COPY internal/migrations /app/migrations
USER oroya
EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--", "/app/oroya"]
