<template>
  <v-container fluid>
    <v-data-table
            :headers.sync="headers"
            :items.sync="items"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
            :search="search"
            :items-per-page="15"
            expand-icon="fal fa-angle-down"
            show-expand
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Пользователи</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="dialogImportUsers = true"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                  x-small
                  class="ml-1 mr-3"
              >far fa-upload</v-icon>
              Импорт
            </v-btn>

            <v-divider
                class="mx-2 mt-0"
                inset
                vertical
            ></v-divider>

            <v-btn @click="updateLists(true)"
                   color="secondary"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>

            <v-btn @click="dialogCreateItem.open = true" color="cyan darken-4" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-plus</v-icon>
              Добавить пользователя
            </v-btn>

          </v-card>

        </v-toolbar>

        <v-row flat class="ml-2">
          <v-col cols="6">
            <v-select
                    v-model="rolesForSearch"
                    :items="listHumanRoles"
                    item-value="id" item-text="name"
                    return-object
                    chips
                    label="Роли пользователей"
                    multiple
                    @change="getDataFromApi(true)"
            ></v-select>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по данным клиента"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

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

        <v-btn @click="openDeleteDialog(item)" small depressed text fab :disabled="isIssuerAccount(item)">
          <v-icon small class="red--text text--lighten-2">fas fa-user-times</v-icon>
        </v-btn>
        <v-btn @click="openDeleteDialog(item)" small depressed text fab :disabled="!isIssuerAccount(item)">
          <v-icon small class="red--text text--lighten-2">fas fa-trash-alt</v-icon>
        </v-btn>


      </template>

      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-row>
            <v-col cols="3">
              <v-text-field
                      label="username" dense class="body-1"
                      :disabled="isDisableEditItem(item)"
                      v-model.trim="item['username']"
              />
              <v-text-field
                      label="email" dense class="body-1"
                      :disabled="isDisableEditItem(item)"
                      v-model.trim="item['email']"
              />
              <v-text-field
                      label="phone" dense class="body-1"
                      :disabled="isDisableEditItem(item)"
                      v-model.trim="item['phone']" />
            </v-col>
            <v-col cols="3" offset="1">
              <v-text-field
                      label="name" dense class="body-1"
                      v-model.trim="item['name']"
              />
              <v-text-field
                      label="surname" dense class="body-1"
                      v-model.trim="item['surname']"
              />
              <v-text-field
                      label="patronymic" dense class="body-1"
                      v-model.trim="item['patronymic']"
              />
            </v-col>

            <v-col cols="3" offset="1">
              <v-select
                      v-if="!isOwnerUser(item)"
                      :disabled="isOwnerUser(item)"
                      label="Роль пользователя" dense
                      :items="listHumanRoles.filter(v => v.tag !== 'owner')"
                      item-text="name" return-object
                      :value="getRoleByRoleId(item['account_user']['role_id'])"
                      @input="$set(item, '_role',$event)"
                      style="max-width: 350px;"
              ></v-select>
              <v-text-field v-else disabled value="Владелец аккаунта"></v-text-field>

              <v-switch
                  class="mt-1 mb-6" dense hide-details light
                  label="Статус подписки"
                  v-model="item['subscribed']"
              ></v-switch>

            </v-col>

          </v-row>

          <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
            <!--<v-btn @click="openDeleteDialog(item)"
                   color="red darken-3"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-user-slash</v-icon>
              Удалить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>-->

            <v-btn
                    @click="updateItemData(item)"
                    dark outlined
                    color="teal darken-3"
                    elevation="1"
                    small
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-save</v-icon>
              Сохранить
            </v-btn>
          </v-card>
        </td>
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
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name || dialog.surname">
              <v-list-item-title class="headline">Имя пользователя: &laquo;{{dialog.name}} {{dialog.surname}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.username">
              <v-list-item-title class="headline">Логин: &laquo;{{dialog.username}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.email">
              <v-list-item-title class="headline">Email: &laquo;{{dialog.email}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.phone">
              <v-list-item-title class="headline">Phone: &laquo;{{dialog.phone}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-checkbox v-model="dialog.softDelete" label="Удалить аккаунт пользователя?" :disabled="!isIssuerAccount(dialog)" class="mt-1"></v-checkbox>
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

    <!-- Создание элемента -->
    <v-dialog
            v-model="dialogCreateItem.open"
            @keydown.esc="dialogCreateItem.open = false"
            max-width="440"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Создание пользователя</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Username"  class="body-1" v-model.trim="dialogCreateItem['username']" aria-required="true" />
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Password"  class="body-1" v-model.trim="dialogCreateItem['password']" aria-required="true" />
            </v-list-item>
           <!-- <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-checkbox
                      v-model="dialogCreateItem['enabledAuthFromApp']"
                      label="Разрешить вход через web-интерфейс BhojpurCRM"
                      class="mt-0"
              ></v-checkbox>
            </v-list-item>-->

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Имя"  class="body-1" v-model.trim="dialogCreateItem['name']" aria-required="true" />
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Фамилия"  class="body-1" v-model.trim="dialogCreateItem['surname']" aria-required="true" />
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field label="Email"  class="body-1" v-model.trim="dialogCreateItem['email']" aria-required="true" />
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-list-item-action>
                <v-select
                        label="Роль пользователя"
                        :items="listHumanRoles"
                        item-value="id" item-text="name"
                        dense
                        v-model="dialogCreateItem['role_id']"
                        style="max-width: 350px"
                ></v-select>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="red darken-1" small text @click.native="dialogCreateItem.open = false">
            Отменить
          </v-btn>

          <!-- Создать -->
          <v-btn color="green darken-1" small text @click.stop="createItem"
                 :disabled="!checkNewUserData(dialogCreateItem)">
            Создать
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Импорт подписчиков -->
    <import-users :open="dialogImportUsers" @import="importUsers" @close="dialogImportUsers = false"></import-users>

  </v-container>
</template>

<script>
  import draggable from "vuedraggable";

  export default {
    data: () => ({
      headers: [
        { text: 'Id', align: 'start', value: 'account_users.public_id'},
        { text: 'ФИО',  value: 'name'},
        { text: 'Email', value: 'email'},
        { text: 'Телефон', value: 'phone'},
        { text: 'Роль',  value: 'role', sortable: false},
        { text: 'Добавлен', value: 'created_at'},
        { text: 'Удалить/Исключить', align: 'center', value: 'actions', sortable: false},
      ],
      items: [],
      totalItems: 0,
      options: {
        sortDesc: [true],
        sortBy:['account_users.public_id'],
        itemsPerPage: 15,
      },
      loading: true,
      search: '',
      roles: [],
      rolesForSearch: [], // clients
      dialog:{
        id: null,
        softDelete: false, // удалять ли самого пользователя
        open: false,
      },
      dialogCreateItem:{
        role_id: 0,
        username: null,
        email: null, //'mail-test@bhojpur.net',
        password: null,
        enabledAuthFromApp: true,
        open: false,
      },
      dialogImportUsers: false,
    }),
    async mounted() {
      this.loading = true
      await Promise.all([
        this.getRoleList(false)
      ])
          // .then(()=> this.getDataFromApi(false))

    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      }
    },
    methods: {
      getDataFromApi: async function(showNotification, searchStr = '') {

        this.loading = true

        let rolesArr = [] // default
        if (this.rolesForSearch && this.rolesForSearch.length > 0) {
          rolesArr = this.rolesForSearch.map((item)=>{
            return item['id']
          })
          rolesArr = _.join(rolesArr, ',')
        }

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "account_users.public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
          roles: rolesArr.length > 0 ? rolesArr : null
        }

        await this.$api.user.getListPagination(payload)
                .then(resp => {
                  this.items = resp.users
                  this.totalItems = resp.total
                  if (showNotification) {
                    this.$notify({
                      group: 'main',
                      title: 'Данные обновлены',
                      text: 'Список пользователей обновлен',
                      type: 'success'
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
                .finally(()=>this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(true, this.search)
      }, 1400),
      getRoleList: async function (showNotification) {
        const payload = {
          accountHashId: this.$store.state.account.hash_id
        }

        await this.$api.userRole.getListPagination(payload)
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
                .catch( (err) => {
                  // Если ошибка - указываем пользователю на это
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      updateLists: async function (showNotification) {
        await this.getDataFromApi(showNotification)
      },
      updateItemData: async function (item) {

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: item.id
        }
        payload = Object.assign(payload, item)

        await this.$api.user.update(payload)
                .then((resp) => {
                  if (resp['user'] !== undefined) {
                    let v = this.items.find(el => el.id === resp['user'].id)
                    if (v === undefined || v === null) return
                    delete v['_role']
                    Object.assign(v, resp['user'])

                    this.$notify({
                      group: 'main',
                      title: 'Пользователь обновлен',
                      type: 'success'
                    });

                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [user]',
                      type: 'warring'
                    });
                  }
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка сохранения',
                    text: err['message'],
                    type: 'error'
                  });
                })

      },
      createItem: function (){

        let payload = {}
        Object.assign(payload, this.dialogCreateItem)
        delete payload['open']
        Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

        this.$api.user.create(payload)
                .then((resp) => {
                  if (resp['user'] !== undefined) {
                    this.items.unshift(resp['user'])
                    // this.$set(this,'items', resp['user'])
                    this.$notify({
                      group: 'main',
                      title: 'Пользователь создан',
                      type: 'success'
                    });
                    Object.assign(this.dialogCreateItem, {open:false})
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [user]',
                      type: 'warring'
                    });
                  }
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
      openDeleteDialog(item) {
        Object.assign(this.dialog, item)

        Object.assign(this.dialog, { open: true,  softDelete: this.isIssuerAccount(item)})
      },
      deleteItemDialog: async function() {

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id,
          softDelete: this.dialog.softDelete
        }
        await this.$api.user.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Пользователь обновлен',
              type: 'success'
            });
            Object.assign(this.dialog, {id: null,softDelete: false,open: false})
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=>this.dialog.open = false)
      },
      isDisableEditItem(item) {
        return item.role_id < 3
      },
      isIssuerAccount(item) {
        return item['issuerAccountId'] === this.$store.state.account.id
      },
      isOwnerUser(user) {
        if (user['account_user']['role_id'] === undefined ) return false
        let role_id = user['account_user']['role_id']
        return role_id === 1
      },
      checkNewUserData(data) {
        if (( (data['username'] === '' || data['username'] === undefined) &&
                (data['email'] === '' || data['email'] === undefined) &&
                (data['phone'] === '' || data['phone'] === undefined))) return false;

        return parseInt(data['role_id']) >= 1;
      },
      getRoleNameByCurrentAccount(roles) {
        // !!! roleAccountId === , т.к. это системные (как прваило) роли
        if (!roles || roles.length < 1) return ''
        let role = roles.find(el=>el.account_id === 1 || el.account_id === this.$store.state.account.id)
        if (role === undefined) return ''
        console.log("this.$store.state.account.id: ", this.$store.state.account.id, role.account_id)
        // console.log(role)
        return role.name
      },
      importUsers(importData){
        let users = []
        for (let i =0; i < importData.length;i++) {
          let user = importData[i]
          if (!user) continue
          if(user[0] && user[1]) {
            users.push({
              email:user[0].toLowerCase(),
              name: user[1][0].toUpperCase() + user[1].slice(1)
            })
          }
        }

        if (users.length < 1) return
        if (users.length > 1000) return
        // console.log(users)

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          users: users,
          role_id: 7
        }
        this.$api.user.upload(payload)
            .then((resp) => {
              this.$notify({
                group: 'main',
                title: 'Пользователи добавлены',
                type: 'success'
              });
              this.getDataFromApi(true)
            })
            .catch( (err) => {
              this.$notify({
                group: 'main',
                title: 'Ошибка создания',
                text: err['message'],
                type: 'error'
              });
            }).finally(()=> this.dialogImportUsers = false)
        // console.log("Import: ", users)
      }

    },
    computed: {
      listHumanRoles() {
        return this.roles.filter(item => {
          return item.type === 'gui'
        })
      },
    },
    components: {
      draggable,
      ImportUsers: () => import('@/components/views/tpl/organizations/users/layouts/DialogImportUsers'),
    },

  }
</script>