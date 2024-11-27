
FROM golang:1.23

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

CMD ["./main"]
