<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items.sync="items"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
            :search="search"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>
            <h1 style="font-size: 26px;">Статьи</h1>
          </v-toolbar-title>
          <v-select
              :items="webSites"
              item-value="id"
              item-text="name"
              return-object
              outlined
              dense
              :value="webSiteId"
              @change="switchCurrentShop"
              style="max-width: 380px"
              class="ml-8 mt-7"
          ></v-select>
          <v-spacer></v-spacer>

          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>

            <v-btn @click="createItem"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-plus</v-icon>
              Создать статью
            </v-btn>

          </v-card>

        </v-toolbar>

        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по данным"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.public="{ item }">
        <span v-text="item['public'] ? 'виден' : 'скрыт'"></span>
      </template>

      <template v-slot:item.created_at="{ item }">
        <span >{{parseDate(item['created_at'])}}</span>
      </template>

      <template v-slot:item.updated_at="{ item }">
        <span >{{parseDate(item['updated_at'])}}</span>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-btn :to="{name:'content.articles.edit', params: {public_id:item['public_id']}}" depressed text small class="white--text">
          <v-icon
                  small
                  class="mr-2 blue--text text--lighten-2"
          >
            fas fa-edit
          </v-icon>
        </v-btn>

        <v-icon
                small
                class="red--text text--lighten-2"
                @click="openDeleteDialog(item)"
        >
          fas fa-trash
        </v-icon>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialogDelete.open"
            @keydown.esc="dialogDelete.open = false"
            max-width="290"
    >
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id</strong>: {{dialogDelete.id}}</span>
          <span><strong>Имя: </strong>{{dialogDelete.name}}</span>
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
                  @click.stop="deleteItemDialog"
          >
            Удалить
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
        // { text: 'id', align: 'start', value: 'id' },
        { text: 'ID', align: 'start', value: 'public_id'},
        { text: 'Имя файла', align: 'start', value: 'name'},
        { text: 'Отображение', value: 'public' },
        { text: 'Дата загрузки', value: 'created_at'},
        { text: 'Дата редактирования', value: 'updated_at'},
        { text: 'Actions', align: 'center', value: 'actions'},
      ],
      dialogDelete:{
        open: false,
        id: '',
        name: '',
      },
      uploadedFile: null,
      options: {},
      items: [],
      totalItems: 0,
      search: '',
      loading: false,
      webSites: [],
      webSiteId: null,
    }),

    async mounted() {
      this.loading = true

      await Promise.all([
        this.loadWebSites(false).then(()=>{
          this.webSiteId = this.webSites.length > 0 ? _.minBy(this.webSites,'id').id : null;
        })
            .then(() => this.getDataFromApi(false)),
      ])
    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      },
    },
    methods: {
      getDataFromApi: async function (showNotification) {

        if (this.loading) return
        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: this.search,
          webSiteId: this.webSiteId,
        }

        return this.$api.article.getListPagination(payload)
          .then((resp) => {

            if (resp['articles'] !== undefined) {
              this.items = resp['articles']
              this.totalItems = resp.total
            }
            if (showNotification) {
              this.$notify({
                group: 'main',
                title: 'Список статей обновлен',
                type: 'success'
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
            return err
          })
          .finally(()=> this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        if (this.loading) return
        this.getDataFromApi(false, this.search)
      }, 1400),
      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новая статья",
          body: "Текст статьи...",
          url: "url_articles",
          meta_title: "meta_title of Article",
          web_site_id: this.webSiteId,
          public: false,
        }
        await this.$api.article.create(payload)
                .then(async (data) => {
                  if (data['article']) this.items.push(data['article'])
                  await this.getDataFromApi(false)
                  this.$notify({
                    group: 'main',
                    title: 'Статья создана',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка создания',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      openDeleteDialog(item) {

        this.dialogDelete.id = item['id']
        this.dialogDelete.name = item['name']
        this.dialogDelete.open = true
      },
      removeItemFromItems(payload) {
        let index = this.items.findIndex(item => item.id === payload.id, payload)
        if (index !== -1) this.$delete(this.items, index)
      },
      async deleteItemDialog() {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialogDelete.id
        }

        await this.$api.article.delete(payload)
                .then(()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
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
        await this.getDataFromApi(false)

      },
      switchCurrentShop(shop) {
        this.webSiteId = shop.id
        this.getDataFromApi(true)
      },

    },
    components: {

    },


  }
</script>