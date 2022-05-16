<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">

    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор карточки товара</h2>
        <v-divider
            class="ml-6 mr-4"
            dark
            vertical
        ></v-divider>
        <v-text-field
            class="mx-2 my-0 white--text"
            style="max-width: 380px;"
            hide-details
            dense dark
            placeholder="поиск..."
            prepend-inner-icon="fal fa-search"
            single-line
            v-model.trim="search">
        </v-text-field>
        <v-spacer></v-spacer>
        <v-btn :block="false" @click="getDataFromApi(true)" class="mt-0 blue--text text--lighten-4" elevation="0" outlined small>
          <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
          Обновить
        </v-btn>
        <v-divider
            class="mx-4"
            dark
            vertical
        ></v-divider>
        <v-btn @click="$emit('close')" class="deep-orange--text" elevation="0" outlined small>
          <v-icon class="mr-3">far fa-times</v-icon>
          Закрыть
        </v-btn>
      </v-card-title>

      <v-data-table
          :headers="[
            { text: 'Id', align: 'start', value: 'public_id'},
            { text: 'Label карточки', value: 'label'},
            { text: 'Second label', value: 'second_label'},
            { text: 'Теги', value: 'product_tags', sortable: false},
            { text: 'Path',  value: 'path'},
            { text: 'Route name',  value: 'route_name'},
            { text: 'Товаров', align: 'end', value: '_product_count'},
            { text: 'Статус Розн.', align: 'end', value: 'enable_retail_sale'},
            { text: 'Статус Опт.',  align: 'end',value: 'enable_wholesale_sale'},
            { text: 'Выбрать', align: 'center', value: 'actions', sortable: false},
          ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          :items-per-page="15"
          selectable-key="id"
          single-select
      >
        <template v-slot:item.enable_retail_sale="{ item }">
          <show-item-yes-no :show="item['enable_retail_sale']"></show-item-yes-no>
        </template>

        <template v-slot:item.enable_wholesale_sale="{ item }">
          <show-item-yes-no :show="item['enable_wholesale_sale']"></show-item-yes-no>
        </template>

        <template v-slot:item.product_tags="{ item }">
          <show-tags :tags="item.product_tags"></show-tags>
        </template>

        <template v-slot:item.product_count="{ item }">
          <span v-text="item['products'] && item['products'].length > 0 ? item['products'].length : '0'" :class="{ 'red--text': !(item['products'] && item['products'].length > 0) }"></span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn @click="choiceItem(item)" depressed text class="px-6">
            <v-icon class="green--text text--lighten-2" small>fas fa-file-export</v-icon>
          </v-btn>
        </template>

      </v-data-table>
    </v-card>
    
  </v-dialog>

</template>

<script>
export default {
  data: () => ({
    items: [],
    totalItems: 0,
    options: {
        itemsPerPage: 15,
      },
    search: '',
    loading: true,
    dialog: {open: false},
  }),
  props: ['open'],
  async mounted() {
    // await this.getDataFromApi(false)
  },
  watch: {
    search() {
      this.options.page = 1
      this.searchDataFromApi();
    }
  },
  methods: {
    getDataFromApi: async function (showNotification, searchStr = '') {

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.web_site_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        preloads: "ProductTags"
      }

      await this.$api.productCard.getListPagination(payload)
          .then(resp => {
            if (resp['product_cards'] !== undefined) {
              this.items = resp['product_cards']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список карточек товаров обновлен',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                type: 'warring'
              });
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
          .finally(() => this.loading = false)
    },
    searchDataFromApi: _.throttle(function () {
      this.getDataFromApi(true, this.search)
    }, 1400),
    choiceItem: async function (item) {
      this.$emit('choice', item)
    }
  },
  components: {
    ProductParameters: () => import('@/components/views/tpl/shops/products/product-tab/Parameters'),
    ProductImages: () => import('@/components/views/tpl/shops/products/product-tab/Images'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
    // ProductVideo: () => import('@/components/views/tpl/shops/product-tab/Video'),
  },

}
</script>
