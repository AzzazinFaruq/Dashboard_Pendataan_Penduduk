# Audit Kelemahan — Project Pendataan Penduduk

Hasil scan menyeluruh terhadap `backend_golang/` dan `frontend_vuetify/` (folder `BACKEND LARAVEL/` diabaikan). Total ~50 temuan, diprioritaskan dari paling kritis.

> Tanggal audit: 2026-05-08

---

## CRITICAL — wajib diperbaiki sebelum lanjut

### Auth & otorisasi

1. **IDOR di semua endpoint resource** — `keluargaCon.go:212-240, 434-598, 600-621` & `pendudukCon.go:151-185, 255-290, 292-313`
   GET/PUT/DELETE `:id` tidak cek pemilik. User A bisa baca/edit/hapus keluarga & penduduk milik user B cuma dengan tebak ID.
   **Fix:** sebelum operasi, load record → cek `record.UserId == ctxUserId || level == "admin"`.

2. **Mass assignment via `ShouldBindJSON` ke struct GORM** — `pendudukCon.go:268`, `userCon.go:53-55`
   Client bisa kirim field `user_id`, `level`, atau apapun → langsung overwrite kolom. Termasuk user biasa promosikan diri jadi `admin` lewat `PUT /api/update/:id`.
   **Fix:** pakai DTO terpisah, whitelist field yang boleh diubah.

3. **`localStorage.auth = 'true'` sebagai sumber kebenaran** — `pages/index.vue:16`, `pages/home.vue:14`, `pages/login.vue:118`
   Buka DevTools → set `auth = 'true'` → masuk dashboard tanpa login. Backend tetap reject API call (cookie JWT tidak ada), tapi UI keburu render & flash data sensitif.
   **Fix:** router `beforeEach` panggil `GET /api/user`; kalau 401 → redirect login. Buang flag localStorage.

4. **Tidak ada router guard** — `src/router/index.js` (file praktis kosong)
   Akses `/dashboard`, `/penduduk`, `/keluarga` langsung via URL tidak diproteksi di level router. Cek auth cuma di `mounted()` masing-masing page → flash konten.
   **Fix:** tambah `router.beforeEach` global.

