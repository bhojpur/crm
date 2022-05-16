<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Доставка товаров
            <v-chip color="grey lighten-3">№ {{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'delivery.orders'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat>

          <!-- Сохранение и создание шаблона -->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" color="cyan darken-4" dark small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>
      </v-toolbar>

      <v-form class="mx-4 mt-8 mb-4" model="valid" v-if="item">

        <!-- Доставка -->
        <v-row v-if="item">
          <v-col cols="5">

            <v-select
                :items=[item.web_site]
                class="mt-1" dense disabled
                item-text="hostname" item-value="id" label="Магазин (сайт)"
                v-if="item && item['web_site']" :value="item['web_site_id']">
            </v-select>

            <!-- Стоимость -->
            <v-text-field
                v-if="item && item.amount && item.amount.value"
                v-model="item.amount.value"
                class="mt-0" dense
                disabled label="Стоимость доставки"
            ></v-text-field>

            <!-- Способ доставки -->
            <v-text-field
                v-if="item && item.delivery && item.delivery.name"
                v-model="item.delivery.name" disabled
                class="mt-0" label="Тип доставки"
            ></v-text-field>

            <v-text-field
                v-if="item && item.delivery && item.delivery.code && item.delivery.code !== 'pickup'"
                v-model="item.address"
                class="mt-0"
                label="Адрес доставки"
            ></v-text-field>

          </v-col>
          <v-col cols="5" offset="2" v-if="order">
            <h2 class="text-h4 pb-3">Комментарий клиента</h2>
            <v-textarea class="rounded body-2" outlined v-model="order['customer_comment']"></v-textarea>
          </v-col>

        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>

        <!-- Заказчик -->
        <h2 class="text-h3 pb-3">Заказчик</h2>
        <v-row v-if="item.customer">
          <v-col cols="5">
            <!-- Фамилия -->
            <v-text-field
                class="mt-0"
                dense
                label="Фамилия" v-model="item.customer.surname"
            ></v-text-field>

            <!-- Имя -->
            <v-text-field
                class="mt-0"
                dense
                label="Имя" v-model="item.customer.name"
            ></v-text-field>

            <!-- Отчество -->
            <v-text-field
                class="mt-0"
                dense
                label="Отчетство" v-model="item.customer.patronymic"
            ></v-text-field>
          </v-col>
          <v-col cols="5" offset="2">

            <!-- Телефон -->
            <v-text-field
                class="mt-0"
                label="Телефон" v-model="item.customer.phone"
            ></v-text-field>

            <v-text-field
                class="mt-0"
                label="Email" v-model="item.customer.email"
            ></v-text-field>
          </v-col>

        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>

        <!-- Список товаров -->
        <div class="d-flex justify-lg-space-between" v-if="order">
          <h2 class="text-h3 pb-3 mb-0" v-if="order.cart_items">Список товаров ({{ order.cart_items.length }})</h2>
          <div class="d-flex flex-column" v-if="order">
            <strong class="my-0 pr-2" style="font-weight: 500;" v-if="order.amount">Сумма заказа: {{ order.amount.value }}
              руб.</strong>
            <span :class="order['payment'].status === 'succeeded' ? 'green white--text px-2':''" class="body-2"
                  v-if="order.payment">статус оплаты: {{ order['payment'].status }}</span>
          </div>
        </div>

        <v-data-table
            :headers="[
                            { text: 'Id', align: 'start', value: 'product.id'},
                            { text: 'Артикул', value: 'product.article'},
                            { text: 'SKU', value: 'product.sku'},
                            { text: 'Model', value: 'product.model'},
                            { text: 'Имя', value: 'name'},
                            // { text: 'Признак способа расчета', value: 'payment_mode.name', sortable: false},
                            // { text: 'Базовая цена', value: 'product.retail_price'},
                            // { text: 'Скидка', value: 'product.retail_discount'},
                            // { text: 'Отпускная цена', value: 'product.releasePrice'},
                            { text: 'Ед.шт.', value: 'quantity'},
                            // { text: 'Удалить', value: 'actions'},
                        ]"
            :items="order.cart_items"
            sort-by="id"
            v-if="order"
        >

          <template v-slot:item.name="{ item }">
            {{ item.product['name'] ? item.product['name'] : item['description'] }}
          </template>

          <template v-slot:item.payment_subject="{ item }">
            {{ item.product['payment_subject'].name }}
          </template>

          <template v-slot:item.payment_mode.name="{ item }">
                        <span
                            :class="$api.payment.getColorByCode(item.payment_mode.code)"
                            style="padding: 2px 6px 1px;border-radius: 2px;"
                        >{{ item.payment_mode.name }}</span>
          </template>

          <template v-slot:item.product.retail_price="{ item }">
            {{ item.product['retail_price'] }}
          </template>
          <template v-slot:item.product.retail_discount="{ item }">
            {{ item.product['retail_discount'] }}
          </template>
          <template v-slot:item.product.releasePrice="{ item }">
            {{ getReleasePrice(item) }}
          </template>

          <template v-slot:item.actions="{ item }">

            <v-btn
                @click="openDeleteCartItemDialog(item)" depressed small
                style="background-color: transparent">
              <v-icon
                  class="red--text text--lighten-2"
                  small
              >
                fas fa-trash
              </v-icon>
            </v-btn>

          </template>

        </v-data-table>

      </v-form>

      <!-- Прогресс загрузки -->
      <section class="d-flex align-center text-center justify-center pa-4 flex-column"
               v-if="loadingItem"
      >
        <v-progress-circular
            :size="40"
            color="teal"
            indeterminate
        ></v-progress-circular>

        <span class="mt-2 body-2">Подождите, шаблон загружается...</span>

      </section>
    </v-card>

    <!-- Подтверждение удаления -->
    <v-dialog
        @keydown.esc="dialogCartItem.open = false"
        max-width="380"
        v-model="dialogCartItem.open"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Удалить из списка товаров?</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list v-if="dialogCartItem && dialogCartItem.product">
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.id">
              <v-list-item-title class="headline">Id: &laquo;{{ dialogCartItem.id }}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.name">
              <v-list-item-title class="headline">Имя: &laquo;{{ dialogCartItem.product.name }}&raquo;
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.model">
              <v-list-item-title class="headline">Модель: &laquo;{{ dialogCartItem.product.model }}&raquo;
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.sku">
              <v-list-item-title class="headline">SKU: &laquo;{{ dialogCartItem.product.sku }}&raquo;
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
              @click.stop="dialogCartItem.open = false"
              color="green darken-1"
              small
              text
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
              @click.stop="deleteCartItemDialog"
              color="red darken-1"
              small
              text
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
    valid: true, // form
    item: {},
    order: null, // сам заказ
    loadingItem: false, // абстрактный статус загрузки..
    orderChannels: [],
    dialogCartItem: {
      open: false,
    },
  }),
  async mounted() {
    await Promise.all([
      // this.uploadOrderChannels(),
    ]).then(() => this.uploadItem().
      then((resp)=>{
        if (resp['delivery_order'] !== undefined) {
          this.uploadOrder()
        }

    }))
  },
  methods: {

    uploadItem: async function (showNotification = false) {
      this.loading = true
      return this.$api.deliveryOrder.get({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id
      })
          .then(resp => {
            if (resp['delivery_order'] !== undefined) {
              this.item = resp['delivery_order']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });
              }
              return resp
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
            return err
          })
          .finally(() => this.loading = false)

    },
    updateItemData: async function () {

      let payload = this.item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      await this.$api.deliveryOrder.update(payload)
          .then((resp) => {
            if (resp['delivery_order'] !== undefined) {
              this.item = resp['delivery_order']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [delivery_order]',
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
    uploadOrder: async function (showNotification = false) {

      console.log("Order ID: ", this.item.order_id)

      if (!this.item.order_id) return
      await this.$api.order.get({
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.order_id
      })
          .then(resp => {
            if (resp['order'] !== undefined) {
              this.order = resp['order']
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


    },
    openDeleteCartItemDialog(item) {
      Object.assign(this.dialogCartItem, item)
      this.dialogCartItem.open = true
    },
    async deleteCartItemDialog() {
      this.dialogCartItem.open = false
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialogCartItem.id
      }

      this.$api.deliveryOrder.update(payload)
          .then(async () => {
            await this.uploadItem(false)
            Object.assign(this.dialogCartItem, {open: false})
            this.$notify({
              group: 'main',
              title: 'Товар удален из списка',
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


    },
    getReleasePrice(item) {
      if (item.product_id > 0) {
        return item.product['retail_price'] - item.product['retail_discount']
      } else {
        return item.amount['value']
      }

    },

  },

}
</script>
