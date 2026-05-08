// Generator SQL seed untuk tabel `keluargas` & `penduduks`
// Menghasilkan 500 keluarga + 500 penduduk (1 kepala keluarga per KK).
// Jalankan: node generate_seed.js  -> output: seed_data.sql

const fs = require('fs');
const path = require('path');

// ---------- PRNG deterministik (Mulberry32) ----------
let _seed = 0xC0FFEE;
function rand() {
  _seed |= 0;
  _seed = (_seed + 0x6D2B79F5) | 0;
  let t = Math.imul(_seed ^ (_seed >>> 15), 1 | _seed);
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
}
const ri = (a, b) => Math.floor(rand() * (b - a + 1)) + a;
const pick = (arr) => arr[ri(0, arr.length - 1)];

// ---------- Data dasar ----------
const maleFirst = ['Budi','Agus','Slamet','Bambang','Sutrisno','Andi','Joko','Hendra','Rudi','Eko','Dedi','Rendi','Yanto','Hadi','Wahyu','Imam','Faisal','Ahmad','Hasan','Iwan','Fajar','Reza','Dimas','Bayu','Aditya','Galih','Yusuf','Iqbal','Arif','Surya','Anton','Hanif','Rama','Dani','Gilang','Hari','Iman','Krisna','Lukman','Mahendra','Nanda','Oscar','Pramana','Riski','Setiawan','Toni','Untung','Vito','Wibowo','Yogi','Zaki','Bagas','Candra','Dharma','Erlangga','Firman','Gilbert','Hario','Ilham','Jefri','Kurnia'];
const femaleFirst = ['Siti','Ayu','Indah','Dewi','Rina','Lina','Sari','Tuti','Wati','Yuni','Ani','Endang','Fitri','Gita','Hesti','Ira','Juli','Kartika','Lestari','Maya','Nita','Olivia','Putri','Qori','Rani','Susi','Tina','Umi','Vera','Wulan','Yanti','Zahra','Anisa','Bella','Cinta','Dinda','Eka','Farah','Galuh','Hana','Intan','Jamilah','Kamila','Lia','Mei','Novi','Octa','Prita','Rita','Sinta','Tika','Ulin','Vita','Winda','Yana','Zaskia','Asri','Bunga','Citra','Diah','Erika'];
const lastNames = ['Santoso','Wijaya','Hartanto','Susilo','Pratama','Saputra','Setiawan','Hidayat','Permana','Kusuma','Nugroho','Putra','Iskandar','Rahardjo','Mulyadi','Wibowo','Hartono','Subagyo','Sutanto','Hermawan','Hakim','Salim','Yusuf','Utomo','Surya','Maulana','Pradana','Rifai','Anggara','Pradipta','Adiwijaya','Suryadi','Setyawan','Mahendra','Anggraini','Ramadhan','Aprilia','Lestari','Wicaksono','Pranoto','Suharto','Sumarno','Cahyono','Darmawan','Effendi','Gunawan','Halim','Irawan','Jatmiko','Kurniawan','Latif','Marwan','Nasution','Oktavianus','Purnomo','Qodri','Rosyid','Sukamto','Taufik','Utama','Vivaldi','Widodo'];
const cities = ['Malang','Surabaya','Jakarta','Bandung','Yogyakarta','Semarang','Solo','Medan','Makassar','Denpasar','Pasuruan','Probolinggo','Lumajang','Kediri','Madiun','Magetan','Trenggalek','Tulungagung','Blitar','Banyuwangi','Jember','Bondowoso','Situbondo','Nganjuk','Mojokerto','Sidoarjo','Gresik','Lamongan','Tuban','Bojonegoro','Pamekasan','Sumenep','Bangkalan','Sampang','Ngawi','Ponorogo','Pacitan','Batu','Kepanjen','Singosari'];
const streetNames = ['Merdeka','Sudirman','Diponegoro','Veteran','Ahmad Yani','Gajah Mada','Pahlawan','Kartini','Cut Nyak Dien','Imam Bonjol','Ir. Soekarno','Pattimura','Hasanuddin','Tentara Pelajar','Brigjen Slamet Riyadi','Kawi','Semeru','Bromo','Arjuno','Welirang','Panderman','Buring','Ijen','Kalpataru','Soekarno-Hatta','Kebalen Wetan','Bandulan','Galunggung','Bukirsari','Wilis'];
const villages = ['Sumberporong','Plaosan','Sukoharjo','Sumberkembar','Tegalweru','Karangbesuki','Tlogomas','Lowokwaru','Dinoyo','Merjosari','Tasikmadu','Tunggulwulung','Mojolangu','Tunjungsekar','Polowijen','Blimbing','Pandanwangi','Purwantoro','Bunulrejo','Polehan','Kesatrian','Jodipan','Kotalama','Mergosono','Bandungrejosari','Sukun','Bakalankrajan','Mulyorejo','Pisangcandi','Karangbesuki'];
const agamaWeights = [{v:1,w:75},{v:2,w:10},{v:3,w:5},{v:4,w:5},{v:5,w:3},{v:6,w:1},{v:7,w:1}];
const pekerjaanWeights = [{v:1,w:15},{v:2,w:15},{v:3,w:15},{v:4,w:5},{v:5,w:8},{v:6,w:2},{v:7,w:2},{v:8,w:15},{v:9,w:13},{v:10,w:5},{v:89,w:5}];
const pendidikanWeights = [{v:1,w:5},{v:2,w:5},{v:3,w:15},{v:4,w:20},{v:5,w:30},{v:6,w:5},{v:7,w:5},{v:8,w:12},{v:9,w:2},{v:10,w:1}];

