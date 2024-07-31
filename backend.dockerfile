FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

COPY .aws /root/.aws

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]