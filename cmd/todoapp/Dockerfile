FROM golang:1.25.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app
COPY --from=builder /app/web/public /app/web/public
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]