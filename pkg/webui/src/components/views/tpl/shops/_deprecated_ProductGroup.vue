<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items.sync="items"
            :loading="loading"
            :options.sync="options"
            :server-items-length="totalItems"
            :single-expand="false"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
            class="elevation-1"
            show-expand
            expand-icon="fal fa-angle-down"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Категории товаров</h1></v-toolbar-title>

          <v-select
                  :items="webSites"
                  item-value="id"
                  item-text="name"
                  return-object
                  outlined
                  dense
                  :value="currentWebSiteId"
                  @change="switchCurrentShop"
                  style="max-width: 380px"
                  class="ml-8 mt-7"
          ></v-select>

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
              >far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>

            <v-btn
                    @click="createItem"
                    color="cyan darken-4"
                    elevation="1"
                    small
                    dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-plus</v-icon>
              Добавить группу
            </v-btn>

          </v-card>

        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по данным"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.parent="{ item }">
        <span v-text="item['parent_id'] ? item['parent_id'] : 'root'" :class="{ 'red--text': !item['parent_id'] }"></span>
      </template>

      <template v-slot:item.shop="{ item }">
        {{ getFieldFromShops(item, 'name') }}
      </template>

      <template v-slot:item.actions="{ item }">

        <v-btn
                small depressed style="background-color: transparent"
                @click="openDeleteDialog(item)">
          <v-icon
                  small
                  class="red--text text--lighten-2"
          >
            fas fa-trash
          </v-icon>
        </v-btn>

      </template>

      <template v-slot:expanded-item="{ headers, item }" v-if="currentWebSiteId">
        <td :colspan="headers.length">
          <h3 class="mb-8">Редактируемые данные</h3>
          <v-row dense class="mt-4">
            <v-col cols="3">

              <v-text-field
                      label="Id группы" dense class="body-1"
                      v-model="item['id']"
                      disabled
              />

              <v-text-field
                      label="Название группы" dense class="body-1"
                      v-model.trim="item['name']"
              />
              <v-text-field
                      label="icon_name" dense class="body-1"
                      v-model.trim="item['icon_name']"
              />
              <v-text-field
                      label="route_name" dense class="body-1"
                      v-model.trim="item['route_name']"
              />

              <v-select
                      :items="items"
                      item-text="name"
                      item-value="id"
                      clearable
                      v-model="item['parent_id']"
                      label="Родитель"
              >
                <template slot="selection" slot-scope="data">
                  {{ data.item.id }} - {{ data.item.name }}
                </template>

                <template slot="item" slot-scope="data">
                  {{ data.item.id }} - {{ data.item.name }}
                </template>

              </v-select>

              <v-select
                      :items="webSites"
                      item-value="name"
                      item-text="name"
                      persistent-hint
                      return-object
                      :value="webSites.find(el=>el.id ===item['web_site_id'])['name']"
                      @change="$store.commit('shops/GROUP_UPDATE_STATE', {item, field: 'web_site_id', value: $event['id']})"

                      label="Магазин"
              ></v-select>

            </v-col>
            <v-col cols="3" >
              <v-text-field
                      label="Code" dense class="body-1"
                      v-model.trim="item['code']"
              />
              <v-text-field
                      label="URL" dense class="body-1"
                      v-model.trim="item['url']"
              />
              <v-text-field
                      label="Breadcrumb" dense class="body-1"
                      v-model.trim="item['breadcrumb']"
              />

            </v-col>
            <v-col cols="6" >
              <v-text-field
                      label="Meta Title" dense class="body-1"
                      v-model.trim="item['meta_title']"
              />
              <v-text-field
                      label="Meta Keywords" dense class="body-1"
                      v-model.trim="item['meta_keywords']"
              />
              <v-textarea
                      label="Meta description" dense class="body-1"
                      rows="3"
                      v-model.trim="item['meta_description']"
              />
              <v-textarea
                      label="Short Description" dense class="body-1"
                      rows="3"
                      v-model.trim="item['short_description']"
              />

            </v-col>

            <v-col cols="12">
              <h3 class="mb-3">Description</h3>
              <prism-editor v-model="item['description']" class="myPrismEditor mb-2 pb-2" language="html" max-height></prism-editor>
            </v-col>

          </v-row>

          <v-row class="mt-4">
            <v-col cols="12">
              <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
                <v-btn @click="openDeleteDialog(item)"
                       color="red darken-3"
                       elevation="1"
                       small
                       dark
                >
                  <v-icon
                          x-small
                          class="ml-1 mr-3"
                  >far fa-user-slash</v-icon>
                  Удалить
                </v-btn>

                <v-divider
                        class="mx-2 mt-0"
                        inset
                        vertical
                ></v-divider>

                <v-btn
                        @click="updateItemData(item)"
                        dark
                        color="teal darken-3"
                        elevation="1"
                        small
                >
                  <v-icon
                          x-small
                          class="ml-1 mr-3"
                  >far fa-save</v-icon>
                  Сохранить
                </v-btn>
              </v-card>
            </v-col>
          </v-row>
        </td>
      </template>

    </v-data-table>
    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialog.open"
            @keydown.esc="dialog.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name">
              <v-list-item-title class="headline">Имя: &laquo;{{dialog.name}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.url">
              <v-list-item-title class="headline">URL: &laquo;{{dialog.url}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.web_site_id">
              <v-list-item-title class="headline">Shop: &laquo;{{ getFieldFromShops({web_site_id: dialog.web_site_id}, 'name') }}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
                  color="green darken-1"
                  small
                  text
                  @click.stop="dialog.open = false"
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
                  color="red darken-1"
                  small
                  text
                  @click.stop="deleteItemDialog"
          >
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
        { text: 'Id', align: 'start', value: 'id'},
        { text: 'Parent', value: 'parent'},
        { text: 'Имя', value: 'name'},
        { text: 'URL', value: 'url'},
        { text: 'Code', value: 'code'},
        // { text: 'Магазин', value: 'shop'},
        { text: 'Удалить', value: 'actions'},
      ],
      items: [],
      webSites: [],
      totalItems: 0,
      options: {},
      loading: true,
      search: '',
      dialog: {
        id: null,
        web_site_id: null,
        name: '',
        url: '',
        open: false,
      },
      currentWebSiteId: null,
    }),
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      },
    },
    async mounted() {
      this.loading = true
      await this.loadWebSites(false)
              .then(async ()=>{

                this.currentWebSiteId = this.webSites.length > 0 ? this.webSites[this.webSites.length-1].id : null;
                if (this.currentWebSiteId !== null) {
                  // await this.getDataFromApi(false)
                }

              }).finally(()=>this.loading = false)


    },
    methods: {
      getDataFromApi: async function(showNotification, searchStr = "") {

        if (this.currentWebSiteId < 1) {
          if (showNotification) {
            this.$notify({
              group: 'main',
              title: 'Список магазинов пуст',
              type: 'warn'
            });
          }
          return null
        }

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
          web_site_id: this.currentWebSiteId,
        }

        await this.$api.productGroup.getListPagination(payload)
                .then(resp => {
                  if (resp['productGroups'] !== undefined ) {
                    this.items = resp['productGroups']
                    this.totalItems = resp.total
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        type: 'success'
                      });
                    }
                  } else {
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Ошибка обновления',
                        text: 'Ошибка в ответе сервера: productGroups - not found',
                        type: 'warring'
                      });
                    }
                  }
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(false, this.search)
      }, 1400),
      updateItemData: async function(item) {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: item.web_site_id,
          id: item.id,
        }
        Object.assign(payload, item)

        await this.$store.dispatch('shops/UPDATE_GROUP', payload)
                .then((resp) => {

                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Группа обновлена',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: 'Не удалось обновить данные',
                    type: 'error'
                  });
                })
      },
      loadWebSites: async function(showNotification) {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
        }

        return this.$api.webSite.getList(payload)
                .then(resp => {
                  if (resp['webSites'] !== undefined ) {
                    this.webSites = resp['webSites']
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        type: 'success'
                      });
                    }
                  } else {
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Ошибка обновления',
                        text: 'Ошибка в ответе сервера: webSites - not found',
                        type: 'warring'
                      });
                    }
                  }
                })
                .catch((err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      },

      createItem: async function (){
        if (this.currentWebSiteId < 1) {
          this.$notify({
            group: 'main',
            title: 'Список магазинов пуст',
            type: 'warn'
          });
          return
        }

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.currentWebSiteId, // для url и для БД
          name: "Новая группа",
          parent_id: null,
          code: '',
          url:  '',
        }
        await this.$api.productGroup.create(payload)
                .then(async (data) => {
                  if (data['productGroup']) this.items.push(data['productGroup'])
                  await this.getDataFromApi(false)
                  this.$notify({
                    group: 'main',
                    title: 'Группа создана',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка создания',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.web_site_id = item['web_site_id']
        this.dialog.name = item['name']
        this.dialog.url = item['url']
        this.dialog.open = true
      },
      deleteItemDialog() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.dialog.web_site_id,
          id: this.dialog.id
        }
        this.$store.dispatch('shops/DELETE_GROUP', payload)
                .then(async ()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'Группа удалена',
                    type: 'success'
                  });
                })
                .catch((err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка удаления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>{
                  this.dialog.id = 0
                  this.dialog.name = ''
                  this.dialog.token = ''
                })

      },
      switchCurrentShop(shop) {
        this.currentWebSiteId = shop.id
        this.getDataFromApi(true)
      },
      getGroupById(id) {
        return  this.items.find(item => item.id === id, id)
      }
    },
  }
</script>