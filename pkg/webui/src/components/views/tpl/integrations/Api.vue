<template>
  <v-container fluid>
    <v-data-table
        :headers="headers"
        :items.sync="items"
        :options.sync="options"
        :server-items-length="totalItems"
        :loading="loading"
        class="elevation-1"
        :search="search"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Api keys</h1></v-toolbar-title>

          <!-- Статус API -->
          <v-card class="d-flex flex-column ml-8 justify-center align-center px-6 py-1" flat style="border: 1px solid #ccc">
            <h4>Статус API</h4>
            <v-switch
                    dense
                    hide-details
                    :loading="loadingApiEnabled"
                    :success="$store.state.account.api_enabled"
                    :error="!$store.state.account.api_enabled"
                    :input-value="$store.state.account.api_enabled"
                    @change="updateApiStatus"
            ></v-switch>

          </v-card>

          <div class="ml-6 d-flex flex-column justify-center pa-3">
            <span style="font-size: 13px;" class="mb-1">Account hash id</span>
            <span style="font-weight: 500;">{{$store.state.account.hash_id}}</span>
          </div>

          <div class="ml-6 d-flex flex-column justify-center pa-3">
            <a href="https://api.crm.bhojpur.net" target="_blank" style="font-size: 13px;" class="mb-1">https://api.crm.bhojpur.net</a>
          </div>

          <v-spacer></v-spacer>

          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi"
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
                    dark
                    color="cyan darken-4"
                    elevation="1"
                    small
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-plus</v-icon>
              Добавить
            </v-btn>

          </v-card>

        </v-toolbar>
      </template>

      <template v-slot:item.enabled="{ item }">
        <v-switch
                dense
                :input-value.sync="item['enabled']"
                @change="updateItem(item, $event, item['name'])"
        ></v-switch>
      </template>

      <template v-slot:item.name="{ item }">
        <v-text-field
                dense
                hide-details
                flat
                single-line
                @change="updateItem(item, item['enabled'], $event)"
                @keydown.enter.stop="updateItem(item, item['enabled'], $event)"
                class="body-1"
                :value="item['name']"
        ></v-text-field>
      </template>

      <template v-slot:item.token="{ item }">
        <v-row>
          <span style="width: 230px;font-size: 13px;">
            {{item['token']}}
          </span>
          <v-btn @click="copyToken(item['token'])" small depressed style="background-color: transparent">
            <v-icon class="pr-1" >fal fa-copy</v-icon>
          </v-btn>
        </v-row>

      </template>

      <template v-slot:item.created_at="{ item }">
        {{parseDate(item['created_at'])}}
      </template>

      <template v-slot:item.updated_at="{ item }">
        {{parseDate(item['updated_at'])}}
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
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-list-item-title class="headline">Имя: &laquo;{{dialog.name}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0">
              <v-list-item-title class="headline">Ключ: &laquo;{{dialog.token}}&raquo;</v-list-item-title>
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
        { text: 'Дата создания', value: 'created_at'},
        { text: 'Имя ключа', align: 'start', value: 'name'},
        { text: 'Ключ', align: 'start', value: 'token' },
        { text: 'Статус ключа', value: 'enabled' },
        { text: 'Удалить', value: 'actions'},
      ],
      items: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        open: false,
      },
      loadingApiEnabled: false,
    }),
    async mounted() {
      // await this.getDataFromApi(false)
    },

    methods: {
      getDataFromApi: async function(showNotification, searchStr = '') {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
        }

        this.$api.apiKey.getListPagination(payload)
            .then(resp => {
              if (resp['api_keys'] !== undefined ) {
                this.items = resp['api_keys']
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
                    text: 'Ошибка в ответе сервера: api_keys - not found',
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
            .finally(() => this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(false, this.search)
      }, 1400),
      updateLists: async function (showNotification) {
        await this.getDataFromApi(showNotification)
      },
      updateItem: async function (item, status, name) {
        if (name['type'] !==undefined && name['type']=== 'keydown') {
          return
        }
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: item['id'],
          name: name,
          enabled: status,
        }
        await this.$api.apiKey.update(payload)
                .then((data) => {
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Данные Api-ключа успешно обновлены',
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
      updateApiStatus: async function (status) {
        this.loadingApiEnabled = true
        const payload = {
          hash_id: this.$store.state.account.hash_id,
          api_enabled: status,
        }

        await this.$api.apiKey.update(payload)
                .then((data) => {
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Статус API обновлен',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: 'Не удалось обновить статус API',
                    type: 'error'
                  });
                })
        .finally(()=> this.loadingApiEnabled = false)


      },
      copyToken: function(token) {
        navigator.clipboard.writeText(token).then(()=> console.log("Скопировано в буфер"))
        this.$notify({
          group: 'main',
          title: 'Скопирован в буфер',
          text: token,
          type: 'success'
        });
      },

      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "New Api key",
          enabled: true
        }
        await this.$api.apiKey.create(payload)
                .then((resp) => {
                  if (resp['api_key']) {
                    this.items.unshift(resp['api_key'])
                    this.$notify({
                      group: 'main',
                      title: 'ApiKey создан',
                      // text: 'Данные Api-ключа успешно обновлены',
                      type: 'success'
                    });
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Получены не вс  е данные',
                      // text: 'Данные Api-ключа успешно обновлены',
                      type: 'warring'
                    });
                  }

                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка создания',
                    text: 'Не удалось создать APiKey',
                    type: 'error'
                  });
                })
      },

      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {

        this.dialog.id = item['id']
        this.dialog.name = item['name']
        this.dialog.token = item['token']
        this.dialog.open = true
      },
      deleteItemDialog() {
        this.dialog.open = false

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }

        this.$api.apiKey.delete(payload)
                .then(()=> {
                  let index = this.items.findIndex(el => el.id === this.dialog.id)
                  if (index !== -1 ) {
                    this.$delete(this.items, index)
                  }
                  console.log("index: ", index)
                  this.$notify({
                    group: 'main',
                    title: 'Ключ удален',
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

    },

    computed: {
     /* api_enabled() {
        // console.log(this.$store.state.account['api_enabled'])
        return true
      }*/
    }

  }
</script>