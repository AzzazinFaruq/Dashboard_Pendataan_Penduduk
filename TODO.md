# Checklist Perbaikan — Dashboard Pendataan Penduduk

## 🔴 P0 — Keamanan (kritikal)
- [x] Pindahkan JWT secret ke `os.Getenv("JWT_SECRET")` (`utils/token.go`)
- [x] Tambah cek kepemilikan (`user_id`) di `GetPendudukByID`, `GetKeluargaByID`
- [x] Tambah cek kepemilikan di `UpdateKeluarga`, `DeleteKeluarga`, `UpdatePenduduk`, `DeletePenduduk`
- [x] Ambil `user_id` dari token (bukan form/body) di `AddKeluarga` & `AddPenduduk`
- [x] Upload: generate nama file unik (UUID/timestamp) + validasi MIME (bukan ekstensi string) — `keluargaCon.go`, `userCon.go`
- [x] Tambah `gin.Recovery()` (atau `gin.Default()`) di `main.go`
- [x] Batasi `GetAllUser` hanya untuk admin

## 🟠 P1 — Bug logika
### Backend
- [x] Fix latitude/longitude di `UpdateKeluarga` (shadow `:=`, salah assign ke `statusInt`, pesan error salah)
- [x] Fix `RememberMe`: `GenerateJWT` harus terima durasi (exp saat ini selalu 24 jam)
- [x] Fix `/api/jumlah` & `/api/alive` (kembalikan JSON, bukan `c.Keys`)
- [x] Fix logika tahun kabisat di `RangeData` (Februari salah hitung hari)
- [x] Fix `GetAllUser`: `Order(...).Find(...)` (Order setelah Find no-op)
- [x] Hapus blok `if isComplete {}` kosong & perbaiki validasi `AddKeluarga`
- [x] Hapus file foto saat `DeleteKeluarga`/`DeletePenduduk`
- [x] Perbaiki `total` agar per-user untuk non-admin (`GetPenduduk`, `Index`)
### Frontend
- [x] Perbaiki konflik dua store ber-`id: "options"` (`nav.js` & `constant.js`)
- [x] Tambah route guard global (`beforeEach`) di `router/index.js`
- [x] Ganti pola `try/catch` sinkron → `.catch()`/`async-await` (navbar, login, dll)
- [x] Perbaiki avatar dobel di `navbar.vue`

## 🟡 P2 — Kebersihan & config
- [x] Pindahkan base URL ke `VITE_API_URL` (`main.js`, hardcode di `navbar.vue`)
- [x] Buat `.env.example` (backend & frontend)
- [x] Lengkapi map `Pekerjaan` (11–88) di `config/constant.go`
- [x] Ganti global mutable export (`export var name/success`) → Pinia store

## 🗑️ Hapus (berlebihan)
### Frontend
- [x] `pages/test.vue`, `components/HelloWorld.vue`, `src/api/authAPI.js`, `src/script/const.js`
- [x] `stores/nav.js` (tak dipakai + buggy)
- [x] `vite.config.mjs.timestamp-*.mjs` + tambahkan ke `.gitignore`
- [x] Hapus (sudah dicek): `pages/about.vue`, `pages/home.vue`, `pages/edit.vue`
- [x] README scaffold di `components/`, `pages/`, `plugins/`, `styles/`
- [x] Dependency: `javascript`, `vuex`, `cors`, `vue-axios`, `chart.js`, `vue-chartjs`, `vue-number-input`, `@chenfengyuan/vue-number-input`
### Backend
- [x] Map `StatusSurat`, `StatusSuratKet`, `Golongan` di `config/constant.go` (di luar scope)
- [x] Fungsi `UserData()` (dead code, tak diroute + salah cast)

## ➕ Yang kurang (opsional/nice-to-have)
- [x] Unit test backend (config + utils) — *frontend belum (butuh vitest)*
- [x] Pagination opsional pada list penduduk/keluarga (`?page=&limit=`)
- [x] README backend
- [x] Dockerfile (backend + frontend) + docker-compose + CI (GitHub Actions)
