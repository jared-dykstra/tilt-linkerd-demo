FROM golang:1.22-alpine

COPY . /app
WORKDIR /app

RUN go build -o baz main.go

CMD ["./baz"]