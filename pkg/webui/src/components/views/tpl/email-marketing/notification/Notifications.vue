<template>
  <v-container fluid>

    <v-data-table
            :headers="headers"
            :items="items"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
            :search="search"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Email-уведомления</h1></v-toolbar-title>

          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>
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
                    elevation="1"
                    small
                    dark
            >
              <v-icon x-small class="ml-1 mr-3">far fa-plus</v-icon>
              Добавить Email Notification
            </v-btn>

          </v-card>

        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск по данным"
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.email_template="{ item }">
        <router-link :to="{ name:'email-marketing.templates.edit', params: {public_id:item['email_template']['public_id']}}" class="mx-2 mb-4" color="teal darken-3" style="text-decoration: none;" target="_blank">
          #{{item['email_template'].id}} - {{item['email_template'].name}}
          <v-icon small class="ml-4 blue-grey--text text--lighten-2">
            fas fa-edit
          </v-icon>
        </router-link>

      </template>

      <template v-slot:item.status="{ item }">
        <span
            :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2"
        >
        {{$api.helper.getHumanStatus(item.status)}}
      </span>
      </template>

      <template v-slot:item.actions="{ item }">
        <div style="width: 40px;" class="d-flex justify-space-between">
          <v-icon
              small
              class="mr-4 blue--text text--lighten-2"
              @click="editItem(item)"
          >
            fas fa-edit
          </v-icon>
          <v-icon
              small
              class="red--text text--lighten-2"
              @click="openDeleteDialog(item)"
          >
            fas fa-trash
          </v-icon>
        </div>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog v-model="dialog.open" @keydown.esc="dialog.open = false" max-width="380">
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.id">
              <v-list-item-title class="headline">Id: &laquo;{{dialog.id}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name">
              <v-list-item-title class="headline">Имя : &laquo;{{dialog.name}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="green darken-1" small text @click.stop="dialog.open = false">
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn color="red darken-1" small text @click.stop="deleteItemDialog">
            Удалить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Добавление пользователя -->
    <v-dialog v-model="dialogAddUser" @keydown.esc="dialogAddUser = false" fullscreen  transition="dialog-bottom-transition">
      <v-card>

        <v-toolbar dark color="primary">
          <v-btn icon dark @click="dialogAddUser = false">
            <v-icon >far fa-times</v-icon>
          </v-btn>
          <v-toolbar-title>Settings</v-toolbar-title>
          <v-spacer></v-spacer>
          <span class="mr-8">Выбрано: {{userSelected.length}}</span>
            <v-btn small @click="removeUsersRecipient" class="red darken-3 ">
              <v-icon class="mr-4">far fa-trash-alt</v-icon>
              Убарть из оповещения</v-btn>
          <v-divider
                  class="mx-2 mt-0"
                  inset
                  vertical
          ></v-divider>

          <v-btn small @click="saveAddedUsersRecipient" class="teal darken-3">
            <v-icon class="mr-4">far fa-save</v-icon>
            Добавить в оповещения
          </v-btn>

        </v-toolbar>

        <v-data-table
                :headers="[
                  { text: 'Id', align: 'start', value: 'id'},
                  { text: 'Логин',  value: 'username'},
                  { text: 'Имя / Фамилия',  value: 'name'},
                  { text: 'Email', value: 'email'},
                  { text: 'Телефон', value: 'phone'},
                  { text: 'Роль',  value: 'role', sortable: false},
                ]"
                :items.sync="userItems"
                :options.sync="userOptions"
                show-select
                :server-items-length="userTotalItems"
                :loading="userLoading"
                :search.sync="userSearch"
                v-model="userSelected"
                dense
                sort-by="id"
                class="mt-4"
                @update:options="getDataUserFromApi(false)"
                @update:search="searchDataUserFromApi"
        >
          <template v-slot:top>
            <v-toolbar flat color="white">
              <v-toolbar-title><h1 style="font-size: 26px;">Поиск по пользователям</h1></v-toolbar-title>
              <v-spacer></v-spacer>
              <v-card class="d-flex align-center" flat>
                <v-btn @click="updateUserLists(true)"
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
                        @change="getDataUserFromApi(true)"
                ></v-select>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="6" class="ml-2">
                <v-text-field
                        v-model="userSearch"
                        label="Поиск по данным пользователей"
                        prepend-inner-icon="fal fa-search"
                        class="mx-2"
                        single-line
                        hide-details>
                </v-text-field>
              </v-col>
            </v-row>
          </template>

          <template v-slot:item.id="{ item }">
            <span class="px-2 py-1 rounded" :class="{'green lighten-4': existUserInArr(item)}">{{ item.id }}</span>
          </template>

          <template v-slot:item.username="{ item }">
            <span :class="{'green--text':existUserInArr(item)}">{{ item.username }}</span>
          </template>

          <template v-slot:item.name="{ item }">
            <span :class="{'green--text': existUserInArr(item)}">{{ getFullName(item) }}</span>
          </template>

          <template v-slot:item.role="{ item }">
            {{item['_role'] ? item['_role'].name : getRoleNameByCurrentAccount(item['roles'])}}
          </template>
          <template v-slot:item.created_at="{ item }">
            {{item['created_at'] | moment("Do MMMM YYYY, h:mm:ss")  }}
          </template>

          <template v-slot:item.updated_at="{ item }">
            {{item['updated_at'] | moment("Do MMMM YYYY, h:mm:ss")}}
          </template>

        </v-data-table>

      </v-card>
    </v-dialog>

    <!-- Выбор шаблона письма -->
    <v-dialog v-model="dialogChoiceEmailTemplate" @keydown.esc="dialogChoiceEmailTemplate = false" fullscreen  transition="dialog-bottom-transition">
      <v-card>

        <v-toolbar dark color="primary">
          <v-spacer></v-spacer>
          <v-card-actions>
            <v-btn dark @click="dialogChoiceEmailTemplate = false" outlined>
              <v-icon class="mr-3">far fa-times</v-icon>
              Закрыть
            </v-btn>
          </v-card-actions>


        </v-toolbar>

        <v-data-table
                :headers="[
                  { text: 'Id', align: 'start', value: 'id'},
                  { text: 'Имя шаблона Фамилия',  value: 'name'},
                  { text: 'Описание',  value: 'description'},
                  { text: 'Дата обновления',  value: 'updated_at'},
                  { text: 'Дата создания',  value: 'created_at'},
                  { text: 'Выбрать шаблон', align: 'center', value: 'actions', sortable: false},
                ]"
                :items.sync="templateItems"
                :options.sync="templateOptions"
                :server-items-length="templateTotalItems"
                :loading="templateLoading"
                :search.sync="templateSearch"
                dense
                sort-by="id"
                class="mt-4"
                @update:options="getDataTemplateFromApi(false)"
                @update:search="searchDataTemplateFromApi"
        >
          <template v-slot:top>
            <v-toolbar flat color="white">
              <v-toolbar-title><h1 style="font-size: 26px;">Поиск по пользователям</h1></v-toolbar-title>
              <v-spacer></v-spacer>
              <v-card class="d-flex align-center" flat>
                <v-btn @click="updateTemplateLists(true)"
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
              </v-card>

            </v-toolbar>

            <v-row>
              <v-col cols="6" class="ml-2">
                <v-text-field
                        v-model="templateSearch"
                        label="Поиск по данным"
                        prepend-inner-icon="fal fa-search"
                        class="mx-2"
                        single-line
                        hide-details>
                </v-text-field>
              </v-col>
            </v-row>
          </template>

          <template v-slot:item.id="{ item }">
            <span class="px-2 py-1 rounded" :class="{'green lighten-4': item.id === items.find(el=>el.id === templateChoiceNotificationId).email_template_id}">{{ item.id }}</span>
          </template>

          <template v-slot:item.name="{ item }">
            <span class="px-2 py-1 rounded" :class="{'green lighten-4': item.id === items.find(el=>el.id === templateChoiceNotificationId).email_template_id}">{{ item.name }}</span>
          </template>

          <template v-slot:item.created_at="{ item }">
            {{item['created_at'] | moment("Do MMMM YYYY, h:mm:ss")  }}
          </template>

          <template v-slot:item.updated_at="{ item }">
            {{item['updated_at'] | moment("Do MMMM YYYY, h:mm:ss")}}
          </template>

          <template v-slot:item.actions="{ item }">
            <!-- Переборка по шаблонам item = email_template -->
