<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Данные платежа
            <v-chip color="grey lighten-3">№ {{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'orders.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку заказов</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="!loadingItem">

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

      <v-form class="mx-4 mt-6 mb-4" v-model="valid" v-if="item">

        <!-- Основные данные -->
        <v-row class="pt-4">
          <v-col cols="5">
            <v-select
                :items=[item.web_site]
                @click.stop="changeWebSite"
                class="mt-1" dense disabled
                item-text="hostname" item-value="id" label="Магазин (сайт)"
                v-if="item && item['web_site']" v-model="item['web_site_id']"></v-select>

            <v-select
                :items=[item.manager]
                @click.stop="changeManager"
                class="mt-1" dense disabled
                item-text="name" item-value="id" label="Менеджер"
                v-if="item && item['manager']" v-model="item['manager_id']"></v-select>
          </v-col>
          <v-col cols="5" offset="2">
            <v-select :items="orderChannels" class="mt-1" dense item-text="name"
                      item-value="id" label="Канал" v-model="item['order_channel_id']"></v-select>

            <v-select
                :items="[{name:'Частное лицо', status:true},{name:'Юридическое лицо', status:false}]" class="mt-1"
                dense item-text="name" item-value="status"
                label="Тип заказчика"
                v-if="item && item['individual']" v-model="item['individual']"></v-select>
          </v-col>
        </v-row>
        <v-divider class="mb-8 mt-4"></v-divider>

        <!-- Доставка -->
        <h2 class="text-h3 pb-3 d-inline">Доставка</h2>
        <v-icon
            small
            class="mr-4 ml-1 pb-1 blue--text text--lighten-2 d-inline-block"

            @click="$router.push({name:'delivery.order.edit', params: {id:item.delivery_order.public_id}})"
        >
          fas fa-edit
        </v-icon>
        <v-row v-if="item.delivery_order">
          <v-col cols="5">

            <!-- Стоимость -->
            <v-text-field
                class="mt-0"
                dense disabled
                label="Стоимость доставки" :value="item.delivery_order.amount.value"
            ></v-text-field>

            <!-- Способ доставки -->
            <v-text-field
                class="mt-0"
                label="Тип доставки" :value="item.delivery_order.delivery.name"
            ></v-text-field>

          </v-col>
          <v-col cols="5" offset="2">

          </v-col>

        </v-row>

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

    loadingItem: false, // абстрактный статус загрузки..
    orderChannels: [],
    dialogCartItem: {
      open: false,
    },
  }),
  async mounted() {
    await Promise.all([
      this.uploadOrderChannels(),

    ]).then(() => this.uploadItem())
  },
  methods: {

    getReleasePrice(item) {
      if (item.product_id > 0) {
        return item.product['retail_price'] - item.product['retail_discount']
      } else {
        return item.amount['value']
      }

    },
    uploadItem: async function (showNotification = false) {
      this.loading = true
      return this.$api.order.get({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id
      })
          .then(resp => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
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
    uploadOrderChannels: async function (showNotification = false) {
      this.loading = true
      return this.$api.orderChannel.getListPagination({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['orderChannels'] !== undefined) {
              this.orderChannels = resp['orderChannels']
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
    updateItemData: async function () {

      let payload = this.item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              this.$notify({
                group: 'main',
                title: 'Даныне заказа обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
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

      this.$api.order.update(payload)
          .then(async () => {
            // this.removeItemFromItems(payload)
            // await this.getDataFromApi(false)
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

    changeManager: async function () {
      console.log("changeManager")
    },
    changeWebSite: async function () {
      console.log("changeWebSite")
    },


  },

}
</script>
