# Stage 1: Build
FROM golang:1.23 
# AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.* ./
RUN go mod download

# Copy source code and build the binary
COPY . . 
RUN go build -o main main.go

# Stage 2: Deploy
# FROM debian:rc-buggy-20241111

# WORKDIR /app

# # Install locales and configure UTF-8
# RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
# 	&& localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
# ENV LANG en_US.utf8


# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

EXPOSE 8080

# Start the application
CMD ["./main"]