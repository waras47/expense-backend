# Expense Tracer — Backend Go

Backend REST API untuk aplikasi pencatatan keuangan pribadi.  

Tech stack:
**Go**, **Gin**, **PostgreSQL**, dan **Clean Architecture**.



## Fitur

| Fitur | Endpoint Prefix | Keterangan |
|-------|----------------|------------|
| Kategori | `/api/categories` | Manajemen kategori pengeluaran |
| Pengeluaran | `/api/expenses` | CRUD + ringkasan bulanan & per kategori |
| Pemasukan | `/api/incomes` | Catat gaji, freelance, bisnis, dll |
| Hutang | `/api/debts` | Catat hutang kita ke orang lain & sebaliknya |
| Transfer | `/api/transfers` | Transfer antar rekening / e-wallet |
| Scan Struk | `/api/receipts/scan` | Upload foto struk → AI baca otomatis |

---

## Arsitektur
```
backend-go/
├── cmd/
│   └── main.go              ← Entry point + Dependency Injection
├── internal/
│   ├── domain/              ← Layer 1: Entity & Interface (aturan bisnis)
│   │   ├── category.go
│   │   ├── expense.go
│   │   ├── income.go
│   │   ├── debt.go
│   │   ├── transfer.go
│   │   └── receipt.go
│   ├── repository/          ← Layer 2: Implementasi database (PostgreSQL)
│   │   ├── category_repository.go
│   │   ├── expense_repository.go
│   │   ├── income_repository.go
│   │   ├── debt_repository.go
│   │   └── transfer_repository.go
│   ├── usecase/             ← Layer 3: Logika bisnis
│   │   ├── category_usecase.go
│   │   ├── expense_usecase.go
│   │   ├── income_usecase.go
│   │   ├── debt_usecase.go
│   │   ├── transfer_usecase.go
│   │   └── receipt_usecase.go
│   ├── handler/             ← Layer 4: HTTP Handler (Gin)
│   │   ├── category_handler.go
│   │   ├── expense_handler.go
│   │   ├── income_handler.go
│   │   ├── debt_handler.go
│   │   ├── transfer_handler.go
│   │   └── receipt_handler.go
│   └── service/
│       └── claude_ocr.go    ← OCR service (Claude AI Vision)
├── pkg/
│   └── apperror/errors.go   ← Error handling terpusat
├── docs/                    ← Auto-generated Swagger docs
├── migrations/              ← SQL migration files
├── uploads/                 ← Folder penyimpanan foto struk
└── .env
```

**Alur request:**  
`HTTP Request → Handler → Usecase → Repository → PostgreSQL`

---

## Tech Stack

