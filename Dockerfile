FROM golang:1.23.6-alpine3.20@sha256:484906fa392c79201ee01325cb42c0414aa9ad836afacf49a0e98a69fd252f78 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o main /app/cmd/main.go

FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "/app/main" ]
