<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items.sync="items"
            :loading="loading"
            item-key="created_at"
            show-expand
            expand-icon="fal fa-angle-down"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Методы доставки</h1></v-toolbar-title>
          <v-select
                  :items="webSites"
                  item-value="id"
                  item-text="name"
                  return-object
                  outlined
                  dense
                  :value="currentWebSiteId"
                  @change="switchCurrentShop"
                  style="max-width: 250px"
                  class="ml-8 mt-7"
          ></v-select>

          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)"
                   color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>

            <v-btn @click="dialogCreateItem.open = true" color="cyan darken-4" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-plus</v-icon>
              Добавить метод
            </v-btn>

          </v-card>

        </v-toolbar>
      </template>

      <template v-slot:item.code="{ item }">
        {{delivery_code_list[item['code']]}}
      </template>

      <template v-slot:item.enabled="{ item }">
        <span v-text="item['enabled'] ? 'в работе' : 'недоступен'" :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"></span>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-btn small depressed style="background-color: transparent" @click="openDeleteDialog(item)">
          <v-icon small class="red--text text--lighten-2">
            fas fa-trash
          </v-icon>
        </v-btn>
      </template>

      <template v-slot:expanded-item="{ headers, item }">
        <india-post-expanded :itemProp.sync="item" :headersLength="headers.length"
                               :currentWebSiteId="currentWebSiteId"
                               @removeItem="openDeleteDialog"
                              @switchCurrentShop="switchCurrentShop">

        </india-post-expanded>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialog.open"
            @keydown.esc="dialog.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.id">
              <v-list-item-title class="headline">Id: &laquo;{{dialog.id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog['name']">
              <v-list-item-title class="headline">Имя: &laquo;{{dialog['name']}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog['code']">
              <v-list-item-title class="headline">Тип: &laquo;{{deliveryListOptions[dialog['code']]}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="green darken-1" small text @click.stop="dialog.open = false">
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn color="red darken-1" small text @click.stop="deleteItemDialog">
            Удалить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Создание элемента -->
    <v-dialog
            v-model="dialogCreateItem.open"
            @keydown.esc="dialogCreateItem.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Создание метода доставки</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Name"  class="body-1" v-model.trim="dialogCreateItem['name']" aria-required="true" />
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-list-item-action>
                <v-select
                        label="Тип метода"
                        :items="delivery_list"
                        item-value="code" item-text="name"
                        dense
                        v-model="dialogCreateItem['code']"
                        style="max-width: 350px"
                ></v-select>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="red darken-1" small text @click.stop="dialogCreateItem.open = false">
            Отменить
          </v-btn>

          <!-- Создать -->
          <v-btn color="green darken-1" small text @click.stop="createItem" :disabled="dialogCreateItem['name'] === '' || dialogCreateItem['code'] === ''">
            Создать
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
        // { text: 'Id', align: 'start', value: 'id'},
        { text: 'Name',  value: 'name'},
        { text: 'Способ доставки',  value: 'code'},
        { text: 'Статус',  align: "center", value: 'enabled'},
        { text: 'Удалить', align: "center", value: 'actions'},
      ],
      items: [],
      webSites: [],
      currentWebSiteId: null,
      totalItems: 0,
      options: {},
      delivery_code_list: [],
      delivery_list: [
        { code: "indiaPost", name: "India Post"},
        { code: "courier", name: "Курьерская доставка"},
        { code: "pickup", name: "Самовывоз"},
      ],
      loading: true,
      search: '',
      dialog:{
        id: null,
        name: '',
        url: '',
        open: false,
      },
      dialogCreateItem:{
        name: '',
        code: '',
        open: false,
      },
    }),
    async mounted() {
      this.loading = true

      // Получаем список магазинов
      await this.updateShopList(false)
              .then(async (resp)=> {
                let webSites = this.webSites

                if (webSites && webSites.length > 0 && webSites[webSites.length-1]['id'] !== undefined) {
                  this.currentWebSiteId = webSites[webSites.length-1]['id']
                  await Promise.all([
                    this.loadDeliveryCodeList(),
                  ]).then(()=> this.getDataFromApi(false))

                }
                else this.currentWebSiteId = null

              })
              .catch((err)=>{
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: err['message'],
                  type: 'error'
                });
              })
              .finally(()=>this.loading = false)
    },
    methods: {
      getDataFromApi: async function(showNotification) {

        if (this.currentWebSiteId === null || this.currentWebSiteId === undefined) {
          if (showNotification) {
            this.$notify({
              group: 'main',
              title: 'Необходимо указать магазин',
              type: 'warring'
            });
          }
          return
        }
        this.loading = true
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.currentWebSiteId,
        }

        await this.$api.deliveries.getList(payload)
          .then(resp => {
            if (resp['deliveries'] !== undefined && resp['deliveries'] !== null) {
              // this.$set(this, 'items', resp['deliveries'])
              this.items = resp['deliveries']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список обновлен',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [deliveries]',
                type: 'warring'
              });
            }
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=>this.loading = false)

      },
      switchCurrentShop(web_site) {
        let needUpdateData = this.currentWebSiteId !== web_site.id
        this.currentWebSiteId = web_site.id
        if (needUpdateData) this.getDataFromApi(true)
      },
      createItem: async function () {

        if (this.currentWebSiteId === null || this.currentWebSiteId === undefined) {
          this.$notify({
            group: 'main',
            title: 'Необходимо выбрать магазин',
            type: 'warring'
          });
          return
        }

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.currentWebSiteId,
          code: this.dialogCreateItem.code,
          name: this.dialogCreateItem.name
        }

        await this.$api.deliveries.create(payload)
          .then(resp => {
            if (resp['delivery'] !== undefined) {
              if (resp['delivery'] !== undefined) this.items.push(resp['delivery'])
              this.$notify({
                group: 'main',
                title: 'Метод успешно добавлен',
                type: 'success'
              });
              this.dialogCreateItem.code = ""
              this.dialogCreateItem.name = ""
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [delivery]',
                type: 'warring'
              });
            }
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=> this.dialogCreateItem.open = false)

      },
      updateShopList: async function (showNotification) {
        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: 0,
          limit: 100,
          sortBy: "id",
          search: "",
          sortDesc: false,
        }

        await this.$api.webSite.getListPagination(payload)
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
                        text: 'Ошибка в ответе сервера: webSites - not found',
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
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.name = item['name']
        this.dialog.code = item['code']
        this.dialog.open = true
      },
      async deleteItemDialog() {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.currentWebSiteId,
          id: this.dialog.id,
          code: this.dialog.code
        }

        await this.$api.deliveries.delete(payload)
                .then(() => {
                  this.$notify({
                    group: 'main',
                    title: 'Объект удален',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка удаления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>{
                  this.dialog.open = false
                  this.dialog.id = ''
                  this.dialog.name = ''
                  this.dialog.code = ''
                })

        // обновляем текущий список
        await this.getDataFromApi(false)

      },

      loadDeliveryCodeList: async function() {
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          web_site_id: this.currentWebSiteId,
        }

        await this.$api.deliveries.getListOptions(payload)
                .then(resp => {
                  if (resp['delivery_code_list'] !== undefined) {
                    this.$set(this, 'delivery_code_list', resp['delivery_code_list'])
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [delivery_code_list]',
                      type: 'warring'
                    });
                  }
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })

      }
    },

    components: {
      IndiaPostExpanded: () => import('@/components/views/tpl/deliveries/layouts/DeliveryExpanded'),
    },


  }
</script>