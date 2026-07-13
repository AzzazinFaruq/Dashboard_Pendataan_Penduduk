<template>
  <div class="pa-5">
    <v-card class="mt-5 pa-5" elevation="0">
      <!-- Bagian Foto Profil -->
      <div class="d-flex justify-center">
        <div class="rounded-circle" style="width: 150px; height: 150px;">
          <v-img v-if="!userdata.profile_picture" src="../assets/avatar.png" width="150" height="150"></v-img>
          <v-img :src="foto" width="150" height="150"></v-img>
        </div>
      </div>

      <!-- Riwayat Login -->
      <div class="mt-8">
        <v-row class="mt-2">
          <v-col cols="12" md="9" class="pr-md-2">
            <v-card variant="outlined" class="pa-4 d-flex align-center justify-center" style="border-color: #e0e0e0;">
              <v-img src="../assets/logo-pendataan-black.png" max-height="40" max-width="120" contain class="mr-4"></v-img>
            </v-card>
          </v-col>
          <v-col cols="12" md="3" class="pl-md-2">
            <v-card variant="outlined" class="pa-4 d-flex align-center justify-center" height="100%" style="border-color: #e0e0e0;">
              <v-icon color="#4285F4" size="x-large" class="mr-2">mdi-google</v-icon>
              <span class="text-subtitle-2">Google</span>
            </v-card>
          </v-col>
        </v-row>
      </div>

      <!-- Form Email dan Password -->
      <div class="mt-8">
        <v-row>
          <v-col cols="12" md="6">
            <label>Email</label>
            <v-text-field
              class="mt-2"
              v-model="userdata.email"
              variant="outlined"
              density="comfortable"
              @change="saveChanges"
            ></v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <label>Kata Sandi</label>
            <v-text-field
              class="mt-2"
              v-model="password"
              variant="outlined"
              density="comfortable"
              :type="showPassword ? 'text' : 'password'"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="showPassword = !showPassword"
              @change="saveChanges"
            ></v-text-field>
          </v-col>
        </v-row>
      </div>

      <!-- Tombol Simpan -->
      <div class="d-flex justify-end mt-4">
        <v-btn 
          color="#01A65D" 
          prepend-icon="mdi-content-save" 
          class="text-white"
          @click="saveChanges"
          :loading="isSaving"
          height="48"
          width="150"
        >
          Simpan
        </v-btn>
      </div>
    </v-card>
  </div>
</template>

<script>
import axios from 'axios';
import { test } from '@/stores/restrict';
const use = test();

export default {
  setup() {
    use.setup();
  },
  data() {
    return {
      foto: "",
      userdata: { email: '' },
      password: "admin123",
      showPassword: false,
      isSaving: false
    }
  },
  mounted() {
    this.user();
  },
  methods: {
    user() {
      try {
        axios.get('/api/user')
          .then((res) => {
            const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
            this.foto = apiUrl + "/" + res.data.data.profile_picture;
            this.userdata = res.data.data;
          })
      } catch {
        // Error handling
      }
    },
    async saveChanges() {
      this.isSaving = true;
      try {
        // Simpan perubahan email
        await axios.post('/api/update-email', { email: this.userdata.email });
        
        // Simpan perubahan password jika diubah
        if (this.password !== "32010101012345") {
          await axios.post('/api/update-password', { password: this.password });
        }
        
        // Notifikasi sukses
        this.$toast.success('Perubahan berhasil disimpan');
      } catch (error) {
        console.error('Error saving changes:', error);
        this.$toast.error('Gagal menyimpan perubahan');
      } finally {
        this.isSaving = false;
      }
    }
  }
}
</script>

<style scoped>
label {
  font-size: 16px;
  font-weight: 600;
}

.rounded-circle {
  overflow: hidden;
  transition: opacity 0.3s;
  cursor: pointer;
}

.rounded-circle:hover {
  opacity: 0.8;
}

.v-card {
  border-radius: 8px;
}

.text-h6 {
  font-weight: 600;
  color: #333;
}

.pr-md-2 {
  padding-right: 8px !important;
}

.pl-md-2 {
  padding-left: 8px !important;
}
</style>