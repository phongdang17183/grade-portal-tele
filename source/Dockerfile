FROM golang:1.23

WORKDIR /app

# Copy mã nguồn và build ứng dụng
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main main.go

# Thiết lập Entrypoint Script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 8080


# Sử dụng script khởi động trước khi chạy ứng dụng
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
