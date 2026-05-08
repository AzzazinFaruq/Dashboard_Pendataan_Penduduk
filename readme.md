# Project Pendataan Penduduk

Sistem pendataan penduduk & keluarga (kelurahan/desa). Repo berisi dua sub-proyek aktif:

- `backend_golang/` — REST API (Go + Gin + GORM + MySQL)
- `frontend_vuetify/` — SPA (Vue 3 + Vuetify + Pinia + Vite)

> Folder `BACKEND LARAVEL/` ada di repo tapi **tidak digunakan**; dokumen ini sengaja mengabaikannya.

---

## Quickstart

### Backend (Go)

```bash
cd backend_golang
cp .env.example .env   # isi DB_* dan FE_URL
go mod download
go run main.go         # listen di :8080
```

`.env` yang dibutuhkan:

```
FE_URL=http://localhost:3000
DB_USERNAME=root
DB_PASSWORD=
DB_HOST=localhost
DB_PORT=3306
DB_NAME=project_pendataan_penduduk
```

Saat startup, `setup.ConnectDatabase` (`backend_golang/setup/setup.go`) menjalankan `AutoMigrate` untuk tabel `users`, `keluargas`, `penduduks`, lalu `SeedersUser` membuat admin default:

- email: `admin@gmail.com`
- password: `admin123`

### Frontend (Vuetify)

```bash
cd frontend_vuetify
npm install
npm run dev            # listen di :3000
```

Base URL backend di-hardcode di `frontend_vuetify/src/main.js:42` (`http://localhost:8080`). Ubah di sana kalau backend pindah host.

---

## Backend — `backend_golang/`

### Stack

| Komponen | Versi | Keterangan |
|---|---|---|
| Go | 1.23.2 | dari `go.mod` |
| Gin | v1.10.0 | HTTP framework |
| GORM | v1.25.12 | ORM, driver MySQL v1.5.7 |
| jwt/v5 | v5.2.1 | JWT HS256 |
| godotenv | v1.5.1 | load `.env` |
| bcrypt | `golang.org/x/crypto` | hashing password |

### Struktur

```
backend_golang/
├── main.go                 # entrypoint, register route
├── setup/setup.go          # GORM connect + AutoMigrate + Seed
├── seeders/userSeeders.go  # admin default (idempotent)
├── config/constant.go      # enum + helper Get*() (gender, agama, dll)
├── middlewares/
│   ├── cors.go             # CORS dari env FE_URL
│   └── jwt.go              # AuthMiddleware (cookie "Authorization")
├── models/
│   ├── User.go             # id, name, email, password, level, profile_picture
│   ├── Keluarga.go         # no_kk, kk_nama, alamat, RT/RW, foto_kk, foto_rumah, lat/long, user_id
│   └── Penduduk.go         # nik, nama, kelamin, agama, pendidikan, pekerjaan, kels_id, user_id, dll
├── controllers/
│   ├── authCon.go          # Register, Login, GetCurrentUser, Logout
│   ├── userCon.go          # GetAllUser, UpdateUser, PasswordUpdate
│   ├── keluargaCon.go      # CRUD keluarga + upload foto KK & rumah
│   ├── pendudukCon.go      # CRUD penduduk
│   └── universalCon.go     # statistik: AllData, DataCount, AliveCount, MarryCount, GenderCount, RangeData
├── utils/token.go          # GenerateJWT, ValidateJWT
└── public/uploads/         # foto-kk/, foto-rumah/, profile_pictures/
```

### Endpoint

Public:

| Method | Path | Handler |
|---|---|---|
| POST | `/register` | `authCon.Register` |
| POST | `/login` | `authCon.Login` (set cookie `Authorization`) |
| GET | `/public/*` | static files (uploads) |

Protected (`/api/*`, butuh cookie `Authorization` valid):

