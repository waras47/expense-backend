# Duitku — Frontend

Aplikasi pencatatan keuangan pribadi berbasis web.  
Dibangun dengan **Vue 3**, **Pinia**, **Vite**, dan **Chart.js**.

> Backend API: lihat [`../backend-go/README.md`](../backend-go/README.md)

---

## Fitur

| Halaman | URL | Keterangan |
|---------|-----|-----------|
| Dashboard | `/` | Ringkasan keuangan, grafik, saldo bulan ini |
| Pengeluaran | `/expenses` | CRUD + filter bulan & kategori |
| Pemasukan | `/incomes` | Catat gaji, freelance, bisnis, dll |
| Hutang | `/debts` | Catat hutang/piutang, tandai lunas |
| Transfer | `/transfers` | Transfer antar rekening / e-wallet |
| Scan Struk | `/receipt` | 📷 Upload foto → AI baca otomatis |
| Kategori | `/categories` | Kelola kategori + warna |

---

## Struktur Folder

```
frontend/
├── src/
│   ├── api/
│   │   └── index.js          ← Semua fungsi HTTP (axios)
│   ├── stores/
│   │   ├── expenseStore.js   ← State pengeluaran & kategori
│   │   ├── incomeStore.js    ← State pemasukan
│   │   ├── debtStore.js      ← State hutang/piutang
│   │   └── transferStore.js  ← State transfer
│   ├── views/
│   │   ├── DashboardView.vue    ← Halaman utama + grafik
│   │   ├── ExpensesView.vue     ← CRUD pengeluaran
│   │   ├── IncomesView.vue      ← CRUD pemasukan
│   │   ├── DebtsView.vue        ← Kelola hutang
│   │   ├── TransfersView.vue    ← Catat transfer
│   │   ├── ReceiptScanView.vue  ← Scan struk AI
│   │   └── CategoriesView.vue   ← Kelola kategori
│   ├── router/index.js       ← Konfigurasi routing
│   ├── App.vue               ← Root component + navbar
│   ├── main.js               ← Entry point
│   └── style.css             ← Design system global
├── index.html
├── vite.config.js            ← Proxy ke backend :8081
├── Dockerfile
└── nginx.conf
```

**Alur data:**  
`View → Store Action → API (axios) → Backend Go → PostgreSQL`

---

## Tech Stack

| Library | Versi | Kegunaan |
|---------|-------|---------|
| [Vue 3](https://vuejs.org) | 3.4 | Framework UI (Composition API) |
| [Pinia](https://pinia.vuejs.org) | 2.1 | State management |
| [Vue Router](https://router.vuejs.org) | 4.3 | Client-side routing |
| [Axios](https://axios-http.com) | 1.6 | HTTP requests ke backend |
| [Chart.js](https://chartjs.org) + vue-chartjs | 4.4 | Grafik bar & donut |
| [Vite](https://vitejs.dev) | 5.0 | Build tool + dev server |

---

## Cara Menjalankan

### Prasyarat
- Node.js 18+
- Backend Go sudah berjalan di port **8081**

### 1. Install dependency

```bash
cd frontend
npm install
```

### 2. Jalankan dev server

```bash
npm run dev
```

Aplikasi buka di: **http://localhost:5173**

> Semua request ke `/api/*` akan di-proxy ke `http://localhost:8081` secara otomatis.

### 3. Build untuk production

```bash
npm run build
# Output di folder dist/
```

---

## Koneksi ke Backend

Konfigurasi proxy ada di [`vite.config.js`](vite.config.js):

```js
proxy: { '/api': { target: 'http://localhost:8081', changeOrigin: true } }
```

Untuk production, ubah `nginx.conf`:

```nginx
location /api/ {
  proxy_pass http://backend:8081;  # sesuaikan dengan host backend
}
```

---

## Halaman & Komponen

### Dashboard (`/`)
- Kartu saldo: **Pemasukan** vs **Pengeluaran** vs **Saldo**
- Kartu hutang: total hutang belum lunas & piutang belum dibayar
- Bar chart: perbandingan pengeluaran vs pemasukan 6 bulan
- Donut chart: pengeluaran per kategori bulan ini
- Shortcut ke semua fitur

### Pengeluaran (`/expenses`)
- Tabel pengeluaran dengan filter **bulan** & **kategori**
- CRUD lengkap (buat, edit, hapus)
- Ringkasan bulanan & per kategori

### Pemasukan (`/incomes`)
- Tabel pemasukan dengan filter bulan & kategori
- Kategori: `salary`, `freelance`, `business`, `investment`, `gift`, `other`
- CRUD + ringkasan 6 bulan

### Hutang (`/debts`)
- Tab: **Semua / Hutang saya / Piutang saya**
- Toggle: tampilkan belum lunas saja
- Tombol **Tandai Lunas** — langsung update status
- Indikator jatuh tempo (merah jika sudah lewat)

### Transfer (`/transfers`)
- Catat perpindahan saldo: BCA → GoPay, Mandiri → OVO, dll
- Filter per bulan & akun asal
- Tampilan arah transfer yang jelas

### Scan Struk (`/receipt`)
1. Upload foto struk (drag & drop atau klik)
2. AI (Claude Vision) membaca struk → hasilkan data terstruktur
3. Review hasil: nama toko, tanggal, list item, total
4. Simpan langsung sebagai pengeluaran — sudah terisi otomatis

> **Prasyarat scan:** `ANTHROPIC_API_KEY` harus diset di backend `.env`

### Kategori (`/categories`)
- Grid kartu kategori dengan warna
- Color picker + 10 preset warna
- Hapus kategori

---

## Design System

Warna utama (CSS variables di `style.css`):

| Variable | Warna | Digunakan untuk |
|----------|-------|----------------|
| `--red` | `#E8442A` | Tombol utama, pengeluaran |
| `--teal` | `#2ABFA3` | Pemasukan, lunas |
| `--yellow` | `#F5C842` | Aksen |
| `--blue` | `#3D8EE8` | Info |
| `--cream` | `#F5F0E8` | Background |

Font: **Syne** (heading) + **DM Sans** (body)

---

## Docker

```bash
# Build image
docker build -t duitku-frontend .

# Jalankan (dengan backend jalan di host)
docker run -p 3000:80 duitku-frontend
```

Atau pakai `docker-compose` bersama backend:

```yaml
version: '3.8'
services:
  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
  backend:
    build: ./backend-go
    ports:
      - "8081:8081"
    environment:
      DATABASE_URL: postgres://postgres:secret@db:5432/expense_tracer
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: expense_tracer
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
```

---

## Link Terkait

- Backend Go API: [`../backend-go/`](../backend-go/)
- Backend Rust (original): [`../backend/`](../backend/)
- Swagger UI (saat backend jalan): http://localhost:8081/swagger/index.html
- GitHub Project: https://github.com/users/waras47/projects/5
