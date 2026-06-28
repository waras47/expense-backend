# API Backend Expense

A lightweight backend RESTful API designed for tracking and managing personal or organizational expenses. This API provides secure endpoints to log daily expenses, categorize spending, and generate basic financial summaries.

## 🚀 Spesifikasi API (Swagger Info)
* **Version:** 1.0
* **Base Path:** `/api`
* **Title:** API Backend Expense

---

## 🛠️ Cara Menjalankan Aplikasi

Ikuti langkah-langkah di bawah ini untuk menjalankan database dan aplikasi di lingkungan lokal:

### 1. Jalankan Infrastruktur (Database/Docker)
Pastikan Docker desktop Anda sudah aktif, lalu jalankan semua *services* pendukung di *background*:
```bash
docker compose up -d
```

### 2. Generate Dokumentasi Swagger
Setiap kali ada perubahan pada komentar anotasi API, perbarui file dokumentasi dengan perintah:
```bash
swag init -d .,./internal/handler -g cmd/main.go
```

### 3. Jalankan Aplikasi Golang
Setelah kontainer Docker berjalan dan Swagger diperbarui, jalankan aplikasi Go Anda:
```bash
go run ./cmd
```

---

## 📖 Mengakses Dokumentasi API
Setelah aplikasi berjalan, Anda dapat melihat dan menguji langsung semua endpoint melalui Swagger UI di browser:
* **URL:** `http://localhost:8080/swagger/index.html` *(Sesuaikan port jika berbeda)*