5. **JWT secret hardcode `"secret"`** — [backend_golang/utils/token.go:11](backend_golang/utils/token.go#L11)
   Kalau repo ini pernah public/leak, semua token forgeable.
   **Fix:** baca dari env (`JWT_SECRET`), validasi non-empty di startup, generate random ≥32 byte.

6. **Cookie JWT tanpa flag aman** — [backend_golang/controllers/authCon.go:91](backend_golang/controllers/authCon.go#L91)
   `Secure=false`, `SameSite` tidak di-set. Rentan CSRF + bisa ke-MITM di HTTP.
   **Fix:** `Secure=true` di prod, `SameSite=Lax` (atau `Strict`), `HttpOnly=true` (sudah aktif).

### File upload & path

7. **Upload pakai filename asli** — `userCon.go:75`, `keluargaCon.go:365, 389, 468, 490`
   Path traversal (`../../something.jpg`), overwrite file user lain dengan nama sama, eksekusi shell-like extension.
   **Fix:** rename ke `uuid.New().String() + ext`, validasi MIME asli (`http.DetectContentType` 512 byte pertama), bukan cuma ekstensi.

8. **File tidak di-rollback saat insert DB gagal** — `keluargaCon.go:351-397`
   File foto sudah tersimpan di disk, lalu DB Create error → orphan file menumpuk.
   **Fix:** `defer` cleanup, atau pakai `tx` lalu hapus file kalau rollback.

### Bug langsung

9. **Bug parsing lat/long → tertulis ke `statusInt`** — `keluargaCon.go:528-543`
   Error parsing latitude ditulis ke variable salah. Geo data rusak senyap.
   **Fix:** assign ke variable yang benar.

10. **`jsconfig.json` alias salah arah** — [frontend_vuetify/jsconfig.json:9-11](frontend_vuetify/jsconfig.json#L9-L11)
    `"@/*": ["https://apikk.spora.id/api"]` → URL eksternal, bukan path. Vite tetap jalan (alias dari `vite.config.mjs`), tapi IDE intellisense rusak.
    **Fix:** `["./src/*"]`.

---

## HIGH — perbaiki segera setelah CRITICAL

### Backend

11. **Statistik tidak filter `user_id` untuk non-admin di sebagian endpoint** — `universalCon.go:37-38`
    Cek konsistensi semua endpoint statistik. Bisa membocorkan total agregat data milik orang lain.

12. **`AutoMigrate` jalan setiap startup** — [setup/setup.go:40-44](backend_golang/setup/setup.go#L40-L44)
    Di production, perubahan model → tabel berubah otomatis tanpa review.
    **Fix:** ganti dengan tool migrasi (`golang-migrate`, `goose`) + jangan jalankan `AutoMigrate` di env prod.

13. **Tidak ada validasi unik untuk NIK & NoKK** — `keluargaCon.go:243-293`, `pendudukCon.go`
    NIK & NoKK secara hukum unik. Saat ini bisa duplikat.
    **Fix:** `gorm:"uniqueIndex"` di model + handle error 23000 saat insert.

14. **Brute force login tanpa rate limit** — `authCon.go:55-98`
    Tidak ada delay/lockout.
    **Fix:** rate limiter per IP (mis. `ulule/limiter`) atau exponential backoff per email.

15. **Error message bocor ke client** — banyak controller (`return err.Error()`)
    GORM/MySQL error berisi nama tabel & kolom.
    **Fix:** map ke pesan generik untuk client; log detail di server.

### Frontend

16. **Hardcode `baseURL: http://localhost:8080`** — [src/main.js:42](frontend_vuetify/src/main.js#L42)
    Tidak bisa deploy production tanpa rebuild.
    **Fix:** `import.meta.env.VITE_API_BASE_URL` + `.env.example`.

17. **Halaman raksasa** — `pages/penduduk/inputPenduduk.vue` (~542 baris), `pages/keluarga/inputKeluarga.vue`
    Susah di-maintain & di-test.
    **Fix:** pecah jadi sub-component (form identitas, form alamat, form berkas).

18. **Submit button tidak di-disable saat loading** — `pages/login.vue:34`, `pages/penduduk/inputPenduduk.vue:290`, `register.vue`
    Double-submit menghasilkan duplikat data.
    **Fix:** `:loading` + `:disabled="submitting"`.

19. **Hardcoded default password `"admin123"` di kode profile** — `pages/profile.vue:89`
    Verifikasi manual. Kalau betul, harus dihapus.

---

## MEDIUM — improve kualitas

### Backend

20. **CORS via env `FE_URL` saja** — `middlewares/cors.go` — pastikan tidak fallback ke `*` saat env kosong.
21. **Tidak ada index DB di `user_id`** di `keluargas` & `penduduks` — query filter `user_id` sering. Tambah `gorm:"index"`.
22. **Route naming tidak RESTful** — `/addkeluarga`, `/editkeluarga`, `/deletekeluarga`. Konsisten ke REST: `POST /api/keluarga`, `PUT /api/keluarga/:id`, dll.
23. **Magic numbers** untuk enum status/kelamin di `universalCon.go:58, 161-169` — bikin `const` di package config.
24. **Hard delete** — keluarga/penduduk dihapus permanen, riwayat hilang. Pertimbangkan soft delete (`gorm.DeletedAt`).
25. **Tidak ada audit log** untuk Create/Update/Delete data sensitif.
26. **`UserData()` dead code** — [userCon.go:141-217](backend_golang/controllers/userCon.go#L141-L217), tidak terdaftar di route.
27. **Time zone `loc=Local`** — kalau server pindah region, tanggal lahir/created_at miss. UTC lebih aman.
28. **Endpoint `/api/data` validasi `year`/`month`** — terima string mentah ke `YEAR()/MONTH()`. Walau parameterized, validasi range tetap perlu.
29. **Password seeder admin `admin123`** — idealnya force-change pada first login.

### Frontend

30. **Pinia `restrict.js` side effect di action `setup`**, bukan reactive state. Refactor ke `useAuthStore` dengan `state.user`, computed `isAuthenticated`.
31. **Duplikasi enum** — backend `config/constant.go`, frontend mirror di `stores/constant.js`. Tambah enum baru → dua tempat. **Fix:** ekspos endpoint `/api/constants` & fetch sekali.
32. **Inline mapping ID → label** tersebar (`dtTable.vue`, `dashboard/index.vue`) — pindah ke util.
33. **Tidak ada cleanup `axios` saat unmount** — pakai `AbortController` untuk fetch di `onMounted`.
34. **Search navbar tanpa debounce** — kalau nanti dipakai panggil API.
35. **`vite.config.mjs:47` proxy `apikk.spora.id`** dead config — hapus.
36. **ESLint cuma `essential`** — naikkan ke `recommended`, tambah `no-console` untuk production build.
37. **Bundle besar** — Chart.js + ApexCharts redundan; Leaflet hanya dipakai kalau ada peta. Audit & buang yang tidak dipakai.
38. **`v-html`** — cek manual, tidak boleh ada `v-html` dengan input user.
39. **Logout** hanya hapus localStorage; pastikan juga panggil `POST /api/logout` agar cookie di-clear server-side.

### Cross-cutting

40. **Tidak ada test sama sekali** — backend (Go test) maupun frontend (Vitest/Vue Test Utils). Minimal cover auth + IDOR check + form validasi.
41. **Tidak ada CI** — lint + test + build harusnya jalan otomatis di PR.
42. **Tidak ada `.env.example`** di backend & frontend → onboarding susah.
43. **`db_project_pendataan_penduduk.sql`** dump pernah ter-commit ke history — kalau ada data sensitif (NIK asli), perlu rewrite history (`git filter-repo`). Sekarang sudah di-untrack via `.gitignore`.

---

## LOW — polish

- `console.log` data sensitif di `pages/login.vue:114`, `inputPenduduk.vue:382`.
- Toast SweetAlert 3000ms terlalu pendek untuk error panjang.
- Form `<label>` tanpa `for` link → screen reader fail.
- Tidak ada TypeScript / type safety.
- Logo & avatar default tanpa `alt` text di beberapa tempat.
- Source map: pastikan tidak terbawa ke production build.

---

## Saran urutan pengerjaan

1. **Hari 1-2 (security darurat):** #1 IDOR + #2 mass assignment + #5 JWT secret + #7 file upload rename. Kalau aplikasi sudah online dengan user nyata, ini paling urgent.
2. **Hari 3-4 (auth bersih):** #3-4 router guard + ganti `localStorage.auth` jadi cek `GET /api/user`, #6 cookie flag, #15 generic error.
3. **Minggu 1:** #11-13 (statistik filter, migrasi tool, unique NIK), #16 env-driven baseURL, #18 disable button, #14 rate limit.
4. **Selanjutnya:** refactor halaman besar (#17), test (#40), CI (#41), cleanup dead code & duplikasi.

**Rekomendasi mulai dari:** #1 IDOR + #2 mass assignment. Dengan kondisi sekarang, satu user bisa baca/edit/hapus data semua user lain hanya dengan tebak `:id` — itu masalah paling besar di sini.

---

## Checklist progress

### Critical
- [ ] #1 IDOR endpoint resource (keluarga + penduduk)
- [ ] #2 Mass assignment (DTO + whitelist)
- [ ] #3 Buang `localStorage.auth` sebagai source of truth
- [ ] #4 Router guard `beforeEach`
- [ ] #5 JWT secret dari env
- [ ] #6 Cookie flag aman (Secure, SameSite)
- [ ] #7 Rename file upload + MIME validation
- [ ] #8 Rollback file saat DB gagal
- [ ] #9 Bug parsing lat/long
- [ ] #10 Fix `jsconfig.json` alias

### High
- [ ] #11 Filter user_id konsisten di statistik
- [ ] #12 Ganti AutoMigrate dengan migration tool
- [ ] #13 Unique NIK & NoKK
- [ ] #14 Rate limit login
- [ ] #15 Generic error message
- [ ] #16 baseURL dari env
- [ ] #17 Pecah halaman besar
- [ ] #18 Disable submit saat loading
- [ ] #19 Hapus hardcoded password di profile

### Medium
- [ ] #20 CORS fallback
- [ ] #21 Index DB di user_id
- [ ] #22 RESTful route naming
- [ ] #23 Const untuk magic number enum
- [ ] #24 Soft delete
- [ ] #25 Audit log
- [ ] #26 Hapus dead code `UserData()`
- [ ] #27 UTC time zone
- [ ] #28 Validasi range year/month
- [ ] #29 Force change admin password
- [ ] #30 Refactor `restrict.js` ke proper Pinia
- [ ] #31 Endpoint `/api/constants`
- [ ] #32 Util mapping ID→label
- [ ] #33 AbortController cleanup
- [ ] #34 Debounce search
- [ ] #35 Hapus dead proxy `apikk.spora.id`
- [ ] #36 ESLint recommended + no-console
- [ ] #37 Buang library tidak terpakai
- [ ] #38 Audit `v-html`
- [ ] #39 Logout panggil backend
- [ ] #40 Test (backend + frontend)
- [ ] #41 CI pipeline
- [ ] #42 `.env.example`
- [ ] #43 Rewrite history bila SQL dump berisi data sensitif

### Low
- [ ] Hapus `console.log` sensitif
- [ ] Perpanjang durasi toast
- [ ] Aksesibilitas form `<label for>`
- [ ] Migrasi ke TypeScript
- [ ] `alt` text untuk gambar
- [ ] Disable source map di prod
