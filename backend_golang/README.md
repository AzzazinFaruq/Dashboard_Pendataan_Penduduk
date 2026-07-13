# Backend — Pendataan Penduduk

REST API untuk sistem pendataan penduduk & keluarga.

**Stack:** Go 1.23 · Gin · GORM · MySQL · JWT (cookie)

---

## Menjalankan

```bash
cp .env.example .env     # isi kredensial DB & JWT_SECRET
go mod download
go run main.go           # listen di :8080
```

Tabel dibuat otomatis via `AutoMigrate` saat startup, dan seeder membuat akun admin default (`admin@gmail.com` / `admin123` — **segera ganti**).

## Environment variable

| Variabel | Wajib | Keterangan |
|---|---|---|
| `DB_USERNAME` | ✅ | User MySQL |
| `DB_PASSWORD` | – | Password MySQL |
| `DB_HOST` | ✅ | Host DB (mis. `127.0.0.1`) |
| `DB_PORT` | ✅ | Port DB (mis. `3306`) |
| `DB_NAME` | ✅ | Nama database |
| `JWT_SECRET` | ✅ | Secret penandatangan JWT. Backend **menolak jalan** jika kosong. Generate: `openssl rand -hex 32` |
| `FE_URL` | – | Origin frontend untuk CORS (default `http://localhost:3000`) |

## Struktur

```
backend_golang/
├── main.go              # entrypoint + daftar route
├── config/              # konstanta enum (agama, pekerjaan, dll)
├── controllers/         # handler HTTP
│   └── helper.go        # getCurrentUser, isAdmin, canAccessResource
├── middlewares/         # CORS & JWT auth
├── models/              # skema GORM (User, Keluarga, Penduduk)
├── seeders/             # seed admin awal
├── setup/               # koneksi DB & migrasi
└── utils/               # JWT & validasi upload gambar
```

## Autentikasi & otorisasi

- Login mengembalikan JWT dalam cookie `Authorization` (httpOnly). Masa berlaku 1 hari (7 hari bila `remember_me`).
- Semua route di grup `/api` diproteksi `AuthMiddleware`.
- **Role:** `admin` mengakses semua data; `user` hanya data miliknya sendiri.
- Endpoint by-ID & mutasi memverifikasi kepemilikan (`canAccessResource`) — non-admin tidak bisa mengakses data user lain.

## Endpoint

### Publik
| Method | Path | Fungsi |
|---|---|---|
| POST | `/register` | Daftar user baru (role `user`) |
| POST | `/login` | Login, set cookie JWT |

### Terproteksi (`/api`, butuh cookie)
| Method | Path | Fungsi |
|---|---|---|
| POST | `/logout` | Hapus cookie |
| GET | `/user` | User yang sedang login |
| GET | `/user/all` | Semua user (**admin**) |
| PUT | `/update/:id` | Update profil user |
| PUT | `/update/password/:id` | Ganti password |
| GET | `/keluarga` · `/keluarga/:id` | List / detail keluarga |
| GET | `/latestkel` · `/latestkelinput` | Keluarga terbaru |
| POST | `/addkeluarga` | Tambah keluarga (multipart, upload foto) |
| PUT | `/editkeluarga/:id` | Update keluarga |
| DELETE | `/deletekeluarga/:id` | Hapus keluarga (+ file foto) |
| GET | `/penduduk` · `/penduduk/:id` | List / detail penduduk |
| GET | `/latestpend` | Penduduk terbaru |
| POST | `/addpenduduk` | Tambah penduduk (JSON) |
| PUT | `/updatependuduk/:id` | Update penduduk |
| DELETE | `/deletependuduk/:id` | Hapus penduduk |
| GET | `/alldata` · `/jumlah` · `/alive` | Statistik jumlah & status |
| GET | `/marry` · `/gender` | Statistik kawin & gender |
| GET | `/data?year=&month=` | Tren input per hari dalam bulan |

> `GET /keluarga` & `/penduduk` mendukung pagination opsional via `?page=&limit=` (tanpa parameter = semua data).

## Testing

```bash
go test ./...
```

## Build

```bash
go build -o app        # binary produksi
```
