<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">

    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор пользователя</h2>
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
            { text: 'Id', align: 'start', value: 'account_users.public_id'},
            { text: 'ФИО',  value: 'name'},
            { text: 'Email', value: 'email'},
            { text: 'Телефон', value: 'phone'},
            { text: 'Роль',  value: 'role', sortable: false},
            { text: 'Добавлен', value: 'created_at'},
            { text: 'Действие', align: 'center', value: 'actions', sortable: false},
          ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          :items-per-page="15"
          class="elevation-1"
      >
        <template v-slot:item.account_users.public_id="{ item }">
          {{ item['account_user'] ? item['account_user'].public_id : ''}}
        </template>
        <template v-slot:item.name="{ item }">
          {{$api.user.getFullName(item)}}
        </template>

        <template v-slot:item.role="{ item }">
          {{item['_role'] ? item['_role'].name : (getRoleByRoleId(item['account_user']['role_id']) ? getRoleByRoleId(item['account_user']['role_id']).name : '')}}
        </template>

        <template v-slot:item.created_at="{ item }">
          {{item['created_at'] | moment("Do MMMM YYYY, h:mm:ss")  }}
        </template>

        <template v-slot:item.updated_at="{ item }">
          {{item['updated_at'] | moment("Do MMMM YYYY, h:mm:ss")}}
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
    getDataFromApi: async function (showNotification) {

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.web_site_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        preloads: 'ProductCards,PaymentSubject,ProductTags'
      }

      await this.$api.user.getListPagination(payload)
          .then(resp => {
            if (resp['users'] !== undefined) {
              this.items = resp['users']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список пользователей обновлен',
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
