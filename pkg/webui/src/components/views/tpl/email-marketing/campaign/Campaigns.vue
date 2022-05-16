<template>
  <v-container fluid>
    <v-data-table
        :headers="headers"
        :items.sync="items"
        :loading="loading"
        :options.sync="options"
        :search="search"
        :server-items-length="totalItems"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
        class="elevation-1"
        sort-by="id"
        sort-desc
    >
      <template v-slot:top>
        <v-toolbar color="white" flat>
          <v-toolbar-title><h1 style="font-size: 26px;">Email-рассылки</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>
            <v-btn @click="getDataFromApi(true)" color="secondary" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="goToCreateItem" color="cyan darken-4" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-external-link</v-icon>
              Создать кампанию
            </v-btn>
          </v-card>
        </v-toolbar>
        <v-row>
          <v-col class="ml-2" cols="6">
            <v-text-field
                class="mx-2"
                hide-details
                label="Поиск по кампаниям"
                prepend-inner-icon="fal fa-search"
                single-line
                v-model.trim="search">
            </v-text-field>
          </v-col>
        </v-row>
      </template>
      <template v-slot:item.public_id="{ item }">
        {{ item.public_id }}
      </template>

      <template v-slot:item._recipients="{ item }">
        {{ item._recipients }}
      </template>
      <template v-slot:item._open_rate="{ item }">
        {{ item._open_rate.toFixed(2) + '%' }}
      </template>
      <template v-slot:item._clickRate="{ item }">
        n/a
      </template>
      <template v-slot:item._unsubscribe_rate="{ item }">
        {{ item._unsubscribe_rate.toFixed(2) + '%' }}
      </template>
      <template v-slot:item.status="{ item }">
        <span
            :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2"
        >
        {{$api.helper.getHumanStatus(item.status)}}
      </span>
      </template>
      <template v-slot:item.schedule_run="{ item }">
        {{ item['schedule_run'] | moment("HH:MM DD.MM.YYYY") }}
      </template>
      <template v-slot:item.users_segment_id="{ item }">
        {{item['users_segment_id'] > 0 ? item['users_segment'].name : 'n/a'}}
      </template>
      <template v-slot:item.updated_at="{ item }">
        {{ item['updated_at'] | moment("HH:MM DD.MM.YYYY ") }}
      </template>

      <template v-slot:item.actions="{ item }">
        <div class="d-flex justify-space-between" style="width: 40px;">
          <v-icon
              @click="$router.push({name:'email-marketing.campaign.edit', params: {public_id:item['public_id']}})"
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
    <v-dialog @keydown.esc="dialog.open = false" max-width="380" v-model="dialog.open">
      <v-card>
        <v-card-title class="headline">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <span class="mb-1"><strong>Id №</strong> {{ dialog.id }}</span>
          <span><strong>Имя: </strong>{{ dialog.name }}</span>
        </v-card-text>

        <!-- Удалить / Отменить -->
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click.stop="dialog.open = false" color="green darken-1" small text>
            Отменить
          </v-btn>
          <v-btn @click.stop="deleteItemDialog" color="red darken-1" small text>
            Удалить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-container>
</template>

<script>
import email_campaign from "@/api/modules/email-campaign";

export default {
  data: () => ({
    headers: [
      {text: 'ID', align: 'start', value: 'public_id'},
      {text: 'Название серии', align: 'start', value: 'name'},
      {text: 'Время запуска', align: 'start', value: 'schedule_run'},
      {text: 'Сегмент', align: 'center', value: 'users_segment_id'},
      {text: 'В очереди на отправку', align: 'center', value: '_queue', sortable: false},
      {text: 'Отправлено писем', align: 'center', value: '_recipients', sortable: false},
      {text: 'Открытий', align: 'center', value: '_open_rate', sortable: false},
      // { text: 'Кликов', align: 'center', value: '_clickRate', sortable: false},
      {text: 'Отписок', align: 'center', value: '_unsubscribe_rate', sortable: false},
      {text: 'Статус работы', align: 'center', value: 'status'},

      // {text: 'Редактировано', align: 'center', value: 'updated_at'},
      {text: 'Actions', value: 'actions', sortable: false},
    ],
    items: [],
    totalItems: 0,
    options: {},
    search: '',
    loading: true,
    dialog: {
      open: false,
    },
  }),
  async mounted() {
    // await this.getDataFromApi(false)
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
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: searchStr,
      }

      this.$api.emailCampaign.getListPagination(payload)
          .then(resp => {
            if (resp['email_campaigns'] !== undefined) {
              this.items = resp['email_campaigns']
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
                  text: 'Ошибка в ответе сервера: email_campaigns - not found',
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
    searchDataFromApi: _.throttle(function () {
      this.getDataFromApi(false, this.search)
    }, 1400),
    updateLists: async function (showNotification) {
      await this.getDataFromApi(showNotification)
    },
    goToCreateItem: async function () {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        name: "Новая Email-кампания",
        enabled: false,
        executed: false,
        schedule_run: new Date(),
        subject: "Тема письма",
        preview_text: "",
        // users_segment_id: 1,
        // email_template_id: 1,
      }

      await this.$api.emailCampaign.create(payload)
          .then((resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.items.unshift(resp['email_campaign'])
              this.totalItems += 1
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_campaign]',
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

      // return this.$router.push({name:'email-marketing.campaign.create'})
    },

    // ниже переписать в общее диалоговое оконо
    openDeleteDialog(item) {
      Object.assign(this.dialog, item)
      this.dialog.open = true
    },
    deleteItemDialog: async function () {
      this.dialog.open = false

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialog.id,
      }
      await this.$api.emailCampaign.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Рассылка удалена',
              type: 'success'
            });
            Object.assign(this.dialog, {open: false})
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error'
            });
          })
    },
  },

}
</script>