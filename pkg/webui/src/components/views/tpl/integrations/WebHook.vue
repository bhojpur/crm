<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items="items"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
            :search="search"
            expand-icon="fal fa-angle-down"
            :single-expand="false"
            show-expand
            @update:options="getDataFromApi(false)"
            @update:search="searchDataFromApi"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">WebHooks</h1></v-toolbar-title>

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
              Добавить webHook
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

      <template v-slot:item.enabled="{ item }">
        <span v-text="item['enabled'] ? 'в работе' : 'остановлен'"></span>
      </template>

      <template v-slot:item.actions="{ item }">

        <v-btn
                small depressed style="background-color: transparent"
                @click="execute(item)">
          <v-icon
                  small
                  class="green--text text--lighten-2"
          >
            fas fa-play
          </v-icon>
        </v-btn>

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

      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-toolbar flat>
            <v-toolbar-title><h4 class="mt-4 mb-6" style="font-weight: 400;">Параметры WebHook</h4></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
              <v-btn @click="openDeleteDialog(item)" color="red lighten-1" elevation="0" small dark outlined>
                <v-icon x-small class="ml-1 mr-3">far fa-user-slash</v-icon>
                Удалить
              </v-btn>
              <v-divider class="mx-2 mt-0" inset vertical></v-divider>
              <v-btn @click="updateItemData(item)" dark color="teal darken-1" elevation="0" small outlined>
                <v-icon x-small class="ml-1 mr-3">far fa-save</v-icon>
                Сохранить
              </v-btn>
            </v-card>
          </v-toolbar>
          <v-row>
           <v-col cols="4">
              <v-text-field
                      label="Id карточки" dense class="body-1"
                      :value="item.id"
                      disabled
              />
              <v-switch
                      class="mb-6 mt-0"
                      dense
                      :label="item['enabled'] ? 'В работе' : 'Остановлен'"
                      hide-details
                      light
                      :success="item['enabled']"

                      v-model="item['enabled']"
              ></v-switch>

             <v-text-field
                     label="Name"  class="body-1"
                     v-model.trim="item['name']"
             />
             <v-text-field
                     label="Code"  class="body-1"
                     v-model.trim="item['code']"
             />



            </v-col>

            <v-col cols="8">
              <v-text-field label="URL"  class="body-1" v-model.trim="item['url']"/>

              <v-select :items="['GET','POST','PATCH','PUT','DELETE']" label="HTTP Method" v-model="item['http_method']" ></v-select>

              <v-textarea label="Описание" dense class="body-1" rows="3" v-model.trim="item['description']" />

            </v-col>
          </v-row>
        </td>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog v-model="dialog.open" @keydown.esc="dialog.open = false" max-width="380">
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.id">
              <v-list-item-title class="headline">Id: &laquo;{{dialog.id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name">
              <v-list-item-title class="headline">Имя : &laquo;{{dialog.name}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="green darken-1" small text @click.stop="dialog.open = false">
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn color="red darken-1" small text @click.stop="deleteItemDialog">
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
        { text: 'Id', align: 'start', value: 'public_id'},
        { text: 'Name',  value: 'name'},
        { text: 'Code',  value: 'code'},
        { text: 'URL',  value: 'url'},
        { text: 'HTTP',  value: 'http_method'},
        { text: 'Описание',  value: 'description'},
        { text: 'Статус',  value: 'enabled'},
        { text: 'Удалить', value: 'actions', sortable: false},
      ],
      items: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        id: null,
        name: '',
        url: '',
        open: false,
      },
    }),
    async mounted() {
      /*this.loading = true

      await this.getDataFromApi(false)
              .finally(()=>this.loading = false)*/
    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      },
    },
    methods: {
      getDataFromApi: async function(showNotification, searchStr = '') {

        this.loading = true
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
        }

        await this.$api.webHook.getListPagination(payload)
                .then(resp => {
                  if (resp['web_hooks'] !== undefined ) {
                    this.items = resp['web_hooks']
                    this.totalItems = resp.total
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        text: 'Список ВебХуков обновлен',
                        type: 'success'
                      });
                    }
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Ошибка обновления',
                      type: 'warring'
                    });
                  }

                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка получения списка',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(true, this.search)
      }, 1400),
      updateItemData: async function (item) {

        let payload = item
        Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

        await this.$api.webHook.update(payload)
                .then((resp) => {
                  if (resp['web_hook'] !== undefined) {
                    // this.getItemFromItems(item.id)
                    item = resp['web_hook']

                  }
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Карточка обнолвена',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },

      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новый вебхук",
        }
        await this.$api.webHook.create(payload)
                .then(async (data) => {
                  if (data['web_hook']) this.items.push(data['web_hook'])
                  await this.getDataFromApi(false) // wtf?
                  this.$notify({
                    group: 'main',
                    title: 'WebHook создан',
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

      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.name = item['name']
        this.dialog.url = item['url']
        this.dialog.open = true
      },
      deleteItemDialog: async function() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }

        this.$api.webHook.delete(payload)
                .then(async ()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'WebHook удален',
                    type: 'success'
                  });

                  // обновляем текущий список
                  await this.getDataFromApi(false)
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
                  this.dialog.url = ''
                })

      },
      getItemFromItems(itemId) {
        let index = this.items.findIndex(item => item.id === itemId, itemId)
        if (index !== -1) return this.items[index]
        else return {}
      },
      removeItemFromItems(payload) {
        let index = this.items.findIndex(item => item.id === payload.id, payload)
        if (index !== -1) this.$delete(this.items, index)
      },

      async execute(item) {
        this.loading = true
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: item.id
        }

        await this.$api.webHook.execute(payload)
                .then(() => {
                  this.$notify({
                    group: 'main',
                    title: 'WebHook выполнен',
                    type: 'success'
                  });
                })
                .catch(err => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      }
    },


  }
</script>