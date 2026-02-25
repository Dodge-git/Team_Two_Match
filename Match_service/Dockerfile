FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/main

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main /app/main

EXPOSE 8000

ENV GIN_MODE=release

ENTRYPOINT ["./match_service"]