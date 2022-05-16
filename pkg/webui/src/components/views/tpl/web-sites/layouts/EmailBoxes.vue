<template>
  <v-data-table
      :headers="[
                  { text: 'Id', align: 'start', value: 'id'},
                  { text: 'box',  value: 'box'},
                  { text: 'Имя отправителя',  value: 'name'},
                  { text: 'Ящик по умолчанию',  align: 'center', value: 'default'},
                  { text: 'Отправка',  align: 'center',  value: 'allowed'},
                  { text: 'Удалить', align: 'center', value: 'actions', sortable: false},
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
      expand-icon="fal fa-angle-down"
      item-key="id"
      show-expand
  >
    <template v-slot:top>
      <v-toolbar class="mb-0 pb-0" color="white" flat>
        <v-toolbar-title class="d-flex flex-column justify-center">
          Почтовые ящики
          <v-btn :block="false" @click="getDataFromApi(true)" class="mt-0" color="secondary" elevation="1" outlined small>
            <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
            Обновить
          </v-btn>
        </v-toolbar-title>
        <v-divider
            class="mx-4"
            inset
            vertical
        ></v-divider>

        <v-toolbar-items class="mt-12">
          <v-text-field :error-messages="!$v.form.name.$error ? null:'Укажите корректное имя'" class="mb-1"
                        label="Имя отправителя"
                        messages="до 32 символов"
                        placeholder="Company name"
                        v-model.trim="$v.form.name.$model"></v-text-field>
          <v-text-field :error-messages="!$v.form.box.$error ? null:'Только [a-zA-Z0-9], 32 символов'" class="ml-4 mb-5"
                        label="Имя почтового ящика"
                        placeholder="info, sale, help"
                        v-model.trim="$v.form.box.$model">

          </v-text-field>
          <v-checkbox label="По-умолчанию" v-model="form.default"></v-checkbox>
          <div>
            <v-btn :block="false" :disabled="(!$v.form.$anyDirty && !$v.form.$error) || $v.form.$anyError" @click="createItem()" class="mt-0 ml-4" color="cyan darken-4" elevation="1"
                   outlined small>
              <v-icon class="ml-1 mr-3" x-small>far fa-save</v-icon>
              Добавить почтовый ящик
            </v-btn>
          </div>
        </v-toolbar-items>
        <v-spacer></v-spacer>
      </v-toolbar>

      <v-divider
          class="mx-2 mt-0"
          inset
          vertical
      ></v-divider>
      <!-- Подтверждение удаления -->
      <v-dialog @keydown.esc="dialog.open = false" max-width="380" v-model="dialog.open">
        <v-card>
          <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

          <v-card-text class="d-flex flex-column">
            <v-list>
              <v-list-item class="px-0 mx-0 mb-0 pb-0">
                <v-list-item-title class="headline">Email Box : &laquo;{{ dialog.box }}@{{ hostname }}&raquo;</v-list-item-title>
              </v-list-item>
              <v-list-item class="px-0 mx-0 mb-0 pb-0">
                <v-list-item-title class="headline">Отправитель : &laquo;{{ dialog.name }}&raquo;</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card-text>

          <v-card-actions>
            <v-spacer></v-spacer>

            <!-- Отменить -->
            <v-btn @click.stop="dialog.open = false" color="green darken-1" small text>
              Отменить
            </v-btn>

            <!-- Удалить -->
            <v-btn @click.stop="deleteItemDialog" color="red darken-1" small text>
              Удалить
            </v-btn>

          </v-card-actions>
        </v-card>
      </v-dialog>
    </template>

    <template v-slot:item.box="{ item }">
      {{ item.box }}@{{ hostname }}
    </template>

    <template v-slot:item.default="{ item }">
      <span :class="item['default'] ? 'green--text' : 'deep-orange--text darken-1'" v-text="item['default'] ? 'да' : 'нет'"></span>
    </template>

    <template v-slot:item.allowed="{ item }">
      <span :class="item['allowed'] ? 'green--text' : 'deep-orange--text darken-1'" v-text="item['allowed'] ? 'разрешена' : 'запрещена'"></span>
    </template>

    <template v-slot:item.actions="{ item }">
      <v-btn @click="openDeleteDialog(item)" depressed small style="background-color: transparent">
        <v-icon class="red--text text--lighten-2" small>
          fas fa-trash
        </v-icon>
      </v-btn>
    </template>

    <template v-slot:expanded-item="{ headers, item }">
      <td :colspan="headers.length">
        <h3 class="mt-4 mb-6">Параметры</h3>
        <v-form
            :ref="'form_'+item.id"
            lazy-validation
            v-model="valid"
        >
          <v-text-field
              :rules="[rules.maxLength32]" class="body-1"
              label="Имя отправителя"
              v-model.trim="item.name"
          />
          <v-switch
              :label="item['default'] ? 'По-умолчанию: да' : 'По умолчанию: нет'"
              :success="item['default']"
              class="mb-6 mt-0"
              dense
              hide-details
              light
              v-model="item['default']"
          ></v-switch>
          <v-text-field
              :rules="[rules.required, rules.box, rules.maxLength32]" class="body-1"
              label="Имя почтового ящика"
              v-model.trim="item.box"
          />
        </v-form>


        <v-card class="d-flex justify-end mb-3" flat style="background-color: transparent;">
          <v-btn @click="openDeleteDialog(item)" color="red darken-3" dark elevation="1" small>
            <v-icon class="ml-1 mr-3" x-small>far fa-user-slash</v-icon>
            Удалить
          </v-btn>

          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData(item)" color="teal darken-3" dark elevation="1" small>
            <v-icon class="ml-1 mr-3" x-small>far fa-save</v-icon>
            Сохранить
          </v-btn>
        </v-card>
      </td>
    </template>

  </v-data-table>