function weighted(arr) {
  const total = arr.reduce((s, x) => s + x.w, 0);
  let r = rand() * total;
  for (const x of arr) { if ((r -= x.w) <= 0) return x.v; }
  return arr[arr.length - 1].v;
}

// ---------- Helpers ----------
const sqlEsc = (s) => "'" + String(s).replace(/'/g, "''") + "'";
const pad = (n, len) => String(n).padStart(len, '0');

function fmtDate(d) {
  return `${d.getFullYear()}-${pad(d.getMonth()+1,2)}-${pad(d.getDate(),2)} ${pad(d.getHours(),2)}:${pad(d.getMinutes(),2)}:${pad(d.getSeconds(),2)}`;
}

function randomKK() {
  // 12 digit, sesuai dengan style data lama (mis: 123412341234)
  let s = String(ri(1, 9));
  for (let i = 0; i < 11; i++) s += String(ri(0, 9));
  return s;
}
function randomNIK() {
  // 16 digit
  let s = String(ri(1, 9));
  for (let i = 0; i < 15; i++) s += String(ri(0, 9));
  return s;
}
function randomHP() {
  const prefixes = ['0812','0813','0821','0822','0852','0853','0856','0857','0858','0878','0895','0896'];
  let s = pick(prefixes);
  for (let i = 0; i < 8; i++) s += String(ri(0, 9));
  return s;
}
function randomBirth(minYear, maxYear) {
  const y = ri(minYear, maxYear);
  const m = ri(1, 12);
  const d = ri(1, 28);
  return new Date(y, m - 1, d, ri(0, 23), ri(0, 59), ri(0, 59));
}

const USER_ID = 1; // admin user id (hasil SeedersUser)
const NUM_KELUARGA = 500;
const NUM_PENDUDUK = 500;

// ---------- Generate KELUARGA ----------
const keluargaRows = [];
const kkInfo = []; // simpan info KK utk relasi dengan penduduk

for (let i = 1; i <= NUM_KELUARGA; i++) {
  const noKk = randomKK();
  const kkNama = `${pick(maleFirst)} ${pick(lastNames)}`;
  const villageName = pick(villages);
  const alamat = `Jl. ${pick(streetNames)} No. ${ri(1, 250)}, Kel. ${villageName}`;
  const rt = pad(ri(1, 15), 2);
  const rw = pad(ri(1, 12), 2);
  const kodePos = String(ri(60000, 69999));
  const status = ri(1, 100) <= 92 ? 1 : 2;
  const fotoKk = `public/uploads/foto-kk/sample-kk-${(i % 10) + 1}.jpg`;
  const fotoRumah = `public/uploads/foto-rumah/sample-rumah-${(i % 10) + 1}.jpg`;
  // Latitude/longitude sekitar Malang Raya
  const lat = (-7.95 + (rand() * 0.35)).toFixed(6);   // -7.95 .. -7.60
  const lon = (112.55 + (rand() * 0.40)).toFixed(6);  // 112.55 .. 112.95
  const created = fmtDate(new Date(2024, ri(0, 11), ri(1, 28), ri(0, 23), ri(0, 59), ri(0, 59)));

  kkInfo.push({ id: i, kkNama });
  keluargaRows.push(
    `(${i}, ${noKk}, ${sqlEsc(kkNama)}, ${sqlEsc(alamat)}, ${sqlEsc(rt)}, ${sqlEsc(rw)}, ${sqlEsc(kodePos)}, ${status}, ${USER_ID}, ${sqlEsc(fotoKk)}, ${sqlEsc(fotoRumah)}, ${lat}, ${lon}, ${sqlEsc(created)}, ${sqlEsc(created)})`
  );
}

// ---------- Generate PENDUDUK ----------
// Distribusi: setiap penduduk ditugaskan ke salah satu keluarga (1..500).
// Penduduk pertama untuk setiap kels_id menjadi KEPALA KELUARGA (hub_kel = 1, kelamin = 1).
const pendudukRows = [];
const seenKepala = new Set();

for (let i = 1; i <= NUM_PENDUDUK; i++) {
  // Ambil kels_id berurutan dulu agar 500 keluarga pertama dapat kepala,
  // kemudian sisanya (jika NUM_PENDUDUK > NUM_KELUARGA) jadi anggota acak.
  const kelsId = i <= NUM_KELUARGA ? i : ri(1, NUM_KELUARGA);
  const kk = kkInfo[kelsId - 1];

  let hubKel, kelamin, nama, statKawin;
  if (!seenKepala.has(kelsId)) {
    seenKepala.add(kelsId);
    hubKel = 1; // KEPALA KELUARGA
    kelamin = 1; // laki-laki
    nama = kk.kkNama;
    statKawin = ri(1, 100) <= 85 ? 2 : 1; // 85% kawin
  } else {
    const role = ri(1, 100);
    if (role <= 25) {           // istri
      hubKel = 3; kelamin = 2;
      nama = `${pick(femaleFirst)} ${pick(lastNames)}`;
      statKawin = 2;
    } else if (role <= 80) {    // anak
      hubKel = 4;
      kelamin = ri(1, 2);
      nama = `${kelamin === 1 ? pick(maleFirst) : pick(femaleFirst)} ${kk.kkNama.split(' ').slice(-1)[0]}`;
      statKawin = 1;
    } else if (role <= 90) {    // orang tua
      hubKel = 7;
      kelamin = ri(1, 2);
      nama = `${kelamin === 1 ? pick(maleFirst) : pick(femaleFirst)} ${pick(lastNames)}`;
      statKawin = ri(1, 100) <= 60 ? 2 : 4;
    } else {                    // famili lain
      hubKel = 9;
      kelamin = ri(1, 2);
      nama = `${kelamin === 1 ? pick(maleFirst) : pick(femaleFirst)} ${pick(lastNames)}`;
      statKawin = ri(1, 4);
    }
  }

  const nik = randomNIK();
  const tmpLahir = pick(cities);
  let birth;
  if (hubKel === 1) birth = randomBirth(1960, 1990);
  else if (hubKel === 3) birth = randomBirth(1965, 1995);
  else if (hubKel === 4) birth = randomBirth(1995, 2020);
  else if (hubKel === 7) birth = randomBirth(1935, 1965);
  else birth = randomBirth(1950, 2005);
  const tglLahir = fmtDate(birth);

  const wargaNeg = ri(1, 100) <= 99 ? 1 : 2;
  const agama = weighted(agamaWeights);
  const pendidikan = weighted(pendidikanWeights);
  const pekerjaan = (statKawin === 1 && birth.getFullYear() >= 2008) ? 3 : weighted(pekerjaanWeights);
  const ayah = `${pick(maleFirst)} ${pick(lastNames)}`;
  const ibu = `${pick(femaleFirst)} ${pick(lastNames)}`;
  const noHp = randomHP();
  const domisili = ri(1, 100) <= 90 ? 1 : 2;
  const status = ri(1, 100) <= 95 ? 1 : 2;
  const created = fmtDate(new Date(2024, ri(0, 11), ri(1, 28), ri(0, 23), ri(0, 59), ri(0, 59)));

  pendudukRows.push(
    `(${i}, ${kelsId}, ${nik}, ${sqlEsc(nama)}, ${sqlEsc(tmpLahir)}, ${sqlEsc(tglLahir)}, ${kelamin}, ${statKawin}, ${hubKel}, ${wargaNeg}, ${agama}, ${pendidikan}, ${pekerjaan}, ${sqlEsc(ayah)}, ${sqlEsc(ibu)}, ${sqlEsc(kk.kkNama)}, ${sqlEsc(noHp)}, ${domisili}, ${status}, ${USER_ID}, ${sqlEsc(created)}, ${sqlEsc(created)})`
  );
}

// ---------- Compose SQL ----------
const header = `-- =====================================================================
-- Seed data: 500 keluarga + 500 penduduk
-- Database  : db_project_pendataan_penduduk
-- Asumsi    : user dengan id = ${USER_ID} (admin) sudah ada (lihat SeedersUser).
-- Cara pakai:
--   mysql -u root -p db_project_pendataan_penduduk < seed_data.sql
-- =====================================================================

SET FOREIGN_KEY_CHECKS = 0;
SET NAMES utf8mb4;

-- Bersihkan data lama (opsional) ---------------------------------------
-- Catatan: TRUNCATE tidak bisa dipakai karena penduduks memiliki FK ke keluargas
-- (InnoDB menolak TRUNCATE pada tabel yang direferensikan FK).
DELETE FROM penduduks;
DELETE FROM keluargas;
ALTER TABLE penduduks  AUTO_INCREMENT = 1;
ALTER TABLE keluargas AUTO_INCREMENT = 1;

`;

const kelInsert =
  `INSERT INTO keluargas (id, no_kk, kk_nama, alamat, rt, rw, kode_pos, status, user_id, foto_kk, foto_rumah, latitude, longtitude, created_at, updated_at) VALUES\n` +
  keluargaRows.join(',\n') + ';\n\n';

const pendInsert =
  `INSERT INTO penduduks (id, kels_id, nik, nama, tmp_lahir, tgl_lahir, kelamin, stat_kawin, hub_kel, warga_neg, agama, pendidikan, pekerjaan, ayah, ibu, kepala_kel, no_hp, domisili, status, user_id, created_at, updated_at) VALUES\n` +
  pendudukRows.join(',\n') + ';\n\n';

const footer = `SET FOREIGN_KEY_CHECKS = 1;\n`;

const sql = header + kelInsert + pendInsert + footer;
const outPath = path.join(__dirname, 'seed_data.sql');
fs.writeFileSync(outPath, sql, 'utf8');
console.log(`Generated ${keluargaRows.length} keluarga + ${pendudukRows.length} penduduk`);
console.log(`Output: ${outPath}`);
