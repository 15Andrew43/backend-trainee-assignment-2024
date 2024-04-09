FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o app ./banner_service

EXPOSE 8080

CMD ["./app"]