<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Категория товаров</h1>
          <v-btn @click="$router.push({name: 'product-categories.list'})" class="ml-0 pl-0" outlined small text>
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
          <v-col cols="5">

            <!-- Label -->
            <v-text-field
                class="mt-0"
                label="Наименование в ед. числе"
                counter="128"
                required
                v-model.trim="item['label']"
            ></v-text-field>

            <!-- label_plural -->
            <v-text-field
                class="mt-0"
                label="Наименование во множ. числе"
                counter="128"
                required
                v-model.trim="item['label_plural']"
            ></v-text-field>

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Code"
                counter="128"
                v-model.trim="item['code']"
            ></v-text-field>

            <!-- Priority -->
            <v-text-field
                class="mt-0"
                label="Приоритет"
                v-model.number="item['priority']"
            ></v-text-field>

          </v-col>
          <v-col cols="4" offset="1">
            <!-- Отображать ли категорию в свойствах товара -->
            <v-checkbox
                v-model="item['show_property']"
                class="my-0 py-0"
                label="Отображать категорию в свойствах товара"
            ></v-checkbox>

            <!-- Родительская категория товара -->
            <div class="">
              <h4 class="mt-4 mb-1 primary--text mr-3" style="font-weight: 500;">Родительская категория: </h4>
              <v-btn text small
                     :class="{'red--text': !item.parent_id}"
                     @click="dialogChoiceCategory = true"
                     class="item-editable-popup">
                <v-icon class="mr-2" x-small>
                  far fa-newspaper
                </v-icon>
                {{ (item.parent && item.parent_id) ? item.parent.label : 'its root' }}
              </v-btn>
              <v-icon
                  @click="deleteParentCategory()"
                  class="blue-grey--text text--lighten-2 ml-3"
                  small
              >
                fas fa-trash
              </v-icon>

              <!--<v-btn text small depressed
                     v-if="item.parent && item.parent_id"
                     :class="{'blue-grey&#45;&#45;text': !item.parent_id}"
                     link
                     :to="{name:'product-category.edit',params: {public_id:item.parent.public_id}}"
                     class="item-editable-popup">
                <v-icon
                    class="blue-grey&#45;&#45;text text&#45;&#45;lighten-2 ml-3"
                    x-small
                >
                  far fa-link
                </v-icon>
              </v-btn>-->

            </div>
          </v-col>

        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>
      <v-data-table
          v-if="item"
          :headers="[
          { text: 'Id', align: 'start', value: 'public_id'},
          { text: 'Название товара', value: 'label'},
          { text: 'Артикул', value: 'article'},
          { text: 'Model', value: 'model'},
          { text: 'Предмета расчета', align: 'end', value: 'payment_subject', sortable: false},
          { text: 'Ед. измерения', align: 'end', value: 'measurement_unit', sortable: false},
          { text: 'Сборный товар', align: 'end', value: 'is_kit'},
          { text: 'Цена розничная', align: 'end', value: 'retail_price'},
          { text: 'Карточек товара', align: 'end', value: 'product_cards', sortable: false},
          { text: 'Розница', align: 'end',value: 'enable_retail_sale'},
          { text: 'Опт', align: 'end', value: 'wholesale_sale'},
          { text: 'Действие', align: 'center', value: 'actions', sortable: false},
        ]"
          :items.sync="item.products"
          class="elevation-1"
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
                  @click="openDialogChoiceProduct"
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

        <template v-slot:item.is_kit="{ item }">
          <show-item-yes-no :show="item['is_kit']" :bright="false"></show-item-yes-no>
        </template>
        <template v-slot:item.enable_retail_sale="{ item }">
          <show-item-yes-no :show="item['enable_retail_sale']"></show-item-yes-no>
        </template>
        <template v-slot:item.wholesale_sale="{ item }">
          <show-item-yes-no :show="item['wholesale_sale']"></show-item-yes-no>
        </template>

        <template v-slot:item.payment_subject="{ item }">
          <span v-text="item['payment_subject'] ? item['payment_subject'].name : '-'"></span>
        </template>

        <template v-slot:item.measurement_unit="{ item }">
          <span v-text="item['measurement_unit'] ?  (item['measurement_unit'].name  + ' (' + item['measurement_unit'].short_name + ')') : '-'"></span>
        </template>

        <template v-slot:item.retail_price="{ item }">
          <show-price :price="item['retail_price']"></show-price>
        </template>

        <template v-slot:item.product_cards="{ item }">
          <show-count :items="item['product_cards']"></show-count>
        </template>

        <template v-slot:item.actions="{ item }">
          <div class="d-flex justify-space-between" style="width: 40px;">
            <v-icon
                @click="$router.push({name:'product.edit', params: {public_id:item['public_id']}})"
                class="mr-4 blue--text text--lighten-2"
                small
            >
              fas fa-edit
            </v-icon>
            <v-icon
                @click="removeProduct(item)"
                class="red--text text--lighten-2"
                small
            >
              fas fa-trash
            </v-icon>
          </div>
        </template>

      </v-data-table>
    </v-card>

    <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>
    <choice-category :open.sync="dialogChoiceCategory" @choice="choiceCategory" @close="dialogChoiceCategory = false" :category="item"></choice-category>

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
    dialogChoiceCategory: false, // dialog template

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
      // this.getEmailBoxes(),
      // this.getWebSiteList(),
    ])
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.productCategory.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads: 'Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,Parent' // может надо отдельно подгружать карточки товаров через filter[]
      })
          .then(resp => {
            if (resp['product_category'] !== undefined) {
              this.item = resp['product_category']
              if (!resp['product_category']['products']) this.item['products'] = []
              else this.item['products'] = resp['product_category']['products']

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
        preloads: 'Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,Parent' // может надо отдельно подгружать карточки товаров через filter[]
      })

      await this.$api.productCategory.update(payload)
          .then((resp) => {
            if (resp['product_category'] !== undefined) {
              this.item = resp['product_category']
              if (!resp['product_category']['products']) this.item['products'] = []
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_category]',
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
        preloads: 'Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,Parent' // может надо отдельно подгружать карточки товаров через filter[]
      }

      this.$api.productCategory.appendProduct(payload)
          .then(async (resp) => {
            if (resp['product_category']) {

              if (!resp['product_category']['products']) this.item['products'] = []
              else this.item['products'] = resp['product_category']['products']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_category',
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
    choiceCategory: async function (category) {

      if (!category || !category.id || !this.item) return

      if (this.item.id === category.id) {
        this.$notify({
          group: 'main',
          title: 'Категория не может быть сама себе родителем',
          type: 'warring',
        });
        return
      }

      this.item.parent_id = category.id

      await this.updateItemData().then(()=>this.dialogChoiceCategory = false)

    },
    deleteParentCategory: async function () {

      this.item.parent_id = null

      await this.updateItemData()

    },
    removeProduct: async function (product) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads: 'Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,Parent' // может надо отдельно подгружать карточки товаров через filter[]
      }

      this.$api.productCategory.removeProduct(payload)
          .then(async (resp) => {
            if (resp['product_category']) {

              if (!resp['product_category']['products']) this.item['products'] = []
              else this.item['products'] = resp['product_category']['products']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_category',
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
    openDialogChoiceProduct(){
      this.dialogChoiceProduct = true
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
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ChoiceCategory: () => import('@/components/views/tpl/shops/product-category/ChoiceCategory'),
  },


}
</script>
