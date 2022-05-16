<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card flat class="rounded-0">

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Почтовые ящики отправления</h2>
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
                  { text: 'Id', align: 'start', value: 'id'},
                  { text: 'box',  value: 'box'},
                  { text: 'Hostname',  value: 'web_site', sortable: false},
                  { text: 'Имя отправителя',  value: 'name'},
                  { text: 'Ящик по умолчанию',  align: 'center', value: 'default'},
                  { text: 'Отправка',  align: 'center',  value: 'allowed'},
                  { text: 'Добавить', align: 'center', value: 'actions', sortable: false},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          :single-expand="false"
          class="elevation-1"
          selectable-key="id"
      >
        <template v-slot:item.box="{ item }">
          {{ item.box }}@{{ item.hostname }}
        </template>

        <template v-slot:item.default="{ item }">
          <span :class="item['default'] ? 'green--text' : 'deep-orange--text darken-1'"
                v-text="item['default'] ? 'да' : 'нет'"></span>
        </template>

        <template v-slot:item.allowed="{ item }">
          <span :class="item['allowed'] ? 'green--text' : 'deep-orange--text darken-1'"
                v-text="item['allowed'] ? 'разрешена' : 'запрещена'"></span>
        </template>

        <template v-slot:item.web_site="{ item }">
          <span v-if="webSites.length > 0">{{webSites.find(el=>el.id === item.web_site_id, item).hostname}}</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn @click="choiceEmailBox(item)" depressed fab small text>
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
    rules: {
      required: value => !!value || 'Required.',
      maxLength32: value => value.length <= 32 || 'Max 32 characters',
      box: value => {
        const pattern = /^[a-zA-Z0-9.]+$/
        return pattern.test(value) || 'Только разрешенные символы a-zA-Z0-9..'
      },
    },
    selected: [],
    webSites: [],
  }),
  props: ['open'],
  async mounted() {
    await Promise.all([
      this.updateWebSites()
    ]).then(()=>this.getDataFromApi(false))
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
      }

      await this.$api.emailBox.getList(payload)
          .then(resp => {
            if (resp['email_boxes'] !== undefined) {
              this.items = resp['email_boxes']
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
    updateWebSites: async function (showNotification) {
      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.webSite.getList(payload)
          .then(resp => {
            if (resp['web_sites'] !== undefined) {
              this.webSites = resp['web_sites']
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
                  text: 'Ошибка в ответе сервера: webSites - not found',
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
          .finally(() => this.loading = false)

    },
    choiceEmailBox: async function (email_box) {
      this.$emit('choice', email_box)
    }
  },

}
</script>
