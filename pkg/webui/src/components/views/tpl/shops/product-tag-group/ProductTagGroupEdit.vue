<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;"># Группа тегов</h1>
          <v-btn @click="$router.push({name: 'product-tag-groups.list'})" class="ml-0 pl-0" outlined small text>
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

            <!-- filter_label -->
            <v-text-field
                class="mt-0"
                label="filter_label"
                counter="128"
                required
                v-model.trim="item['filter_label']"
            ></v-text-field>

            <!-- filter_code -->
            <v-text-field
                class="mt-0"
                label="filter_code"
                counter="128"
                v-model.trim="item['filter_code']"
            ></v-text-field>

            <v-textarea
                label="Краткое описание тега"  class="body-1 ml-0 pl-0" rows="2" flat
                counter="255"
                v-model="item['description']"
            />


          </v-col>
          <v-col cols="4" offset="1">
            <!-- Отображать ли категорию в фильтрах -->
            <v-checkbox
                v-model="item['enable_viewing']"
                class="my-0 py-0"
                label="Enable view"
            ></v-checkbox>

            <!-- enable_sorting -->
            <v-checkbox
                v-model="item['enable_sorting']"
                class="my-0 py-0"
                label="Enable sorting"
            ></v-checkbox>

            <!-- enable_many_of -->
            <v-checkbox
                v-model="item['enable_many_of']"
                class="my-0 py-0"
                label="Enable many of"
            ></v-checkbox>

            <!-- min_range -->
            <v-text-field
                class="mt-0"
                label="Min range"
                v-model.number="item['min_range']"
            ></v-text-field>

            <!-- max_range -->
            <v-text-field
                class="mt-0"
                label="Max range"
                v-model.number="item['max_range']"
            ></v-text-field>
          </v-col>
          <v-col cols="4">
          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>

      <v-data-table
          v-if="productTags"
          :headers="[
          { text: 'Id', align: 'start', value: 'public_id'},
          { text: 'Имя категории', value: 'label'},
          { text: 'Группа тегов',  value: 'product_tag_group_id'},
          { text: 'Код',  value: 'code'},
          { text: 'Числится товаров',  align: 'center', value: '_product_count'},
          { text: 'Действие', align: 'center', value: 'actions', sortable: false},
        ]"
          :items.sync="productTags"
          :options.sync="optionsTag"
          :server-items-length="totalItemTags"
          :search="searchTag"
          @update:search="searchTagDataFromApi"
          @update:options="getTagDataFromApi(false)"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h2 style="font-size: 21px;">Теги в группе</h2></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card class="d-flex align-center" flat>

              <v-btn @click="getTagDataFromApi(true)"
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
                  @click="dialogChoiceTag = true"
                  outlined
                  color="cyan darken-4"
                  elevation="1"
                  small disabled
              >
                <v-icon
                    class="ml-1 mr-3"
                    x-small
                >far fa-plus</v-icon>
                Присоединить тег
              </v-btn>

            </v-card>

          </v-toolbar>

          <v-row class="mb-4 mt-0">
            <v-col class="ml-2 mt-0 pt-0" cols="6">
              <v-text-field
                  class="mx-2"
                  hide-details
                  label="Поиск по тегам"
                  prepend-inner-icon="fal fa-search"
                  single-line
                  v-model="searchTag">
              </v-text-field>
            </v-col>
          </v-row>
        </template>

        <template v-slot:item.product_tag_group_id="{ item }">
          <div>
            {{item['product_tag_group_id'] && item['product_tag_group'] ? item['product_tag_group'].label : ''}}
          </div>
        </template>

        <template v-slot:item.product_count="{ item }">
        <span :class="{ 'red--text': !(item['product_cards'] && item['product_cards'].length > 0) }">
          {{item['product_cards'] && item['product_cards'].length > 0 ? item['product_cards'].length : '0'}}
        </span>
        </template>

        <template v-slot:item.actions="{ item }">
          <show-action :route="{name:'product-tag.edit', params: {public_id:item['public_id']}}" @delete="removeTag(item)" removeIcon></show-action>
        </template>

      </v-data-table>

    </v-card>

<!--    <choice-product :open.sync="dialogChoiceTag" @choice="choiceTag" @close="dialogChoiceTag = false"></choice-product>-->

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  data: () => ({
    formValid: true, // form
    item: null,
    loading: true, // абстрактный статус загрузки..

    productTags: [],
    totalItemTags: 0,
    searchTag: '',
    optionsTag: {},
    dialogChoiceTag: false,
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
    ]).then(()=>this.getTagDataFromApi())
  },
  watch: {
    searchTag() {
      this.searchTagDataFromApi();
    },
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.productTagGroup.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads:"Tags",
      })
        .then(resp => {
          if (resp['product_tag_group'] !== undefined) {
            this.item = resp['product_tag_group']

            if (!resp['product_tag_group']['product_tags']) this.item['product_tags'] = []

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

      await this.$api.productTagGroup.update(payload)
          .then((resp) => {
            if (resp['product_tag_group'] !== undefined) {
              this.item = resp['product_tag_group']
              if (!resp['product_tag_group']['product_tags']) this.item['product_tags'] = []
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_tag_group]',
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

    searchTagDataFromApi: _.throttle( function () {
      this.getTagDataFromApi(false, this.searchTag)
    }, 1400),
    getTagDataFromApi: async function (showNotification = false) {

      if(this.loading) return
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        offset: (this.optionsTag.page-1)*(this.optionsTag.itemsPerPage),
        limit: this.optionsTag.itemsPerPage,
        sortBy: (this.optionsTag.sortBy !== undefined && typeof this.optionsTag.sortBy[0] === 'string') ? this.optionsTag.sortBy[0] : "public_id",
        sortDesc: this.optionsTag.sortDesc[0] !== undefined ? this.optionsTag.sortDesc[0] : true,
        search: this.searchTag,
        preloads:"ProductTagGroup",
      }

      await this.$api.productTagGroup.getTagListPagination(payload)
          .then(resp => {
            if (resp['product_tags']) {
              this.productTags =  resp['product_tags']
              this.totalItemTags = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список тегов обновлен',
                  type: 'success'
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Не хватает переменной product_tags',
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

    choiceTag: async function (tag) {

      if (!tag || !tag.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_tag_id: tag.id,
        // preloads: 'ProductCards,PaymentSubject'
      }

      this.$api.productTagGroup.appendProductTag(payload)
          .then(async (resp) => {
            if (resp['product_tag_group']) {

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product_tag_group',
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
            await this.getTagDataFromApi()
            this.dialogChoiceTag = false
          })

    },
    removeTag: async function (tag) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_tag_id: tag.id,
      }

      this.$api.productTagGroup.removeProductTag(payload)
          .then(async (resp) => {
            if (resp['product_tag_group']) {

              await this.getTagDataFromApi()

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
    openDialogChoiceTag(){
      this.dialogChoiceTag = true
    },

  },
  components: {
    draggable,
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    // ChoiceProduct: () => import('@/components/views/tpl/shops/product_tags/ChoiceProduct'),
  },


}
</script>