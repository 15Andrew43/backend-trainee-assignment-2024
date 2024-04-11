FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./banner_service

EXPOSE 8080

CMD ["sh", "-c", "sleep 5 && ./app"]
