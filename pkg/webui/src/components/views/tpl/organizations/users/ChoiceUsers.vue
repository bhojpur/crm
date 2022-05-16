<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card class="rounded-0" flat>

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор пользователей</h2>
        <v-divider
            class="ml-6 mr-4"
            dark
            vertical
        ></v-divider>
        <v-text-field
            class="mx-2 my-0 white--text"
            dark
            dense
            hide-details placeholder="поиск..."
            prepend-inner-icon="fal fa-search"
            single-line
            style="max-width: 380px;"
            v-model.trim="search">
        </v-text-field>
        <v-spacer></v-spacer>
        <span class="mr-4 body-1 grey--text">Выбрано: {{ userSelected.length }}</span>
        <v-btn :disabled="userSelected.length< 1" @click="emitUsers(false)" class="orange--text mr-3" dark elevation="0"
               outlined small>
          <v-icon class="mr-4">far fa-trash-alt</v-icon>
          Убарть из оповещения
        </v-btn>

        <v-btn :disabled="userSelected.length< 1" @click="emitUsers(true)" class="teal--text" dark elevation="0"
               outlined small>
          <v-icon class="mr-4">far fa-save</v-icon>
          Добавить в оповещения
        </v-btn>

        <v-divider
            class="mx-5"
            dark
            vertical
        ></v-divider>

        <v-btn :block="false" @click="getDataFromApi(true)" class="mt-0 blue--text text--lighten-4 mr-3" elevation="0"
               outlined small>
          <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
          Обновить
        </v-btn>

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
                  { text: 'Роль', value: 'role', sortable: false},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          item-key="id"
          selectable-key="id"
          show-select
          v-model="userSelected"
      >
        <template v-slot:item.data-table-select="{ isSelected, select }">
          <v-simple-checkbox :value="isSelected" @input="select($event)" :ripple="false"></v-simple-checkbox>
        </template>

        <template v-slot:item.account_users.public_id="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{ item.account_user.public_id }}</span>
        </template>

        <template v-slot:item.name="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{$api.user.getFullName(item)}}</span>
        </template>

        <template v-slot:item.email="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{item.email  }}</span>
        </template>

        <template v-slot:item.public="{ item }">
          <span :class="item['allowed'] ? 'green--text' : 'deep-orange--text darken-1'"
                v-text="item['public'] ? 'отображен' : 'скрыт'"></span>
        </template>

        <template v-slot:item.created_at="{ item }">
          {{item['created_at'] | moment("Do MMMM YYYY")  }}
        </template>

        <template v-slot:item.role="{ item }">
          {{ getRoleByRoleId(item['account_user']['role_id']).name }}
        </template>

      </v-data-table>

    </v-card>
  </v-dialog>

</template>

<script>
export default {
  data: () => ({
    items: [],
    roles: [],
    totalItems: 0,
    userSelected: [],
    options: {
      sortDesc: [true],
      sortBy:['account_users.public_id']
    },
    search: '',
    loading: true,
    dialog: {open: false},
  }),
  props: {
    open: {
      type: Boolean,
      default: false
    },
    usersIds: {
      type: Array,
      default: []
    }
  },
  async mounted() {
    await Promise.all([
      this.getRoles(),

    ])
        // .then(()=>this.getDataFromApi(false))

  },
  watch: {
    search() {
      this.options.page = 1;
      this.searchDataFromApi();
    },
    open() {
      this.getDataFromApi(false)
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
        roles: ''
      }

      await this.$api.user.getListPagination(payload)
          .then(resp => {
            if (resp['users'] !== undefined) {
              this.items = resp['users']
              this.totalItems = resp.total


              if (this.usersIds && this.usersIds.length > 0) {
                this.userSelected = _.intersectionWith(this.items, this.usersIds, function (obj, other) {
                  let v1, v2
                  if (typeof obj === 'object')
                    v1 = parseInt(obj.id)
                  else
                    v1 = parseInt(obj)

                  if (typeof other === 'object')
                    v2 = parseInt(other.id)
                  else
                    v2 = parseInt(other)

                  return v1 === v2;
                })
              }

              // console.log(this.userSelected)
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
    getRoles: async function (showNotification = false) {
      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id
      }

      await this.$api.userRole.getList(payload)
          .then((resp) => {
            if (resp['roles'] !== undefined) {
              this.roles = resp['roles']
            }
            if (showNotification) {
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });
            }

          })
          .catch((err) => {
            // Если ошибка - указываем пользователю на это
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    emitUsers: function (status) {
      this.$emit('choice', {
        users: this.userSelected.map(el => el.id),
        append: status
      })
      this.userSelected = []
    },
    isAdded(item) {
      if (!item) return false
      return !!_.find(this.usersIds, (el => el === item.id))
    }
  }
}
</script>
