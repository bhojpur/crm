<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Web page
            <v-chip color="grey lighten-3" ripple v-if="item">{{ item.label }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'web-pages.list'})" class="ml-0 pl-0" outlined small text>
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

            <v-select
                :items="webSites"
                item-value="id"
                item-text="name"
                outlined
                dense
                v-model="item['web_site_id']"
                style="max-width: 250px"
                class="ml-0 mt-2"
            ></v-select>

            <!-- Label -->
            <v-text-field
                class="mt-0"
                label="Label"
                required
                v-model.trim="item['label']"
            ></v-text-field>

            <!-- Order -->
            <v-text-field
                class="mt-0"
                label="Порядок в категории"
                required
                v-model.number="item['order']"
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
                label="Path"
                v-model.trim="item['path']"
            ></v-text-field>

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Code (group of)"
                v-model.trim="item['code']"
            ></v-text-field>



          </v-col>
          <v-col cols="5">
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
                counter="180" rows="4"
                v-model.trim="item['meta_description']"
            ></v-textarea>

            <!-- ShortDescription -->
            <v-textarea
                label="Краткое описание страницы"
                counter="255" rows="4" class="mt-0"
                v-model.trim="item['short_description']"
            ></v-textarea>

          </v-col>
          <v-col cols="3">
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

        <v-expansion-panels class="px-0 mx-0" multiple hover tile>
          <v-expansion-panel class="light-green lighten-5">
            <v-expansion-panel-header>Описание страницы</v-expansion-panel-header>
            <v-expansion-panel-content>
              <section class="pb-6">
                <v-btn small text outlined class="mb-1" @click="item.description='<p>Описание страницы..</p>'">Set HTML tpl</v-btn>
                <prism-editor v-model="item.description" class="myPrismEditor" language="html"></prism-editor>
              </section>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>

      <!-- Categories -->
      <v-data-table
          v-if="item && item.product_categories"
          :headers="[
            { text: 'Id', align: 'start', value: 'public_id'},
            { text: 'Имя категории', value: 'label'},
            { text: 'Код',  value: 'code'},
            { text: 'Приоритет',  align: 'end', value: 'priority'},
            { text: 'Товаров', align: 'end', value: '_products_count', sortable: false},
            { text: 'Действие', align: 'center', value: 'actions', sortable: false}
          ]"
          :items.sync="item.product_categories"
          class="elevation-1"
          item-key="id"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h2 style="font-size: 21px;line-height: 1.1;">Список категорий товаров<br>
              <small style="font-weight: 400;font-size: 14px;">Для отображения на странице</small>
            </h2></v-toolbar-title>
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
                  @click="dialogChoiceCategories = true"
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
                Выбрать категории
              </v-btn>

            </v-card>

          </v-toolbar>
        </template>

        <template v-slot:item._products_count="{ item }">
          {{item['_products_count']}}
        </template>

        <template v-slot:item.actions="{ item }">
          <show-action :route="{name:'product-category.edit', params: {public_id:item['public_id']}}" @delete="removeCategory(item)" removeIcon></show-action>
        </template>

      </v-data-table>
    </v-card>

    <choice-categories v-if="item" :open="dialogChoiceCategories" @choice="choiceCategories" @close="dialogChoiceCategories = false"
                       :catIds.sync="item.product_categories"></choice-categories>

  </v-container>
</template>

<script>
import _ from 'lodash'

