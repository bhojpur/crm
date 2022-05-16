<template>
  <v-container fluid>
    <v-data-table
        :headers="headers"
        :items.sync="items"
        :options.sync="options"
        :server-items-length="totalItems"
        :loading="loading"
        class="elevation-1"
        :items-per-page="15"
        :search="search"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
        show-expand
        expand-icon="fal fa-angle-down"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>
            <h1 style="font-size: 26px;">Публичные файлы</h1>
            <v-progress-linear
                height="8"
                :title="valueSpaceDiskString"
                :value="valueSpaceDisk"
                striped
                color="light-green darken-4"
            >
            </v-progress-linear>

          </v-toolbar-title>

          <v-spacer></v-spacer>

          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                  x-small
                  class="ml-1 mr-3"
              >far fa-sync
              </v-icon>
              Обновить
            </v-btn>

            <v-divider
                class="mx-2 mt-0"
                inset
                vertical
            ></v-divider>

            <v-btn @click="dialogAddFile.open = true"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                  x-small
                  class="ml-1 mr-3"
              >far fa-plus
              </v-icon>
              Добавить
            </v-btn>

          </v-card>

        </v-toolbar>
      </template>

      <template v-slot:item.name="{ item }">
        <a
            :href="item['_url']"
            target="_blank"
            style="text-decoration: none;">
          <v-icon
              small
              class="mr-2 text--primary"
          >
            fal fa-link
          </v-icon>
          {{ item['name'] }}</a>
      </template>

      <template v-slot:item.size="{ item }">
        <span class="grey--text text--darken-2">{{ formatBytes(item['size']) }}</span>
      </template>

      <template v-slot:item.mime="{ item }">
        <span class="grey--text text--darken-2">{{ (item['mime']) }}</span>
      </template>

      <template v-slot:item.enabled="{ item }">
        <span v-text="item['enabled'] ? 'виден' : 'скрыт'"></span>
      </template>

      <template v-slot:item.type="{ item }">
        <span v-text="item['product_id'] > 0 ? 'товар' : item['emailId'] > 0 ? 'email' : '-'"></span>
      </template>

      <template v-slot:item.created_at="{ item }">
        <span>{{ parseDate(item['created_at']) }}</span>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-icon
            small
            class="red--text text--lighten-2"
            @click="openDeleteDialog(item)"
        >
          fas fa-trash
        </v-icon>
      </template>

      <template v-slot:expanded-item="{ headers, item }">
        <image-expand-template :file.sync="item" :headersLength="headers.length" @removeItem="removeItemFromItems"></image-expand-template>
      </template>

    </v-data-table>

    <!-- Диалог создания файла -->
    <v-dialog
        v-model="dialogAddFile.open"
        @keydown.esc="dialogAddFile.open = false"
        max-width="340"
    >
      <v-card>
        <v-card-title class="headline">Выберите файл (max 32mb)</v-card-title>

        <v-card-text>
          <v-file-input show-size label="Файла для загрузки" v-model="uploadedFile"></v-file-input>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- закрыть -->
          <v-btn
              color="primary"
              class="mr-2"
              small
              text
              outlined
              @click.stop="dialogAddFile.open = false, file=null"
          >
            Закрыть
          </v-btn>

          <!-- Загрузить -->
          <v-btn
              color="green"
              small
              text
              outlined
              @click.stop="createItem"
          >
            Загрузить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Подтверждение удаления -->
    <v-dialog
        v-model="dialogDelete.open"
        @keydown.esc="dialogDelete.open = false"
        max-width="290"
    >
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id</strong>: {{ dialogDelete.id }}</span>
          <span><strong>Имя: </strong>{{ dialogDelete.name }}</span>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
              color="primary"
              class="mr-2"
              small
              text
              outlined
              @click.stop="dialogDelete.open = false"
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
              color="red darken-1"
              small
              text
              outlined
              @click.stop="deleteFileDialog"
          >
            Удалить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-container>
</template>

<script>

import axios from 'axios'
import Vue from "vue";

export default {
  data: () => ({
    headers: [
      {text: 'id', align: 'start', value: 'id'},
      {text: 'Имя файла', align: 'start', value: 'name'},
      {text: 'Размер', align: 'start', value: 'size'},
      {text: 'MIME-тип', value: 'mime'},
      // { text: 'Отображение (#)', value: 'public' },
      {text: 'Назначение', value: 'type'},
      {text: 'Видимость', value: 'enabled'},
      {text: 'Дата загрузки', value: 'created_at'},
      // { text: 'Дата редактирования', value: 'updated_at'},
      {text: 'Actions', value: 'actions'},
    ],
    dialogDelete: {
      open: false,
      id: '',
      name: '',
    },
    dialogAddFile: {
      open: false,
      file: null,
    },
    uploadedFile: null,
    options: {
      itemsPerPage: 15,
    },
    items: [],
    totalItems: 0,
    search: '',
    loading: false,
    diskSpaceUsed: 0,
  }),
  watch: {},
  async mounted() {
    await this.getDataFromApi(false)
  },

  methods: {
    getDataFromApi: async function (shopNotification) {

      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        // preloads: 'ProductCards,PaymentSubject,MeasurementUnit,ProductTags'
      }

      await this.$api.storage.getListPagination(payload)
          .then((resp) => {
            if (resp['files'] !== undefined && resp['files'] !== null) {
              this.totalItems = parseInt(resp['total'], 10)
              this.diskSpaceUsed = resp['disk_space_used']
              this.items = resp['files']
              if (shopNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Список файлов обновлен',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Не полные данные от сервера',
                type: 'warring',
              });
            }

          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },

    searchDataFromApi: _.throttle(function () {
      this.getDataFromApi(false, this.search)
    }, 1400),
    openDeleteDialog(item) {

      this.dialogDelete.id = item['id']
      this.dialogDelete.name = item['name']
      this.dialogDelete.open = true
    },
    removeItemFromItems(payload) {
      let index = this.items.findIndex(item => item.id === payload.id, payload)
      if (index !== -1) this.$delete(this.items, index)
    },
    async deleteFileDialog() {
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialogDelete.id,
      }
      await this.$api.storage.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) Vue.delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Файл удален',
              type: 'success',
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogDelete.open = false
            this.dialogDelete.hash_id = ''
            this.dialogDelete.name = ''
          });


      // обновляем текущий список
      await this.getDataFromApi(false)

    },
    createItem() {

      this.uploadFile()
          .then((resp) => {

            if (resp['file'] !== undefined && resp['file'] !== null) {
              this.$notify({
                group: 'main',
                title: 'Файл создан',
                type: 'success',
              });
              this.totalItems++
              this.diskSpaceUsed = resp['disk_space_used']
              this.items.unshift(resp['file'])
            } else {
              this.$notify({
                group: 'main',
                title: 'Не полные данные от сервера',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка при создании файла',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogAddFile.open = false
            this.uploadedFile = null
          });


    },

    /////////////////
  },
  computed: {
    valueSpaceDiskString() {
      return this.formatBytes(this.diskSpaceUsed) + ' / ' + this.formatBytes(this.$store.state.account.disk_space_available)
      // return this.formatBytes(this.diskSpaceUsed) + ' / ' + this.formatBytes(524288000)
    },
    valueSpaceDisk() {
      return (this.diskSpaceUsed / parseInt(this.$store.state.account.disk_space_available)) * 100
    },
  },
  components: {
    ImageExpandTemplate: () => import('@/components/views/tpl/content/layouts/ImagesExpandTemplate'),
  },


}
</script>