| Komponen | Library/Tools |
|----------|--------------|
| Framework | [Gin](https://github.com/gin-gonic/gin) v1.10 |
| Database | PostgreSQL + [pgx/v5](https://github.com/jackc/pgx) |
| Desimal | [shopspring/decimal](https://github.com/shopspring/decimal) |
| CORS | [gin-contrib/cors](https://github.com/gin-contrib/cors) |
| Env config | [godotenv](https://github.com/joho/godotenv) |
| OCR / AI | [Claude claude-haiku-4-5 Vision](https://www.anthropic.com) (Anthropic) |
| Dokumentasi | [Swagger (swaggo)](https://github.com/swaggo/swag) |

---

## Cara Menjalankan

### 1. Prasyarat

- Go 1.22+
- PostgreSQL berjalan

### 2. Clone & masuk folder

```bash
cd backend-go
```

### 3. Konfigurasi `.env`

Salin dan edit file `.env`:

```env
DATABASE_URL=postgres://postgres:PASSWORD@localhost:5432/expense_tracer
SERVER_HOST=0.0.0.0
SERVER_PORT=8081
UPLOAD_DIR=uploads

# Untuk fitur scan struk (daftar gratis di https://console.anthropic.com)
ANTHROPIC_API_KEY=sk-ant-xxxxx
```

### 4. Migrasi database

```bash
# Tabel awal (categories, expenses) — jika belum ada
psql -U postgres -d expense_tracer -f migrations/001_initial.sql

# Tabel baru (incomes, debts, transfers)
psql -U postgres -d expense_tracer -f migrations/002_income_debt_transfer.sql
```

### 5. Install dependency & jalankan

```bash
go mod tidy
go run ./cmd/main.go
```

Server berjalan di: `http://localhost:8081`  
Swagger UI: `http://localhost:8081/swagger/index.html`

### 6. Update Swagger docs (setelah edit anotasi handler)

```bash
swag init -g cmd/main.go --output docs
```

---

## Dokumentasi API (Swagger)

Setelah server jalan, buka:

```
http://localhost:8081/swagger/index.html
```

## API Reference

### Kategori

| Method | URL | Keterangan |
|--------|-----|-----------|
| `GET` | `/api/categories` | Daftar semua kategori |
| `POST` | `/api/categories` | Buat kategori baru |
| `DELETE` | `/api/categories/:id` | Hapus kategori |

**POST body:**
```json
{
  "name": "Makanan",
  "color": "#f59e0b"
}
```

---

### Pengeluaran

| Method | URL | Keterangan |
|--------|-----|-----------|
| `GET` | `/api/expenses` | Daftar (filter: `?month=2025-05&category_id=1`) |
| `GET` | `/api/expenses/:id` | Detail |
| `POST` | `/api/expenses` | Buat baru |
| `PUT` | `/api/expenses/:id` | Update |
| `DELETE` | `/api/expenses/:id` | Hapus |
| `GET` | `/api/expenses/summary/monthly` | Ringkasan 6 bulan |
| `GET` | `/api/expenses/summary/category` | Ringkasan per kategori (bulan ini) |

**POST body:**
```json
{
  "title": "Makan siang",
  "amount": 25000,
  "category_id": 1,
  "expense_date": "2025-05-30",
  "note": "Warteg dekat kantor"
}
```

---

### Pemasukan

| Method | URL | Keterangan |
|--------|-----|-----------|
| `GET` | `/api/incomes` | Daftar (filter: `?month=2025-05&category=salary`) |
| `GET` | `/api/incomes/:id` | Detail |
| `POST` | `/api/incomes` | Catat pemasukan |
| `PUT` | `/api/incomes/:id` | Update |
| `DELETE` | `/api/incomes/:id` | Hapus |
| `GET` | `/api/incomes/summary/monthly` | Ringkasan 6 bulan |

**Kategori yang tersedia:** `salary`, `freelance`, `business`, `investment`, `gift`, `other`

**POST body:**
```json
{
  "title": "Gaji Mei 2025",
  "amount": 5000000,
  "category": "salary",
  "income_date": "2025-05-01",
  "note": "Transfer ke BCA"
}
```

---

### Hutang

| Method | URL | Keterangan |
|--------|-----|-----------|
| `GET` | `/api/debts` | Daftar (filter: `?type=owe&is_paid=false`) |
| `GET` | `/api/debts/:id` | Detail |
| `POST` | `/api/debts` | Catat hutang baru |
| `PUT` | `/api/debts/:id` | Update |
| `PATCH` | `/api/debts/:id/paid` | Tandai sudah lunas |
| `DELETE` | `/api/debts/:id` | Hapus |

**Tipe hutang:**
- `owe` → **kita** yang berhutang ke orang lain
- `lent` → **orang lain** yang berhutang ke kita

**POST body:**
```json
{
  "person_name": "Budi",
  "amount": 150000,
  "type": "lent",
  "due_date": "2025-06-30",
  "note": "Pinjam uang buat bensin"
}
```

---

### Transfer

| Method | URL | Keterangan |
|--------|-----|-----------|
| `GET` | `/api/transfers` | Daftar (filter: `?month=2025-05&from_account=BCA`) |
| `GET` | `/api/transfers/:id` | Detail |
| `POST` | `/api/transfers` | Catat transfer |
| `PUT` | `/api/transfers/:id` | Update |
| `DELETE` | `/api/transfers/:id` | Hapus |

**POST body:**
```json
{
  "title": "Top up GoPay",
  "amount": 200000,
  "from_account": "BCA",
  "to_account": "GoPay",
  "transfer_date": "2025-05-30",
  "note": "Untuk transportasi minggu ini"
}
```

---

### Scan Struk (AI)

| Method | URL | Keterangan |
|--------|-----|-----------|
| `POST` | `/api/receipts/scan` | Upload foto struk |

**Request:** `multipart/form-data`, field `file` (jpg/png/webp, maks 10MB)

```bash
# Contoh dengan curl:
curl -X POST http://localhost:8081/api/receipts/scan \
  -F "file=@struk_indomaret.jpg"
```

**Response:**
```json
{
  "store_name": "Indomaret Jl. Sudirman",
  "date": "2025-05-30",
  "items": [
    { "name": "Indomie Goreng", "qty": 3, "price": "3500", "subtotal": "10500" },
    { "name": "Aqua 600ml",    "qty": 2, "price": "4000", "subtotal": "8000"  }
  ],
  "subtotal": "18500",
  "tax": "0",
  "total": "18500",
  "payment_type": "qris",
  "image_path": "/uploads/receipt_1748597234.jpg"
}
```

> Hasil ini bisa langsung dipakai untuk mengisi form `POST /api/expenses`.

---

## Skema Database

```sql
-- Tabel yang sudah ada
categories  (id, name, color)
expenses    (id, title, amount, category_id, note, expense_date, created_at)

-- Tabel baru
incomes     (id, title, amount, category, note, income_date, created_at)
debts       (id, person_name, amount, type, due_date, is_paid, note, created_at, paid_at)
transfers   (id, title, amount, from_account, to_account, note, transfer_date, created_at)
```

---

## Docker

```bash
# Build image
docker build -t expense-tracer-go .

# Jalankan (pastikan PostgreSQL sudah jalan)
docker run -p 8081:8081 \
  -e DATABASE_URL=postgres://postgres:pass@host.docker.internal:5432/expense_tracer \
  -e ANTHROPIC_API_KEY=sk-ant-xxxxx \
  expense-tracer-go
```

---

## Referensi

- Rust backend (referensi awal): `../backend/`
- Swagger UI: http://localhost:8081/swagger/index.html
- Anthropic Console (API Key): https://console.anthropic.com
