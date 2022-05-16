<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">

    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">

        <h2 style="font-weight: 500;">Склады товаров</h2>

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
          <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>Обновить</v-btn>

        <v-divider class="mx-4" dark vertical></v-divider>

        <v-btn @click="$emit('close')" class="deep-orange--text" elevation="0" outlined small>
          <v-icon class="mr-3">far fa-times</v-icon>Закрыть</v-btn>

      </v-card-title>

      <v-data-table
          :headers="[
            { text: 'Id', align: 'start', value: 'public_id'},
            { text: 'Наименование', value: 'label'},
            { text: 'Краткое наименование', value: 'short_label'},
            { text: 'Код склада',  value: 'code'},
            { text: 'Адрес',  value: 'address'},
            { text: 'Телефон',  value: 'phone'},
            { text: 'Email',  value: 'email'},
            { text: 'Товаров', align: 'center', value: '_product_count'},
            { text: 'Действия', value: 'actions'},
          ]"
          :items="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          :items-per-page="15"
          class="elevation-1"
          selectable-key="id"
          single-select
      >
        <template v-slot:item.name="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item.name}}</span>
        </template>
        <template v-slot:item.code="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item.code}}</span>
        </template>
        <template v-slot:item.address="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item.address}}</span>
        </template>
        <template v-slot:item.phone="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item.phone}}</span>
        </template>
        <template v-slot:item.email="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item.email}}</span>
        </template>
        <template v-slot:item._product_count="{ item }">
          <span :class="{'grey--text':productExist(item.id)}">{{item._product_count}}</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn @click="choiceItem(item)" depressed text class="px-6" :disabled="productExist(item.id)">
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
    options: {
        itemsPerPage: 15,
      },
    search: '',
    loading: true,
    dialog: {open: false},
    totalItems: 0
  }),
  // props: ['open','wIdx'],
  props: {
    open: {
      type: Boolean,
      default: false,
    },
    wIdx: {
      type: [Array, Object],
      default: function () {
        return []
      },
    },
    fIdx: {
      type: [Array, Object],
      default: function () {
        return null
      },
    }
  },
  watch: {
    search() {
      this.options.page = 1
      this.searchDataFromApi();
    },
    open() {
      if (this.open) {
        this.getDataFromApi()
      }
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
        sortDesc: this.options.sortDesc ? this.options.sortDesc[0] : true,
        search: this.search,
        // preloads: 'ProductCards,PaymentSubject,ProductTags,MeasurementUnit'
      }

      await this.$api.warehouse.getListPagination(payload)
          .then(resp => {
            if (resp['warehouses'] !== undefined) {

              if (!this.fIdx) {

                // console.log("!fIdx")

                this.items = resp['warehouses']
                this.totalItems = resp['total']

              } else {
                // console.log("this.fIdx: ", this.fIdx)
                let items = resp['warehouses']

                for (let key in items) {
                  let id = items[key].id

                  // удалять или нет
                  let ex = true

                  for (let k2 in this.fIdx) {
                    let id2 = this.fIdx[k2]
                    if (id === id2) ex = false
                  }

                  if (ex) items.splice(key, 1);
                }

                this.items = items
                this.totalItems = items.length

                // console.log(this.items)
              }

              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список складов обновлен',
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
    },
    productExist(warehouseId) {
      for(let i = 0; i < this.wIdx.length; i++){
        if (this.wIdx[i] === warehouseId) return true
      }

      return false
    },
    warehouses() {
      if (!this.fIdx) {
        this.totalItems = items.length
        return this.items
      } else {
        let items = this.items

        for (let key in items) {
          let id = items[key].id

          // удалять или нет
          let ex = true

          for (let k2 in this.fIdx) {
            let id2 = this.fIdx[k2]
            if (id === id2) {
              console.log("Совпадение: ", id, " - ", id2)
              ex = false
            }
          }

          if (ex) {
            console.log('удаляем: ', key)
            items.splice(key, 1);
          }
        }

        this.items = items
        this.totalItems = items.length
      }
    },
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
