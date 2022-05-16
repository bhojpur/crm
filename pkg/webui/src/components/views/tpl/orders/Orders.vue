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
        :items-per-page="15"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Заказы</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" dark color="cyan darken-4" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-external-link</v-icon>
              Новый заказ
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

      <template v-slot:item.amount="{ item }">
        <span class="body-2">{{ item['amount']['value'].toFixed(2) }}</span>
      </template>

      <template v-slot:item.cart_items="{ item }">
        <span class="body-2">{{ item['delivery_order'] && item['delivery_order'].id > 0 ? item['cart_items'].length - 1 : item['cart_items'].length}}</span>
      </template>

      <template v-slot:item.customer="{ item }">
        <span class="body-2">{{ item['customer_id'] ? '#' + item['customer']['id'] + ' - ' + item['customer']['name'] + " " + item['customer']['surname'] : "-" }}</span>
      </template>

      <template v-slot:item.customer.phone="{ item }">
        <span class="body-2">{{ item['customer'].phone }}</span>
      </template>

      <template v-slot:item.order_channel.name="{ item }">
        <span class="body-2">{{ item['order_channel'].name }}</span>
      </template>

      <template v-slot:item.payment.paid="{ item }">
        <span class="body-2" :class="item['payment'].status === 'succeeded' ? 'green white--text px-2':''">{{ item['payment'].status }}</span>
      </template>

      <template v-slot:item.cost="{ item }">
              <span v-if="item.product_id !== 0">
              {{item['cost'].toLocaleString('ru-RU',{currency:"RUB"})}} руб.
              </span><span v-else>-</span>
      </template>

      <template v-slot:item.status="{ item }">
        <span
            :class="getCssClass(item.status_id)"
            style="padding: 2px 6px 1px;border-radius: 2px;"
        >
          {{ orderStatuses.find(el => el.id === item.status_id) ? orderStatuses.find(el => el.id === item.status_id).name : ''}}
        </span>
      </template>

      <template v-slot:item.delivery_order="{ item }">
        <show-item-yes-no :show="item['delivery_order'] !== null" ></show-item-yes-no>

      </template>

      <template v-slot:item.created_at="{ item }">
        <span class="body-2">{{ getMomentJSData(item['created_at']) }}</span>
      </template>

      <template v-slot:item.actions="{ item }">
        <show-action :route="{name:'order.edit', params: {public_id:item['public_id']}}" @delete="openDeleteDialog(item)" editOnly></show-action>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog v-model="dialog.open" @keydown.esc="dialog.open = false" max-width="380">
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id №</strong> {{ dialog.id }}</span>
          <span><strong>Имя: </strong>{{ dialog.name }}</span>
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
const preloads = 'CartItems.Product,Customer,OrderChannel,Payment,PaymentMethod,Status,DeliveryOrder';
export default {
  data: () => ({
    headers: [
      {text: 'p.ID', align: 'start', value: 'public_id'},
      {text: 'Клиент', align: 'center', value: 'customer'},
      {text: 'Телефон', align: 'center', value: 'customer.phone'},
      // { text: 'Комментарий клиента', align: 'start', value: 'customer_comment'},
      {text: 'Канал', align: 'center', value: 'order_channel.name'},
      {text: 'Способ оплаты', align: 'start', value: 'payment_method.name'},
      {text: 'Статус оплаты', align: 'center', value: 'payment.paid', sortable: false},
      {text: 'Позиций', align: 'center', value: 'cart_items', sortable: false},
      {text: 'Сумма заказа', align: 'center', value: 'cost', sortable: false},
      {text: 'С доставкой', align: 'center', value: 'delivery_order', sortable: false},
      {text: 'Статус', align: 'center', value: 'status', sortable: false},
      {text: 'Дата создания', value: 'created_at'},
      {text: 'Подробнее', align: 'center', value: 'actions', sortable: false},
    ],
    items: [],
    orderStatuses: [],
    totalItems: 0,
    options: {
        itemsPerPage: 15,
      },
    search: '',
    loading: true,
    dialog: {
      open: false,
    },
  }),
  async mounted() {
    await Promise.all([
      this.uploadOrderStatuses(),
    ]).then(()=>this.getDataFromApi(false))
  },
  watch: {
    search() {
      this.options.page = 1;
      this.searchDataFromApi();
    },
  },
  methods: {
    getDataFromApi: async function (showNotification, searchStr = '') {

      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        preloads
      }

      this.$api.order.getListPagination(payload)
          .then(resp => {
            if (resp['orders'] !== undefined) {
              this.items = resp['orders']
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
                  text: 'Ошибка в ответе сервера: orders - not found',
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
      this.getDataFromApi(false, this.search)
    }, 1400),
    updateLists: async function (showNotification) {
      await this.getDataFromApi(showNotification)
    },
    editItem: function (item) {
      this.$router.push({name: 'order.edit', params: {public_id: item['public_id']}})
    },
    goToCreateItem: async function () {
      return this.$router.push({name: 'order.create'})
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
      await this.$api.order.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Заказ удален',
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


    uploadOrderStatuses: async function (showNotification = false) {

      return this.$api.orderStatus.getList({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['order_statuses'] !== undefined) {
              this.orderStatuses = resp['order_statuses']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Заказ не найден',
              text: err['message'],
              type: 'error',
            });
          })
    },
    updateOrderStatus: async function (item, statusId) {
      let id = item.statusId
      let _itemStatus = this.orderStatuses.find(el => el.id === id, id)
      if (_itemStatus) item.status = _itemStatus

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: item.id,
        status_id: statusId,
        preloads: 'CartItems,Customer,OrderChannel,Payment,PaymentMethod,Status',
      }


      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {

              let _item = this.items.find(el => el.id === item.id)
              if (_item) {
                // console.log(_item)
                item = resp['order']
                // _item = resp.order['status']
                // _item = resp.order['status_id']
              }

              this.$notify({
                group: 'main',
                title: 'Статус обновлен',
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
    },
    getCssClass(status_id){
      if (!status_id) return ''
      let status = this.orderStatuses.find(el=>el.id===status_id)
      if (!status) return 'grey'

      return this.$api.orderStatus.getStatusColorByGroup(status.group)
    }
  },
  components: {
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),


  },
}
</script>