FROM golang:1.22-alpine

COPY . /app
WORKDIR /app

RUN go build -o bar main.go

CMD ["./bar"]
