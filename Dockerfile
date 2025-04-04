FROM golang:1.24.2-alpine3.20@sha256:00f149d5963f415a8a91943531b9092fde06b596b276281039604292d8b2b9c8 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o main /app/cmd/main.go

FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "/app/main" ]
