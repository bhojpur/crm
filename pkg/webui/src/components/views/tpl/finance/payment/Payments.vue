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
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Список платежей</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по шаблонам"
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

      <template v-slot:item.status="{ item }">
        <span :class="item['status'] === 'succeeded' ? 'green white--text px-2':'-'">{{item.status}}</span>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("Do MMMM YYYY")}}
      </template>

      <template v-slot:item.updated_at="{ item }">
        {{item['updated_at'] | moment("Do MMMM YYYY")}}
      </template>

      <template v-slot:item.actions="{ item }">
        <!--<v-icon
                small
                class="mr-4 blue&#45;&#45;text text&#45;&#45;lighten-2"
                @click="editItem(item)"
        >
          fas fa-edit
        </v-icon>-->
        <v-icon
                small
                class="red--text text--lighten-2"
                :disabled="item.status==='succeeded'"
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
          <span><strong>Имя шаблона: </strong>{{dialog.name}}</span>
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
        { text: 'Описание', align: 'start', value: 'description'},
        { text: 'status ', value: 'status' },
        { text: 'Дата создания', value: 'created_at'},
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
        name: '',
      },
    }),
    async mounted() {
      // await this.getDataFromApi(false)
    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi(false);
      }
    },
    methods: {
      getDataFromApi: async function(showNotification = false, searchStr = '') {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
        }

        this.$api.payment.getListPagination(payload)
                .then(resp => {
                  if (resp['payments'] !== undefined ) {
                    this.items = resp['payments']
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
                        text: 'Ошибка в ответе сервера: payments - not found',
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

      editItem: function(item) {
        // this.$router.push({name:'payment.edit', params: {public_id:item['public_id']}})
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
        await this.$api.payment.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Платеж удален',
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