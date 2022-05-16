<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Автоматическая серия
            <v-chip color="grey lighten-3" v-if="item">№ {{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'email-marketing.queues'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку серий</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="!loading">

          <!--<template v-if="$store.state.account.id === 7">
            <div class="d-flex align-center">
              <span v-if="countBro" class="mr-2 blue-grey&#45;&#45;text">добавлено: {{countBro}}</span>
              <v-btn @click="specBro" color="red" small class="white&#45;&#45;text">
                <v-icon class="ml-1 mr-3" small>far fa-upload</v-icon>
                spec. BroUser
              </v-btn>
            </div>
            <v-divider class="mx-2 mt-0" inset vertical></v-divider>
          </template>-->
          <!-- Обновление и сохранение серии -->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" color="cyan darken-4" dark small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>

      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">
          <v-col cols="5">
            <!-- Имя серии писем -->
            <v-text-field
                :rules="[v => !!v || 'Укажите имя серии']"
                class="mt-0"
                label="Имя серии писем"
                required
                v-model="item['name']"
            ></v-text-field>

          </v-col>
          <v-col cols="6" offset="1">

            <!-- В работе или нет-->
            <work-status v-if="item"
                         :status.sync="item.status"
                         :canRun.sync="errors.length === 0"
                         :errors.sync="errors"
                         ownerType="email_queue"
                         @change="setNewStatusNotification">
            </work-status>

          </v-col>
        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>

      <div class="d-flex flex-column mb-2" v-if="item">
        <div class="queue-stat">
          <div class="item">
            <span>В очереди</span>
            <span>{{item._queue}}</span>
          </div>
          <div class="item">
            <span>Завершено</span>
            <span>{{item._completed}}</span>
          </div>
          <div class="item">
            <span>Отправлено писем</span>
            <span>{{item._recipients}}</span>
          </div>
          <div class="item">
            <span>Открытий</span>
            <span>{{item._open_rate.toFixed(2)+'%'}}</span>
          </div>
          <div class="item">
            <span>Отписок</span>
            <span>{{item._unsubscribe_rate.toFixed(2)+'%'}}</span>
          </div>
          <div class="item">
            <span>Кликов</span>
            <span>n/a</span>
          </div>
        </div>
      </div>

      <!-- Письма в серии -->
      <v-data-table
          :headers="[
            { text: 'Порядок отправки', align: 'start', value: 'step'},
            { text: 'Задержка', align: 'start', value: 'delay_time'},
            { text: 'Шаблон',  value: 'email_template.name'},
            { text: 'Тема письма',  value: 'email_template.subject'},
            { text: 'Preview text',  value: 'email_template.preview_text'},
            { text: 'Статус работы',  value: 'enabled'},
            { text: 'From email box',  align: 'start', value: 'email_box_id'},
            { text: 'Удалить', align: 'center', value: 'actions'},
          ]"
          :items="emailQueueEmailTemplates"
          :show-select="false"
          expand-icon="fal fa-angle-down"
          show-expand
          sort-by="step"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h3 style="font-size: 26px;">Цепочка писем</h3></v-toolbar-title>

            <v-spacer></v-spacer>
            <v-card class="d-flex align-center" flat>

              <v-btn
                  @click="dialogChoiceEmailTemplate = true"
                  color="secondary" elevation="1" small>
                <v-icon class="mr-3 ml-2" small>fal fa-plus-circle</v-icon>
                Добавить письмо
              </v-btn>
              <v-divider
                  class="mx-2 mt-0"
                  inset
                  vertical
              ></v-divider>
              <v-btn
                  :disabled="!carried"
                  @click="saveNewOrders"
                  color="primary" elevation="1" small>
                <v-icon class="mr-3 ml-2" small>fal fa-save</v-icon>
                Сохранить порядок
              </v-btn>
            </v-card>

          </v-toolbar>
        </template>

        <template v-slot:body="props" v-if="emailQueueEmailTemplates && emailQueueEmailTemplates.length > 0">
          <draggable
              :list="emailQueueEmailTemplates"
              @end="carried = true"
              handle=".handle"
              tag="tbody"
          >
            <tr :key="item.id" v-for="(item, index) in props.items">
              <td>
                <v-icon class="handle" small>
                  fas fa-align-justify
                </v-icon>
              </td>
              <td><strong>№ {{ item.step }}</strong></td>
              <td><div @click="openDialogSetDelayTime(item)" class="item-editable-popup">{{getDelay(item.delay_time)}}</div></td>
              <td>
                <div @click="changeEmailTemplateForStep(item)" class="item-editable-popup" :class="{'red--text': !item.email_template_id}">
                  <v-icon x-small class="mr-2">
                    far fa-newspaper
                  </v-icon>
                  {{getEmailTemplateName(item)}}
                </div>
              </td>
              <td>
                <div class="item-editable-popup" @click="openDialogChangeSubject(item)">{{ item.subject }}</div>
              </td>
              <td>{{ item.preview_text ? item.preview_text.slice(0,50) : '-' }}{{item.preview_text && item.preview_text.length > 50 ? '...' : ''}}</td>
              <td>
                <span :class="item.enabled ? 'green--text' : 'orange--text'" class="item-editable-popup" @click="changeStatusStep(item)">
                    <v-icon :class="item.enabled ? 'green--text' : 'orange--text'" class="mr-1"
                      small
                      v-text="item.enabled ? 'fal fa-mail-bulk' : 'fal fa-mail-bulk'">

                    </v-icon>
                  {{ item.enabled ? 'в работе' : 'остановлен' }}
                </span>
              </td>
              <td>
                <div @click="openDialogEmailBox(item)" class="d-flex align-center item-editable-popup" :class="{'red--text': !item.email_box_id}">
                  <v-icon x-small class="mr-2">
                    far fa-mailbox
                  </v-icon>
                  {{getEmailBoxName(item)}}
                </div>
              </td>
              <td align="center">
                <div class="d-flex justify-space-between" style="width: 40px;margin-right: 12px;">
                  <v-icon
                      disabled
                      @click=""
                      class="mr-4 blue--text text--lighten-2"
                      small
                  >
                    fas fa-edit
                  </v-icon>
                  <v-icon
                      @click="removeEmailTemplate(item)"
                      class="red--text text--lighten-2"
                      small
                  >
                    fas fa-trash
                  </v-icon>
                </div>
              </td>
            </tr>
          </draggable>
        </template>

      </v-data-table>

      <v-divider
          class="mx-2 mt-2"
          inset
      ></v-divider>

      <!-- История отправки -->
      <v-data-table
          :headers="headers"
          :items.sync="mta_histories"
          :options.sync="options"
          :server-items-length="totalItems"
          class="elevation-1 mt-6"
          :items-per-page="15"
          sort-by="id" :sort-desc="true"
          :search="search"
          @update:search="searchMTAHistoriesFromApi"
          @update:options="getMTAHistoriesFromApi(false)"
      >
        <template v-slot:top>
          <v-toolbar color="white" flat>
            <v-toolbar-title><h1 style="font-size: 26px;">История отправки</h1></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card class="d-flex align-center" flat>
              <v-btn @click="getMTAHistoriesFromApi(true)"
                     color="secondary"
                     dark
                     elevation="1"
                     small
              >
                <v-icon
                    class="ml-1 mr-3"
                    x-small
                >far fa-sync</v-icon>
                Обновить
              </v-btn>
            </v-card>

          </v-toolbar>
          <v-row>
            <v-col class="ml-2" cols="6">
              <v-text-field
                  class="mx-2"
                  hide-details
                  label="Поиск по данным..."
                  prepend-inner-icon="fal fa-search"
                  single-line
                  v-model.trim="search">
              </v-text-field>
            </v-col>
          </v-row>
        </template>

        <template v-slot:item.user.username="{ item }">
          {{item.user.username > 0 ? item.user.username.length  : '-'}}
        </template>
        <template v-slot:item.user.name="{ item }">
          {{$api.user.getFullName(item.user)}}
        </template>

        <template v-slot:item.unsubscribed="{ item }">
          <show-item-yes-no :show="item['unsubscribed']" :bright="false"></show-item-yes-no>
        </template>

        <template v-slot:item.queue_completed="{ item }">
          <show-item-yes-no :show="item['queue_completed']"></show-item-yes-no>
        </template>

        <template v-slot:item.email_template.name="{ item }">
         {{item.email_template.name}}
        </template>

        <template v-slot:item.opens="{ item }">
          <span
              :class="item.opens > 0 ? 'blue--text': 'pink--text text--lighten-1'"
              v-text=" item.opens > 0  ? item.opens : '0'"></span>
        </template>

        <template v-slot:item.unsubscribed_at="{ item }">
          <span class="body-2">{{ item['unsubscribed_at'] ? getMomentJSData(item['unsubscribed_at']) : '-' }}</span>
        </template>

        <template v-slot:item.opened_at="{ item }">
          <span class="body-2">{{ item['opened_at'] ? getMomentJSData(item['opened_at']) : '-' }}</span>
        </template>

      </v-data-table>

    </v-card>

    <choice-email-template :open="dialogChoiceEmailTemplate" @close="dialogChoiceEmailTemplate = false" @choice="choiceEmailTemplate"></choice-email-template>
    <choice-email-box :open="dialogChoiceEmailBox" @close="dialogChoiceEmailBox = false" @choice="choiceEmailBox"></choice-email-box>
    <v-dialog v-model="dialogSubjectEmailTemplate" persistent max-width="600px" @keydown.esc="dialogSubjectEmailTemplate = false" >
      <v-card class="mb-0 pb-0">
        <v-card-title>
          <span class="headline">Изменение темы письма</span>
        </v-card-title>
        <v-card-text class="mb-0 pb-0">
          <v-container>
            <v-row>
              <v-col cols="12" v-if="emailTemplateChangeSubjectId">
                <v-text-field label="Тема письма" required v-model="emailQueueEmailTemplates.find(el=>el.id === emailTemplateChangeSubjectId).subject"></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions class="mt-0 pt-0 pb-3">
          <v-spacer></v-spacer>
          <v-btn color="secondary" dark elevation="1" small outlined @click="dialogSubjectEmailTemplate = false;emailTemplateChangeSubjectId = null" class="mr-3">Отмена</v-btn>
          <v-btn color="teal darken-3" dark elevation="1" small @click="changeSubject">Сохранить</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <delay-time :open.sync="dialogDelayTimeEmailTemplate" :time="emailTemplateChangeDelayTime.delay_time" @close="dialogDelayTimeEmailTemplate = false" @choice="setDelayTime"></delay-time>

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'

export default {
  data: () => ({
    formValid: true, // form
    item: null,
    emailQueueEmailTemplates: [],
    email_boxes: [],
    web_sites: [],
    mta_histories: [],
    search: '',
    totalItems: 0,
    options: {
        itemsPerPage: 15,
      },
    headers: [
      { text: 'Id', align: 'start', value: 'id'},
      { text: 'User id', align: 'start', value: 'user_id'},
      { text: 'Username', align: 'center', value: 'user.username', sortable: false},
      { text: 'ФИО', align: 'center', value: 'user.name', sortable: false},
      { text: 'Email (факт.)', align: 'center', value: 'email'},
      { text: 'Шаблон', align: 'center', value: 'email_template.name', sortable: false},
      { text: 'Шаг', align: 'center', value: 'queue_step_id'},
      { text: 'Завершение', align: 'center', value: 'queue_completed', sortable: false},
      { text: 'Открытий', align: 'center', value: 'opens'},
      { text: 'Дата 1-го открытия', align: 'center', value: 'opened_at'},
      { text: 'Отписка', align: 'center', value: 'unsubscribed'},
      { text: 'Дата отписки', align: 'center', value: 'unsubscribed_at'},
      // { text: 'Наименование', value: 'label'},
      // { text: 'Действие', align: 'center', value: 'actions', sortable: false},
    ],
    loading: false, // абстрактный статус загрузки..
    carried: false,
    dialogChoiceEmailTemplate: false, // dialog box
    dialogChoiceEmailBox: false, // dialog box
    emailQueueEmailTemplateId: 0,
    emailTemplateForStepId: null,

    emailTemplateChangeSubjectId: null,
    dialogSubjectEmailTemplate: false,

    dialogDelayTimeEmailTemplate: false,
    emailTemplateChangeDelayTimeId: null,
    emailTemplateChangeDelayTime: { delay_time: 0},

    errors: [],
    countBro: null,
    // moment: moment,
  }),
  computed: {
    colorEnabled(){
      return this.item.enabled ? 'green' : 'orange'
    }
  },
  async mounted() {
    await Promise.all([
      this.getEmailBoxes(),
      this.getWebSiteList(),

    ]).then(()=> {
      this.uploadItem().then(()=>{
        this.getMTAHistoriesFromApi()
      })
    })
  },
  watch: {
   search() {
        this.options.page = 1;
        this.searchMTAHistoriesFromApi();
      },
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.emailQueue.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads:"",
      })
          .then(resp => {
            if (resp['email_queue'] !== undefined) {
              this.item = resp['email_queue']
              this.uploadEmailQueueEmailTemplates()
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring'
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'История не найдена',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    uploadEmailQueueEmailTemplates: async function (showNotification = false) {
      return this.$api.emailQueueEmailTemplate.getListPagination({
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        preloads:"EmailTemplate"
      })
          .then(resp => {
            if (resp['email_queue_email_templates'] !== undefined) {

              this.emailQueueEmailTemplates = this.sortByKey(resp['email_queue_email_templates'],'step')
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring'
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Серия писем не найдена',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    getMTAHistoriesFromApi: async function (showNotification = false) {
      if (!this.item || !this.item.id) return

      return this.$api.mtaHistory.getListPagination({
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page-1)*(this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        email_queue_id: this.item.id,
        preloads:"EmailTemplate,User",
        ownerType: 'email_queues',
        ownerId: this.item.id,
        search: this.search,
      })
          .then(resp => {
            if (resp['mta_histories'] !== undefined) {

              this.mta_histories = resp['mta_histories']
              this.totalItems = resp['total']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring'
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Серия писем не найдена',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },
    searchMTAHistoriesFromApi: _.throttle( function () {
      this.getMTAHistoriesFromApi(false, this.search)
    }, 1400),
    updateItemData: async function () {

      if (!this.$refs.form.validate()) {
        this.$notify({
          group: 'main',
          title: 'Проверьте заполненные поля',
          type: 'warring'
        });
        return
      }

      let payload = this.item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

      await this.$api.emailQueue.update(payload)
          .then((resp) => {
            if (resp['email_queue'] !== undefined) {
              this.item = resp['email_queue']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [step]',
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
    onMoveCallback: async function (evt, e) {

      let keys = Object.keys(evt)
      if (keys.length < 1) {
        this.$notify({
          group: 'main',
          title: 'Не удалось обновить порядок',
          type: 'error'
        });
        return
      }

      let key = keys[0]
      let element = evt[key].element
      let newIndex = evt[key].newIndex

      switch (key) {
        case "added":
          break
        case "removed":
          break
        case "moved":
          let oldIndex = evt[key].oldIndex
          let stepMemory = parseInt(this.emailQueueEmailTemplates[newIndex].step)
          console.log("orderMemory:", stepMemory)
          this.emailQueueEmailTemplates[newIndex].step = this.emailQueueEmailTemplates[oldIndex].step
          this.emailQueueEmailTemplates[oldIndex].step = stepMemory

          // console.log("element id:", element.id)
          // console.log("newIndex:", newIndex)
          // console.log("oldIndex:", oldIndex)
          break
      }

    },
    saveNewOrders: async function (showNotification = false, force = false) {

      let arr = []

      this.emailQueueEmailTemplates.forEach((item, i) => {
        if (force) {
          let v = {
            id: item.id,
            step: i + 1
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].step = (i + 1)
          } else {
            arr.push(v)
          }
        } else if(item.step !== (i + 1)) {
          let v = {
            id: item.id,
            step: i + 1
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].step = (i + 1)
          } else {
            arr.push(v)
          }
        }
      })

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        email_queue_email_templates: arr
      }

      this.$api.emailQueueEmailTemplate.massUpdates(payload)
          .then(resp => {
            if (resp['email_queue_email_templates'] !== undefined) {
              this.emailQueueEmailTemplates = this.sortByKey(resp['email_queue_email_templates'],'step')
              this.carried = false
              // this.uploadEmailQueueEmailTemplates()
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success'
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring'
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Серия писем не найдена',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)
    },

    wKeyInItem(key, msg) {
      if (!this.item[key] || this.item[key].length < 1) {
        let index = this.errors.findIndex(el => el === msg)
        if (index === -1)
          this.errors.push(msg)
        this.enabledToEdit = false
      } else {
        let index = this.errors.findIndex(el => el === msg)
        if (index !== -1)
          this.errors.splice(index, 1);
      }
    },

    changeEmailTemplateForStep: async function(email_template) {
      this.emailTemplateForStepId = email_template.id
      this.dialogChoiceEmailTemplate = true

    },
    choiceEmailTemplate: async function(email_template) {

      // узнаем это для новых или текущему изменяем
      if (this.emailTemplateForStepId == null) {
        // начальное значение
        let nextStep = 1
        if (this.emailQueueEmailTemplates.length > 0) {
          let _item = _.maxBy(this.emailQueueEmailTemplates, function(o) { return o.step; })
          if (_item) {
            nextStep = _item.step + 1
          }
        }

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          email_queue_id: this.item.id,
          enabled: false,
          step: nextStep,
          email_template_id: email_template.id,
          delay_time: 0, // без задержки
          // EmailBoxId: 1, // !!!
          subject: "Сообщение для нового письма в серии",
          preloads:"EmailTemplate"
        }
        this.$api.emailQueueEmailTemplate.create(payload)
            .then(async (resp) => {
              if (resp['email_queue_email_template']) this.emailQueueEmailTemplates.push(resp['email_queue_email_template'])
              this.$notify({
                group: 'main',
                title: 'Новое письмо добавлено в очередь',
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
            .finally(()=> {
              this.dialogChoiceEmailTemplate=false
      })
      } else if (this.emailTemplateForStepId > 0 ) {

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          email_queue_id: this.item.id,
          email_template_id: email_template.id,
          id: this.emailTemplateForStepId,
          preloads:"EmailTemplate"
        }
        this.$api.emailQueueEmailTemplate.update(payload)
            .then(async (resp) => {
              if (resp['email_queue_email_template']) {

                let index = this.emailQueueEmailTemplates.findIndex(el=> el.id === this.emailTemplateForStepId, this.emailTemplateForStepId)
                if (index !== -1) {
                  this.emailQueueEmailTemplates[index] = resp['email_queue_email_template']
                }
              }
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
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
            .finally(()=> {
              this.dialogChoiceEmailTemplate = false
              this.emailTemplateForStepId = null
            })

      }


    },
    removeEmailTemplate: async function(email_template) {

      if (email_template.enabled) {
        this.$notify({
          group: 'main',
          title: 'Нельзя удалить активное письмо',
          text: 'Сначала остановите отправку письма',
          type: 'warring'
        });
        return
      }

      const res = await this.$confirm("Подтвердите удаление письма из очереди")
      if (!res) { return;}

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        id: email_template.id,
      }
      this.$api.emailQueueEmailTemplate.delete(payload)
          .then(async () => {
            let index = this.emailQueueEmailTemplates.findIndex(el=>el.id === payload.id, payload)
            if (index !== -1) {
              this.emailQueueEmailTemplates.splice(index, 1);
            }
            await this.saveNewOrders(false, true)

            this.$notify({
              group: 'main',
              title: 'Шаг удален из очереди',
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
      .finally(()=>this.dialogChoiceEmailTemplate=false)

    },

    getEmailBoxes: async function (showNotification) {

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.emailBox.getList(payload)
          .then(resp => {
            if (resp['email_boxes'] !== undefined) {
              this.email_boxes = resp['email_boxes']
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
    openDialogEmailBox(item) {
      if (item.enabled) {
        this.$notify({
          group: 'main',
          title: 'Нельзя изменить box для активного письма',
          text: 'Сначала остановите отправку письма',
          type: 'warring'
        });
        return
      }
      this.dialogChoiceEmailBox = true
      this.emailQueueEmailTemplateId = item.id
    },
    choiceEmailBox: async function(email_box) {

      // поиск по notification
      let item = this.emailQueueEmailTemplates.find(el=>el.id === this.emailQueueEmailTemplateId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error'
        });
        return
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        email_box_id: email_box.id,
        id: item.id
      }
      this.$api.emailQueueEmailTemplate.update(payload)
          .then(async (resp) => {
            item.email_box_id = email_box.id
            this.$notify({
              group: 'main',
              title: 'Email успешно обновлен',
              type: 'success'
            });
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=> {
            this.dialogChoiceEmailBox=false
            this.emailQueueEmailTemplateId = null
    })

    },
    getWebSiteList: async function (showNotification) {
      this.loading = true

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: 0,
        limit: 100,
        sortBy: "id",
        search: "",
        sortDesc: false,
      }

      await this.$api.webSite.getListPagination(payload)
          .then(resp => {
            if (resp['web_sites'] !== undefined ) {
              this.web_sites = resp['web_sites']
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
                  text: 'Ошибка в ответе сервера: web_sites - not found',
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

    changeStatusStep: async function(item) {

      if (!item.email_template_id || item.email_template_id < 1) {
        this.$notify({
          group: 'main',
          title: 'Сначала необходимо выбрать шаблон',
          type: 'warring'
        });
        return
      }

      if (!item.email_box_id || item.email_box_id < 1) {
        this.$notify({
          group: 'main',
          title: 'Сначала необходимо выбрать почтовый ящик',
          type: 'warring'
        });
        return
      }


      let msg = ""
      if (item.enabled) {
        msg = "Подтвердите остановку письма"
      } else {
        msg = "Подтвердите запуск письма"
      }

      const res = await this.$confirm(msg)
      if (!res) { return;}

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        id: item.id,
        enabled: !item.enabled

      }
      this.$api.emailQueueEmailTemplate.update(payload)
          .then(async (resp) => {
            if (resp['email_queue_email_template'] !== undefined) {
              item.enabled = !item.enabled
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success'
            });
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=>this.dialogChoiceEmailBox=false)
    },

    openDialogChangeSubject(item) {
      this.emailTemplateChangeSubjectId = item.id
      this.dialogSubjectEmailTemplate = true
    },
    changeSubject() {
      // поиск по notification
      let item = this.emailQueueEmailTemplates.find(el=>el.id === this.emailTemplateChangeSubjectId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error'
        });
        return
      }


      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        id: this.emailTemplateChangeSubjectId,
        subject: item.subject
      }
      this.$api.emailQueueEmailTemplate.update(payload)
          .then(async (resp) => {
            if (resp['email_queue_email_template'] !== undefined) {
              item.subject = resp['email_queue_email_template'].subject
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success'
            });
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=> {
            this.dialogSubjectEmailTemplate = false
            this.emailTemplateChangeSubjectId = null
          })

    },

    openDialogSetDelayTime(item) {
      this.emailTemplateChangeDelayTime = item
      this.emailTemplateChangeDelayTimeId = item.id
      this.dialogDelayTimeEmailTemplate = true
    },
    async setDelayTime(durationTime) {

      // поиск по notification
      let item = this.emailQueueEmailTemplates.find(el=>el.id === this.emailTemplateChangeDelayTimeId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error'
        });
        return
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        id: this.emailTemplateChangeDelayTimeId,
        delay_time: durationTime
      }
      await this.$api.emailQueueEmailTemplate.update(payload)
          .then(async (resp) => {
            if (resp['email_queue_email_template'] !== undefined) {
              item.delay_time = resp['email_queue_email_template'].delay_time
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success'
            });
          })
          .catch( (err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(()=> {
            this.dialogDelayTimeEmailTemplate = false
            this.emailTemplateChangeDelayTimeId = null
            this.emailTemplateChangeDelayTime = {}
          })

    },
    setNewStatusNotification: async function (status) {

      if (status === '') {
        return
      }
      this.loading = true
      await this.$api.emailQueue.changeStatus({
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        status: status,
        preloads:"",
      })
          .then(resp => {
            if (resp['email_queue'] !== undefined) {
              this.item = resp['email_queue']
              this.$notify({
                group: 'main',
                title: 'Статус изменен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_queue]',
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

    specBro: async function() {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id
      }

      await this.$api.emailQueue.appendAll(payload)
          .then((resp) => {
            if (resp['count'] !== undefined) {
              this.countBro = parseInt(resp['count'])
            }
            if (resp['email_queue'] !== undefined) {
              this.item = resp['email_queue']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [step]',
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
    }
  },
  components: {
    draggable,
    ChoiceEmailTemplate: () => import('@/components/views/tpl/email-marketing/templates/ChoiceEmailTemplate'),
    ChoiceEmailBox: () => import('@/components/views/tpl/web-sites/layouts/ChoiceEmailBox'),
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    WorkStatus: () => import('@/components/views/tpl/layouts/WorkStatus'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
  },

}
</script>

<style lang="scss">
.queue-stat {
  display: flex;
  justify-content: space-around;
  background-color: #f1f5f6;
  padding: 5px 10px;
  & > div.item {

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    & > span:first-child {
      font-weight: 500;
    }
    & > span:last-child {
      color: tomato;
    }
  }
}
.tdItem {
  padding: 12px 12px 10px;

  &:first-child {
    border-top: 1px solid grey;
  }

  border-left: 1px solid grey;
  border-right: 1px solid grey;
  border-bottom: 1px solid grey;

  &:hover {
    cursor: pointer;
    background-color: ghostwhite;
  }

  transition: .2s ease-in;
}


.flip-list-move {
  transition: transform 0.5s;
}
</style>
