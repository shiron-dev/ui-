FROM golang:1.23.4-alpine3.20@sha256:2314d93ea3899b8118d11ec70714e96c4bb2d52ff46ee919622269a5c6207077 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o main /app/cmd/main.go

FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "/app/main" ]
