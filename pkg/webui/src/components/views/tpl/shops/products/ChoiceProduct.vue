<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">

    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Товары и услуги</h2>
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
            { text: 'pId', align: 'start', value: 'public_id'},
            { text: 'Наименование', value: 'label'},
            { text: 'Second label', value: 'second_label'},
            { text: 'Артикул', value: 'article'},
            { text: 'Теги', align: 'end', value: 'product_tags', sortable: false},
            { text: 'Предмета расчета', align: 'end', value: 'payment_subject', sortable: false},
            { text: 'Ед. измерения', align: 'end', value: 'measurement_unit.name', sortable: false},
            { text: 'Сборный товар', align: 'end', value: 'is_kit'},
            { text: 'Цена розничная', align: 'end', value: 'retail_price'},
            { text: 'Карточек товара', align: 'end', value: 'product_cards', sortable: false},
            { text: 'Розница', align: 'end',value: 'enable_retail_sale'},
            { text: 'Опт', align: 'end', value: 'wholesale_sale'},
            { text: 'Действие', align: 'center', value: 'actions', sortable: false},
          ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search" sort-desc
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          :items-per-page="15"
          class="elevation-1"
          selectable-key="id"
          single-select
      >
        <template v-slot:item.is_source="{ item }">
          <show-item-yes-no :show="item['is_source']" :bright="false"></show-item-yes-no>
        </template>

        <template v-slot:item.enable_retail_sale="{ item }">
          <show-item-yes-no :show="item['enable_retail_sale']"></show-item-yes-no>
        </template>

        <template v-slot:item.is_kit="{ item }">
          <show-item-yes-no :show="item['is_kit']"></show-item-yes-no>
        </template>

        <template v-slot:item.wholesale_sale="{ item }">
          <show-item-yes-no :show="item['wholesale_sale']"></show-item-yes-no>
        </template>

        <template v-slot:item.product_tags="{ item }">
          <show-tags :tags="item.product_tags"></show-tags>
        </template>

        <template v-slot:item.measurement_unit="{ item }">
          {{item['measurement_unit']}}
          <span v-text="item['measurement_unit'] ?  (item['measurement_unit'].name.toLowerCase().trim()  + ' (' + item['measurement_unit'].short_name + ')') : '-'"></span>
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
          <v-btn @click="choiceItem(item)" small outlined class="px-4 grey--text text--lighten-2">
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
  props: {
    open: {
      type: Boolean,
      default: false,
    },
    isKit: { // исключить isKit
      type: Boolean,
      default: null,
    },
  },
  watch: {
    search() {
      this.options.page = 1
      this.searchDataFromApi();
    },
    open() {
      this.getDataFromApi(false)
    }
  },
  methods: {
    getDataFromApi: async function (showNotification) {

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.web_site_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        // sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        sortDesc: this.options.sortDesc ? this.options.sortDesc[0] : false,
        search: this.search,
        preloads: 'ProductCards,PaymentSubject,ProductTags,MeasurementUnit',
        isKit: this.isKit
      }
      /*if (this.isKit)
        Object.assign(payload,{isKit:this.isKit})*/

      await this.$api.product.getListPagination(payload)
          .then(resp => {
            if (resp['products'] !== undefined) {
              this.items = resp['products']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список товаров обновлен',
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
