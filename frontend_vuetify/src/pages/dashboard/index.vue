<template>
  <v-container fluid class="">
  <div class="pa-5 mr-2">
    <div class="mt-5">
    <v-row class="">
      <v-col cols="12" md="6" lg="3">
     <card
      bgColor="#F1EDF5"
      icon-color="#6E4D99"
      icon="mdi-account-multiple"
      title="Total Keluarga"
      :count="datacardjumlah.keluargaCount"
      />
    </v-col>
      <v-col cols="12" md="6" lg="3">

      <card
      bgColor="#E9F3FB"
      icon-color="#2184D8"
      icon="mdi-account-multiple"
      title="Total Penduduk"
      :count="datacardjumlah.pendudukCount"
      />
    </v-col>
      <v-col cols="12" md="6" lg="3">
      <card
      bgColor="#E6F6EF"
      icon-color="#01A65D"
      icon="mdi-check-bold"
      title="Penduduk Aktif"
      :count="datacardstatus.alivepend"
      />
    </v-col>
      <v-col cols="12" md="6" lg="3">
      <card
      bgColor="#FEF0EF"
      icon-color="#F76C5E"
      icon="mdi-minus-circle"
      title="Penduduk Inaktif"
      :count="datacardstatus.nopen"
      />
    </v-col>
    </v-row>
    </div>
    <div class="mt-5">
        <v-row elevation="5">
          <v-col cols="12" md="8">
          <h3  class="mb-5">Status Penduduk</h3>
          <v-card class="pa-5" elevation="0" height="90%" rounded="lg">
          <marrychart/>
          </v-card>
          </v-col>
          <v-col cols="12" md="4">
          <h3 class="mb-5">Jenis Kelamin</h3>
          <v-card class="pa-5" elevation="0" height="90%" rounded="lg">
          <genderchart/>
          </v-card>
          </v-col>
      </v-row>
    </div>
    <div class="mt-5" >
    <h3  class="mb-5">Riwayat Input</h3>
    <v-card elevation="0" class="pa-5">
        <totalchart/>
    </v-card>

    </div>
    <div class="mt-10">
        <dtTable
        :dthead="headpen"
        :dtbody="dtpen"
        :toCapitalize="toCapitalize"
        />
    </div>
  </div>
  </v-container>
</template>

<script>
import genderchart from '@/components/chart/genderchart.vue';
import marrychart from '@/components/chart/marrychart.vue';
import totalchart from '@/components/chart/totalchart.vue';
import axios from 'axios'
import card from '@/components/card.vue';
import dtTable from '@/components/dtTable.vue';
import { test } from '@/stores/restrict';
import { useTitle } from '@/stores/title'
const use = test();

export default {
  setup(){
    use.setup();
    const title = useTitle()
    title.setTitle('Dashboard')  // Set title halaman
  },
  components:{
  },
  data() {
    return {
      loginPoUp:false,
      headpen:[
        {id:0, name:'Nomor KK / Kepala Keluarga'},
        {id:1, name:'Nomor NIK'},
        {id:2, name:'Nama'},
        {id:3, name:'Status'},
        {id:4, name:'User'}
      ],

      datacardjumlah:[],
      datacardstatus:[],
      tab: 'one',
      data: [],
      level:"",
      dtkel:[],
      dtpen:[]
    };
  },
  mounted() {
    this.fetchData();
    this.popup();
  },
  methods: {
    popup(){
      if (localStorage.getItem('loginAlert')==='true') {
        this.loginPoUp=true;
        localStorage.removeItem('loginAlert');
      }
    },
    async fetchData() {
        this.status();
        this.penlast();
        this.dataForCard();
    },
    dataForCard(){
      axios.get('/api/alldata')
      .then((res)=>{
        this.datacardjumlah = res.data.jumlah;
        this.datacardstatus = res.data.status;

      })
    },
    status() {
      // try {
        axios.get("/api/user").then((res) => {
          this.data = res.data.data;
          const data = res.data.data.level
          if (data == 'admin' || data == 'superAdmin') {
            this.level = 'Admin';
          }
          else{
            this.level = '';
          }
        });
    },
    penlast(){
      axios.get('/api/latestpend')
        .then((res)=>{
          this.dtpen = res.data.data.map(item => ({
            ...item,
            nama: this.toCapitalize(item.nama),
            kepala: this.toCapitalize(item.kepala),
            status: this.toCapitalize(item.status),
            user: this.toCapitalize(item.user)
          }));
        })
    },
    toCapitalize(text) {
      if (!text) return '';
      return text.toString().toLowerCase().split(' ').map(word => {
        return word.charAt(0).toUpperCase() + word.slice(1);
      }).join(' ');
    },
  },
  created() {

  },
};

</script>

<style>

</style>
