<div align="center">

# 🏘️ Dashboard Pendataan Penduduk

**Sistem pendataan penduduk & keluarga skala kelurahan/desa**

[![CI](https://img.shields.io/badge/CI-GitHub%20Actions-2088FF?logo=githubactions&logoColor=white)](.github/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)](backend_golang/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vuedotjs&logoColor=white)](frontend_vuetify/)
[![Vuetify](https://img.shields.io/badge/Vuetify-3.6-1867C0?logo=vuetify&logoColor=white)](frontend_vuetify/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)](#database)
[![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)](docker-compose.yml)

Aplikasi web dua sub-proyek: REST API **Go + Gin + GORM** dan SPA **Vue 3 + Vuetify + Pinia**.

</div>

---

## ✨ Fitur

| | Fitur |
|---|---|
| 🔐 | Autentikasi multi-user berbasis **JWT httpOnly cookie** (login, register, logout, remember-me 7 hari) |
| 👥 | Role **`admin`** (akses semua data) & **`user`** (hanya data miliknya — dicek di setiap endpoint) |
| 🏠 | Manajemen **Keluarga**: CRUD + upload foto KK & rumah + titik koordinat GPS (Leaflet) |
| 🧑 | Manajemen **Penduduk**: CRUD lengkap dengan field demografis (NIK, agama, pendidikan, pekerjaan, dll) |
| 📊 | Dashboard statistik: total keluarga/penduduk, distribusi gender, status kawin, tren input harian (ApexCharts) |
| 🖼️ | Upload foto profil — validasi **MIME asli** + nama file acak (anti path-traversal) |
| 📄 | Pagination opsional `?page=&limit=` pada endpoint list |

---

## 🚀 Quickstart

### Opsi A — Docker (paling cepat)

```bash
# edit dulu JWT_SECRET & password di docker-compose.yml
docker compose up --build
```

Frontend: `http://localhost:3000` · API: `http://localhost:8080`

### Opsi B — Manual

**Prasyarat:** Go 1.23+ · Node.js 18+ · MySQL 8.0+ · npm 9+

**1. Database**

```sql
CREATE DATABASE project_pendataan_penduduk;
```

Tabel dibuat otomatis via `AutoMigrate` saat backend pertama jalan.

**2. Backend**

```bash
cd backend_golang
cp .env.example .env      # isi kredensial DB + JWT_SECRET
go mod download
go run main.go            # listen di :8080
```

> ⚠️ `JWT_SECRET` **wajib diisi** (mis. `openssl rand -hex 32`) — backend menolak jalan tanpa ini.

**3. Frontend**

```bash
cd frontend_vuetify
cp .env.example .env      # opsional; default API http://localhost:8080
npm install
npm run dev               # listen di :3000
```

**4. Login**

Seeder membuat akun admin default:

| Email | Password |
|---|---|
| `admin@gmail.com` | `admin123` |

> 🔒 **Segera ganti password admin** setelah login pertama.

---

## 📁 Struktur Repo

```
Dashboard_Pendataan_Penduduk/
├── backend_golang/            # 🟦 REST API (Go)
│   ├── main.go                #    entrypoint + daftar route
│   ├── config/                #    enum demografis + helper Get*()
│   ├── controllers/           #    handler (auth, user, keluarga, penduduk, statistik)
│   │   └── helper.go          #    getCurrentUser / isAdmin / canAccessResource
│   ├── middlewares/           #    CORS + JWT auth
│   ├── models/                #    User, Keluarga, Penduduk (GORM)
│   ├── seeders/               #    admin default
│   ├── setup/                 #    koneksi DB + AutoMigrate
│   ├── utils/                 #    JWT + validasi upload gambar
│   └── Dockerfile
├── frontend_vuetify/          # 🟩 SPA (Vue 3)
│   ├── src/
│   │   ├── pages/             #    auto-route: login, dashboard, penduduk/, keluarga/, profile
│   │   ├── components/        #    navbar, card, dtTable, chart/, modal/
│   │   ├── stores/            #    Pinia: constant, restrict, title
│   │   ├── router/            #    auto-router + global auth guard
│   │   └── plugins/           #    registrasi Vuetify
│   ├── Dockerfile             #    build Vite → nginx
│   └── nginx.conf
├── docker-compose.yml         # 🐳 MySQL + backend + frontend
├── .github/workflows/ci.yml   # ⚙️ CI: vet/build/test Go + build Vite
└── TODO.md                    # ✅ checklist audit & perbaikan
```

---

## 🔌 API

Dokumentasi endpoint lengkap ada di **[backend_golang/README.md](backend_golang/README.md)**. Ringkasan:

| Grup | Endpoint |
|---|---|
| **Publik** | `POST /register` · `POST /login` · `GET /public/*` (file upload) |
| **Auth & User** | `POST /api/logout` · `GET /api/user` · `GET /api/user/all` (admin) · `PUT /api/update/:id` · `PUT /api/update/password/:id` |
| **Keluarga** | `GET/POST/PUT/DELETE /api/{keluarga,addkeluarga,editkeluarga/:id,deletekeluarga/:id}` · `GET /api/latestkel` |
| **Penduduk** | `GET/POST/PUT/DELETE /api/{penduduk,addpenduduk,updatependuduk/:id,deletependuduk/:id}` · `GET /api/latestpend` |
| **Statistik** | `GET /api/{alldata,jumlah,alive,marry,gender}` · `GET /api/data?year=&month=` |

---

## 🔑 Alur Autentikasi

```
login.vue ──POST /login──▶ verifikasi bcrypt ──▶ JWT HS256 (sub, exp)
                                                    │
   axios (withCredentials) ◀── Set-Cookie: Authorization (httpOnly) ──┘
```

- Masa berlaku token **1 hari**, atau **7 hari** dengan `remember_me` — cookie & klaim `exp` sinkron.
- Semua route `/api/*` lewat `AuthMiddleware`; endpoint by-ID/mutasi memverifikasi **kepemilikan resource** (`canAccessResource`) — non-admin tidak bisa membaca/mengubah data user lain.
- `user_id` data baru selalu diambil **dari token**, tidak pernah dari input client.
- Router frontend punya guard global (`beforeEach`); pengaman sesungguhnya tetap di sisi API.

---

## 🗄️ Database

Schema dibangun otomatis via `AutoMigrate` ([setup/setup.go](backend_golang/setup/setup.go)).

| Tabel | Kolom kunci | Relasi |
|---|---|---|
| `users` | `email` (unique), `level`, `profile_picture` | — |
| `keluargas` | `no_kk`, `foto_kk`, `foto_rumah`, `latitude/longtitude` | `user_id → users.id` |
| `penduduks` | `nik`, field demografis (`int8` enum) | `user_id → users.id` · `kels_id → keluargas.id` |

Enum demografis dipetakan ke string via helper di [config/constant.go](backend_golang/config/constant.go); frontend punya mirror di [stores/constant.js](frontend_vuetify/src/stores/constant.js) — **sinkronkan keduanya** bila mengubah enum.

---

## 🧪 Testing & CI

```bash
cd backend_golang && go test ./...     # unit test: config enum, JWT, upload helper
```

CI (GitHub Actions) berjalan otomatis di setiap push/PR ke `main`:
**backend** → `go vet` + `go build` + `go test` · **frontend** → `npm run build`

---

## 🛡️ Catatan Keamanan

Hasil audit keamanan (lihat [TODO.md](TODO.md)) sudah ditindaklanjuti:

- ✅ JWT secret dari environment (`JWT_SECRET`) — tidak lagi hardcode
- ✅ IDOR ditutup: cek kepemilikan di semua endpoint by-ID & mutasi
- ✅ Upload divalidasi berdasarkan MIME asli + nama file acak
- ✅ `user_id` diambil dari token, bukan body request
- ✅ `gin.Recovery()` — panic tidak mematikan server
- ⚠️ Sisa: ganti password admin default & tambahkan rate-limiting sebelum produksi

---

## 📜 Lisensi

Belum ditentukan — tambahkan `LICENSE` sebelum publikasi.
