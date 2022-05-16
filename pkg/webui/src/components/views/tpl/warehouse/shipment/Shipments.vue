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
            :search="search"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>
            <h1 style="font-size: 26px;">Поставки</h1>

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
              Новая поставка
            </v-btn>
          </v-card>

        </v-toolbar>

        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по поставкам"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item._payment_amount_order="{ item }">
        <span>{{item['_payment_amount_order'].toLocaleString('ru-RU',{currency:"RUB"})}} руб.</span>
        <span
            v-if="item.status === 'posting' || item.status === 'completed'|| item.status === 'completed'"
            :class="differentAmount < 0 ? 'red--text': 'green--text'">
          ({{MathGetSign(differentAmount)}}{{Math.abs(differentAmount).toLocaleString('ru-RU',{currency:"RUB"})}}<span v-if="differentAmount !== 0"> руб.</span>)
        </span>
      </template>

      <template v-slot:item.status="{ item }">
        <span :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2">
        {{$api.helper.getHumanStatus(item.status, 'shipment')}}
      </span>
      </template>

      <template v-slot:item._product_units="{ item }">
        {{item['_product_units']}} шт.
      </template>

      <template v-slot:item.actions="{ item }">
        <div class="d-flex justify-space-between" style="width: 40px;">
          <v-icon
              @click="$router.push({name:'shipment.edit', params: {public_id:item['public_id']}})"
              class="mr-4 blue--text text--lighten-2"
              small
          >
            fas fa-edit
          </v-icon>
          <v-icon
              @click="openDeleteDialog(item)"
              class="red--text text--lighten-2"
              small
          >
            fas fa-trash
          </v-icon>
        </div>

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
  import ChoiceProduct from '@/components/views/tpl/shops/products/ChoiceProduct'
  export default {
    components: {ChoiceProduct},
    data: () => ({
      headers: [
        { text: 'Id', align: 'start', value: 'public_id'},
        { text: 'Название', value: 'name'},
        { text: 'Ex. code',  value: 'code'},
        // { text: 'GroupId',  value: 'product_group_id'},
        // { text: 'WebPage',  value: 'web_page_id'},
        { text: 'Позиций', align: "center", value: '_product_units'},
        { text: 'Сумма закупки',align: 'center', value: '_payment_amount_order'},
        { text: 'Статус',align: 'center', value: 'status'},
        { text: 'Действия', value: 'actions'},
      ],
      items: [],

      totalItems: 0,
      options: {},
      loading: true,
      search: '',
      dialog:{
        id: null,
        open: false,
      },
    }),
    async mounted() {
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
          webSiteId: this.webSiteId,
        }

        this.$api.shipment.getListPagination(payload)
                .then(resp => {
                  if (resp['shipments'] !== undefined ) {
                    this.items = resp['shipments']
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
                        text: 'Ошибка в ответе сервера: shipments - not found',
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
          name: "Новая поставка",
        }

        await this.$api.shipment.create(payload)
          .then(async (data) => {
            if (data['shipment']) {
              if (!this.items) {
                this.items = []
              }
              this.items.push(data['shipment'])
            }
            await this.getDataFromApi(false)
            this.$notify({
              group: 'main',
              title: 'Поставка создан',
              type: 'success'
            });
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: 'Не удалось создать поставка',
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
        this.$api.shipment.delete(payload)
          .then(async ()=> {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)

            this.$notify({
              group: 'main',
              title: 'Поставка удалена',
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
    computed: {
      differentAmount() {
        if (!this.item || !this.item['_payment_amount_fact'] || !this.item['_payment_amount_order']) return 0
        return this.item['_payment_amount_fact']-this.item['_payment_amount_order']
      }
    }
  }
</script>