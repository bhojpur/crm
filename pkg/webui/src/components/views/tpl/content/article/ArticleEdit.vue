<template>
  <v-card style="border-top-left-radius: 2px;border-top-right-radius: 2px;" >

    <v-toolbar flat color="white">
      <v-toolbar-title class="d-flex flex-column pt-6">
        <h1 style="font-size: 26px;">Редактировать статью</h1>
        <v-btn
                @click="$router.push({name: 'content.articles'})"
                class="ml-0 pl-0"
                small text
                outlined
        >
          <v-icon
                  small
                  class="mr-3"
          >fal fa-arrow-circle-left</v-icon>
          <span style="padding-top: 2px">Назад к списку</span>
        </v-btn>
      </v-toolbar-title>
      <v-spacer></v-spacer>

      <v-btn
              @click="uploadItem"
              color="secondary"
              style="font-weight: 500"
              small dark
      >
        <v-icon
                small
                class="ml-1 mr-3"
        >
          far fa-sync
        </v-icon>
        Обновить
      </v-btn>
      <v-divider
              class="mx-2 mt-2"
              inset
              vertical
      ></v-divider>
      <v-btn
              @click="updateItem"
              color="cyan darken-4"
              style="font-weight: 500"
              small dark
      >
        <v-icon
                small
                class="mr-3"
        >fal fa-save</v-icon>
        Сохранить
      </v-btn>

    </v-toolbar>
    <v-row class="mt-9 ml-1" aria-multiline="true" v-if="item">
      <v-col cols="4" >
        <v-text-field
                v-model.trim="item.name"
                class="pt-0" label="Имя статьи"
                prepend-inner-icon="fal fa-file-code"
                placeholder="Введите имя статьи"
        ></v-text-field>

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
      </v-col>
      <v-col cols="3" offset="1">
        <v-card flat class="d-flex flex-column">
          <v-checkbox
                  v-model="item.shared"
                  :label="'Предпросмотр: ' + item.id"
                  class="mb-0 pb-0"
          ></v-checkbox>
          <div v-if="item.shared">
            <a
                    :href="$api.storage.getCDN_ARTICLES_RAW_URL() + '/'+item.hash_id"
                    target="_blank"
                    class="d-inline-flex"
                    style="text-decoration: none;color:#3c6eb6 !important;border-bottom: 1px dotted #3c6eb6;"
            >
              link raw
            </a>
          </div>
          <span v-if="!item.shared" class="deep-orange--text darken-3">скрыт</span>
        </v-card>
      </v-col>
      <v-col cols="3">
        <v-checkbox
                v-model="item.public"
                label="Опубликована"
                class="mb-0 pb-0"
        ></v-checkbox>
      </v-col>

      <v-col cols="4">
        <v-text-field
                label="Краткое имя статьи" dense class="body-1"
                v-model.trim="item['short_name']"
        />
        <v-text-field
                label="Path" dense class="body-1"
                v-model.trim="item['path']"
        />
        <v-text-field
                label="Breadcrumb" dense class="body-1"
                v-model.trim="item['breadcrumb']"
        />

        <v-textarea
                label="Описание (255 символов)" dense class="body-1"
                rows="3"
                counter
                v-model.trim="item['description']"
        />

      </v-col>

      <v-col cols="4" offset="1">
        <v-text-field
                label="Meta Title" dense class="body-1"
                v-model.trim="item['meta_title']"
        />
        <v-text-field
                label="Meta Keywords" dense class="body-1"
                v-model.trim="item['meta_keywords']"
        />
        <v-text-field
                label="Meta Description" dense class="body-1"
                v-model.trim="item['meta_description']"
        />

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

    <!-- Редактор HTML -->
    <section v-show="!loading" class="px-2 pb-4" v-if="item">

      <v-card class="d-flex justify-space-between mb-3 mx-1" flat>
        <v-card-title>
          <h2 class="mb-1">Код статьи</h2>
        </v-card-title>

        <v-card-actions>
          <v-btn
                  small text outlined @click="setBaseHtmlTemplate">Set HTML tpl</v-btn>
        </v-card-actions>

      </v-card>

      <prism-editor v-model="item.body" class="myPrismEditor" language="html"></prism-editor>
    </section>

    <!-- Прогресс загрузки -->
    <section v-if="loading"
            class="d-flex align-center text-center justify-center pa-4 flex-column"
    >
      <v-progress-circular
              :size="40"
              color="teal"
              indeterminate
      ></v-progress-circular>

      <span class="mt-2 body-2">Подождите, статья загружается...</span>

    </section>
    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialogDelete.open"
            @keydown.esc="dialogDelete.open = false"
            max-width="290"
    >
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id</strong>: {{dialogDelete.id ? dialogDelete.id : ''}}</span>
          <span><strong>Имя: </strong>{{dialogDelete ? dialogDelete.name : ''}}</span>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
                  color="primary"
                  class="mr-2"
                  small
                  text
                  outlined
                  @click.stop="dialogDelete.open = false"
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
                  color="red darken-1"
                  small
                  text
                  outlined
                  @click.stop="deleteItem"
          >
            Удалить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
  export default {
    data: () => ({
      item: null,
      uploadedFile: null, // загружаемый файл изображения
      dialogDelete:{
        open: false,
      },
      // newTpl: true, // работаем с новым шаблоном или нет
      loading: false, // абстрактный статус загрузки..

      webSites: [],
    }),
    async mounted() {
      await Promise.all([
        this.getWebSiteList(),
      ]).then(()=>this.uploadItem())
    },
    methods: {

      updateItem: async function () {

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          preloads:"Image"
        }
        Object.assign(payload, this.item)

        // delete arrays not used
        delete payload['image'];

        await this.$api.article.update(payload)
                .then((resp) => {
                  this.item = resp['article']
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Карточка обнолвена',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      async uploadImage($emit) {
        await this.uploadFile(this.item.id,"articles")
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

      openDeleteDialog() {

        this.dialogDelete.id = this.item['id']
        this.dialogDelete.name = this.item['name']
        this.dialogDelete.open = true
      },
      async deleteItem() {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialogDelete.id
        }
        await this.$api.article.delete(payload)
                .then(()=> {
                  this.$notify({
                    group: 'main',
                    title: 'Статья удалена',
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
                .finally(()=>{
                  this.dialogDelete.open = false
                  this.dialogDelete.id = ''
                  this.dialogDelete.name = ''
                });


        // обновляем текущий список
        return this.$router.push({name:'content.articles'})

      },
      uploadItem: async function () {
        this.loading = true

        // 1. Загружаем данные шаблона
        this.$api.article.getByPublicId({
          accountHashId: this.$store.state.account.hash_id,
          public_id: this.$route.params.public_id,
          preloads:"Image",
        })
        .then(resp => this.item = resp['article'])
        .catch(err => {
          this.$notify({
            group: 'main',
            title: 'Шаблон не найден',
            text: err['message'],
            type: 'warn'
          });
        })
        .finally(()=> this.loading = false)
      },
      setBaseHtmlTemplate: function () {
        this.item.body =
                `<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Bhojpur CRM - Новая статья</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
</head>
<body>
<div class="container">
    <h1>Заголовок статьи</h1>
    <p>Текст статьи...</p>
</div>
</body>
</html>`
      },

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
    },
  }
</script>
