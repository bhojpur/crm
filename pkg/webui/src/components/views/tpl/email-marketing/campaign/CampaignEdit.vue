<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Email-рассылка
            <v-chip color="grey lighten-3" v-if="item">№ {{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'email-marketing.campaigns'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="item">

          <!-- Обновление и сохранение серии -->
          <!--<v-btn @click="checkDouble" color="secondary" dark elevation="1" small>
            <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
            Проверить дубли
          </v-btn>-->

<!--          <v-divider class="mx-2 mt-0" inset vertical></v-divider>-->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" :disabled="item['executed']" color="cyan darken-4" class="white--text" small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>
      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">
          <v-col cols="5">

            <!-- Название уведомления -->
            <v-text-field
                :rules="[v => !!v || 'Укажите имя кампании']"
                class="mt-0"
                label="Имя кампании"
                :disabled="item['status'] !=='pending'"
                required
                v-model="item['name']"
            ></v-text-field>

            <!-- Тема сообщения -->
            <v-text-field
                :rules="[v => !!v || 'Укажите тему сообщения']"
                class="mt-0"
                :disabled="item['status'] !=='pending'"
                label="Тема сообщения*"
                v-model="item['subject']"
            ></v-text-field>

            <!-- PreviewText -->
            <v-text-field
                :rules="[v => v.length <= 255 || 'Максимум 255 символов']"
                class="mt-0"
                label="Preview text"
                :disabled="item['status'] !=='pending'"
                required
                v-model="item['preview_text']"
            ></v-text-field>

            <h4 class="mt-4 mb-1 primary--text" style="font-weight: 500;">Сегмент пользователей</h4>
            <v-btn
                :disabled="item['status'] !=='pending'"
                small text :class="{'red--text': item.users_segment_id < 1}" @click="dialogChoiceUsersSegment = true"
                class="d-flex align-center item-editable-popup">
              <v-icon class="mr-2" x-small>
                far fa-mailbox
              </v-icon>
              {{ item.users_segment_id ? item.users_segment.name : 'n/a' }}
            </v-btn>

          </v-col>
          <v-col cols="6" offset="1">

            <!-- В работе или нет-->
            <work-status v-if="item"
                         :status.sync="item.status"
                         :canRun.sync="errors.length === 0"
                         :errors.sync="errors"
                         ownerType="email_campaign"
                         @change="setNewStatus">
            </work-status>

            <!--<v-card class="pl-2 py-2" outlined>
              <v-toolbar flat class="mb-0 pb-0">

                <h4 class="mr-3" style="font-weight: 400;">Статус рассылки: </h4>
                <span :class="$api.helper.getCSSColorByStatus(item.status)" class="status-table mr-2">
                  {{$api.helper.getHumanStatus(item.status)}}
                </span>

                <v-divider class="mx-2 mr-2" inset vertical v-show="item['status'] !== 'completed'&&item['status'] !== 'failed'&&item['status'] !== 'cancelled'"></v-divider>
                <v-btn
                    @click="setNewStatus('planned')"
                    v-show="item['status'] === 'pending'"
                    :disabled="item['status'] !== 'pending'"
                    small outlined class="teal&#45;&#45;text mx-2"
                >
                  <v-icon class="mr-2" small>fal fa-calendar-alt</v-icon>
                  Запланировать рассылку
                </v-btn>
                <v-btn
                    @click="setNewStatus('pending')"
                    v-show="item['status'] === 'planned'"
                    :disabled="item['status'] !== 'planned'"
                    small outlined class="brown&#45;&#45;text mx-2"
                >
                  <v-icon class="mr-2" small>fal fa-undo</v-icon>
                  Вернуть на доработку
                </v-btn>
                <v-btn
                    @click="setNewStatus('paused')"
                    v-show="item['status'] === 'active'"
                    :disabled="item['status'] !== 'active'"
                    small outlined
                    class="blue-grey&#45;&#45;text mx-2"
                >
                  <v-icon class="mr-2" small>fal fa-pause-circle</v-icon>
                  Приостановить
                </v-btn>
                <v-btn
                    @click="setNewStatus('active')"
                    v-show="item['status'] === 'paused'"
                    :disabled="item['status'] !== 'paused'"
                    small outlined
                    class="teal&#45;&#45;text mx-2"
                >
                  <v-icon class="mr-2" small>fal fa-play-circle</v-icon>
                  Возобновить отправку
                </v-btn>
                <v-btn
                    @click="setNewStatus('cancelled')"
                    v-show="item['status'] === 'planned' ||item['status'] === 'active' || item['status'] === 'paused'"
                    :disabled="item['status'] !== 'planned' && item['status'] !== 'active' && item['status'] !== 'paused'"
                    small outlined
                    class="red&#45;&#45;text mx-2"
                >
                  <v-icon class="mr-2" small>fal fa-ban</v-icon>
                  Отменить кампанию
                </v-btn>
              </v-toolbar>
              <div class="text-center" style="max-width: 300px;">
                <v-sheet class="white&#45;&#45;text" color="blue-grey lighten-2" v-if="item['status'] === 'failed'">{{item['failed_status']}}</v-sheet>
              </div>
            </v-card>-->

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Дата запуска</h4>
            <v-datetime-picker
                :disabled="item['status'] !=='pending'"
                v-model="timeToStart" date-format="dd-MM-yyyy"
                :time-picker-props="{format: '24hr'}"
                clearText="Отмена"
                okText="Установить">
              <template slot="dateIcon">
                <v-icon>fas fa-calendar</v-icon>
              </template>
              <template slot="timeIcon">
                <v-icon>fas fa-clock</v-icon>
              </template>
            </v-datetime-picker>

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Шаблон email-сообщения</h4>
            <v-btn text small
                   :disabled="item['status'] !=='pending'"
                   :class="{'red--text': !item.email_template_id}"
                   @click="dialogChoiceEmailTemplate = true"
                   class="item-editable-popup">
              <v-icon class="mr-2" x-small>
                far fa-newspaper
              </v-icon>
              {{ getEmailTemplateName(item) }}
            </v-btn>

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Отправитель</h4>
            <v-btn
                :disabled="item['status'] !=='pending'"
                small text :class="{'red--text': !item.email_box_id}" @click="dialogChoiceEmailBox = true"
                class="d-flex align-center item-editable-popup">
              <v-icon class="mr-2" x-small>
                far fa-mailbox
              </v-icon>
              {{ getEmailBoxName(item) }}
            </v-btn>

          </v-col>

        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>
    </v-card>

    <choice-email-template :open="dialogChoiceEmailTemplate" @choice="choiceEmailTemplate" @close="dialogChoiceEmailTemplate = false"></choice-email-template>
    <choice-email-box :open="dialogChoiceEmailBox" @choice="choiceEmailBox" @close="dialogChoiceEmailBox = false"></choice-email-box>
    <delay-time :open.sync="dialogDelayTime" :time="item.delay_time" @choice="setDelayTime"
                @close="dialogDelayTime = false" v-if="item"></delay-time>
    <choice-users-segment :open="dialogChoiceUsersSegment" @choice="choiceUsersSegment"
                          @close="dialogChoiceUsersSegment = false"></choice-users-segment>

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'
import emailQueueEmailTemplate from "@/api/modules/email-queue-email-template";

export default {
  data: () => ({
    formValid: true, // form
    item: null,
    enabledEdit: true,
    emailQueueEmailTemplates: [],
    email_boxes: [],
    web_sites: [],
    loading: false, // абстрактный статус загрузки..
    carried: false,
    dialogChoiceEmailTemplate: false, // dialog template
    dialogChoiceEmailBox: false, // dialog box
    dialogChoiceUsersSegment: false, // dialog box
    dialogSubjectEmailTemplate: false,
    dialogDelayTime: false,

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
      this.getEmailBoxes(),
      this.getWebSiteList(),
    ])
  },
  watch: {
    item: {
      deep: true,
      handler() {
        this.wKeyInItem('subject', 'Необходимо установить тему сообщения')
        this.wKeyInItem('email_template_id', 'Необходимо установить шаблон письма')
        this.wKeyInItem('email_box_id', 'Необходимо установить отправителя')
        this.wKeyInItem('users_segment_id', 'Необходимо выбрать сегмент')
      },
    },
  },
  methods: {
    wKeyInItem(key, msg) {
      if (!this.item[key] || this.item[key].length < 1) {
        let index = this.errors.findIndex(el => el === msg)
        if (index === -1)
          this.errors.push(msg)
        this.enabledEdit = false
      } else {
        let index = this.errors.findIndex(el => el === msg)
        if (index !== -1)
          this.errors.splice(index, 1);
      }
    },
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.emailCampaign.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads: "EmailTemplate,UsersSegment"
      })
          .then(resp => {
            if (resp['email_campaign'] !== undefined) {
              this.item = resp['email_campaign']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Кампания не найдена',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)

    },

    updateItemData: async function () {
      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        preloads: "EmailTemplate,UsersSegment"
      })

      await this.$api.emailCampaign.update(payload)
          .then((resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.item = resp['email_campaign']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
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
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })

    },
    choiceEmailTemplate: async function (email_template) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_template_id: email_template.id,
        id: this.item.id,
        preloads: "EmailTemplate,UsersSegment"
      }

      this.$api.emailCampaign.update(payload)
          .then(async (resp) => {
            if (resp['email_campaign']) {
              this.item = resp['email_campaign']
              // this.item['email_template'] = resp['email_campaign']['email_template']
              // this.item['email_template_id'] = resp['email_campaign']['email_template_id']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе email_campaign',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogChoiceEmailTemplate = false
          })


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
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },
    choiceEmailBox: async function (email_box) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_box_id: email_box.id,
        id: this.item.id,
        preloads: "EmailTemplate,UsersSegment"
      }
      this.$api.emailCampaign.update(payload)
          .then(async (resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.item.email_box_id = resp['email_campaign'].email_box_id
              this.item.email_box = resp['email_campaign'].email_box
              this.$notify({
                group: 'main',
                title: 'Отправитель успешно обновлен',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Не хватает данных',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogChoiceEmailBox = false
          })

    },
    choiceUsersSegment: async function (segment) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        users_segment_id: segment.id,
        id: this.item.id,
        preloads: "EmailTemplate,UsersSegment"
      }
      this.$api.emailCampaign.update(payload)
          .then(async (resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.item.users_segment_id = resp['email_campaign'].users_segment_id
              this.item.users_segment = resp['email_campaign'].users_segment
              this.$notify({
                group: 'main',
                title: 'Сегмент успешно обновлен',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Не хватает данных',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogChoiceUsersSegment = false
          })

    },

    /////////////
    removeUsersRecipient: async function () {

      this.dialogAddUser = false

      let notificationId = this.addedUsersNotificationId
      let removeListUsers = this.userSelected.map(el => el.id) // список id

      let item = this.items.find(el => el.id === notificationId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error',
        });
        return
      }
      item['recipient_users_list'] = _.difference(item['recipient_users_list'], removeListUsers);
      await this.updateItemData(item)
      this.userSelected = []
    },
    saveAddedUsersRecipient: async function () {

      this.dialogAddUser = false

      let notificationId = this.addedUsersNotificationId
      let newListUsers = this.userSelected.map(el => el.id) // список id

      let item = this.items.find(el => el.id === notificationId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error',
        });
        return
      }
      item['recipient_users_list'] = _.union(item['recipient_users_list'], newListUsers);
      await this.updateItemData(item)
      this.userSelected = []
    },
    /////////////////////

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
            if (resp['web_sites'] !== undefined) {
              this.web_sites = resp['web_sites']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Ошибка в ответе сервера: web_sites - not found',
                  type: 'warring',
                });
              }
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)

    },
    openDialogSetDelayTime() {
      this.dialogDelayTime = true
    },
    async setDelayTime(durationTime) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        delay_time: durationTime,
        preloads: "EmailTemplate,UsersSegment"
      }
      await this.$api.emailCampaign.update(payload)
          .then(async (resp) => {
            if (resp['email_campaign'] !== undefined) {
              this.item.delay_time = resp['email_campaign'].delay_time
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success',
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.dialogDelayTime = false
          })

    },

    requiredForRunNotification() {
      let arr = []
      if (!this.item.subject)
        arr.push('Необходимо указать тему сообщения')

      if (!this.item.email_box_id)
        arr.push('Необходимо установить отправителя')

      if (!this.item.email_template_id)
        arr.push('Необходимо установить шаблон')

      return arr
    },
    setNewStatus: async function (status) {
      if (status ==='') {
        return
      }
      await this.$api.emailCampaign.changeStatus({
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        status: status,
        preloads: "EmailTemplate,UsersSegment"
      })
          .then(resp => {
            if (resp['email_campaign'] !== undefined) {
              this.item = resp['email_campaign']
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

    checkDouble: async function () {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id
      }
      this.$api.emailCampaign.checkDouble(payload)
          .then(resp => {
            console.log("Дублей: ", resp['count'])
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
  },
  components: {
    draggable,
    ChoiceEmailTemplate: () => import('@/components/views/tpl/email-marketing/templates/ChoiceEmailTemplate'),
    ChoiceEmailBox: () => import('@/components/views/tpl/web-sites/layouts/ChoiceEmailBox'),
    ChoiceUsers: () => import('@/components/views/tpl/organizations/users/ChoiceUsers'),
    ChoiceUsersSegment: () => import('@/components/views/tpl/organizations/users/segment/layouts/ChoiceUsersSegment'),
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    WorkStatus: () => import('@/components/views/tpl/layouts/WorkStatus'),
  },
  computed: {
    timeToStart: {
      get() {
        // let t = this.$moment(this.item.schedule_run).format("MM-DD-YYYY hh:mm")
        let t = this.$moment(this.item.schedule_run).format("DD-MM-YYYY hh:mm")
        // let t = this.$moment(this.item.schedule_run).format("dd/MM/YYYY hh:mm")
        // console.log(t)
        return t
      }, set($val) {
        // console.log($val)
        // 2020-01-01T03:00:00+03:00
        // let t = this.$moment($val).format("MM:DD:YYYY hh:mm:ss")
        let t = this.$moment($val).toISOString();
        // let t = this.$moment($val).tz().format();
        // console.log(t)
        this.item.schedule_run = t
      },
    },
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

.item-ch {
  border: 1px dotted darkgrey;
  padding: 2px 6px;

  &:hover {
    cursor: pointer;
  }
}

.flip-list-move {
  transition: transform 0.5s;
}
</style>
