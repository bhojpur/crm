<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор группы Тегов</h2>
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
                  { text: 'Имя группы', value: 'label'},
                  { text: 'Цвет',  align: 'center', value: 'color'},
                  { text: 'Код',  value: 'code'},
                  { text: 'filter_label',  value: 'filter_label'},
                  { text: 'filter_code',  value: 'filter_code'},
                  { text: '_tag_count',  value: '_tag_count'},

                  { text: 'Действия', align: 'center', value: 'actions', sortable: false},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          selectable-key="id"
          single-select
      >

        <template v-slot:item.public_id="{ item }">
          {{item.public_id}}
        </template>

        <template v-slot:item.actions="{ item }">

          <router-link :to="{ name:'email-marketing.templates.edit', params: {public_id:item['public_id']}}" class="mx-2 mb-4" color="teal darken-3" style="text-decoration: none;"
                       target="_blank">
            <v-icon small class="blue--text text--lighten-2">
              fas fa-edit
            </v-icon>
          </router-link>
          <v-btn @click="choiceItem(item)" depressed fab small text>
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
    options: {},
    search: '',
    loading: true,
    dialog: {open: false},
  }),
  props: ['open','choiceId'],
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
        search: searchStr,
      }

      await this.$api.productTagGroup.getListPagination(payload)
          .then(resp => {
            if (resp['product_tag_groups'] !== undefined) {
              this.items = resp['product_tag_groups']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список типов обновлен',
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

}
</script>
