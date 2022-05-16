<template>
  <v-dialog @keydown.esc="$emit('close')" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Импорт подписчиков</h2>
        <v-spacer></v-spacer>
        <v-btn :block="false" @click="parse" class="mt-0 blue--text text--lighten-4" elevation="0" outlined small>
          <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
          Загрузить
        </v-btn>
        <v-divider
            class="mx-4"
            dark
            vertical
        ></v-divider>
        <v-btn @click="$emit('close')" class="deep-orange--text" elevation="0" outlined small>
          <v-icon class="mr-3">far fa-times</v-icon>
          Закрыть
        </v-btn>
      </v-card-title>
<!--      <input ref="inputUpload" type="file" @change="uploadFile" >-->
      <v-file-input v-model="file" @change="uploadFile" accept=".csv"></v-file-input>

<!--      <v-text-field v-model="inputData" outlined color="success" auto-grow></v-text-field>-->

    </v-card>
  </v-dialog>
</template>

<script>
// const parse = require('@fast-csv/parse');
// const fs = require("fs");
// const fastcsv = require("fast-csv");
 import papaparse from 'papaparse'

export default {
  data: () => ({
    inputData: '',
    file: null,
    // fileName: '',
    // fileUrl: '',
    // fileData: '',
  }),
  props: ['open'],
  methods: {
    uploadFile() {

      var reader = new FileReader();
      reader.readAsText(this.file);

      reader.onload = function() {
        let result = reader.result
        let csvObj = papaparse.parse(result, {});
        this.$emit('import', csvObj.data)
        this.file = null
      }.bind(this);

      reader.onerror = function() {
        console.log(reader.error);
        this.$emit('close')
        this.file = null
      }.bind(this);


      /*reader.readAsText(this.file)
      reader.addEventListener('load', () => {
        let result = reader.result
        let csvObj = papaparse.parse(result, {});
        this.$emit('import', csvObj.data)
        this.file = null
        // console.log(csvObj.data)
      })*/



      /*fr.readAsDataURL(this.file)
      fr.addEventListener('load', () => {
        // this.fileUrl = fr.result
        // this.fileData = files[0] // this is an image file that can be sent to server...
        console.log()
      })*/




    },
    parse: async function () {
      console.log('parse')


      // this.$emit('upload', email_box)
    }
  },

}
</script>
