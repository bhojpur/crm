<template>
  <v-container fluid>
    <v-card style="border-top-left-radius: 2px;border-top-right-radius: 2px;">

      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Редактировать сайт #{{ item ? item.public_id : ''}}</h1>
          <v-btn @click="$router.push({name: 'web-sites.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку сайтов</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="!loading && item">


          <v-btn :href="'http://' + item.hostname" link small target="_blank" text>{{ item.hostname }}</v-btn>

          <!-- Сохранение и создание -->
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>
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

      <v-row class="mt-9 mx-1" v-if="item">
        <v-col cols="6" sm="6">
          <v-text-field class="mb-1" label="Имя сайта" placeholder="Имя сайта" v-model.trim="item.name"></v-text-field>
          <v-text-field class="mb-5" label="Домен" messages="без протокола http/https: example.com"
                        placeholder="Hostname" v-model.trim="item.hostname"></v-text-field>
          <v-text-field class="mb-2" label="URL" messages="полный url: https://example.com"
                        placeholder="URL сайта" v-model.trim="item.url"></v-text-field>

          <h4 class="mt-4 mb-1 headline">Склад по-умолчанию</h4>
          <div class="item-editable-popup" @click="dialogChoiceWarehouse = true">
            <span v-if="item.warehouse" class="blue-grey--text text--darken-3">{{item.warehouse.label}}</span>
            <span v-else class="deep-orange--text">n/a</span>
          </div>
        </v-col>

        <v-col cols="6">
          <v-textarea class="" label="Адрес компании, магазина)"
                      placeholder="Адрес компании или магазина" rows="3" v-model.trim="item.address"></v-textarea>
          <v-text-field class="mb-2" label="Email" placeholder="Контактный email"
                        v-model.trim="item.email"></v-text-field>
          <v-text-field class="mb-2" label="Телефон" placeholder="Контактный телефон"
                        v-model.trim="item.phone"></v-text-field>
        </v-col>
      </v-row>

      <!-- Доставка / Заказчик -->
      <v-expansion-panels class="px-0 mx-0 mt-3" multiple hover v-if="item">
        <v-expansion-panel>
          <v-expansion-panel-header class="lime lighten-4">Почтовые ящики
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <email-boxes v-if="item" :hostname.sync="item['hostname']" :webSiteId.sync="item.id"></email-boxes>
          </v-expansion-panel-content>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-header class="blue-grey lighten-4">Настройка DKIM подписи</v-expansion-panel-header>
          <v-expansion-panel-content>
            <h3 class="ml-3 mt-8" style="font-weight: 500;">Настройка DKIM подписи</h3>
            <a class="ml-3" href="https://easydmarc.com/tools/dkim-record-generator" style="text-decoration: none;"
               target="_blank">Link to generate</a>
            <v-row class="mt-5 mx-1">
              <!-- DKIM Public key -->
              <v-col cols="6">
                <v-text-field class="mb-2" label="DKIM Selector" placeholder="dkim selector" style="max-width: 200px;"
                              v-model.trim="item.dkim_selector"></v-text-field>


                <h4 class="pa-2">Generated DKIM Record</h4>
                <v-textarea
                    :value="'v=DKIM1;t=s;p=' + item['dkim_public_rsa_key']"
                    @click:append="copyBuffer('v=DKIM1;t=s;p=' + item['dkim_public_rsa_key'])"
                    append-icon="fal fa-copy"
                    auto-grow
                    background-color="#e2fdef"
                    filled
                    readonly
                    rows="3"
                    style="font-size: 12px;line-height: 1.45 !important;"
                >
                </v-textarea>
                <p>Please, publish above DNS TXT record on:
                  <strong>{{ item['dkim_selector'] + '._domainkey.' + item['hostname'] }}</strong></p>
                <p class="mt-2">SPF для Яндекса: v=spf1 ip4:45.84.226.178 include:_spf.yandex.net ~all</p>
              </v-col>

              <!-- DKIM Private key -->
              <v-col cols="6">
                <div class="pa-1">
                  <h4 class="pa-2 d-inline">Public key</h4>
                  <v-btn @click="copyBuffer(item['dkim_public_rsa_key'])" class="d-inline" depressed small style="background-color: aliceblue"
                         text>
                    <v-icon class="pr-1">fal fa-copy</v-icon>
                    Copy
                  </v-btn>
                  <v-btn @click="readonlyDKIMPublicRsaKey = !readonlyDKIMPublicRsaKey" class="d-inline ml-2" color="secondary" depressed outlined small
                         text>
                    <v-icon class="pr-1">fal fa-edit</v-icon>
                    Edit
                  </v-btn>
                </div>
                <v-textarea
                    :disabled="readonlyDKIMPublicRsaKey"
                    auto-grow
                    filled
                    outlined
                    style="font-size: 12px;line-height: 1.45 !important;"
                    v-model="item['dkim_public_rsa_key']"
                >
                </v-textarea>

                <div class="pa-1">
                  <h4 class="pa-2 d-inline">Private key</h4>
                  <v-btn @click="copyBuffer(item['dkim_public_rsa_key'])" class="d-inline" depressed small style="background-color: aliceblue"
                         text>
                    <v-icon class="pr-1">fal fa-copy</v-icon>
                    Copy
                  </v-btn>
                  <v-btn @click="readonlyDKIMPrivateRsaKey = !readonlyDKIMPrivateRsaKey" class="d-inline ml-2" color="secondary" depressed outlined small
                         text>
                    <v-icon class="pr-1">fal fa-edit</v-icon>
                    Edit
                  </v-btn>
                </div>
                <v-textarea

                    v-model="item['dkim_private_rsa_key']"
                    @click:append="copyBuffer(item['dkim_private_rsa_key'])"
                    append-icon="fal fa-copy"
                    auto-grow
                    background-color="#f7f9f8"
                    filled
                    outlined
                    style="font-size: 12px;line-height: 1.45 !important;"
                >
                </v-textarea>
              </v-col>

            </v-row>
          </v-expansion-panel-content>
        </v-expansion-panel>

      </v-expansion-panels>



      <choice-warehouse v-if="item" :open.sync="dialogChoiceWarehouse" @choice="choiceWarehouse"
                        @close="dialogChoiceWarehouse = false" :wIdx.sync="warehouseIdArr"></choice-warehouse>

    </v-card>
  </v-container>
