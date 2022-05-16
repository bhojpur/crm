<template>
  <v-container fluid>
    <v-data-table
        :headers="headers"
        :items="items"
        :loading="loading"
        :options.sync="options"
        :search="search"
        :server-items-length="totalItems"
        :single-expand="false"
        @update:search="searchDataFromApi"
        @update:options="getDataFromApi(false)"
        class="elevation-1"
        expand-icon="fal fa-angle-down"
        show-expand
    >
      <template v-slot:top>
        <v-toolbar color="white" flat>
          <v-toolbar-title><h1 style="font-size: 26px;">Listeners</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="updateLists(true)" color="secondary" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
              Обновить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="dialogCreateItem.open = true" color="cyan darken-4" dark elevation="1" small>
              <v-icon class="ml-1 mr-3" x-small>far fa-plus</v-icon>
              Добавить обработчик
            </v-btn>

          </v-card>

        </v-toolbar>
        <v-row>
          <v-col class="ml-2" cols="6">
            <v-text-field
                class="mx-2"
                hide-details
                label="Поиск по данным клиента"
                prepend-inner-icon="fal fa-search"
                single-line
                v-model.trim="search">
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.event_id="{ item }">
        {{ getEventById(item['event_id']) ? getEventById(item['event_id']).name : '' }}
      </template>

      <template v-slot:item.handler_id="{ item }">
        {{ getHandlerById(item['handler_id']) ? getHandlerById(item['handler_id']).name : '' }}
      </template>

      <template v-slot:item.enabled="{ item }">
        <span :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"
              v-text="item['enabled'] ? 'в работе' : 'остановлен'"></span>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{ item['created_at'] | moment("Do MMMM YYYY") }}
      </template>

      <template v-slot:item.updated_at="{ item }">
        {{ item['updated_at'] | moment("Do MMMM YYYY") }}
      </template>

      <template v-slot:item.actions="{ item }">

        <v-btn @click="openDeleteDialog(item)" depressed fab small text>
          <v-icon class="red--text text--lighten-2" small>fas fa-trash-alt</v-icon>
        </v-btn>
      </template>

      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-toolbar flat>
            <v-toolbar-title><h4 class="mt-4 mb-6" style="font-weight: 400;">Параметры слушателя</h4></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
              <v-btn @click="openDeleteDialog(item)" color="red lighten-1" elevation="0" small dark outlined>
                <v-icon x-small class="ml-1 mr-3">far fa-user-slash</v-icon>
                Удалить
              </v-btn>
              <v-divider class="mx-2 mt-0" inset vertical></v-divider>
              <v-btn @click="updateItemData(item)" dark color="teal darken-1" elevation="0" small outlined>
                <v-icon x-small class="ml-1 mr-3">far fa-save</v-icon>
                Сохранить
              </v-btn>
            </v-card>
          </v-toolbar>
          <v-row>
            <v-col cols="4">

              <v-switch
                  :error="!item['enabled']" :success="item['enabled']" class="mt-1 mb-4" dense
                  hide-details
                  light
                  v-model="item['enabled']"
              >
                <template v-slot:name>
                  <span :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"
                        class="input__label">{{ item['enabled'] ? 'в работе' : 'остановлен' }}</span>
                </template>
              </v-switch>
              <v-text-field label="Имя слушателя" v-model.trim="item['name']"></v-text-field>

              <!-- Выбор события -->
              <v-select
                  label="Выберите событие"
                  :items="events" persistent-hint style="max-width: 380px;"
                  class="mb-4" item-text="name" item-value="id"
                  v-model.number="item['event_id']"
              ></v-select>

              <v-card flat>
                <v-card-text class="teal lighten-5"
                             v-text="getEventById(item['event_id'])!==null ? getEventById(item['event_id']).description : ''"/>
              </v-card>
            </v-col>
            <v-col cols="7" offset="1">

              <!-- Выбор действия -->
              <v-select
                  :hint="getHandlerById(item['handler_id']) ? getHandlerById(item['handler_id']).description : ''"
                  :items="handlerItems"
                  class="mb-4" item-text="name"
                  item-value="id" label="Выберите действие" persistent-hint
                  style="max-width: 380px;"
                  v-model.number="item['handler_id']"
              ></v-select>
              <v-card flat>
                <v-card-text class="teal lighten-5"
                             v-text="getHandlerById(item['handler_id']) ? getHandlerById(item['handler_id']).description : ''"/>
              </v-card>

              <v-select
                  v-if="isWebHookHandler(item)"
                  :hint="getWebHookById(item['entity_id']) ? getWebHookById(item['entity_id']).code : ''"
                  :items="webHooks" class="mt-6 mb-4"
                  item-text="name"
                  item-value="id" label="Выберите WebHook"
                  persistent-hint style="max-width: 380px;"
                  v-model.number="item['entity_id']"
              ></v-select>

              <v-select
                  v-if="isEmailQueueHandler(item)"
                  :hint="getEmailQueueById(item['entity_id']) ? 'Статус запуска: '+ $api.helper.getHumanStatus(item['entity_id'].status) : ''"
                  :items="emailQueues" class="mt-6 mb-4"
                  item-text="name"
                  item-value="id" label="Выберите автоматическую серию"
                  persistent-hint style="max-width: 380px;"
                  v-model.number="item['entity_id']"
              ></v-select>

              <v-btn @click="openDialogChoiceEmailNotification(item)" class="mt-2" color="secondary" outlined small
                     v-if="isEmailNotificationHandler(item)">
                #{{ item['entity_id'] }} - Посмотреть / выбрать другое уведомление
                <v-icon class="ml-4 blue-grey--text text--lighten-2" small>
                  fas fa-file-upload
                </v-icon>
              </v-btn>

            </v-col>
          </v-row>
        </td>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog @keydown.esc="dialog.open = false" max-width="380" v-model="dialog.open">
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.event && dialog.event.name">
              <v-list-item-title class="headline">Event: <strong>{{ dialog.event.name }}</strong></v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.handler && dialog.handler.name">
              <v-list-item-title class="headline">Handler: <strong>{{ dialog.handler.name }}</strong>
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn @click.stop="dialog.open = false" color="green darken-1" small text>Отменить</v-btn>

          <!-- Удалить -->
          <v-btn @click.stop="deleteItemDialog" color="red darken-1" small text>Удалить</v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Создание элемента -->
    <v-dialog @keydown.esc="dialogCreateItem.open = false" max-width="480" v-model="dialogCreateItem.open">
      <v-card>
        <v-card-title style="font-size: 21px;">Создание Event Listener</v-card-title>
        <p class="ml-6">Заполните поле имя и выберите событие и обработчик</p>

        <v-card-text class="d-flex flex-column">
          <v-list>

            <v-list-item class="px-0 mx-0 mb-2 pb-0">
              <v-text-field label="Имя слушателя" v-model.trim="dialogCreateItem['name']"></v-text-field>
            </v-list-item>

            <v-list-item class="px-0 mx-0 pb-0 mb-6">
              <div class="d-flex flex-column justify-start">
                <v-select
                    :items="events"
                    item-text="name"
                    item-value="id" label="Выберите событие"
                    style="max-width: 480px;" v-model.number="dialogCreateItem['event_id']"
                ></v-select>
                <div class="body-1 cyan lighten-5 pa-2 rounded">
                  {{ getEventById(dialogCreateItem['event_id']) !== null ? getEventById(dialogCreateItem['event_id']).label : '' }}
                </div>
              </div>
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-2 pb-0">
              <div class="d-flex flex-column justify-start">
                <v-select
                    :items="handlerItems"
                    item-text="name"
                    item-value="id" label="Выберите обработчик"
                    style="max-width: 480px;" v-model.number="dialogCreateItem['handler_id']"
                ></v-select>
                <div class="body-1 cyan lighten-5 pa-2 rounded">
                  {{ getHandlerById(dialogCreateItem['handler_id']) !== null ? getHandlerById(dialogCreateItem['handler_id']).name : '' }}
                </div>
              </div>
            </v-list-item>

            <v-list-item class="px-0 mt-4 mb-0 pb-0">
              <v-switch
                  :error="!dialogCreateItem['enabled']" :success="dialogCreateItem['enabled']" class="mt-1 mb-6" dense
                  hide-details
                  light
                  v-model="dialogCreateItem['enabled']"
              >
                <template v-slot:label>
                  <span :class="dialogCreateItem['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"
                        class="input__label">{{ dialogCreateItem['enabled'] ? 'в работе' : 'остановлен' }}</span>
                </template>
              </v-switch>
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field
                  class="mx-2"
                  label="Приоритет [-200;+200]"
                  v-model.number="dialogCreateItem['priority']">
              </v-text-field>
            </v-list-item>

          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn @click.native="dialogCreateItem.open = false" color="red darken-1" small text>
            Отменить
          </v-btn>

          <!-- Создать -->
          <v-btn :disabled="!validateNewItem" @click.stop="createItem" color="green darken-1" small text>
            Создать
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Выбор уведомления -->
    <v-dialog @keydown.esc="dialogChoiceEmailNotification = false" fullscreen transition="dialog-bottom-transition"
              v-model="dialogChoiceEmailNotification">
      <v-card>

        <v-toolbar color="primary" dark>
          <v-spacer></v-spacer>
          <v-card-actions>
            <v-btn @click="dialogChoiceEmailNotification = false" dark outlined>
              <v-icon class="mr-3">far fa-times</v-icon>
              Закрыть
            </v-btn>
          </v-card-actions>
        </v-toolbar>

        <v-data-table
            :headers="[
                  { text: 'Id', align: 'start', value: 'id'},
                  { text: 'Name',  value: 'name'},
                  { text: 'Назначение уведомления',  value: 'description'},
                  { text: 'Статус',  value: 'enabled'},
                  { text: 'Выбрать уведомление', align: 'center', value: 'actions', sortable: false},
                ]"
            :items.sync="emailNotificationItems"
            :loading="emailNotificationLoading"
            :options.sync="emailNotificationOptions"
            :search.sync="emailNotificationSearch"
            :server-items-length="emailNotificationTotalItems"
            @update:options="getDataEmailNotificationFromApi(false)"
            @update:search="searchDataEmailNotificationFromApi"
            class="mt-4"
            dense
            sort-by="id"
        >
          <template v-slot:top>
            <v-toolbar color="white" flat>
              <v-toolbar-title><h1 style="font-size: 26px;">Поиск email-уведомлений</h1></v-toolbar-title>
              <v-spacer></v-spacer>
              <v-card class="d-flex align-center" flat>
                <v-btn @click="updateEmailNotificationLists(true)"
                       color="secondary"
                       dark
                       elevation="1"
                       small
                >
                  <v-icon
                      class="ml-1 mr-3"
                      x-small
                  >far fa-sync
                  </v-icon>
                  Обновить
                </v-btn>
              </v-card>

            </v-toolbar>

            <v-row>
              <v-col class="ml-2" cols="6">
                <v-text-field
                    class="mx-2"
                    hide-details
                    label="Поиск по данным"
                    prepend-inner-icon="fal fa-search"
                    single-line
                    v-model="emailNotificationSearch">
                </v-text-field>
              </v-col>
            </v-row>
          </template>

          <template v-slot:item.id="{ item }">
            <span :class="{'green lighten-4': item.id === items.find(el=>el.id === emailNotificationChoiceListenerId).entity_id}"
                  class="px-2 py-1 rounded">{{item.id}}</span>
          </template>

          <template v-slot:item.name="{ item }">
            <span :class="{'green lighten-4': item.id === items.find(el=>el.id === emailNotificationChoiceListenerId).entity_id}"
                  class="px-2 py-1 rounded">{{item.name}}</span>
          </template>

          <template v-slot:item.created_at="{ item }">
            {{ item['created_at'] | moment("Do MMMM YYYY, h:mm:ss") }}
          </template>

          <template v-slot:item.updated_at="{ item }">
            {{ item['updated_at'] | moment("Do MMMM YYYY, h:mm:ss") }}
          </template>

          <template v-slot:item.actions="{ item }">
            <v-btn :disabled="item.id === items.find(el=>el.id === emailNotificationChoiceListenerId).entity_id"
                   @click="choiceEmailNotification(item)" depressed
                   fab small text>
              <v-icon class="green--text text--lighten-2" small>fas fa-file-export</v-icon>
            </v-btn>
          </template>

        </v-data-table>

      </v-card>
    </v-dialog>

  </v-container>
