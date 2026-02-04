FROM golang:1.25.6-alpine AS builder

COPY . /app/
WORKDIR /app/

RUN go mod download
RUN go build -o ./bin/server cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bin/server .
CMD ["./server"]