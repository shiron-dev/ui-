FROM golang:1.23.6-alpine3.20@sha256:22caeb4deced0138cb4ae154db260b22d1b2ef893dde7f84415b619beae90901 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o main /app/cmd/main.go

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "/app/main" ]
