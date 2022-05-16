<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items.sync="items"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
            item-key="id"
            :items-per-page="15"
            :search="search"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>
            <h1 style="font-size: 26px;">Теги товаров</h1>
          </v-toolbar-title>

          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)" color="secondary" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="createItem" color="cyan darken-4" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-plus</v-icon>
              Добавить тег
            </v-btn>
          </v-card>

        </v-toolbar>

        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по данным в карточке товаров"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.product_tag_group_id="{ item }">
        <div>
          {{item['product_tag_group_id'] && item['product_tag_group'] ? item['product_tag_group'].label : ''}}
          <v-btn
              v-if="item && item['product_tag_group'] && item['product_tag_group']['public_id']"
              :to="{name:'product-tag-group.edit', params: {public_id:item['product_tag_group']['public_id']}}" depressed text class="px-0 ml-2" x-small>
            <v-icon class="blue-grey--text text--lighten-2" small>
              fas fa-edit
            </v-icon>
          </v-btn>
        </div>
      </template>

      <template v-slot:item.product_count="{ item }">
        <span :class="{ 'red--text': !(item['product_cards'] && item['product_cards'].length > 0) }">
          {{item['product_cards'] && item['product_cards'].length > 0 ? item['product_cards'].length : '0'}}
        </span>
      </template>

      <template v-slot:item.code="{ item }">
        <span class="mr-3">{{item['code'] ? item['code'] : '-'}}</span>
      </template>

      <template v-slot:item.actions="{ item }">
        <show-action :route="{name:'product-tag.edit', params: {public_id:item['public_id']}}" @delete="openDeleteDialog(item)"></show-action>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialog.open"
            @keydown.esc="dialog.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.id">
              <v-list-item-title class="headline">Id: &laquo;{{dialog.id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.label">
              <v-list-item-title class="headline">Имя карточки: &laquo;{{dialog.label}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
                  color="green darken-1"
                  small
                  text
                  @click.stop="dialog.open = false"
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
                  color="red darken-1"
                  small
                  text
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
        { text: 'Id', align: 'start', value: 'public_id'},
        { text: 'Имя категории', value: 'label'},
        { text: 'Группа тегов',  value: 'product_tag_group_id'},
        { text: 'Код', align: 'end', value: 'code'},
        { text: 'Числится товаров',  align: 'center', value: '_product_count'},
        { text: 'Действия',  align: 'center', value: 'actions'},
      ],
      items: [],

      totalItems: 0,
      options: {
        itemsPerPage: 15,
      },
      loading: true,
      webSiteId: null, // shop_id
      search: '',
      switchProduct: '',
      expanded: [],
      dialog:{
        id: null,
        open: false,
      },

      webSites: [],
      // productGroups: [],
      webPages: [],

      dialogChoiceProduct: false, // dialog template
    }),
    async mounted() {},
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      },
    },
    methods: {
      getDataFromApi: async function(showNotification = false) {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: this.search,
          preloads: "ProductTagGroup"
        }

        this.$api.productTag.getListPagination(payload)
                .then(resp => {
                  if (resp['product_tags'] !== undefined ) {
                    this.items = resp['product_tags']
                    this.totalItems = resp.total
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        type: 'success'
                      });
                    }
                  } else {
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Ошибка обновления',
                        text: 'Ошибка в ответе сервера: product_categories - not found',
                        type: 'warring'
                      });
                    }
                  }
                })
                .catch((err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        if (this.loading) return
        this.getDataFromApi(false, this.search)
      }, 1400),
      updateLists: async function (showNotification) {
        await this.getDataFromApi(showNotification)
      },
      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          label: "Новый тег",
        }
        await this.$api.productTag.create(payload)
                .then(async (data) => {
                  if (data['product_tag']) {
                    if (!this.items) {
                      this.items = []
                    }
                    this.items.push(data['product_tag'])
                  }
                  await this.getDataFromApi(false)
                  this.$notify({
                    group: 'main',
                    title: 'Тег товара создан',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка создания',
                    text: 'Не удалось создать тег товара',
                    type: 'error'
                  });
                })
      },

      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.open = true
      },
      deleteItemDialog() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }
        this.$api.productTag.delete(payload)
                .then(async ()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'Тег удален',
                    type: 'success'
                  });

                  // обновляем текущий список
                  await this.getDataFromApi(false)
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
                  this.dialog.id = 0
                  this.dialog.label = ''
                })

      },
      switchCurrentShop(shop) {
        this.webSiteId = shop.id
        this.getDataFromApi(true)
      },

    },
    components: {
      ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
      ProductImages: () => import('@/components/views/tpl/shops/products/product-tab/Images'),
      ShowItem: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
      ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
      ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
      ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
      ChoiceProductTagGroup: () => import('@/components/views/tpl/shops/product-tag-group/ChoiceProductTagGroup'),
      // ProductVideo: () => import('@/components/views/tpl/shops/product-tab/Video'),
    },


  }
</script>