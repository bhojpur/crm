<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Listener</h1>
          <v-btn @click="$router.push({name: 'listeners.list'})" class="ml-0 pl-0" outlined small text>
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
          </v-col>
          <v-col cols="4">
          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>
    </v-card>

    <v-data-table
        v-if="item"
        :headers="[
          { text: 'Id', align: 'start', value: 'public_id'},
          { text: 'Артикул', value: 'article'},
          { text: 'Model', value: 'model'},
          { text: 'Имя', value: 'label'},
          { text: 'Розничная цена', value: 'retail_price',align: 'center'},
          { text: 'Отпуск',  align: 'center',value: 'enabled'},
          { text: 'Действия', value: 'actions', sortable: false},
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

      <template v-slot:item.enabled="{ item }">
          <span :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"
                v-text="item['enabled'] ? 'в продаже' : 'скрыт'"></span>
      </template>

      <template v-slot:item.payment_subject="{ item }">
        <span v-text="item['payment_subject'] ? item['payment_subject'].name : '-'"></span>
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

    <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>

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
        preloads:"Products", // может надо отдельно подгружать карточки товаров через filter[]
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
        preloads:"Products",
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
        preloads:"Products",
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
    removeProduct: async function (product) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads:"Products",
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
  },


}
</script>

