FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-service .


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/product-service .
RUN chmod +x product-service

EXPOSE 8081

CMD [ "./product-service web" ]

