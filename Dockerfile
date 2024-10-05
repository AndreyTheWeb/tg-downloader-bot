FROM golang:1.22.1 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

FROM alpine:latest

COPY --from=builder /app/server /.

EXPOSE 80

CMD ["./server", "--port", "80"]
