package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Struct cho thông tin đăng nhập
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Handler cho API đăng nhập
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Kiểm tra tài khoản, mật khẩu (chỉ là ví dụ đơn giản)
    if req.Username == "admin" && req.Password == "password" {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Login successful!")
    } else {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    }
}
