FROM golang:1.23 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update \
    && apt-get install -y --no-install-recommends mysql-client \
    && rm -rf /var/lib/apt/lists/*
    
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