| Method | Path | Handler |
|---|---|---|
| POST | `/api/logout` | `authCon.Logout` |
| GET | `/api/user` | `authCon.GetCurrentUser` |
| GET | `/api/user/all` | `userCon.GetAllUser` |
| PUT | `/api/update/:id` | `userCon.UpdateUser` (form-data, foto profil ≤2MB) |
| PUT | `/api/update/password/:id` | `userCon.PasswordUpdate` |
| GET | `/api/keluarga` | `keluargaCon.Index` |
| GET | `/api/latestkel` | `keluargaCon.Latest` (top 5) |
| GET | `/api/latestkelinput` | `keluargaCon.LatestForInput` (top 1) |
| GET | `/api/keluarga/:id` | `keluargaCon.GetKeluargaByID` |
| POST | `/api/addkeluarga` | `keluargaCon.AddKeluarga` (form-data, foto ≤5MB) |
| PUT | `/api/editkeluarga/:id` | `keluargaCon.UpdateKeluarga` |
| DELETE | `/api/deletekeluarga/:id` | `keluargaCon.DeleteKeluarga` |
| GET | `/api/penduduk` | `pendudukCon.GetPenduduk` |
| GET | `/api/latestpend` | `pendudukCon.GetLatestPenduduk` |
| GET | `/api/penduduk/:id` | `pendudukCon.GetPendudukByID` |
| POST | `/api/addpenduduk` | `pendudukCon.AddPenduduk` (JSON) |
| PUT | `/api/updatependuduk/:id` | `pendudukCon.UpdatePenduduk` |
| DELETE | `/api/deletependuduk/:id` | `pendudukCon.DeletePenduduk` |
| GET | `/api/alldata` | gabungan DataCount + AliveCount |
| GET | `/api/jumlah` | total keluarga + penduduk |
| GET | `/api/alive` | aktif vs tidak aktif |
| GET | `/api/marry` | breakdown status kawin |
| GET | `/api/gender` | breakdown gender |
| GET | `/api/data?year=YYYY&month=MM` | jumlah penduduk per hari |

### Auth & role

JWT HS256, dikirim sebagai HTTP cookie `Authorization`. Claim: `sub` (user id) + `exp` (24 jam, atau 7 hari kalau `remember_me`).

Role di kolom `users.level`:

- `admin` — lihat semua data
- `user` — hanya data milik sendiri (filter `user_id = ?` di setiap query list/statistik)

