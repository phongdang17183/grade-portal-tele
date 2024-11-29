FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Run
FROM alpine:latest

# Thiết lập thư mục làm việc
WORKDIR /root/

# Sao chép binary từ giai đoạn build
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./main"]
