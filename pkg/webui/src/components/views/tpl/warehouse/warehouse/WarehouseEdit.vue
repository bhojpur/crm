<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Склад товаров</h1>
          <v-btn @click="$router.push({name: 'warehouses.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="item">

          <!-- Обновление и сохранение серии -->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" color="cyan darken-4" class="white--text" small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>
      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">
          <v-col cols="6">

            <!-- Label -->
            <v-text-field
                class="mt-0"
                label="Наименование"
                v-model.trim="item['label']"
            ></v-text-field>

            <!-- Name -->
            <v-text-field
                class="mt-0"
                label="Краткое наименование"
                v-model.trim="item['short_label']"
            ></v-text-field>

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Код склада"
                required
                v-model.trim="item['code']"
            ></v-text-field>
          </v-col>

          <v-col cols="6">


            <!-- Phone -->
            <v-text-field
                class="mt-0"
                label="Телефон"
                v-model.trim="item['phone']"
            ></v-text-field>

            <!-- Email -->
            <v-text-field
                class="mt-0"
                label="Email"
                v-model.trim="item['email']"
            ></v-text-field>

            <!-- Description -->
            <v-textarea class="" label="Адрес склада"
                        placeholder="Адрес компании или магазина" rows="3" v-model.trim="item.address"></v-textarea>

          </v-col>

        </v-row>

      </v-form>

      <v-data-table
          v-if="item"
          :headers="[
            { text: 'p.Id', align: 'start', value: 'product.public_id'},
            { text: 'SKU', value: 'sku'},

            { text: 'Имя', value: 'product.label', sortable: false},
            { text: 'Доп. имя', value: 'product.second_label', sortable: false},
            { text: 'Сборный товар', align: 'end', value: 'is_kit'},
            // { text: 'Время хранения', align: 'center', value: 'expired_at'},
            // { text: 'Розничная цена', value: 'product.retail_price',align: 'center', sortable: false},
            { text: 'Ед. измерения', align: 'end', value: 'product.measurement_unit'},
            { text: 'Доступно', align: 'end', value: 'stock'},
            { text: 'Резерв', align: 'end', value: 'reservation'},
            { text: 'Всего', align: 'end', value: 'all_in_stock', sortable: false},

            { text: 'Действие', align: 'center', value: 'actions', sortable: false},
          ]"
          :items.sync="item.warehouse_items"
          class="elevation-0"
          item-key="id"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h2 style="font-size: 21px;">Список товаров</h2></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card class="d-flex align-center" flat>

              <v-btn @click="uploadItem(true)"
                     color="secondary"
                     dark outlined
                     elevation="1"
                     small
                     class="align-center"
              >
                <v-icon
                    class="ml-1 mr-3"
                    x-small
                >far fa-sync</v-icon>
                Обновить
              </v-btn>

              <v-divider
                  class="mx-2 mt-0"
                  inset
                  vertical
              ></v-divider>

              <v-btn
                  @click="dialogChoiceProduct = true"
                  outlined
                  color="cyan darken-4"
                  dark
                  elevation="1"
                  small
              >
                <v-icon
                    class="ml-1 mr-3"
                    x-small
                >far fa-plus</v-icon>
                Добавить товар
              </v-btn>

            </v-card>

          </v-toolbar>
        </template>

        <template v-slot:item.product.public_id="{ item }">
          {{item.product.public_id}}/{{item.product.id}}
        </template>

        <template v-slot:item.product.measurement_unit="{ item }">
          {{item.product['measurement_unit'] ?  (item.product['measurement_unit'].name.toLowerCase().trim()  + ' (' + item.product['measurement_unit'].short_name + ')') : '-'}}
        </template>

        <template v-slot:item.stock="{ item }">
          <span style="color: #076533;">{{item.stock}}</span>
        </template>
        <template v-slot:item.reservation="{ item }">
          <span style="color: #794f09;">{{item.reservation}}</span>
        </template>

        <template v-slot:item.all_in_stock="{ item }">
          <span style="color: #415663;">{{item.stock+item.reservation}}</span>
        </template>


        <template v-slot:item.is_kit="{ item }">
          <show-item-yes-no :show="item['is_kit']" :bright="false"></show-item-yes-no>
        </template>

        <template v-slot:item.product.retail_price="{ item }">
          {{item['product']['retail_price'].toLocaleString('ru-RU',{currency:"RUB"})}} руб.
        </template>

        <template v-slot:item.expired_at="{ item }">
          {{item['expired_at'] | moment("DD.MM.YYYY")}}
        </template>

        <template v-slot:item.actions="{ item }">
          <show-action :route="{name:'product.edit', params: {public_id:item.product['public_id']}}" @delete="removeProduct(item)" removeIcon></show-action>
        </template>

      </v-data-table>

      <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false" ></choice-product>

    </v-card>

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  data: () => ({
    formValid: true, // form
    item: null,
    uploadedFile: null, // загружаемый файл изображения
    tab: null,
    enabledEdit: true,
    products: [],
    email_boxes: [],
    web_sites: [],
    loading: false, // абстрактный статус загрузки..
    carried: false,

    dialogChoiceProduct: false, // dialog template

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
    ])
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.warehouse.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads:"WarehouseItems.Product, WarehouseItems.Product.MeasurementUnit", // может надо отдельно подгружать карточки товаров через filter[]
      })
          .then(resp => {
            if (resp['warehouse'] !== undefined) {
              this.item = resp['warehouse']
              if (!resp['warehouse']['warehouse_items']) this.item['warehouse_items'] = []
              else this.item['warehouse_items'] = resp['warehouse']['warehouse_items']
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
              title: 'Кампания не найдена',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },
    updateItemData: async function () {
      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        preloads:"WarehouseItems.Product",
      })

      await this.$api.warehouse.update(payload)
          .then((resp) => {
            if (resp['warehouse'] !== undefined) {
              this.item = resp['warehouse']

              if (!resp['warehouse']['warehouse_items']) this.item['warehouse_items'] = []
              else this.item['warehouse_items'] = resp['warehouse']['warehouse_items']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [warehouse]',
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
    choiceProduct: async function (product) {

      if (!product || !product.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads:"WarehouseItems.Product",
      }

      this.$api.warehouse.appendProduct(payload)
          .then(async (resp) => {
            if (resp['warehouse']) {

              if (!resp['warehouse']['warehouse_items']) this.item['warehouse_items'] = []
              else this.item['warehouse_items'] = resp['warehouse']['warehouse_items']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе warehouse',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogChoiceProduct = false
          })

    },
    removeProduct: async function (warehouseItem) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: warehouseItem.product_id,
        preloads:"WarehouseItems.Product",
      }

      this.$api.warehouse.removeProduct(payload)
          .then(async (resp) => {
            if (resp['warehouse']) {

              if (!resp['warehouse']['warehouse_items']) this.item['warehouse_items'] = []
              else this.item['warehouse_items'] = resp['warehouse']['warehouse_items']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе warehouse',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error',
            });
          })

    },

    savePriority: async function (showNotification = false, force = false) {

      let arr = []

      this.item.images.forEach((item, i) => {
        if (force) {
          let v = {
            id: item.id,
            priority: i + 1,
            owner_id: item.id,
            owner_type: item.owner_type
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].priority = (i + 1)
          } else {
            arr.push(v)
          }
        } else if(item.priority !== (i + 1)) {
          let v = {
            id: item.id,
            priority: i + 1,
            owner_id: this.item.id,
            owner_type: item.owner_type
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].priority = (i + 1)
          } else {
            arr.push(v)
          }
        }
      })

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        files: arr
      }

      this.$api.storage.massUpdates(payload)
          .then(resp => {
            if (resp['files'] !== undefined) {
              this.item.images = this.sortByKey(resp['files'],'priority')
              this.carried = false
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
              title: 'Ошибка обновления порядка',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)
    },

  },
  components: {
    draggable,
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),

  },


}
</script>