</template>

<script>

export default {
  data: () => ({
    headers: [
      {text: 'Id', align: 'start', value: 'id'},
      {text: 'Имя', value: 'name'},
      {text: 'Событие', value: 'event_id'},
      {text: 'Обработчик', value: 'handler_id'},
      // { text: 'Entity Id',  value: 'entity_id'},
      {text: 'Статус', value: 'enabled'},

      {text: 'Добавлен', value: 'created_at', sortable: false},
      {text: 'Удалить/Исключить', align: 'center', value: 'actions', sortable: false},
    ],
    items: [], //       systemObservers
    totalItems: 0,
    events: [], // Список возможных событий
    handlerItems: [], // Список возможных обработчиков
    webHooks: [], // Список возможных вебхуков для вызова обработчиком
    emailQueues: [], // Список возможных серий писем

    options: {},
    loading: true,
    search: '',
    dialog: {
      id: null,
      softDelete: false, // удалять ли самого пользователя
      open: false,
    },
    dialogCreateItem: {
      name: '',
      event_id: 0,
      handler_id: 0,
      enabled: false,
      priority: 0,

      open: false,
    },

    dialogChoiceEmailNotification: false, // dialog
    emailNotificationItems: [],
    emailNotificationOptions: {},
    emailNotificationTotalItems: 0,
    emailNotificationLoading: false,
    emailNotificationSearch: '',
    emailNotificationChoiceListenerId: null, // id listener Для сохранения
  }),
  async mounted() {
    await Promise.all([
      this.getEventItems(false),
      this.getHandlerItems(false),
      this.getWebHookItems(false),
      this.getEmailQueueItems(false),
    ])

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
        search: this.search,
        // preloads: "Handler,Event"
      }

      await this.$api.eventListener.getListPagination(payload)
          .then(resp => {
            if (resp['event_listeners'] !== undefined) {
              this.items = resp['event_listeners']
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
                  text: 'Нет ответа от серера',
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
    getEventItems: async function (showNotification) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        limit: 120,
        sortBy: "id",
      }

      await this.$api.event.getListPagination(payload)
          .then(resp => {
            if (resp['events'] !== undefined) {
              this.events = resp['events']
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
    },
    getHandlerItems: async function (showNotification) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.eventHandlerItem.getList(payload)
          .then(resp => {
            if (resp['handler_items'] !== undefined) {
              this.handlerItems = resp['handler_items']
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
    },
    getWebHookItems: async function (showNotification) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        // limit: 120,
        // sortBy: "id",
      }

      await this.$api.webHook.getList(payload)
          .then(resp => {

            if (resp['web_hooks'] !== undefined) {
              this.webHooks = resp['web_hooks']
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
    },
    getEmailQueueItems: async function (showNotification) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.emailQueue.getList(payload)
          .then(resp => {

            if (resp['email_queues'] !== undefined) {
              this.emailQueues = resp['email_queues']
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
    },
    searchDataFromApi: _.throttle(function () {
      this.getDataFromApi(true, this.search)
    }, 1400),
    updateLists: async function (showNotification) {
      await this.getDataFromApi(showNotification)
    },
    updateItemData: async function (item) {

      let payload = item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      delete payload['event']
      delete payload['handler']

      await this.$api.eventListener.update(payload)
          .then((resp) => {
            if (resp['event_listener'] !== undefined) {
              let id = resp['event_listener'].id
              let obj = this.items.find(el => el.id === id, id)
              if (obj) obj = resp['event_listener']

              this.$notify({
                group: 'main',
                title: 'EventListener обновлен',
                type: 'success'
              });

              this.updateLists()

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [eventListener]',
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
    createItem: async function () {

      let payload = {}
      Object.assign(payload, this.dialogCreateItem)
      delete payload['open']
      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      await this.$api.eventListener.create(payload)
          .then((resp) => {
            if (resp['event_listener'] !== undefined) {
              this.items.unshift(resp['event_listener'])
              this.totalItems++
              // this.$set(this, 'items', resp['eventListener'])
              this.$notify({
                group: 'main',
                title: 'Event Listener создан',
                type: 'success'
              });
              Object.assign(this.dialogCreateItem, {open: false})
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [event_listener]',
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
      await this.$api.eventListener.delete(payload)
          .then(() => {
            let index = this.items.findIndex(item => item.id === payload.id, payload)
            if (index !== -1) this.$delete(this.items, index)
            this.$notify({
              group: 'main',
              title: 'Слушатель удален',
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

    getHandlerById(id) {
      let el = this.handlerItems.find(v => v.id === id, id)
      return (el === undefined) ? null : el
    },
    getEventById(id) {
      // console.log(id)
      let el = this.events.find(v => v.id === id, id)
      return (el === undefined) ? null : el
    },
    getWebHookById(id) {
      let el = this.webHooks.find(v => v.id === id, id)
      return (el === undefined) ? null : el
    },
    getEmailQueueById(id) {
      let el = this.emailQueues.find(v => v.id === id, id)
      return (el === undefined) ? null : el
    },
    isWebHookHandler(item) {
      let handler = this.handlerItems.find(el=> el.id === item.handler_id)
      if (handler === undefined || handler['entity_type'] === undefined) return false
      return handler.entity_type === "web_hooks"
    },
    isEmailNotificationHandler(item) {
      let handler = this.handlerItems.find(el=> el.id === item.handler_id)
      if (handler === undefined || handler['entity_type'] === undefined) return false
      return handler.entity_type === "email_notification"
    },
    isEmailQueueHandler(item) {
      let handler = this.handlerItems.find(el=> el.id === item.handler_id)
      if (handler === undefined || handler['entity_type'] === undefined) return false
      return handler.entity_type === "email_queue"
    },

    openDialogChoiceEmailNotification(item) {
      this.dialogChoiceEmailNotification = true
      this.emailNotificationChoiceListenerId = item.id
    },
    getDataEmailNotificationFromApi: async function (showNotification, searchStr = '') {

      this.emailNotificationLoading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.emailNotificationOptions.page - 1) * (this.emailNotificationOptions.itemsPerPage),
        limit: this.emailNotificationOptions.itemsPerPage,
        sortBy: typeof this.emailNotificationOptions.sortBy[0] === 'string' ? this.emailNotificationOptions.sortBy[0] : "id",
        sortDesc: this.emailNotificationOptions.sortDesc[0] !== undefined ? this.emailNotificationOptions.sortDesc[0] : false,
        search: searchStr,
      }

      await this.$api.emailNotification.getListPagination(payload)
          .then(resp => {
            if (resp['email_notifications'] !== undefined) {
              this.emailNotificationItems = resp['email_notifications']
              this.emailNotificationTotalItems = resp.total
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
                title: 'Ошибка обновления данных',
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
          .finally(() => this.emailNotificationLoading = false)
    },
    searchDataEmailNotificationFromApi: _.throttle(function () {
      this.getDataEmailNotificationFromApi(true, this.emailNotificationSearch)
    }, 1400),
    updateEmailNotificationLists: async function (showNotification) {
      await this.getDataEmailNotificationFromApi(showNotification)
    },
    choiceEmailNotification: async function (notification) {

      let item = this.items.find(el => el.id === this.emailNotificationChoiceListenerId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error'
        });
        return
      }
      item['entity_id'] = parseInt(notification.id)
      this.dialogChoiceEmailNotification = false
      delete item['email_notification']
      await this.updateItemData(item)
    }
  },
  computed: {
    validateNewItem() {
      let item = this.dialogCreateItem
      if (item['name'].length > 255) return false
      if (item['event_id'] < 1) return false
      if (item['handler_id'] < 1) return false
      return !(item['priority'] < -200 || item['priority'] > 200);

    }
  }
}
</script>