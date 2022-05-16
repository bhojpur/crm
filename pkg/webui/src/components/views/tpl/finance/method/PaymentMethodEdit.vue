<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;" >
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Редактирование метода</h1>
          <v-btn @click="$router.push({name: 'finance.payment-method.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку методов</span>
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

      <v-form model="valid" class="mx-4 mt-6 mb-4" v-if="item">

        <!-- Основные данные -->
        <v-row class="pt-4">
          <v-col cols="6" >
            <!-- Name -->
            <v-text-field
                label="Имя"
                v-model="item.name"
                class="body-2 mt-0"
            ></v-text-field>

            <!-- Label -->
            <v-text-field
                label="Public label"
                v-model="item.label"
                class="body-2 mt-0"
            ></v-text-field>

            <!-- Select Shop -->
            <v-select
                v-if="item && item['web_site']"
                v-model="item['web_site_id']"
                class="mt-1"
                :items="webSites" item-text="hostname" item-value="id"
                label="Магазин (сайт)">

            </v-select>


          </v-col>

          <v-col cols="6" v-if="item.type =='payment_yandexes'">
            <v-text-field
                label="api_key"
                v-model="item.api_key"
                class="body-2 mt-0"
            ></v-text-field>

            <v-text-field
                label="Yandex shop id"
                v-model="item.shop_id"
                class="body-2 mt-0"
            ></v-text-field>

            <v-text-field
                label="Hash ID метода (защита от спама)"
                disabled
                v-model="item.hash_id"
                class="body-2 mt-0"
            ></v-text-field>

          </v-col>

          <v-col cols="6" >
            <v-text-field
                label="Return Url"
                v-model="item.return_url"
                class="body-2 mt-0"
            ></v-text-field>
          </v-col>

          <v-col cols="6" v-if="item.type =='payment_yandexes'">
            <h3 class="body-1">Уведомления от Yandex</h3>
            <div class="mt-6">
              <v-text-field
                  label="HTTP-URL для уведомлений со стороны Яндекса"
                  disabled
                  :value="`https://ui.api.crm.bhojpur.net/accounts/${$store.state.account.hash_id}/yandex-payment/${item['hash_id']}/web-hooks`"
                  class="d-inline body-2 cxt">
              </v-text-field>
            </div>

            <div class="d-flex align-center">
              <v-btn
                  small class="d-inline mr-4" style="background-color: aliceblue" text depressed
                  @click="copyBuffer(`https://ui.api.crm.bhojpur.net/accounts/${$store.state.account.hash_id}/yandex-payment/${item['hash_id']}/web-hooks`)">
                <v-icon class="pr-1">fal fa-copy</v-icon>Copy URL
              </v-btn>
              <v-switch
                  class="mb-2 mt-1 ml-6 mt-0" dense hide-details light
                  :success="item['enabled_incoming_notifications']"
                  :error="!item['enabled_incoming_notifications']"
                  v-model="item['enabled_incoming_notifications']"
              >
                <template v-slot:label>
              <span class="input__label" :class="item['enabled_incoming_notifications'] ? 'green--text' : 'deep-orange--text darken-1'">
                {{item['enabled_incoming_notifications'] ? 'принимать' : 'отключить'}}
              </span>
                </template>
              </v-switch>
            </div>
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
              <v-list-item-title class="headline">Id: &laquo;{{dialogCartItem.id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.name">
              <v-list-item-title class="headline">Имя: &laquo;{{dialogCartItem.product.name}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.model">
              <v-list-item-title class="headline">Модель: &laquo;{{dialogCartItem.product.model}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.sku">
              <v-list-item-title class="headline">SKU: &laquo;{{dialogCartItem.product.sku}}&raquo;</v-list-item-title>
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
    item: [],
    loadingItem: false, // абстрактный статус загрузки..
    webSites:[],
    dialogCartItem:{
      open: false,
    },
  }),
  async mounted() {
    await Promise.all([
      this.uploadWebSites(),
    ]).then(()=> this.uploadItem())
  },
  methods: {

    uploadItem: async function (showNotification = false) {
      this.loading = true
      return this.$api.paymentMethod.get({
        accountHashId: this.$store.state.account.hash_id,
        id: this.$route.params.id,
        code: this.$route.query.code
      })
          .then(resp => {
            if (resp['payment_method'] !== undefined) {
              this.item = resp['payment_method']
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
    uploadWebSites: async function(showNotification = false) {
      await this.$api.webSite.getList({accountHashId: this.$store.state.account.hash_id})
          .then(resp => {
            if (resp['web_sites'] !== undefined ) {
              this.webSites = resp['web_sites']
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

    updateItemData: async function () {

      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
      })

      await this.$api.paymentMethod.update(payload)
          .then((resp) => {
            if (resp['payment_method'] !== undefined) {
              this.item = resp['payment_method']
              this.$notify({
                group: 'main',
                title: 'Даныне заказа обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [payment_method]',
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

      this.$api.paymentMethod.update(payload)
          .then(async ()=> {
            await this.uploadItem(false)
            Object.assign(this.dialogCartItem, {open:false})
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

    changeManager: async function(){
      console.log("changeManager")
    },
    changeWebSite: async function(){
      console.log("changeWebSite")
    }
  },

}
</script>
