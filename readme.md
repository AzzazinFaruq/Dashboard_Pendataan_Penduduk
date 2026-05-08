# Project Pendataan Penduduk

Sistem pendataan penduduk & keluarga skala kelurahan/desa. Web app dua sub-proyek:

- [`backend_golang/`](backend_golang/) — REST API (Go + Gin + GORM + MySQL)
- [`frontend_vuetify/`](frontend_vuetify/) — SPA (Vue 3 + Vuetify + Pinia + Vite)

> Folder `BACKEND LARAVEL/` ada di repo tapi **tidak dipakai** — abaikan.

> ⚠ Sebelum deploy ke production, **wajib baca [AUDIT.md](AUDIT.md)** — ada beberapa kelemahan keamanan kritis (IDOR, mass assignment, JWT secret hardcoded) yang harus diperbaiki dulu.

---

## Fitur

- Autentikasi multi-user (login, register, logout) berbasis JWT cookie
- Role: `admin` (akses semua data) & `user` (hanya data milik sendiri)
- Manajemen **Keluarga**: CRUD + upload foto KK & foto rumah + koordinat GPS
- Manajemen **Penduduk**: CRUD lengkap dengan field demografis (NIK, agama, pendidikan, pekerjaan, dll)
- Dashboard statistik: total keluarga/penduduk, distribusi gender, status kawin, status aktif, tren input per bulan
- Upload & ganti foto profil user

---

## Prasyarat

| Tool | Versi minimum |
|---|---|
| Go | 1.23+ |
| Node.js | 18+ |
| MySQL | 8.0+ |
| npm | 9+ |

---

## Quickstart

### 1. Database

Buat database kosong di MySQL:

```sql
CREATE DATABASE project_pendataan_penduduk;
```

Tabel akan dibuat otomatis lewat `AutoMigrate` saat backend pertama kali jalan.

### 2. Backend (Go)

```bash
cd backend_golang
cp .env.example .env   # buat dulu kalau belum ada
# edit .env, isi kredensial DB
go mod download
go run main.go         # listen di :8080
```

Isi `.env` minimal:

```env
FE_URL=http://localhost:3000
DB_USERNAME=root
DB_PASSWORD=
DB_HOST=localhost
DB_PORT=3306
DB_NAME=project_pendataan_penduduk
```

Saat startup pertama, seeder otomatis bikin akun admin:

| Email | Password |
|---|---|
| `admin@gmail.com` | `admin123` |

> **Ganti password admin segera** setelah login pertama.

### 3. Frontend (Vuetify)

```bash
cd frontend_vuetify
npm install
npm run dev            # listen di :3000
```

Buka `http://localhost:3000` dan login dengan kredensial admin di atas.

