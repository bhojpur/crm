<template>
  <v-container fluid>
    <v-data-table
        :headers="headers"
        :items.sync="items"
        :options.sync="options"
        :server-items-length="totalItems"
        :loading="loading"
        class="elevation-1"
        :items-per-page="15"
        :search="search"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar color="white" flat>
          <v-toolbar-title><h1 style="font-size: 26px;">Товары и услуги</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)"
                   color="secondary"
                   dark
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
                    @click="createItem"
                    color="cyan darken-4"
                    dark elevation="1" small
            >
              <v-icon
                      class="ml-1 mr-3"
                      x-small
              >far fa-plus</v-icon>
              Добавить товар
            </v-btn>

          </v-card>

        </v-toolbar>

        <v-row>
          <v-col class="ml-2" cols="6">
            <v-text-field
                    class="mx-2"
                    hide-details
                    label="Поиск по товарам и услугам"
                    prepend-inner-icon="fal fa-search"
                    single-line
                    v-model.trim="search">
            </v-text-field>
          </v-col>
        </v-row>
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

      <template v-slot:item.payment_subject="{ item }">
        <span v-text="item['payment_subject'] ? item['payment_subject'].name.toLowerCase().trim() : '-'"></span>
      </template>

      <template v-slot:item.measurement_unit="{ item }">
        <span v-text="item['measurement_unit'] ?  (item['measurement_unit'].name.toLowerCase().trim()  + ' (' + item['measurement_unit'].short_name + ')') : '-'"></span>
      </template>

      <template v-slot:item.product_tags="{ item }">
        <show-tags :tags="item.product_tags"></show-tags>
      </template>

      <template v-slot:item.retail_price="{ item }">
        <show-price :price="item['retail_price']"></show-price>
      </template>

      <template v-slot:item.product_cards="{ item }">
        <show-count :items="item['product_cards']"></show-count>
      </template>

      <template v-slot:item.actions="{ item }">
        <show-action :route="{name:'product.edit', params: {public_id:item['public_id']}}" @delete="openDeleteDialog(item)"></show-action>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog
            @keydown.esc="dialog.open = false"
            max-width="380"
            v-model="dialog.open"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.public_id">
              <v-list-item-title class="headline">Id: &laquo;{{dialog.public_id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.label">
              <v-list-item-title class="headline">Имя: &laquo;{{dialog.label}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.model">
              <v-list-item-title class="headline">Модель: &laquo;{{dialog.model}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
                  @click.stop="dialog.open = false"
                  color="green darken-1"
                  small
                  text
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
                  @click.stop="deleteItemDialog"
                  color="red darken-1"
                  small
                  text
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
        { text: 'pId', align: 'start', value: 'public_id'},
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
      ],
      items: [],
      totalItems: 0,
      options:{
        itemsPerPage: 15,
      },
      loading: true,
      search: '',
      tab: null,

      dialog:{
        id: null,
        name: '',
        model: null,
        sku: null,
        open: false,
      },

      payment_subjects: [],
      vat_codes: [],
    }),
    async mounted() {
        await Promise.all([
          this.getPaymentSubjectsFromApi(),
          this.getVatCodesFromApi()
        ])
            // .then(()=>this.getDataFromApi(false))
      // await this.getDataFromApi(false)
    },
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
          preloads: 'ProductCards,PaymentSubject,MeasurementUnit,ProductTags'
        }

        await this.$api.product.getListPagination(payload)
                .then(resp => {
                  if (resp['products']) {
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
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Ошибка обновления',
                        text: 'Не хватает переменной products',
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
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(false, this.search)
      }, 1400),
      getPaymentSubjectsFromApi: async function(showNotification = false) {

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
        }

        await this.$api.paymentSubject.getList(payload)
                .then(resp => {
                  if (resp['payment_subjects'] !== undefined) {
                    this.payment_subjects = resp['payment_subjects']
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Ошибка получения списка payment Subjects',
                      type: 'warring'
                    });
                  }

                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка получения списка Payment Subjects',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      getVatCodesFromApi: async function(showNotification = false) {

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
        }

        await this.$api.vatCode.getList(payload)
                .then(resp => {
                  if (resp['vat_codes'] !== undefined) {
                    this.vat_codes = resp['vat_codes']
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Ошибка получения списка payment vat_codes',
                      type: 'warring'
                    });
                  }

                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка получения списка Vat Codes',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          label: "Новый товар",
          payment_subject_id: 1,
        }
        await this.$api.product.create(payload)
                .then(async (data) => {
                  if (data['product']) this.items.push(data['product'])
                  await this.getDataFromApi(false)
                  this.$notify({
                    group: 'main',
                    title: 'Продукт создан',
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

      removeItemFromItems(payload) {
        let index = this.items.findIndex(item => item.id === payload.id, payload)
        if (index !== -1) this.$delete(this.items, index)
      },
      openDeleteDialog(item) {
        Object.assign(this.dialog, item)
        this.dialog.open = true
      },
      async deleteItemDialog() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }

        await this.$api.product.delete(payload)
                .then(async ()=> {
                  this.removeItemFromItems(payload)
                  await this.getDataFromApi(false)

                  this.$notify({
                    group: 'main',
                    title: 'Товар удален',
                    type: 'success'
                  });

                  // обновляем текущий список

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
                  this.dialog.name = ''
                })

      },
      getFullName(item) {
        if (item['user'] === undefined) return ""

        let fillName = ""
        if (item.user['name'] !== undefined && item.user['surname'] !== undefined) {
          fillName = item.user['name'] + " " + item.user['surname']
        } else {
          if (item.user['name'] !== undefined) {
            fillName = item.user['name']
          } else {
            fillName = "-"
          }
        }
        return fillName
      },
    },
    components: {
      ProductParameters: () => import('@/components/views/tpl/shops/products/product-tab/Parameters'),
      ProductImages: () => import('@/components/views/tpl/shops/products/product-tab/Images'),
      ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
      ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
      ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
      ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
      ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
      // ProductVideo: () => import('@/components/views/tpl/shops/product-tab/Video'),
    },
  }
</script>