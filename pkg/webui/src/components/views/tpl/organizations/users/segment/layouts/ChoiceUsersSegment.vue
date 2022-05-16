<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор сегмента пользователей</h2>
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
                  { text: 'Имя сегмента',  value: 'name'},
                  { text: 'Строгое соответствеи',  align: 'center', value: 'strict_matching'},
                  { text: 'Добавить', align: 'center', value: 'actions', sortable: false},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          :single-expand="false"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          item-key="id"
          selectable-key="id"
      >
        <template v-slot:item.box="{ item }">
          {{ item.box }}@{{ item.hostname }}
        </template>

        <template v-slot:item.public_id="{ item }">
          {{item.public_id}}
        </template>


        <template v-slot:item.actions="{ item }">
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
    webSites: [],
  }),
  props: ['open'],
  async mounted() {
    // await this.searchDataFromApi(false);
  },
  watch: {
    search() {
      this.options.page = 1;
      this.searchDataFromApi();
    }
  },
  methods: {
    getDataFromApi: async function (showNotification, searchStr = '') {

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: searchStr,
      }

      await this.$api.usersSegment.getListPagination(payload)
          .then(resp => {
            if (resp['users_segments'] !== undefined) {
              this.items = resp['users_segments']
              this.totalItems = resp.total
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список ящиков обновлен',
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
    choiceItem: async function (segment) {
      this.$emit('choice', segment)
    }
  },

}
</script>