</template>

<script>
import {required, minLength, maxLength, email} from 'vuelidate/lib/validators'

export default {
  data: () => ({
    items: [],
    totalItems: 0,
    options: {},
    search: '',
    loading: true,
    dialog: {open: false},
    form: {
      name: '',
      box: '',
      default: false,
      allowed: true,
    },
    rules: {
      required: value => !!value || 'Required.',
      maxLength32: value => value.length <= 32 || 'Max 32 characters',
      box: value => {
        const pattern = /^[a-zA-Z0-9.]+$/
        return pattern.test(value) || 'Только разрешенные символы a-zA-Z0-9..'
      },
    },
    valid: true,

  }),
  validations: {
    form: {
      name: {
        maxLength: maxLength(32),
      },
      box: {
        required,
        maxLength: maxLength(32),
        allowedCharacters: (name) => /^[a-zA-Z0-9.]+$/.test(name),
      },
    }
  },
  props: ['email_boxes', 'webSiteId', 'hostname'],
  async mounted() {
    await this.getDataFromApi(false)
  },
  methods: {
    getDataFromApi: async function (showNotification) {

      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        webSiteId: this.webSiteId,
      }

      await this.$api.emailBox.getListPagination(payload)
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
    updateItemData: async function (item) {

      // return
      if (!this.$refs["form_" + item.id].validate()) {
        return
      } else {
        this.$refs["form_" + item.id].resetValidation()
      }

      let payload = item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      await this.$api.emailBox.update(payload)
          .then((resp) => {
            if (resp['email_box'] !== undefined) {
              this.item = resp['email_box']
              this.$notify({
                group: 'main',
                title: 'Почтовый ящик обновлен',
                type: 'success'
              });
              this.$refs["form_" + item.id].resetValidation()
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_box]',
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

    },
    openDeleteDialog(item) {
      Object.assign(this.dialog, item)
      Object.assign(this.dialog, {open: true})
    },
    deleteItemDialog: async function () {
      this.dialog.open = false
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialog.id,
      }
      await this.$api.emailBox.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Пользователь обновлен',
              type: 'success'
            });
            Object.assign(this.dialog, {id: null, softDelete: false, open: false})
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error'
            });
          })
    },
    createItem: async function () {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.webSiteId,
      }
      Object.assign(payload, this.form)

      // console.log(payload)
      // return
      await this.$api.emailBox.create(payload)
          .then(async (resp) => {
            if (resp['email_box'] !== undefined) {
              this.items.unshift(resp['email_box'])
              this.form.name = ""
              this.form.box = ""
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_box]',
                type: 'warring'
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error'
            });
          })
    },
  },

}
</script>
