FROM golang:1.26.7-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o /app/ticket-er /app/cmd/ticketer/main.go 

FROM scratch
COPY --from=builder /app/ticket-er /ticket-er
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/ticket-er"]