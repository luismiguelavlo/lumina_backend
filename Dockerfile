FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/library_back .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/library_back /app/library_back
USER app

ENV PORT=8080
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=20s --retries=5 \
    CMD sh -c 'wget -qO- http://127.0.0.1:${PORT}/health || exit 1'

ENTRYPOINT ["/app/library_back"]
