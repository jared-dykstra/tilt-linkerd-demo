FROM golang:1.22-alpine

COPY . /app
WORKDIR /app

RUN go build -o foo main.go

CMD ["./foo"]