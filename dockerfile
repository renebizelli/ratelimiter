FROM golang:latest AS builder

WORKDIR /app

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/main.go

FROM scratch

COPY --from=builder /app/server .
COPY --from=builder /app/.env .

EXPOSE 8080 6379

CMD ["./server"]
