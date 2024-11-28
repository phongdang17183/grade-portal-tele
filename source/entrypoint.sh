#!/bin/bash

# Tạo file .env từ các biến môi trường
cat <<EOF > /app/.env
API_URL=${API_URL}
DBURL=${DBURL}
BOT_TOKEN=${BOT_TOKEN}
EOF

echo "File .env đã được tạo với nội dung:"
cat /app/.env

# Chạy ứng dụng
exec "$@"
