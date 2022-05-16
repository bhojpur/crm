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
            expand-icon="fal fa-angle-down"
            :single-expand="false"
            :show-expand="isMainAccount"
            @update:search="searchDataFromApi"
            @update:options="getDataFromApi(false)"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Список событий</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="updateLists(true)" color="secondary" elevation="1" small dark>
              <v-icon x-small class="ml-1 mr-3">far fa-sync</v-icon>Обновить</v-btn>

            <v-divider class="mx-2 mt-0" inset vertical v-if="isMainAccount"></v-divider>

            <v-btn @click="dialogCreateItem.open = true" color="cyan darken-4" elevation="1" small dark v-if="isMainAccount">
              <v-icon x-small class="ml-1 mr-3">far fa-plus</v-icon>Добавить обработчик</v-btn>

          </v-card>

        </v-toolbar>

        <v-row class="mt-0">
          <v-col cols="6" class="ml-2 mt-0">
            <v-text-field
                    v-model.trim="search"
                    label="Найти событие..."
                    prepend-inner-icon="fal fa-search"
                    class="mx-2 mt-0"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.enabled="{ item }">
        <span v-text="item['enabled'] ? 'в работе' : 'недоступен'" :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'"></span>
      </template>

      <template v-slot:item.description="{ item }">
        <div style="max-width: 540px;">{{item.description}}</div>
      </template>

      <template v-slot:item.external_call_available="{ item }">
        <show-item-yes-no :show="item['external_call_available']"></show-item-yes-no>
      </template>

      <template v-slot:item.external_call_available="{ item }">
        <show-item-yes-no :show="item['external_call_available']"></show-item-yes-no>
      </template>

      <template v-slot:item.actions="{ item }">
        <show-action :route="{name:'event.edit', params: {id:item['id']}}" @delete="openDeleteDialog(item)"></show-action>
      </template>

      <!-- Только для главного аккаунта -->
      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-row>
            <v-col cols="5">
              <v-switch
                      class="mt-1 mb-6" dense hide-details light
                      :success="item['enabled']"
                      :error="!item['enabled']"
                      v-model="item['enabled']"
              >
                <template v-slot:name>
                  <span class="input__label" :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'">{{item['enabled'] ? 'в работе' : 'недоступен'}}</span>
                </template>
              </v-switch>
              <v-text-field v-model.trim="item['name']" name="Имя события" dense class="body-1" />
            </v-col>
            <v-col cols="5" offset="1">
              <v-textarea v-model.trim="item['description']" name="Описание" dense class="body-1" rows="2"/>
            </v-col>
          </v-row>

          <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
           <!-- <v-btn @click="openDeleteDialog(item)" color="red darken-3" elevation="1" small outlined>
              <v-icon x-small class="ml-1 mr-3">far fa-user-slash</v-icon>
              Удалить
            </v-btn>-->

<!--            <v-divider class="mx-2 mt-0" inset vertical></v-divider>-->

            <v-btn @click="updateItemData(item)" outlined color="teal darken-3" elevation="1" small>
              <v-icon x-small class="ml-1 mr-3">far fa-save</v-icon>
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
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-list-item-title class="headline">Id-Name: [{{dialog.id}}] - {{dialog.name}}</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.description">
              <v-list-item-title class="headline">Описание: {{dialog.description}}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="green darken-1" small text @click.stop="dialog.open = false">Отменить</v-btn>

          <!-- Удалить -->
          <v-btn color="red darken-1" small text @click.stop="deleteItemDialog">Удалить</v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Создание элемента -->
    <v-dialog
            v-model="dialogCreateItem.open"
            @keydown.esc="dialogCreateItem.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Создание Observer</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field name="Имя"  class="body-1" v-model.trim="dialogCreateItem['name']" aria-required="true" />
            </v-list-item>

            <v-list-item class="px-0 mx-0 mb-0 pb-0">
              <v-text-field name="Фамилия"  class="body-1" v-model.trim="dialogCreateItem['description']" aria-required="true" />
            </v-list-item>

            <v-switch
                    class="mt-1 mb-6" dense hide-details light
                    :success="dialogCreateItem['enabled']"
                    :error="!dialogCreateItem['enabled']"
                    v-model="dialogCreateItem['enabled']"
            >
              <template v-slot:name>
                <span class="input__label" :class="dialogCreateItem['enabled'] ? 'green--text' : 'deep-orange--text darken-1'">{{dialogCreateItem['enabled'] ? 'в работе' : 'отключен'}}</span>
              </template>
            </v-switch>

          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn color="red darken-1" small text @click.native="dialogCreateItem.open = false">
            Отменить
          </v-btn>

          <!-- Создать -->
          <v-btn color="green darken-1" small text @click.stop="createItem">
            Создать
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-container>
</template>

