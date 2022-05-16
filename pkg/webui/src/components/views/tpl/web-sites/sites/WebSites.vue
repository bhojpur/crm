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
          <v-toolbar-title><h1 style="font-size: 26px;">Сайты</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreate" dark color="cyan darken-4" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-external-link</v-icon>
              Создать сайт
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по сайтам"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.public_id="{ item }">
        {{item.public_id}}
      </template>

      <template v-slot:item.email_boxes="{ item }">
        <span class="item['email_boxes'] ? 'green--text' : 'blue--text'">{{ item['email_boxes'] ? item['email_boxes'].length : 0 }}</span>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("Do MMMM YYYY")}}
      </template>

      <template v-slot:item.updated_at="{ item }">
        {{item['updated_at'] | moment("Do MMMM YYYY")}}
      </template>

      <template v-slot:item.actions="{ item }">
        <v-icon
                small
                class="mr-4 blue--text text--lighten-2"
                @click="$router.push({name:'web-sites.edit', params: {public_id:item.public_id}})"
        >
          fas fa-edit
        </v-icon>
        <v-icon
                small
                class="red--text text--lighten-2"
                @click="openDeleteDialog(item)"
        >
          fas fa-trash
        </v-icon>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog v-model="dialog.open" @keydown.esc="dialog.open = false" max-width="380">
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id №</strong> {{dialog.id}}</span>
          <span><strong>Имя магазина: </strong>{{dialog.name}}</span>
          <span><strong>Host: </strong>{{dialog.hostname}}</span>
        </v-card-text>

        <!-- Удалить / Отменить -->
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" small text @click.stop="dialog.open = false">
            Отменить
          </v-btn>
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
        { text: 'Имя шаблона', align: 'start', value: 'name'},
        { text: 'Hostname', align: 'start', value: 'hostname'},
        { text: 'Телефон', align: 'start', value: 'phone'},
        { text: 'Email', align: 'start', value: 'email'},
        // { text: 'Отображение (#)', value: 'public' },
        // { text: 'Дата создания', value: 'created_at'},
        // { text: 'Дата редактирования', value: 'updated_at'},
        { text: 'Actions', value: 'actions', sortable: false},
      ],
      items: [], //       systemObservers
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        open: false,
        id: null,
      },
    }),
    async mounted() {
      // await this.getDataFromApi(false)
    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      }
    },
    methods: {
      getDataFromApi: async function(showNotification) {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: this.search,
        }

        this.$api.webSite.getListPagination(payload)
                .then(resp => {
                  if (resp['web_sites'] !== undefined ) {
                    this.items = resp['web_sites']
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
                        text: 'Ошибка в ответе сервера: web_sites - not found',
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
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(false, this.search)
      }, 1400),

      updateLists: async function (showNotification) {
        await this.getDataFromApi(showNotification)
      },
      goToCreate: async function (){

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новый сайт",
          hostname: "example.com",
          url: "http://example.com"
        }

        await this.$api.webSite.create(payload)
                .then((resp) => {
                  if (resp['web_site'] !== undefined) {
                    this.items.unshift(resp['web_site'])
                    this.$router.push({name:'web-sites.edit', params: {public_id:resp['web_site'].public_id}})
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [web_site]',
                      type: 'warring'
                    });
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
      },

      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {
        Object.assign(this.dialog, item)
        this.dialog.open = true
      },
      deleteItemDialog: async function() {
        this.dialog.open = false

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id,
        }
        await this.$api.webSite.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Объект удален',
                    type: 'success'
                  });
                  Object.assign(this.dialog, {open: false})
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка удаления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
    },

  }
</script>