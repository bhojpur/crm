<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Карточка товара</h1>
          <v-btn @click="$router.push({name: 'product-cards.list'})" class="ml-0 pl-0" outlined small text>
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

            <!-- Label -->
            <v-text-field
                class="mt-0"
                label="Label"
                required counter="255"
                v-model.trim="item['label']"
            ></v-text-field>

            <!-- second_label -->
            <v-text-field
                class="mt-0"
                label="Second label"
                required counter="255"
                v-model.trim="item['second_label']"
            ></v-text-field>

            <!-- Breadcrumb -->
            <v-text-field
                class="mt-0"
                label="Breadcrumb label"
                v-model.trim="item['breadcrumb']"
            ></v-text-field>

            <!-- Route name -->
            <v-text-field
                class="mt-0"
                label="Route name"
                v-model.trim="item['route_name']"
            ></v-text-field>

            <!-- Icon name -->
            <v-text-field
                class="mt-0"
                label="Icon name"
                v-model.trim="item['icon_name']"
            ></v-text-field>

            <!-- Path -->
            <v-text-field
                class="mt-0"
                label="Path card"
                v-model.trim="item['path']"
            ></v-text-field>


            <!-- Теги карточки -->
            <h4 class="mt-4 mb-1 primary--text" style="font-weight: 500;">Теги карточки товара</h4>
            <div @click="dialogChoiceTags = true" class="d-flex align-center item-editable-popup">
              <v-chip-group column>
                <template v-if="item.product_tags && item.product_tags.length > 0">
                  <v-chip v-for="tag in item.product_tags" :key="tag.id" class="body-2 blue lighten-5" small>
                    {{ tag.label }}
                  </v-chip>
                </template>
                <span v-else class="deep-orange--text">
                  <v-icon class="mr-3 ml-1" small>
                    far fa-tags
                  </v-icon>
                  n/a
                </span>
              </v-chip-group>
            </div>

          </v-col>
          <v-col cols="4">
            <!-- MetaTitle -->
            <v-text-field
                class="mt-0"
                label="Meta title"
                counter="60"
                v-model.trim="item['meta_title']"
            ></v-text-field>

            <!-- MetaKeywords -->
            <v-textarea
                class="mt-0"
                label="Meta keywords"
                counter="255" rows="2"
                v-model.trim="item['meta_keywords']"
            ></v-textarea>

            <!-- MetaDescription -->
            <v-textarea
                class="mt-0"
                label="Meta description"
                counter="180" rows="3"
                v-model.trim="item['meta_description']"
            ></v-textarea>

            <!-- ShortDescription -->
            <v-textarea
                label="Краткое описание карточки"
                counter="255" rows="3" class="mt-0"
                v-model.trim="item['short_description']"
            ></v-textarea>

            <!-- Показывать для разничной торговли -->
            <v-checkbox
                v-model="item['enable_retail_sale']"
                class="my-0 py-0"
                label="Для розничной торговли"
            ></v-checkbox>

            <!-- Показывать для оптовой торговли -->
            <v-checkbox
                v-model="item['enable_wholesale_sale']"
                class="my-0 py-0"
                label="Для оптовой торговли"
            ></v-checkbox>

          </v-col>
          <v-col cols="4">
            <draggable
                :list="item.images"
                @end="savePriority"
            >
              <transition-group name="flip-list">
                <v-card v-for="image in item.images" :key="image.id" class="mt-1 pl-1">
                  <v-row>
                    <v-col cols="3">
                      <a :href="image['_url']" target="_blank" style="width: 100px;">
                        <v-img :src="image['_url']" aspect-ratio="1" max-width="100px"></v-img>
                      </a>
                    </v-col>
                    <v-col cols="9">
                      <v-text-field v-model="image.alt" placeholder="укажите alt" dense></v-text-field>
                      <section class="d-flex justify-space-between">
                        <span class="body-2 blue-grey--text">{{formatBytes(image['size'])}}</span>
                        <span class="body-2 grey--text">priority: {{image['priority']}}</span>
                        <div>
                          <v-icon @click="updateImage(image)" small class="blue--text text--lighten-2 mr-4">fas fa-save</v-icon>
                          <v-icon @click="deleteImage(image)" small class="red--text text--lighten-2 mr-2">fas fa-trash</v-icon>
                        </div>
                      </section>
                    </v-col>
                  </v-row>
                </v-card>
              </transition-group>
            </draggable>
            <v-divider class="mb-8 mt-4" v-if="item.images && item.images.length > 0"></v-divider>

            <h3 class="headline">Добавить изображение</h3>
            <v-file-input
                   @change.native="uploadImage" :label="item.image ? 'Загрузить новое изображение' : 'Выберите файл'" show-size v-model="uploadedFile" counter></v-file-input>

          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>
        <v-expansion-panels class="px-0 mx-0" multiple hover tile>
          <v-expansion-panel style="background:#f1f2e7;color:white">
            <v-expansion-panel-header class="body-1 grey--text text--darken-4">Описание карточки товара</v-expansion-panel-header>
            <v-expansion-panel-content>
              <section class="pb-6">
                <v-btn small text outlined class="mb-1" @click="item.description='<p>Описание карточки товара..</p>'">Set HTML tpl</v-btn>
                <prism-editor v-model="item.description" class="myPrismEditor" language="html"></prism-editor>
              </section>
            </v-expansion-panel-content>
          </v-expansion-panel>

        </v-expansion-panels>
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>

      <v-data-table
          v-if="item"
          :headers="[
            { text: 'p/Id', align: 'center', value: 'public_id'},
            { text: 'Наименование', value: 'label'},
            { text: 'Second label', value: 'second_label'},
            { text: 'Артикул', value: 'article'},
            // { text: 'Model', value: 'model'},
            { text: 'Теги', align: 'end', value: 'product_tags', sortable: false},
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

        <template v-slot:item.public_id="{ item }">
          {{item.id}}/{{item.public_id}}
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

        <template v-slot:item.product_tags="{ item }">
          <show-tags :tags="item.product_tags"></show-tags>
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
          <show-action :route="{name:'product.edit', params: {public_id:item['public_id']}}" @delete="removeProduct(item)" removeIcon></show-action>
        </template>

      </v-data-table>
    </v-card>

    <choice-product-tags v-if="item" :open="dialogChoiceTags" @choice="choiceTags" @close="dialogChoiceTags = false"
                         :catIds.sync="item.product_tags"></choice-product-tags>
    <choice-product :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>

  </v-container>