Pola filter di [keluargaCon.go:32-36](backend_golang/controllers/keluargaCon.go#L32-L36) dipakai konsisten di `GetPenduduk`, `Latest`, `DataCount`, `AliveCount`, `MarryCount`, `GenderCount`, `RangeData`.

### Catatan keamanan

- JWT secret di-hardcode `"secret"` di [utils/token.go:11](backend_golang/utils/token.go#L11) — pindahkan ke env.
- Validasi upload hanya cek ekstensi, bukan MIME type.
- Password admin default `admin123` — wajib diganti setelah deploy.

---

## Frontend — `frontend_vuetify/`

### Stack

| Komponen | Versi |
|---|---|
| Vue | 3.4.31 |
| Vuetify | 3.6.11 (theme `myCustomTheme`, primary `#795548`) |
| Vue Router | 4.4.0 (auto-routing via `unplugin-vue-router`) |
| Pinia | 2.1.7 |
| Vite | 5.1.5 (dev di port 3000) |
| Axios | 1.7.7 + `axios-retry` 4.5.0 |
| ApexCharts, Chart.js | charting |
| Leaflet, vue-leaflet | peta (lat/long keluarga) |
| vue-sweetalert2 | toast/alert |

### Script

```
npm run dev      # vite dev
npm run build    # production build → dist/
npm run preview  # preview dist
npm run lint     # eslint --fix
```

### Struktur `src/`

```
src/
├── main.js                 # bootstrap: Pinia, Axios, SweetAlert2, ApexCharts
├── App.vue                 # root, sembunyikan navbar di /login /register / /forbidden
├── router/index.js         # createRouter (file-based dari pages/)
├── pages/                  # auto-route
│   ├── index.vue           # redirect awal
│   ├── login.vue, register.vue, forbidden.vue
│   ├── home.vue, about.vue, profile.vue
│   ├── dashboard/index.vue
│   ├── penduduk/index.vue, penduduk/inputPenduduk.vue, penduduk/edit/[id].vue
│   └── keluarga/index.vue, keluarga/inputKeluarga.vue, keluarga/edit/[id].vue
├── components/
│   ├── navbar.vue, AppFooter.vue
│   ├── card.vue            # stat card dashboard
│   ├── dtTable.vue         # data table reusable
│   ├── formPenduduk.vue, inputKeluarga.vue
│   ├── chart/{gender,marry,total}chart.vue
│   └── modal/detailPenduduk.vue
├── stores/
│   ├── constant.js         # useCons — dropdown enum (kelamin, agama, pekerjaan, dll)
│   ├── nav.js              # useNav — link sidebar (links1 standar, links2 + admin)
│   ├── restrict.js         # test — guard: redirect /login kalau !localStorage.auth
│   └── title.js            # useTitle — page title dinamis
├── plugins/
│   ├── index.js            # registerPlugins(app)
│   └── vuetify.js          # theme + lab components (VDateInput, VNumberInput)
├── styles/{settings.scss,global.css}
└── assets/                 # logo, avatar default
```

### Konfigurasi Axios

Di [src/main.js:42-45](frontend_vuetify/src/main.js#L42-L45):

- `baseURL: http://localhost:8080`
- `withCredentials: true`, `withXSRFToken: true` — JWT cookie ikut otomatis
- Retry 3× untuk status 429 (eksponensial)
- Response interceptor: 401 → hapus `localStorage.auth`

### Auth flow

1. `pages/login.vue` → `POST /login` → backend set cookie `Authorization`.
2. Kalau sukses: `localStorage.setItem('auth', 'true')`, toast SweetAlert2, redirect `/dashboard`.
3. Setiap halaman protected memanggil action `setup()` di store `test` (`stores/restrict.js`) — kalau `localStorage.auth !== 'true'`, push ke `/login`.
4. Cookie HTTP-only dikirim oleh browser otomatis di setiap request `axios` karena `withCredentials: true`.
5. Logout: `POST /api/logout` (backend clear cookie) + clear localStorage. Interceptor 401 juga otomatis bersihkan localStorage.

### Env / hal yang perlu diketahui

- Tidak ada `.env` di frontend — `baseURL` di-hardcode. Kalau mau pakai env, ganti `http://localhost:8080` di `main.js` jadi `import.meta.env.VITE_API_BASE_URL`.
- `vite.config.mjs:47` punya proxy ke `https://apikk.spora.id` — sisa konfigurasi lama, tidak dipakai oleh axios karena baseURL di-set absolut.
- `jsconfig.json` punya alias `@` yang salah arah (ke URL eksternal). Vite tetap pakai alias dari `vite.config.mjs:41` (`@ → ./src`), jadi runtime aman; perbaiki jsconfig kalau IDE-nya error resolve import.

---

## Database

MySQL, schema awal dari `db_project_pendataan_penduduk.sql` (sudah di-`.gitignore`). Saat development cukup biarkan `AutoMigrate` GORM yang menyiapkan tabel.

Tabel utama:

- **users** — `id, name, email (unique), password, level, profile_picture`
- **keluargas** — `id, no_kk, kk_nama, alamat, rt, rw, kode_pos, status, foto_kk, foto_rumah, latitude, longtitude, user_id`
- **penduduks** — `id, kels_id, nik, nama, tmp_lahir, tgl_lahir, kelamin, stat_kawin, hub_kel, warga_neg, agama, pendidikan, pekerjaan, ayah, ibu, kepala_kel, no_hp, domisili, status, user_id`

Relasi: `keluargas.user_id → users.id`, `penduduks.user_id → users.id`, `penduduks.kels_id → keluargas.id` (semua `ON UPDATE CASCADE`, `ON DELETE SET NULL`).

Enum disimpan sebagai `int8` di DB, dipetakan ke string via helper di `backend_golang/config/constant.go` (`GetJenisKelamin`, `GetAgama`, `GetPekerjaan`, dll). Frontend punya pemetaan paralel di `frontend_vuetify/src/stores/constant.js` — jaga keduanya tetap sinkron kalau menambah/mengubah enum.

---

## Konvensi & Tips

- **Foto upload** disajikan via `GET /public/uploads/...`. Frontend menyusun URL absolut dari `baseURL + path` yang dikembalikan API.
- **Filter per-user** ada di hampir semua list & statistik backend; cek `user.Level == "admin"` sebelum nambah endpoint baru biar konsisten.
- **CORS** ketat ke satu origin (`FE_URL`). Kalau frontend pindah origin/port, update `.env` backend.
- **Dropdown sumber-tunggal** sebaiknya tetap di `config/constant.go` (backend) — frontend hanya mirror untuk UX.
