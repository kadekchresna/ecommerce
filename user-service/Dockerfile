FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user-service .


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user-service .
RUN chmod +x user-service

EXPOSE 8081

CMD [ "./user-service web" ]