</template>

<script>
const preloads = "Warehouse"
export default {
  data: () => ({
    item: null,
    // newTpl: true, // работаем с новым шаблоном или нет
    loading: false, // абстрактный статус загрузки..
    email_boxes: [],
    readonlyDKIMPublicRsaKey: true,
    readonlyDKIMPrivateRsaKey: true,
    dialogChoiceWarehouse: false,
  }),
  computed: {
    warehouseIdArr() {
      let arr = []
      arr.push(this.item.warehouse_id)
      return arr
    }
  },
  async mounted() {
    Promise.all([
      await this.uploadEmailBoxesItem(false),
    ]).then(() => this.uploadItem(false))
    // this.uploadItem(false)
  },
  methods: {
    uploadItem: async function (showNotification) {
      this.loading = true
      return this.$api.webSite.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads
      })
          .then(resp => {
            if (resp['web_site'] !== undefined) {
              this.item = resp['web_site']
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
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    uploadEmailBoxesItem: function (showNotification) {
      // console.log("WebSite ID: ", this.$route.params.id)
      if (this.$route.params.id === null || this.$route.params.id === undefined) return
      this.loading = true
      return this.$api.emailBox.getList({
        accountHashId: this.$store.state.account.hash_id,
        webSiteId: this.item.id,
      })
          .then(resp => {
            if (resp['email_boxes'] !== undefined) {
              this.emailBoxes = resp['email_boxes']
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
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    updateItemData: async function () {

      let payload = {accountHashId: this.$store.state.account.hash_id,preloads}

      Object.assign(payload, this.item)

      await this.$api.webSite.update(payload)
          .then((resp) => {
            if (resp['web_site'] !== undefined) {
              this.item = resp['web_site']
              this.$notify({
                group: 'main',
                title: 'Данные сохранены',
                type: 'success'
              });

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
              title: 'Ошибка сохранения',
              text: err['message'],
              type: 'error'
            });
          })

    },
    choiceWarehouse: async function (warehouse) {
      if (!warehouse || !warehouse.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        warehouse_id: warehouse.id, // new warehouse id
        preloads
      }

      await this.$api.webSite.update(payload)
          .then((resp) => {
            if (resp['web_site'] !== undefined) {
              this.item = resp['web_site']
              this.$notify({
                group: 'main',
                title: 'Данные сохранены',
                type: 'success'
              });

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
              title: 'Ошибка сохранения',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=>this.dialogChoiceWarehouse = false)



    },
  },

  components: {
    emailBoxes: () => import('@/components/views/tpl/web-sites/layouts/EmailBoxes'),
    ChoiceWarehouse: () => import('@/components/views/tpl/warehouse/warehouse/ChoiceWarehouse'),
  },

}
</script>