<script>
  export default {
    data: () => ({
      headers:[
        { text: 'Id', align: 'start', value: 'id'},
        { text: 'Имя',  value: 'name'},
        { text: 'Код',  value: 'code'},
        { text: 'Описание',  value: 'description'},
        { text: 'Вызов по API',  align: 'end', value: 'external_call_available'},
        { text: 'Удалить', align: 'center', value: 'actions', sortable: false}
      ],
      items: [], //       systemObservers
      totalItems: 0,
      options: {},
      loading: true,
      search: '',
      dialog:{
        id: null,
        softDelete: false, // удалять ли самого пользователя
        open: false,
      },
      dialogCreateItem:{
        name: 'NameEvent',
        description: 'Описание Event',
        enabled: false,
        open: false,
      },
    }),
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
      getDataFromApi: async function(showNotification) {

        this.loading = true

        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          offset: (this.options.page-1)*(this.options.itemsPerPage),
          limit: this.options.itemsPerPage,
          sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
          sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
          search: this.search,
        }

        await this.$api.event.getListPagination(payload)
                .then(resp => {
                  if (resp['events'] !== undefined ) {
                    this.items = resp['events']
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
                .finally(()=>this.loading = false)
      },
      searchDataFromApi: _.throttle( function () {
        this.getDataFromApi(true, this.search)
      }, 1400),

      updateLists: async function (showNotification) {
        await this.getDataFromApi(showNotification)
      },
      updateItemData: async function (item) {

        const payload = item

        Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

        await this.$api.eventItem.update(payload)
                .then((resp) => {
                  if (resp['event'] !== undefined) {
                    item = this.items.find(el => el.id === resp['event'].id)
                    if (item !== undefined && item !== null) item = resp['event']
                    this.$notify({
                      group: 'main',
                      title: 'Данные обновлены',
                      type: 'success'
                    });

                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [event]',
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
      createItem: async function (){

        let payload = {}
        Object.assign(payload, this.dialogCreateItem)
        delete payload['open']
        Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

        await this.$api.eventItem.create(payload)
                .then((resp) => {
                  if (resp['event'] !== undefined) {
                    this.items.push(resp['event'])
                    this.$notify({
                      group: 'main',
                      title: 'Объект создан',
                      type: 'success'
                    });
                    Object.assign(this.dialogCreateItem, {open:false})
                  } else {
                    this.$notify({
                      group: 'main',
                      title: 'Данные от сервера не полные',
                      text: 'не хватает объекта [event]',
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
        this.dialog.open = true
      },
      deleteItemDialog: async function() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id,
        }
        await this.$api.eventItem.delete(payload)
                .then(() => {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)
                  this.$notify({
                    group: 'main',
                    title: 'Обсервер удален',
                    type: 'success'
                  });
                  Object.assign(this.dialog, {open: false})
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка удаления',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },

    },
    components: {
      ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
      ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
      ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
      // ProductVideo: () => import('@/components/views/tpl/shops/product-tab/Video'),
    },
  }
</script>