</template>

<script>
let preloads = "Images,Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,ProductTags,Products.ProductTags"
import draggable from 'vuedraggable'
import _ from 'lodash'
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
    dialogChoiceTags: false, // dialog template

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
      await this.$api.productCard.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads, // может надо отдельно подгружать карточки товаров через filter[]
      })
          .then(resp => {
            if (resp['product_card'] !== undefined) {
              this.item = resp['product_card']
              if (resp['product_card']['images']) this.item.images = this.sortByKey(resp['product_card']['images'],'priority')
              if (!resp['product_card']['products']) this.item['products'] = []
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
    async uploadImage() {
      await this.uploadFile(this.item.id,"product_cards")
          .then((resp)=> {
            if (resp['file'] !== undefined && resp['file']!== null) {
              this.uploadedFile = null
              if (!this.item.images) this.item["images"] = []
              this.item.images.unshift(resp['file'])
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
    async updateImage(image) {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: image.id,
        preloads:"Images,Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit",
      }
      Object.assign(payload,image)
      this.$api.storage.update(payload)
          .then((resp)=> {

            if (resp['file']) {
              let index = this.item.images.findIndex(el=> el.id === image.id)
              if (index !== -1) this.item.images[index]=resp['file']
            }

            this.$notify({
              group: 'main',
              title: 'Изображение сохранено',
              type: 'success'
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка сохранения',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=>{})
    },
    async deleteImage(image) {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: image.id
      }
      this.$api.storage.delete(payload)
          .then(()=> {
            let index = this.item.images.findIndex(el=> el.id === image.id)
            if (index !== -1) {
              this.item.images.splice(index, 1);
            }
            this.item.images = []
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
    updateItemData: async function () {
      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        preloads:"Images,Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,ProductTags,Products.ProductTags",
      })

      await this.$api.productCard.update(payload)
          .then((resp) => {
            if (resp['product_card'] !== undefined) {
              this.item = resp['product_card']
              if (resp['product_card']['images']) this.item.images = this.sortByKey(resp['product_card']['images'],'priority')
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_card]',
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
        preloads:"Images,Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,ProductTags,Products.ProductTags",
      }

      this.$api.productCard.appendProduct(payload)
          .then(async (resp) => {
            if (resp['product_card']) {

              if (!resp['product_card']['products']) this.item['products'] = []
              else this.item['products'] = resp['product_card']['products']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_card',
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
    choiceTags: async function (data) {

      let tags = []
      let ids = []

      if (data.append) {
        for (let _index in this.item['product_tags']) {
          if (this.item['product_tags'].hasOwnProperty(_index) && typeof this.item['product_tags'][_index] === 'object') {
            ids.push(this.item['product_tags'][_index].id)
          } else {
            ids.push(this.item['product_tags'][_index])
          }
        }

        tags = _.union(ids, data.categories);

      } else {
        // надо подготовить список idx
        for (let _index in this.item['product_tags']) {
          if (this.item['product_tags'].hasOwnProperty(_index) && typeof this.item['product_tags'][_index] === 'object') {
            ids.push(this.item['product_tags'][_index].id)
          } else {
            ids.push(this.item['product_tags'][_index])
          }
        }

        tags = _.difference(ids, data.categories);

      }


      let arr = []
      for (let index in tags) {
        arr.push({id: tags[index]})
      }

      // Обновляем данные
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_tags: arr,
        preloads: "ProductTags",
      }

      await this.$api.productCard.syncProductTags(payload)
          .then((resp) => {
            if (resp['product_card'] !== undefined) {
              this.item['product_tags'] = resp['product_card']['product_tags']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_card]',
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
          .finally(() => this.dialogChoiceTags = false)
    },

    removeProduct: async function (product) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads:"Images,Products,Products.ProductCards,Products.PaymentSubject,Products.MeasurementUnit,ProductTags,Products.ProductTags",
      }

      this.$api.productCard.removeProduct(payload)
          .then(async (resp) => {
            if (resp['product_card']) {

              if (!resp['product_card']['products']) this.item['products'] = []
              else this.item['products'] = resp['product_card']['products']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_card',
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

    // Images
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

      if (arr.length < 1) {return}

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
    ChoiceProductTags: () => import('@/components/views/tpl/shops/product-tag/ChoiceTags'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
  },
}
</script>