> Base URL backend di-hardcode di [src/main.js:42](frontend_vuetify/src/main.js#L42). Kalau backend pindah host, edit di sana.

---

## Script

### Backend

```bash
go run main.go         # dev
go build -o app         # build binary
```

### Frontend

```bash
npm run dev      # vite dev server (port 3000)
npm run build    # production bundle ke dist/
npm run preview  # preview hasil build
npm run lint     # eslint --fix
```

---

## Struktur repo

```
Project-one-fixed/
├── backend_golang/         # Go API
│   ├── main.go             # entrypoint + register route
│   ├── setup/              # GORM connect + AutoMigrate + seed
│   ├── seeders/            # admin default
│   ├── config/             # enum + helper Get*()
│   ├── middlewares/        # CORS, JWT auth
│   ├── models/             # User, Keluarga, Penduduk
│   ├── controllers/        # auth, user, keluarga, penduduk, universal (statistik)
│   ├── utils/              # JWT generate/validate
│   └── public/uploads/     # foto-kk, foto-rumah, profile_pictures
├── frontend_vuetify/       # Vue 3 SPA
│   ├── src/
│   │   ├── pages/          # auto-route (login, dashboard, penduduk, keluarga, dll)
│   │   ├── components/     # navbar, card, dtTable, chart/, modal/
│   │   ├── stores/         # Pinia (constant, nav, restrict, title)
│   │   ├── plugins/        # vuetify config + register
│   │   ├── styles/         # SCSS settings + global CSS
│   │   └── main.js         # bootstrap (axios, pinia, sweetalert2, apexcharts)
│   ├── vite.config.mjs
│   └── package.json
├── AUDIT.md                # daftar kelemahan & TODO perbaikan
└── readme.md
```

---

## Tech stack

### Backend

| Komponen | Versi |
|---|---|
| Go | 1.23.2 |
| Gin | v1.10.0 |
| GORM | v1.25.12 (driver MySQL v1.5.7) |
| jwt/v5 | v5.2.1 (HS256) |
| godotenv | v1.5.1 |
| bcrypt | `golang.org/x/crypto` |

### Frontend

| Komponen | Versi |
|---|---|
| Vue | 3.4.31 |
| Vuetify | 3.6.11 |
| Vue Router | 4.4.0 (auto-routing via `unplugin-vue-router`) |
| Pinia | 2.1.7 |
| Vite | 5.1.5 |
| Axios | 1.7.7 + `axios-retry` 4.5.0 |
| ApexCharts, Chart.js | charting |
| Leaflet, vue-leaflet | peta lat/long keluarga |
| vue-sweetalert2 | toast/alert |

---

## API Endpoint

### Public

| Method | Path | Deskripsi |
|---|---|---|
| POST | `/register` | Registrasi user baru |
| POST | `/login` | Login → set cookie `Authorization` |
| GET | `/public/*` | Static file (foto upload) |

### Protected (`/api/*`, butuh cookie `Authorization`)

<details>
<summary>Auth & user</summary>

| Method | Path | Handler |
|---|---|---|
| POST | `/api/logout` | `authCon.Logout` |
| GET | `/api/user` | `authCon.GetCurrentUser` |
| GET | `/api/user/all` | `userCon.GetAllUser` |
| PUT | `/api/update/:id` | `userCon.UpdateUser` (form-data, foto profil ≤2MB) |
| PUT | `/api/update/password/:id` | `userCon.PasswordUpdate` |

</details>

<details>
<summary>Keluarga</summary>

| Method | Path | Handler |
|---|---|---|
| GET | `/api/keluarga` | `keluargaCon.Index` |
| GET | `/api/latestkel` | top 5 keluarga terbaru |
| GET | `/api/latestkelinput` | top 1 keluarga terbaru |
| GET | `/api/keluarga/:id` | detail keluarga |
| POST | `/api/addkeluarga` | tambah (form-data, foto ≤5MB) |
| PUT | `/api/editkeluarga/:id` | update |
| DELETE | `/api/deletekeluarga/:id` | hapus |

</details>

<details>
<summary>Penduduk</summary>

| Method | Path | Handler |
|---|---|---|
| GET | `/api/penduduk` | list |
| GET | `/api/latestpend` | latest |
| GET | `/api/penduduk/:id` | detail |
| POST | `/api/addpenduduk` | tambah (JSON) |
| PUT | `/api/updatependuduk/:id` | update |
| DELETE | `/api/deletependuduk/:id` | hapus |

</details>

<details>
<summary>Statistik</summary>

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/alldata` | gabungan DataCount + AliveCount |
| GET | `/api/jumlah` | total keluarga + penduduk |
| GET | `/api/alive` | aktif vs tidak aktif |
| GET | `/api/marry` | breakdown status kawin |
| GET | `/api/gender` | breakdown gender |
| GET | `/api/data?year=YYYY&month=MM` | jumlah penduduk per hari |

</details>

---

## Auth flow

JWT HS256, dikirim sebagai HTTP-only cookie `Authorization`. Claim: `sub` (user id) + `exp` (24 jam, atau 7 hari kalau `remember_me=true`).

1. `pages/login.vue` → `POST /login` dengan email + password.
2. Backend verify bcrypt → generate JWT → set cookie via `Set-Cookie`.
3. Frontend simpan flag `localStorage.auth = 'true'` (lihat catatan keamanan), redirect `/dashboard`.
4. Setiap request axios otomatis kirim cookie karena `withCredentials: true`.
5. Logout: `POST /api/logout` (backend clear cookie) + clear localStorage.

> **Catatan:** flag `localStorage.auth` saat ini adalah satu-satunya pemicu redirect di guard frontend. Ini rentan — user bisa set manual via DevTools dan masuk ke dashboard tanpa login (walau API call tetap akan 401). Lihat [AUDIT.md](AUDIT.md) #3-4 untuk perbaikan.

### Role / level

Kolom `users.level`:

- `admin` — akses semua data lintas user
- `user` — query difilter `WHERE user_id = ?`

Pola filter konsisten di `Index`, `Latest`, `DataCount`, `AliveCount`, `MarryCount`, `GenderCount`, `RangeData`.

---

## Database

MySQL. Schema dibangun otomatis lewat `AutoMigrate` di [setup/setup.go:40-44](backend_golang/setup/setup.go#L40-L44).

### Tabel utama

- **users** — `id, name, email (unique), password, level, profile_picture`
- **keluargas** — `id, no_kk, kk_nama, alamat, rt, rw, kode_pos, status, foto_kk, foto_rumah, latitude, longtitude, user_id`
- **penduduks** — `id, kels_id, nik, nama, tmp_lahir, tgl_lahir, kelamin, stat_kawin, hub_kel, warga_neg, agama, pendidikan, pekerjaan, ayah, ibu, kepala_kel, no_hp, domisili, status, user_id`

### Relasi

- `keluargas.user_id` → `users.id`
- `penduduks.user_id` → `users.id`
- `penduduks.kels_id` → `keluargas.id`

Semua relasi: `ON UPDATE CASCADE`, `ON DELETE SET NULL`.

### Enum

Field demografis (kelamin, agama, pendidikan, pekerjaan, dll) disimpan sebagai `int8`, dipetakan ke string via helper `GetJenisKelamin`, `GetAgama`, `GetPekerjaan`, dst di [config/constant.go](backend_golang/config/constant.go).

> Frontend punya pemetaan paralel di [stores/constant.js](frontend_vuetify/src/stores/constant.js). Sinkronkan kalau menambah/mengubah enum.

---

## Konvensi

- **Foto upload** disajikan via `GET /public/uploads/...`. Frontend menyusun URL dari `baseURL + path`.
- **Filter per-user** wajib di setiap endpoint baru — cek `user.Level == "admin"` sebelum return data.
- **CORS** ketat ke satu origin (`FE_URL`). Update `.env` backend kalau frontend pindah port/host.
- **Sumber tunggal enum** ada di backend (`config/constant.go`); frontend hanya mirror.

---

## Catatan keamanan

Kelemahan yang sudah teridentifikasi (perlu diperbaiki sebelum produksi):

- JWT secret di-hardcode `"secret"` di [utils/token.go:11](backend_golang/utils/token.go#L11)
- Validasi upload hanya cek ekstensi, bukan MIME
- Password admin default `admin123`
- IDOR di endpoint `:id` (user bisa akses data user lain)
- Mass assignment via `ShouldBindJSON` ke struct GORM
- `localStorage.auth` sebagai satu-satunya pemicu router guard

**Daftar lengkap & rencana fix → [AUDIT.md](AUDIT.md).**

---

## Lisensi

Belum ditentukan. Tambahkan `LICENSE` sebelum publikasi.
