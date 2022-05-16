<template>

  <v-container fluid>

    <v-data-table
        :headers="headers"
        :items.sync="items"
        :loading="loading"
        :options.sync="options"
        :search="search"
        :items-per-page="15"
        :server-items-length="totalItems"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
        class="elevation-1"
    >
      <template v-slot:top>
        <v-toolbar color="white" flat>
          <v-toolbar-title><h1 style="font-size: 26px;">Страницы сайта</h1></v-toolbar-title>
          <v-select
              :items="webSites"
              item-value="id"
              item-text="name"
              return-object
              outlined
              dense
              :value="webSiteId"
              @change="switchCurrentShop"
              style="max-width: 380px"
              class="ml-8 mt-7"
          ></v-select>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" color="cyan darken-4" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-external-link</v-icon>
              Добавить страницу
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col class="ml-2" cols="6">
            <v-text-field
                class="mx-2"
                hide-details
                label="Поиск по страницам"
                prepend-inner-icon="fal fa-search"
                single-line
                v-model.trim="search">
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.actions="{ item }">
        <show-action :route="{name:'web-pages.edit', params: {public_id:item['public_id']}}" @delete="openDeleteDialog(item)"></show-action>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog @keydown.esc="dialog.open = false" max-width="380" v-model="dialog.open">
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id №</strong> {{ dialog.id }}</span>
          <span><strong>Имя: </strong>{{ dialog.label }}</span>
        </v-card-text>

        <!-- Удалить / Отменить -->
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click.stop="dialog.open = false" color="green darken-1" small text>
            Отменить
          </v-btn>
          <v-btn @click.stop="deleteItemDialog" color="red darken-1" small text>
            Удалить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-container>
</template>

<script>
export default {
  data: () => ({
    headers: [
      {text: 'ID', align: 'start', value: 'public_id'},
      {text: 'Label', align: 'start', value: 'label'},
      {text: 'Meta title', align: 'start', value: 'meta_title'},
      {text: 'Code', align: 'start', value: 'code'},
      {text: 'Route name', align: 'start', value: 'route_name'},
      {text: 'Path', align: 'start', value: 'path'},
      {text: 'Категорий', align: 'end', value: '_product_categories_count'},
      { text: 'Действие', align: 'center', value: 'actions', sortable: false},
    ],
    open: ['root'],
    items: [],
    tree: [],
    totalItems: 0,
    options:{
      itemsPerPage: 15,
    },
    search: '',
    loading: true,
    dialog: {
      open: false,
    },
    webSites: [],
    webSiteId: null,
  }),
  async mounted() {
    this.loading = true

    await Promise.all([
      this.loadWebSites(false).then(()=>{
        this.webSiteId = this.webSites.length > 0 ? _.minBy(this.webSites,'id').id : null;
      }).then(() => this.getDataFromApi(false)),
    ])

  },
  watch: {
    search() {
      this.options.page = 1;
      this.searchDataFromApi();
    },
  },
  methods: {
    getDataFromApi: async function (showNotification) {
      if (this.loading) return
      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        webSiteId: this.webSiteId,
      }

      this.$api.webPage.getListPagination(payload)
          .then(resp => {
            if (resp['web_pages'] !== undefined) {
              this.items = resp['web_pages']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Ошибка в ответе сервера: web_pages - not found',
                  type: 'warring',
                });
              }
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },

    searchDataFromApi: _.throttle(function () {
      if (this.loading) return
      this.getDataFromApi(false, this.search)
    }, 1400),
    updateLists: async function (showNotification) {
      await this.getDataFromApi(showNotification)
    },
    goToCreateItem: async function () {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        label: "Новая страница",
        web_site_id: this.webSiteId,
      }

      await this.$api.webPage.create(payload)
          .then((resp) => {
            if (resp['web_page'] !== undefined) {
              this.items.unshift(resp['web_page'])
              this.totalItems += 1
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [web_page]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
    },

    // ниже переписать в общее диалоговое оконо
    openDeleteDialog(item) {
      Object.assign(this.dialog, item)
      this.dialog.open = true
    },
    deleteItemDialog: async function () {
      this.dialog.open = false

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialog.id,
      }
      await this.$api.webPage.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Страница удалена',
              type: 'success',
            });
            Object.assign(this.dialog, {open: false})
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error',
            });
          })
    },
    switchCurrentShop(shop) {
      this.webSiteId = shop.id
      this.getDataFromApi(true)
    },
  },
  components: {
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
  },
}
</script>
<style>
/*@import "~@syncfusion/ej2-base/styles/material.css";
@import "~@syncfusion/ej2-vue-navigations/styles/material.css";
@import "~@syncfusion/ej2-inputs/styles/material.css";
@import "~@syncfusion/ej2-buttons/styles/material.css";
.control_wrapper {
  display: block;
  max-width: 350px;
  max-height: 350px;
  margin: auto;
  overflow: auto;
  border: 1px solid #dddddd;
  border-radius: 3px;
}*/
</style>