export default {
  data: () => ({
    formValid: true, // form
    item: null,
    uploadedFile: null, // загружаемый файл изображения

    webSites: [],
    loading: false, // абстрактный статус загрузки..

    dialogChoiceCategories: false,

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.getWebSiteList(),
    ]).then(()=>this.uploadItem())
  },
  methods: {

    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.webPage.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads: "ProductCategories",
      })
          .then(resp => {
            if (resp['web_page'] !== undefined) {
              this.item = resp['web_page']
              if (!resp['web_page']['product_categories']) this.item['product_categories'] = []
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
      await this.uploadFile(this.item.id,"web_pages")
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

    updateItemData: async function () {
      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        preloads: "ProductCategories",
      })

      await this.$api.webPage.update(payload)
          .then((resp) => {
            if (resp['web_page'] !== undefined) {
              this.item = resp['web_page']
              if (!resp['web_page']['product_categories']) this.item['product_categories'] = []
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [web_page]',
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
    choiceCategories: async function (data) {

      let categories = []
      let ids = []

      if (data.append) {
        for (let _index in this.item['product_categories']) {
          if (this.item['product_categories'].hasOwnProperty(_index) && typeof this.item['product_categories'][_index] === 'object') {
            ids.push(this.item['product_categories'][_index].id)
          } else {
            ids.push(this.item['product_categories'][_index])
          }
        }

        categories = _.union(ids, data.categories);

      } else {
        // надо подготовить список idx
        for (let _index in this.item['product_categories']) {
          if (this.item['product_categories'].hasOwnProperty(_index) && typeof this.item['product_categories'][_index] === 'object') {
            ids.push(this.item['product_categories'][_index].id)
          } else {
            ids.push(this.item['product_categories'][_index])
          }
        }

        categories = _.difference(ids, data.categories);
      }

      let arr = []
      for (let index in categories) {
        arr.push({product_category_id: categories[index]})
      }

      // Обновляем данные
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        web_page_product_categories: arr,
        preloads: "ProductCategories",
      }

      await this.$api.webPage.syncProductCategories(payload)
          .then((resp) => {
            if (resp['web_page'] !== undefined) {
              this.item['product_categories'] = resp['web_page']['product_categories']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [web_page]',
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
          .finally(() => this.dialogChoiceCategories = false)

    },
    removeCategory: async function (category) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_category_id: category.id,
        preloads: "ProductCategories",
      }

      this.$api.webPage.removeProductCategory(payload)
          .then(async (resp) => {
            if (resp['web_page'] !== undefined) {
              this.item['product_categories'] = resp['web_page']['product_categories']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [web_page]',
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

    /////////////
    removeUsersRecipient: async function () {

      this.dialogAddUser = false

      let notificationId = this.addedUsersNotificationId
      let removeListUsers = this.userSelected.map(el => el.id) // список id

      let item = this.items.find(el => el.id === notificationId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error',
        });
        return
      }
      item['recipient_users_list'] = _.difference(item['recipient_users_list'], removeListUsers);
      await this.updateItemData(item)
      this.userSelected = []
    },
    saveAddedUsersRecipient: async function () {

      this.dialogAddUser = false

      let notificationId = this.addedUsersNotificationId
      let newListUsers = this.userSelected.map(el => el.id) // список id

      let item = this.items.find(el => el.id === notificationId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error',
        });
        return
      }
      item['recipient_users_list'] = _.union(item['recipient_users_list'], newListUsers);
      await this.updateItemData(item)
      this.userSelected = []
    },
    /////////////////////

    getWebSiteList: async function (showNotification) {
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
            if (resp['web_sites'] !== undefined) {
              this.webSites = resp['web_sites']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Ошибка в ответе сервера: web_sites - not found',
                  type: 'warring',
                });
              }
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
          .finally(() => this.loading = false)

    },
    openDialogSetDelayTime() {
      this.dialogDelayTime = true
    },
    async setDelayTime(durationTime) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        delay_time: durationTime,
      }
      await this.$api.emailCampaign.update(payload)
          .then(async (resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.item.delay_time = resp['email_campaign'].delay_time
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success',
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogDelayTime = false
          })

    },

    requiredForRunNotification() {
      let arr = []
      if (!this.item.subject)
        arr.push('Необходимо указать тему сообщения')

      if (!this.item.email_box_id)
        arr.push('Необходимо установить отправителя')

      if (!this.item.email_template_id)
        arr.push('Необходимо установить шаблон')

      return arr
    },
    setNewStatus: async function (status) {
      if (status ==='') {
        return
      }
      await this.$api.emailCampaign.changeStatus({
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        status: status,
      })
          .then(resp => {
            if (resp['email_campaign'] !== undefined) {
              this.item = resp['email_campaign']
              this.$notify({
                group: 'main',
                title: 'Статус изменен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_campaign]',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Уведомление не найдено',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },

    checkDouble: async function () {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id
      }
      this.$api.emailCampaign.checkDouble(payload)
          .then(resp => {
            console.log("Дублей: ", resp['count'])
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
  },
  components: {
    ChoiceCategories: () => import('@/components/views/tpl/shops/product-category/ChoiceCategories'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
  },
}
</script>