<!--            <v-btn @click="choiceEmailTemplate(item)" :disabled="item.id === items.find(n=>n.id === templateChoiceNotificationId).email_template_id" small depressed text fab >-->
            <v-btn @click="choiceEmailTemplate(item)" small depressed text fab >
              <v-icon small class="green--text text--lighten-2">fas fa-file-export</v-icon>
            </v-btn>
          </template>

        </v-data-table>

      </v-card>
    </v-dialog>

    <choice-email-box :open="dialogChoiceEmailBox" @close="dialogChoiceEmailBox = false" @choice="choiceEmailBox"></choice-email-box>


  </v-container>
</template>

<script>
  import { required, minLength, maxLength, email} from 'vuelidate/lib/validators'

  export default {
    data: () => ({
      headers: [
        { text: 'Id', align: 'start', value: 'public_id'},
        { text: 'Name',  value: 'name'},
        { text: 'Шаблон письма',  value: 'email_template', sortable: false},
        { text: 'Назначение уведомления',  value: 'description'},
        { text: 'Статус',  value: 'status'},
        { text: 'Actions', align: 'center',value: 'actions', sortable: false},
      ],
      items: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        id: null,
        name: '',
        open: false,
      },
      dialogAddUser: false,
      dialogChoiceEmailTemplate: false,

      recipientAddress: "", // Новый адрес для добавления в список

      /////////

      userItems: [],
      userOptions: {},
      userTotalItems: 0,
      userLoading: false,
      userSearch: '',
      rolesForSearch: [{id:1},{id:2},{id:3},{id:4},{id:5},{id:6},{id:7}],
      userSelected: [],
      addedUsersNotificationId: null, // id notification Для сохранения

      templateItems: [],
      templateOptions: {},
      templateTotalItems: 0,
      templateLoading: false,
      templateSearch: '',
      templateChoiceNotificationId: null, // id notification Для сохранения

      dialogChoiceEmailBox: false,
      emailBoxChoiceNotificationId: null,// id notification Для сохранения

    }),
    validations: {
      recipientAddress: {
        required,
        email,
        maxLength: maxLength(32),
        minLength: minLength(6)
      }
    },
    async mounted() {
      await Promise.all([
        this.getRoleList(false)
      ])
          // .then(()=>this.getDataFromApi(false))

    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
      },
      userSearch() {
        this.searchDataUserFromApi();
      },
    },
    methods: {
      getDataFromApi: async function(showNotification, searchStr = '') {
        this.loading = true
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: searchStr,
          preloads: 'Status,EmailTemplate,EmailBox'
        }

        await this.$api.emailNotification.getListPagination(payload)
                .then(resp => {
                    if (resp['email_notifications'] !== undefined ) {
                      this.items = resp['email_notifications']
                      this.totalItems = resp.total
                      if (showNotification) {
                        this.$notify({
                          group: 'main',
                          title: 'Данные обновлены',
                          text: 'Список уведомлений обновлен',
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
        this.getDataFromApi(true, this.search)
      }, 1400),
      updateItemData: async function (item) {

        if (!this.isEmailBoxSelected(item)) {
          item['enabled'] = false
        }
        let payload = item
        Object.assign(payload, {
          accountHashId: this.$store.state.account.hash_id
        })

        await this.$api.emailNotification.update(payload)
                .then((resp) => {
                  if (resp['email_notification'] !== undefined) {
                    let id = resp['email_notification'].id
                    if (id === undefined) {
                      this.getDataFromApi(true, '')
                      return resolve(resp)
                    } else {
                      let obj = this.items.find(el => el.id === id, id)
                      if (obj) Object.assign(obj, resp['email_notification'])
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        text: 'Карточка обнолвена',
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
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },
      editItem: function(item) {
        this.$router.push({name:'email-marketing.notifications.edit', params: {public_id:item['public_id']}})
      },
      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новое уведомление",
        }
        await this.$api.emailNotification.create(payload)
                .then(async (data) => {
                  if (data['email_notification']) this.items.push(data['email_notification'])
                  await this.getDataFromApi(false) // wtf?
                  this.$notify({
                    group: 'main',
                    title: 'Уведомление создано',
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
      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {
        Object.assign(this.dialog, item)
        this.dialog.open = true
      },
      deleteItemDialog: async function() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }

        this.$api.emailNotification.delete(payload)
                .then(async ()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'Уведомление удалено',
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
                  this.dialog.name = ''
                  this.dialog.url = ''
                })

      },
      getItemFromItems(itemId) {
        if (!this.items || this.items.length < 0) return {}
        let index = this.items.findIndex(item => item.id === itemId, itemId)
        if (index !== -1) return this.items[index]
        else return {}
      },
      removeItemFromItems(payload) {
        let index = this.items.findIndex(item => item.id === payload.id, payload)
        if (index !== -1) this.$delete(this.items, index)
      },
      addRecipient(item) {
        if (this.$v.recipientAddress.$error || (item['recipient_list'] && item['recipient_list'].includes(this.recipientAddress)) || (item['recipient_list'] && item['recipient_list'].length > 10)) {
          return
        }
        if (item['recipient_list'] === null || item['recipient_list'] === undefined) {
          item['recipient_list'] = []
        }
        item['recipient_list'].push(this.recipientAddress)
        this.recipientAddress = ""
        this.$v.recipientAddress.$touch()
      },
      isEmailBoxSelected(item) {
        return !(item['email_box_id'] === null || item['email_box_id'] === undefined)
      },

      removeUsersRecipient: async function() {

        this.dialogAddUser = false

        let notificationId = this.addedUsersNotificationId
        let removeListUsers = this.userSelected.map(el => el.id) // список id

        let item = this.items.find(el=>el.id === notificationId)
        if (!item) {
          this.$notify({
            group: 'main',
            title: 'Ошибка обновления',
            type: 'error'
          });
          return
        }
        item['recipient_users_list'] = _.difference(item['recipient_users_list'],removeListUsers);
        await this.updateItemData(item)
        this.userSelected = []
      },
      saveAddedUsersRecipient: async function() {

        this.dialogAddUser = false

        let notificationId = this.addedUsersNotificationId
        let newListUsers = this.userSelected.map(el => el.id) // список id

        let item = this.items.find(el=>el.id === notificationId)
        if (!item) {
          this.$notify({
            group: 'main',
            title: 'Ошибка обновления',
            type: 'error'
          });
          return
        }
        item['recipient_users_list'] = _.union(item['recipient_users_list'],newListUsers);
        await this.updateItemData(item)
        this.userSelected = []
      },
      openDialogAddUserToNotification(item) {
        this.dialogAddUser = true
        this.addedUsersNotificationId = item.id
      },
      getChipUser(user){
        let name = ""
        if (user['name'])
          name = user['name'] + user['surname'] ? " " + user['surname'] : ""
        else name = user['username']

        return user['email'] !== "" ? name + " - " +user['email'] : name;
      },
      getDataUserFromApi: async function(showNotification, searchStr = '') {
        this.userLoading = true

        let rolesArr = [] // default
        if (this.rolesForSearch && this.rolesForSearch.length > 0) {
          rolesArr = this.rolesForSearch.map((item)=>{
            return item['id']
          })
          rolesArr = _.join(rolesArr, ',')
        }

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.userOptions.page-1)*(this.userOptions.itemsPerPage),
          limit: this.userOptions.itemsPerPage,
          sortBy: (this.userOptions.sortBy !==undefined && typeof this.userOptions.sortBy[0] === 'string') ? this.userOptions.sortBy[0] : "id",
          sortDesc: this.userOptions.sortDesc[0] !== undefined ? this.userOptions.sortDesc[0] : false,
          search: searchStr,
          roles: rolesArr
        }

        await this.$api.user.getListPagination(payload)
                .then(resp => {
                  this.userItems = resp.users
                  this.userTotalItems = resp.total
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
                .finally(()=>this.userLoading = false)
      },
      searchDataUserFromApi: _.throttle( function () {
        this.getDataUserFromApi(true, this.userSearch)
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
                    text: err['mesage'],
                    type: 'error'
                  });
                })
      },
      updateUserLists: async function (showNotification) {
        await this.getDataUserFromApi(showNotification)
      },
      getFullName(item) {

        let fillName = ""
        if (item['name'] !== undefined && item['surname'] !== undefined) {
          fillName = item['name'] + " " + item['surname']
        } else {
          if (item['name'] !== undefined) {
            fillName = item['name']
          } else {
            fillName = "-"
          }
        }
        return fillName
      },
      getRoleNameByCurrentAccount(roles) {
        if (!roles || roles.length < 1) return ''
        let role = roles.find(el=>el.account_id === 1 || el.account_id === this.$store.state.account.id)
        if (!role) return ''
        return role.name
      },
      existUserInArr(user) {
        let item = this.items.find(el=>el.id === this.addedUsersNotificationId)
        if (!item || item['recipient_users_list']===undefined) return

        let rec = item['recipient_users_list']
        return rec.findIndex(el => el === user.id) !== -1
      },

      openDialogChoiceEmailTemplate(item) {

        this.dialogChoiceEmailTemplate = true
        this.templateChoiceNotificationId = item.id
      },
      getDataTemplateFromApi: async function(showNotification, searchStr = '') {

        this.templateLoading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.templateOptions.page-1)*(this.templateOptions.itemsPerPage),
          limit: this.templateOptions.itemsPerPage,
          sortBy: typeof this.templateOptions.sortBy[0] === 'string' ? this.templateOptions.sortBy[0] : "id",
          sortDesc: this.templateOptions.sortDesc[0] !== undefined ? this.templateOptions.sortDesc[0] : false,
          search: searchStr,
        }

        await this.$api.emailTemplate.getListPagination(payload)
                .then(resp => {
                  this.templateItems = resp.email_templates
                  this.templateTotalItems = resp.total
                  if (showNotification) {
                    this.$notify({
                      group: 'main',
                      title: 'Данные обновлены',
                      text: 'Список шаблонов обновлен',
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
                .finally(()=>this.templateLoading = false)
      },
      searchDataTemplateFromApi: _.throttle( function () {
        this.getDataTemplateFromApi(true, this.templateSearch)
      }, 1400),
      updateTemplateLists: async function (showNotification) {
        await this.getDataTemplateFromApi(showNotification)
      },
      choiceEmailTemplate: async function(tpl) {

        // поиск по notification
        let item = this.items.find(el=>el.id === this.templateChoiceNotificationId)
        if (!item) {
          this.$notify({
            group: 'main',
            title: 'Ошибка обновления',
            type: 'error'
          });
          return
        }
        let index = tpl['id']
        let _item = {}
        Object.assign(_item,item)
        delete _item['email_template_id']

        _item['email_template_id'] = parseInt(index)
        this.dialogChoiceEmailTemplate = false
        delete _item['email_template']
        await this.updateItemData(_item)
      },

      openDialogChoiceEmailBox(item) {

        this.dialogChoiceEmailBox = true
        this.emailBoxChoiceNotificationId = item.id
      },
      choiceEmailBox: async function(email_box) {

        this.dialogChoiceEmailBox=false

        // поиск по notification
        let item = this.items.find(el=>el.id === this.emailBoxChoiceNotificationId)
        if (!item) {
          this.$notify({
            group: 'main',
            title: 'Ошибка обновления',
            type: 'error'
          });
          return
        }

        /*item['email_box_id'] = parseInt(email_box['id'])
        delete item['email_box']*/

        let _item = {}
        Object.assign(_item,item)
        delete _item['email_box']

        _item['email_box_id'] = parseInt(email_box['id'])

        await this.updateItemData(_item)
      },

    },
    computed: {
      listHumanRoles() {
        return this.roles.filter(item => {
          return item.type === 'gui'
        })
      },
    },
    components: {
      ChoiceEmailBox: () => import('@/components/views/tpl/web-sites/layouts/ChoiceEmailBox'),
    },
  }
</script>