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
            show-expand
            @update:options="getDataFromApi(false)"
            @update:search="searchDataFromApi"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Вопросы / CallBack</h1></v-toolbar-title>

          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="getDataFromApi(true)"
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

          </v-card>

        </v-toolbar>
        <v-row>
          <v-col cols="6" class="ml-2">
            <v-text-field
                    v-model.trim="search"
                    label="Поиск..."
                    prepend-inner-icon="fal fa-search"
                    class="mx-2"
                    single-line
                    hide-details>
            </v-text-field>
          </v-col>
        </v-row>
      </template>

      <template v-slot:item.user.name="{ item }">
        {{$api.user.getFullName(item.user)}}
      </template>

      <template v-slot:item.status="{ item }">
        <span
            :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2"
        >
        {{$api.helper.getHumanStatus(item.status, 'question')}}
      </span>
      </template>

      <template v-slot:item.expect_an_answer="{ item }">
        <show-item-yes-no :show="item['expect_an_answer']" :bright="false"></show-item-yes-no>
      </template>
      <template v-slot:item.expect_a_callback="{ item }">
        <show-item-yes-no :show="item['expect_a_callback']" :bright="false"></show-item-yes-no>
      </template>

      <template v-slot:item.created_at="{ item }">
        {{item['created_at'] | moment("Do MMMM YYYY, h:mm:ss")  }}
      </template>

      <template v-slot:item.actions="{ item }">

        <v-btn
                small depressed style="background-color: transparent"
                @click="openDeleteDialog(item)">
          <v-icon
                  small
                  class="red--text text--lighten-2"
          >
            fas fa-trash
          </v-icon>
        </v-btn>

      </template>

      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-toolbar flat>
            <v-toolbar-title><h4 class="mt-4 mb-6" style="font-weight: 400;">Параметры Вопроса</h4></v-toolbar-title>
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
              <v-text-field
                      label="Id вопроса" dense class="body-1"
                      :value="item.id"
                      disabled
              />
             <!-- В работе или нет-->
             <work-status v-if="item"
                          :status.sync="item.status"
                          ownerType="question"
                          @change="(e) => setNewStatus(e,item.id)">
             </work-status>

             <v-text-field
                     label="Form name"  class="body-1"
                     v-model.trim="item['form_name']"
             />
             <v-text-field
                     label="From ipv4"  class="body-1" disabled
                     v-model.trim="item['ipv4']"
             />

             <v-text-field
                 label="User Id"  class="body-1" disabled
                 v-model.trim="item.user.id"
             />



            </v-col>

            <v-col cols="7" offset="1">
              <v-textarea
                  label="Message"  class="body-1 ml-0 pl-0" rows="3"
                  counter="255"
                  v-model="item['message']"
              />


            </v-col>
          </v-row>
        </td>
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

  </v-container>
</template>

<script>
  export default {
    data: () => ({
      headers: [
        { text: 'Id', align: 'start', value: 'public_id'},
        { text: 'Name',  value: 'user.name'},
        { text: 'Имя формы',  value: 'form_name'},
        { text: 'Тип запроса',  value: 'question_type'},
        { text: 'Статус',  align: 'center',value: 'status'},
        { text: 'an Answer',  align: 'center', value: 'expect_an_answer'},
        { text: 'call Back',  align: 'center',value: 'expect_a_callback'},
        { text: 'ipv4',  align: 'center',value: 'ipv4'},

        { text: 'Дата создания',  align: 'center',value: 'created_at'},
        { text: 'Удалить', value: 'actions', sortable: false},
      ],
      items: [],
      totalItems: 0,
      options: {},
      search: '',
      loading: true,
      dialog:{
        id: null,
        name: '',
        url: '',
        open: false,
      },
    }),
    async mounted() {
      /*this.loading = true

      await this.getDataFromApi(false)
              .finally(()=>this.loading = false)*/
    },
    watch: {
      search() {
        this.options.page = 1;
        this.searchDataFromApi();
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
          preloads: "User"
        }

        await this.$api.question.getListPagination(payload)
                .then(resp => {
                  if (resp['questions'] !== undefined ) {
                    this.items = resp['questions']
                    this.totalItems = resp.total
                    if (showNotification) {
                      this.$notify({
                        group: 'main',
                        title: 'Данные обновлены',
                        text: 'Список вопросов обновлен',
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
                    title: 'Ошибка получения списка',
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

        let payload = item
        Object.assign(payload, {
          accountHashId: this.$store.state.account.hash_id,
          preloads: "User"
        })

        await this.$api.question.update(payload)
                .then((resp) => {
                  if (resp['web_hook'] !== undefined) {
                    // this.getItemFromItems(item.id)
                    item = resp['web_hook']

                  }
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Карточка обнолвена',
                    type: 'success'
                  });
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

      // ниже переписать в общее диалоговое окно
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.name = item.user.name
        this.dialog.open = true
      },
      deleteItemDialog: async function() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }

        this.$api.question.delete(payload)
                .then(async ()=> {
                  let index = this.items.findIndex(item => item.id === payload.id, payload)
                  if (index !== -1) this.$delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'Вопрос удален',
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
        let index = this.items.findIndex(item => item.id === itemId, itemId)
        if (index !== -1) return this.items[index]
        else return {}
      },
      removeItemFromItems(payload) {
        let index = this.items.findIndex(item => item.id === payload.id, payload)
        if (index !== -1) this.$delete(this.items, index)
      },

      setNewStatus: async function (status, idItem) {
        if (status ==='') {
          return
        }
        await this.$api.question.changeStatus({
          accountHashId: this.$store.state.account.hash_id,
          id: idItem,
          status: status,
          preloads: "EmailTemplate,UsersSegment"
        })
            .then(resp => {
              if (resp['question'] !== undefined) {
                this.item = resp['question']
                this.$notify({
                  group: 'main',
                  title: 'Статус изменен',
                  type: 'success',
                });

              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера не полные',
                  text: 'не хватает объекта [email_campaign]',
                  type: 'warring',
                });
              }
            })
            .catch(err => {
              this.$notify({
                group: 'main',
                title: 'Уведомление не найдено',
                text: err['message'],
                type: 'error',
              });
            })
            .finally(() => this.loading = false)
      },

    },
    components: {
      ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
      WorkStatus: () => import('@/components/views/tpl/layouts/WorkStatus'),
    },
  }
</script>