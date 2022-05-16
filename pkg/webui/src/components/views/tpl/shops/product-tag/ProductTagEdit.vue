<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Тег товара</h1>
          <v-btn @click="$router.push({name: 'product-tags.list'})" class="ml-0 pl-0" outlined small text>
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

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Code"
                counter="128"
                v-model.trim="item['code']"
            ></v-text-field>

            <h4 class="mt-4 mb-1 primary--text" style="font-weight: 500;">Входит в группу тегов</h4>
            <div class="d-flex align-center">
              <v-btn text small
                     :class="{'red--text': !item.product_tag_group_id}"
                     @click="dialogChoiceTagGroup = true"
                     class="item-editable-popup" style="padding-left: 20px;padding-right: 20px;">
                <v-icon class="mr-2" x-small>
                  far fa-newspaper
                </v-icon>
                {{ (item.product_tag_group && item.product_tag_group_id) ? item.product_tag_group.label : 'n/a'}}
              </v-btn>
              <!-- Buttons -->
              <div class="d-flex align-center">
                <v-btn
                    v-if="item && item['product_tag_group'] && item['product_tag_group']['public_id']"
                    :to="{name:'product-tag-group.edit', params: {public_id:item['product_tag_group']['public_id']}}" depressed text class="px-0 ml-2" x-small>
                  <v-icon class="blue--text text--lighten-2" small>
                    fas fa-edit
                  </v-icon>
                </v-btn>
                <v-icon
                    @click="deleteTagGroup()"
                    class="deep-orange--text text--lighten-2 ml-3"
                    small
                >
                  fas fa-trash
                </v-icon>
              </div>
            </div>

          </v-col>
          <v-col cols="4" offset="3">
            <!-- color -->
            <v-color-picker :value="item['color']" @update:color="el => item['color']=el.hexa" flat></v-color-picker>

          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>

      <v-data-table
          v-if="products"
          :headers="[
          { text: 'Id', align: 'start', value: 'public_id'},
          { text: 'Название товара', value: 'label'},
          { text: 'Артикул', value: 'article'},
          { text: 'Model', value: 'model'},
          { text: 'Предмета расчета', align: 'end', value: 'payment_subject', sortable: false},
          { text: 'Сборный', align: 'end', value: 'is_kit'},
          { text: 'Цена розничная', align: 'end', value: 'retail_price'},
          { text: 'Карточек товара', align: 'end', value: 'product_cards', sortable: false},
          { text: 'Розница', align: 'end',value: 'enable_retail_sale'},
          { text: 'Опт', align: 'end', value: 'wholesale_sale'},
          { text: 'Действие', align: 'center', value: 'actions', sortable: false},
        ]"
          :items.sync="products"
          :options.sync="optionsProduct"
          :server-items-length="totalItemProducts"
          :search="searchProduct"
          @update:search="searchProductDataFromApi"
          @update:options="getProductDataFromApi(false)"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h2 style="font-size: 21px;">Товары и услуги</h2></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card class="d-flex align-center" flat>

              <v-btn @click="getProductDataFromApi(true)"
                     color="secondary"
                     dark outlined
                     elevation="1"
                     small
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
                Присоединить товар
              </v-btn>

            </v-card>

          </v-toolbar>

          <v-row class="mb-4 mt-0">
            <v-col class="ml-2 mt-0 pt-0" cols="6">
              <v-text-field
                  class="mx-2"
                  hide-details
                  label="Поиск по товарам и услугам"
                  prepend-inner-icon="fal fa-search"
                  single-line
                  v-model="searchProduct">
              </v-text-field>
            </v-col>
          </v-row>
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

        <template v-slot:item.retail_price="{ item }">
          <show-price :price="item['retail_price']"></show-price>
        </template>

        <template v-slot:item.product_cards="{ item }">
          <show-count :items="item['product_cards']"></show-count>
        </template>

        <template v-slot:item.actions="{ item }">
          <show-action :route="{name:'product.edit', params: {public_id:item['public_id']}}" @delete="removeProduct(item)" removeIcon></show-action>
        </template>

      </v-data-table>
    </v-card>



    <choice-product-tag-group v-if="item" :open="dialogChoiceTagGroup" @choice="choiceTagGroup" @close="dialogChoiceTagGroup = false"
                         :choiceId.sync="item.product_type_id"></choice-product-tag-group>
    <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  data: () => ({
    formValid: true, // form
    item: null,

    products: [],
    totalItemProducts: 0,
    searchProduct: '',
    optionsProduct: {},

    loading: true, // абстрактный статус загрузки..
    carried: false,

    dialogChoiceTagGroup: false, // dialog template
    dialogChoiceProduct: false, // dialog template

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
    ]).then(()=>this.getProductDataFromApi())
  },
  watch: {
    searchProduct() {
      this.searchProductDataFromApi();
    },
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.productTag.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads:"ProductTagGroup",
      })
          .then(resp => {
            if (resp['product_tag'] !== undefined) {
              this.item = resp['product_tag']

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
        preloads:"ProductTagGroup",
      })

      await this.$api.productTag.update(payload)
          .then((resp) => {
            if (resp['product_tag'] !== undefined) {
              this.item = resp['product_tag']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_tag]',
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

    searchProductDataFromApi: _.throttle( function () {
      this.getProductDataFromApi(false, this.searchProduct)
    }, 1400),
    getProductDataFromApi: async function (showNotification = false) {

      if(this.loading) return
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        offset: (this.optionsProduct.page-1)*(this.optionsProduct.itemsPerPage),
        limit: this.optionsProduct.itemsPerPage,
        sortBy: (this.optionsProduct.sortBy !== undefined && typeof this.optionsProduct.sortBy[0] === 'string') ? this.optionsProduct.sortBy[0] : "public_id",
        sortDesc: this.optionsProduct.sortDesc[0] !== undefined ? this.optionsProduct.sortDesc[0] : true,
        search: this.searchProduct,
        preloads: 'ProductCards,PaymentSubject'
      }

      await this.$api.productTag.getProductListPagination(payload)
          .then(resp => {
            if (resp['products']) {
              this.products =  resp['products']
              this.totalItemProducts = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список товаров обновлен',
                  type: 'success'
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Не хватает переменной products',
                  type: 'warring'
                });
              }
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

    choiceProduct: async function (product) {

      if (!product || !product.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads: 'ProductCards,PaymentSubject'
      }

      this.$api.productTag.appendProduct(payload)
          .then(async (resp) => {
            if (resp['product_tag']) {

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_tag',
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
          .finally(async () => {
            await this.getProductDataFromApi()
            this.dialogChoiceProduct = false
          })

    },
    choiceTagGroup: async function (productTagGroup) {

      this.item.product_tag_group_id = parseInt(productTagGroup.id)

      await this.updateItemData(false)
          .finally(() => this.dialogChoiceTagGroup = false)
    },
    removeProduct: async function (product) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
      }

      this.$api.productTag.removeProduct(payload)
          .then(async (resp) => {
            if (resp['product_tag']) {

              await this.getProductDataFromApi()

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_tag',
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

    deleteTagGroup: async function () {
      this.item.product_tag_group_id = null
      await this.updateItemData()
    },

  },
  components: {
    draggable,
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    ChoiceProductTagGroup: () => import('@/components/views/tpl/shops/product-tag-group/ChoiceProductTagGroup'),
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ProductParameters: () => import('@/components/views/tpl/shops/products/product-tab/Parameters'),
    ProductImages: () => import('@/components/views/tpl/shops/products/product-tab/Images'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
  },


}
</script>
