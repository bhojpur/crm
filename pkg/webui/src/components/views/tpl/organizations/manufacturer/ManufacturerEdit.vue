<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Производитель товаров</h1>
          <v-btn @click="$router.push({name: 'manufacturers.list'})" class="ml-0 pl-0" outlined small text>
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
          <v-col cols="4">

            <!-- name -->
            <v-text-field
                class="mt-0"
                label="Имя"
                required
                v-model.trim="item['name']"
            ></v-text-field>

            <!-- short_name -->
            <v-text-field
                class="mt-0"
                label="Краткое имя"
                v-model.trim="item['short_name']"
            ></v-text-field>

            <!-- Основной вид деятельности -->
            <v-text-field
                class="mt-0"
                label="Приоритет"
                v-model.number="item['primary_activity']"
            ></v-text-field>


          </v-col>
          <v-col cols="4">
            <!-- web_site -->
            <v-text-field
                class="mt-0"
                label="Web сайт"
                v-model.trim="item['web_site']"
            ></v-text-field>

            <!-- address -->
            <v-textarea
                label="Адресс"
                counter="255" rows="2" class="mt-0"
                v-model.trim="item['address']"
            ></v-textarea>

            <!-- ShortDescription -->
            <v-textarea
                label="Краткое описание товара"
                counter="255" rows="2" class="mt-0"
                v-model.trim="item['short_description']"
            ></v-textarea>
          </v-col>
          <v-col cols="4">

            <v-file-input  v-if="item.image == null"
                           @change.native="uploadImage" :label="item.image ? 'Загрузить новое изображение' : 'Выберите файл'" show-size v-model="uploadedFile" counter></v-file-input>
            <v-card flat v-else class="pt-2 ma-0 d-flex align-end">
              <a :href="item.image._url" target="_blank" style="width: 100px;">
                <v-img :src="item['image']['_url']" aspect-ratio="1" max-width="100px"></v-img>
              </a>
              <v-btn depressed text @click="deleteImage" class="ml-3 white--text">
                <v-icon small class="red--text text--lighten-2">fas fa-trash</v-icon>
              </v-btn>

            </v-card>

          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>
    </v-card>

    <!--<v-data-table
        v-if="item"
        :headers="[
          { text: 'Id', align: 'start', value: 'public_id'},
          { text: 'Артикул', value: 'item'},
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
          <span :class="item['enabled'] ? 'green&#45;&#45;text' : 'deep-orange&#45;&#45;text darken-1'"
                v-text="item['enabled'] ? 'в продаже' : 'скрыт'"></span>
      </template>

      <template v-slot:item.payment_subject="{ item }">
        <span v-text="item['payment_subject'] ? item['payment_subject'].name : '-'"></span>
      </template>

      <template v-slot:item.actions="{ item }">
        <div class="d-flex justify-space-between" style="width: 40px;">
          <v-icon
              @click="$router.push({name:'product.edit', params: {public_id:item['public_id']}})"
              class="mr-4 blue&#45;&#45;text text&#45;&#45;lighten-2"
              small
          >
            fas fa-edit
          </v-icon>
          <v-icon
              @click="removeProduct(item)"
              class="red&#45;&#45;text text&#45;&#45;lighten-2"
              small
          >
            fas fa-trash
          </v-icon>
        </div>
      </template>

    </v-data-table>

    <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>-->

  </v-container>
</template>

<script>

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
      await this.$api.manufacturer.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads:"Image",
      })
          .then(resp => {
            if (resp['manufacturer'] !== undefined) {
              this.item = resp['manufacturer']
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
        preloads:"Image",
      })

      await this.$api.manufacturer.update(payload)
          .then((resp) => {
            if (resp['manufacturer'] !== undefined) {
              this.item = resp['manufacturer']
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

    async uploadImage($emit) {
      await this.uploadFile(this.item.id,"manufacturers")
          .then((resp)=> {
            if (resp['file'] !== undefined && resp['file']!== null) {
              this.uploadedFile = null
              this.item.image = resp['file']
              this.$notify({
                group: 'main',
                title: 'Файл создан',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Не полные данные от сервера',
                type: 'warring'
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка при создании файла',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=> {
            this.uploadedFile = null
          });
    },
    async deleteImage() {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.image.id
      }
      this.$api.storage.delete(payload)
          .then(()=> {
            this.item.image = null
            this.$notify({
              group: 'main',
              title: 'Изображение удалено',
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
          .finally(()=>{})
    },

  },
  components: {
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
  },


}
</script>

