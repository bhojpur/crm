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
          <v-toolbar-title><h1 style="font-size: 26px;">Автоматические серии писем</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" dark color="cyan darken-4" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-external-link</v-icon>
              Создать серию
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                v-model.trim="search"
                label="Поиск по сериям писем"
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

      <template v-slot:item._recipients="{ item }">
        {{item._recipients}}
      </template>
      <template v-slot:item._open_rate="{ item }">
        {{item._open_rate.toFixed(2)+'%'}}
      </template>
      <template v-slot:item._clickRate="{ item }">
        n/a
      </template>
      <template v-slot:item._unsubscribe_rate="{ item }">
        {{item._unsubscribe_rate.toFixed(2)+'%'}}
      </template>

      <template v-slot:item.status="{ item }">
        <span
            :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2"
        >
        {{$api.helper.getHumanStatus(item.status)}}
      </span>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("DD.MM.YYYY")}}
      </template>

      <template v-slot:item.actions="{ item }">
        <div style="width: 40px;" class="d-flex justify-space-between">
          <v-icon
              small
              class="mr-4 blue--text text--lighten-2"
              @click="$router.push({name:'email-marketing.queues.edit', params: {public_id:item['public_id']}})"
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
        </div>

      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog v-model="dialog.open" @keydown.esc="dialog.open = false" max-width="380">
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id №</strong> {{dialog.id}}</span>
          <span><strong>Имя: </strong>{{dialog.name}}</span>
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
        { text: 'ID', align: 'start', value: 'public_id'},
        { text: 'Название серии', align: 'start', value: 'name'},
        { text: 'Писем в цепочке', align: 'center', value: '_active_email_templates', sortable: false},
        { text: 'Ожидают в очереди', align: 'center', value: '_queue', sortable: false},
        { text: 'Завершили серию', align: 'center', value: '_completed', sortable: false},
        { text: 'Отправлено писем', align: 'center', value: '_recipients', sortable: false},
        { text: 'Открытий', align: 'center', value: '_open_rate', sortable: false},
        { text: 'Кликов', align: 'center', value: '_clickRate', sortable: false},
        { text: 'Отписок', align: 'center', value: '_unsubscribe_rate', sortable: false},
        { text: 'Статус работы', align: 'center', value: 'status'},
        // { text: 'Запущена', align: 'start', value: 'created_at'},
        { text: 'Дата создания', align: 'center', value: 'created_at'},
        { text: 'Actions', value: 'actions', sortable: false},
      ],
      items: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        open: false,
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

        this.$api.emailQueue.getListPagination(payload)
                .then(resp => {
                  if (resp['email_queues'] !== undefined ) {
                    this.items = resp['email_queues']
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
                        text: 'Ошибка в ответе сервера: emailQueues - not found',
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
      goToCreateItem: async function () {
        return this.$router.push({name:'email-marketing.queues.create'})
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
        await this.$api.emailQueue.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Серия удалена',
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