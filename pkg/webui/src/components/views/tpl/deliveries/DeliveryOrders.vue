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
            sort-by="id"
            sort-desc
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Доставки</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" color="cyan darken-4" elevation="1" small disabled>
              <v-icon x-small class="ml-1 mr-3">far fa-external-link</v-icon>
              Создать заявку
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по заявкам"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.public_id="{ item }">
        {{ item.public_id }}/{{ item.id }}
      </template>

      <!--<template v-slot:item.order_id="{ item }">
        {{ item.order.public_id }}/{{ item.order_id }}
      </template>-->

      <!--<template v-slot:item.amount="{ item }">
        <span class="body-2">{{item['amount']['value'].toFixed(2)}}</span>
      </template>-->

      <template v-slot:item.numberOfPositions="{ item }">
        <span class="body-2">{{item['cart_items'].length - 1}}</span>
      </template>

      <template v-slot:item.address="{ item }">
        {{item['address'] ? item['address'] : '-'}}
      </template>

      <template v-slot:item.customer.name="{ item }">
        <span class="body-2">{{getCustomerFullName(item.customer)}}</span>
      </template>

      <template v-slot:item.customer.phone="{ item }">
        <span class="body-2">{{item['customer'].phone}}</span>
      </template>


      <template v-slot:item.delivery.name="{ item }">
        <span class="body-2">{{item['delivery'].name}}</span>
      </template>

      <template v-slot:item.status="{ item }">
        <v-select
            @change="updateDeliveryStatus(item, $event)"
            v-model="item.status_id"
            item-text="name"
            item-value="id"
            dense class="select-no-border"
            style="max-width: 260px"
            :items="deliveryStatuses">
          <template v-slot:selection="{ obj, index }">
            <span
                :class="getDeliveryStatusColorByGroup(item.status.group)"
                style="padding: 2px 6px 1px;border-radius: 2px;"
            >{{ item.status.name }}</span>
          </template>

        </v-select>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("HH:MM - DD.MM.YYYY")}}
      </template>

      <template v-slot:item.actions="{ item }">
        <div style="width: 40px;" class="d-flex justify-space-between">
          <v-icon
              small disabled
              class="d-inline-block mr-4 blue--text text--lighten-2"
              @click="editItem(item)"
          >
            fas fa-edit
          </v-icon>
          <v-icon
              small
              class="d-inline-block red--text text--lighten-2"
              :disabled="item.status.group==='completed'"
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
        { text: 'pID', align: 'start', value: 'public_id'},
        { text: 'Order id', align: 'start', value: 'order_id'},
        { text: 'Дата создания', value: 'created_at'},
        { text: 'ФИО заказчика', align: 'start', value: 'customer.name', sortable: false},
        { text: 'Контактный тел.', align: 'start', value: 'customer.phone', sortable: false},
        { text: 'Адрес доставки', align: 'start', value: 'address', sortable: false},
        { text: 'Тип доставки', align: 'start', value: 'delivery.name', sortable: false},
        { text: 'Стоимость доставки', align: 'center', value: 'cost', sortable: false},

        { text: 'Статус', align: 'start', value: 'status', sortable: false},
        { text: 'Actions', value: 'actions', sortable: false},
      ],
      items: [],
      deliveryStatuses: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        open: false,
      },
    }),
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      }
    },
    async mounted() {
      await Promise.all([
        this.uploadDeliveryOrderStatuses(),
      ]).then(() => this.getDataFromApi(false))
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
          // preloads:'WebSite,Customer',
          preloads:'Order,Status,WebSite,Customer',
        }

        this.$api.deliveryOrder.getListPagination(payload)
                .then(resp => {
                  if (resp['delivery_orders'] !== undefined ) {
                    this.items = resp['delivery_orders']
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
                        text: 'Ошибка в ответе сервера: delivery_orders - not found',
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
      goToCreateItem: async function () {
        return this.$router.push({name:'order.create'})
      },
      uploadDeliveryOrderStatuses: async function (showNotification = false) {
        this.loading = true
        return this.$api.deliveryStatuses.getList({
          accountHashId: this.$store.state.account.hash_id,
        })
            .then(resp => {
              if (resp['delivery_statuses'] !== undefined) {
                this.deliveryStatuses = resp['delivery_statuses']
                if (showNotification) {
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    type: 'success'
                  });
                }
              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера получены не полностью',
                  type: 'warring'
                });
              }
            })
            .catch(err => {
              this.$notify({
                group: 'main',
                title: 'Заказ не найден',
                text: err['message'],
                type: 'error'
              });
            })
            .finally(() => this.loading = false)

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
        await this.$api.deliveryOrder.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Заказ удален',
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
      editItem: function(item) {
            this.$router.push({name:'delivery.order.edit', params: {public_id:item['public_id']}})
      },
      getCustomerFullName(customer) {
        let fullName = ''
        if (customer['name'] !== undefined && customer['name'] !== null && customer['name'] !== '') fullName += customer['name']
        if (customer['surname'] !== undefined && customer['surname'] !== null && customer['surname'] !== '') fullName += " " + customer['surname']
        if (customer['patronymic'] !== undefined && customer['patronymic'] !== null && customer['patronymic'] !== '') fullName += " " + customer['patronymic']

        return fullName
      },
      getDeliveryStatusColorByGroup(group) {
        if (!group) return
        switch (group) {
            // новая доставка
          case "new":
            return "blue  lighten-2 white--text";
            // доставка подтверждена
          case "agreement":
            return "purple  lighten-3 white--text";
            // В процессе
          case "delivery":
            return "brown lighten-2 white--text";
            // Отменена
          case "canceled":
            return "red lighten-3 white--text";
            // Завершена успешно
          case "completed":
            return "green  lighten-1 white--text";
        }
      },
      updateDeliveryStatus: async function(item, status_id) {
        let id = item.status_id
        let _itemStatus = this.deliveryStatuses.find(el => el.id === id, id)
        if (_itemStatus) item.status = _itemStatus

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: item.id,
          status_id: status_id,
          preloads: "Status",
        }

        await this.$api.deliveryOrder.update(payload)
            .then((resp) => {
              if (resp['delivery_order'] !== undefined) {
                item.delivery_order = resp['delivery_order']
                this.$notify({
                  group: 'main',
                  title: 'Статус доставки обновлен',
                  type: 'success',
                });
              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера не полные',
                  text: 'не хватает объекта [delivery_order]',
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

      }
    },
  }
</script>