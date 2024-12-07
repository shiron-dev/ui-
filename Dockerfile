FROM golang:1.23.4-alpine3.20 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o main /app/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "/app/main" ]
