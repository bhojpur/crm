<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items.sync="items"
            :loading="loading"
            class="elevation-1"
            item-key="created_at"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Методы оплат</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" dark color="cyan darken-4" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-external-link</v-icon>
              Добавить метод
            </v-btn>
          </v-card>
        </v-toolbar>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("Do MMMM YYYY")}}
      </template>

      <template v-slot:item.actions="{ item }">
        <v-icon
                small
                class="mr-4 blue--text text--lighten-2"
                @click="editItem(item)"
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
        { text: '#id', align: 'start', value: 'id' },
        { text: 'Имя', align: 'start', value: 'name'},
        { text: 'Тип оплаты', value: 'code' },
        { text: 'Магазин', value: 'web_site.hostname'},
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
      await this.getDataFromApi()
    },
    methods: {
      getDataFromApi: async function(showNotification = false,) {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
        }

        this.$api.paymentMethod.getList(payload)
                .then(resp => {
                  if (resp['payment_methods'] !== undefined ) {
                    this.items = resp['payment_methods']
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
                        text: 'Ошибка в ответе сервера: payment_method - not found',
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
        this.$router.push({name:'finance.payment-method.edit', params: {id:item['id']}, query: { code: item.code }})
      },
      goToCreateItem: async function () {
        /*let payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новый метод",
          enabled: false,
          executed: false,
          schedule_run: new Date(),
          subject: "Тема письма",
          preview_text: "",
          // users_segment_id: 1,
          // email_template_id: 1,
        }

        await this.$api.paymentMethod.create(payload)
            .then((resp) => {
              if (resp['email_campaign'] !== undefined) {
                this.items.unshift(resp['email_campaign'])
                this.totalItems += 1
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });

              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера не полные',
                  text: 'не хватает объекта [email_campaign]',
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
            })*/
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
          code: this.dialog.code
        }
        await this.$api.paymentMethod.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => (item.id === payload.id && item.code === payload.code), payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Метод удален',
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