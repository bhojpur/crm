<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Email-уведомление
            <v-chip color="grey lighten-3" v-if="item">№ {{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'email-marketing.notifications'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="!loading">

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

            <!-- Название уведомления -->
            <v-text-field
                :rules="[v => !!v || 'Укажите название']"
                class="mt-0"
                label="Название уведомления"
                required
                v-model="item['name']"
            ></v-text-field>

            <!-- Тема сообщения -->
            <v-text-field
                :rules="[v => !!v || 'Укажите тему сообщения']"
                class="mt-0"
                label="Тема сообщения*"
                v-model="item['subject']"
            ></v-text-field>

            <!-- PreviewText -->
            <v-text-field
                :rules="[v => (v !== null && v.length <= 255) || 'Максимум 255 символов']"
                class="mt-0"
                label="Preview text"
                required
                v-model="item['preview_text']"
            ></v-text-field>

            <h4 class="mt-4 mb-1 primary--text" style="font-weight: 500;">Пользователи</h4>
            <div :class="{'red--text': !item.email_box_id}" @click="dialogChoiceUsers = true"
                 class="d-flex align-center item-editable-popup">
              <v-chip-group column>
                <template v-if="item._recipient_users && item._recipient_users.length > 0">
                  <v-chip v-for="tag in item._recipient_users" :key="tag.id" class="body-2 blue lighten-5" small>
                    {{$api.user.getFullName(tag) !== "" ? $api.user.getFullName(tag) + ' - ': ''}} {{ tag.email }}
                  </v-chip>
                </template>
                <span v-else class="deep-orange--text">
                  <v-icon class="mr-3 ml-1" small>
                    far fa-users
                  </v-icon>
                  n/a
                </span>

              </v-chip-group>
            </div>

            <h4 class="mt-4 pb-4 pt-0 primary--text" style="font-weight: 500;">Парсинг пользователей в контексте</h4>
            <div class="d-flex flex-column">
              <v-checkbox v-model="item['parse_recipient_user']" class="my-0 py-0" dense hide-details color="success" label="Искать пользователя .user_id" ></v-checkbox>
              <v-checkbox v-model="item['parse_recipient_customer']" class="my-0 py-0" label="Искать клиента .customer_id" dense hide-details color="success"></v-checkbox>
              <v-checkbox v-model="item['parse_recipient_manager']" class="my-0 py-0" label="Искать менеджера .manager_id" dense hide-details color="success"></v-checkbox>
              <v-checkbox v-model="item['parse_recipient_from_event']" class="my-0 py-0" label="Искать Ids получателей в событии" dense hide-details color="success"></v-checkbox>
            </div>
          </v-col>

          <v-col cols="6" offset="1">

            <!-- В работе или нет-->
            <work-status v-if="item"
                         :status.sync="item.status"
                         :canRun.sync="errors.length === 0"
                         :errors.sync="errors"
                         ownerType="email_notification"
                         @change="setNewStatus">
            </work-status>

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Задержка перед отправкой</h4>
            <div @click="openDialogSetDelayTime" class="item-editable-popup">{{ getDelay(item.delay_time) }}</div>

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Шаблон email-сообщения</h4>
            <div :class="{'red--text': !item.email_template_id}" @click="dialogChoiceEmailTemplate = true"
                 class="item-editable-popup">
              <v-icon class="mr-2" x-small>
                far fa-newspaper
              </v-icon>
              {{ getEmailTemplateName(item) }}
            </div>

            <h4 class="mt-4 mb-1" style="font-weight: 400;">Отправитель</h4>
            <div :class="{'red--text': !item.email_box_id}" @click="dialogChoiceEmailBox = true"
                 class="d-flex align-center item-editable-popup">
              <v-icon class="mr-2" x-small>
                far fa-mailbox
              </v-icon>
              {{ getEmailBoxName(item) }}
            </div>

          </v-col>

          <v-col cols="5">

          </v-col>
        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>
    </v-card>

    <choice-email-template v-if="item" :open="dialogChoiceEmailTemplate" @choice="choiceEmailTemplate" @close="dialogChoiceEmailTemplate = false"></choice-email-template>
    <choice-email-box v-if="item" :open="dialogChoiceEmailBox" @choice="choiceEmailBox" @close="dialogChoiceEmailBox = false"></choice-email-box>
    <choice-users v-if="item" :open="dialogChoiceUsers" @choice="choiceUsers" @close="dialogChoiceUsers = false" :usersIds.sync="item.recipient_users_list"></choice-users>
    <delay-time v-if="item" :open.sync="dialogDelayTime" :time="item.delay_time" @choice="setDelayTime" @close="dialogDelayTime = false" ></delay-time>

  </v-container>

</template>

<script>
import draggable from 'vuedraggable'
import emailQueueEmailTemplate from "@/api/modules/email-queue-email-template";

export default {
  data: () => ({
    formValid: true, // form
    item: null,
    enabledToEdit: false,
    emailQueueEmailTemplates: [],
    email_boxes: [],
    web_sites: [],
    loading: false, // абстрактный статус загрузки..
    carried: false,
    dialogChoiceEmailTemplate: false, // dialog template
    dialogChoiceEmailBox: false, // dialog box
    dialogChoiceUsers: false, // dialog box
    dialogSubjectEmailTemplate: false,
    dialogDelayTime: false,

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.getEmailBoxes(),
      this.getWebSiteList(),
    ]).then(()=>this.uploadItem())
  },
  watch: {
    item: {
      deep: true,
      handler() {
        this.wKeyInItem('subject', 'Необходимо установить тему сообщения')
        this.wKeyInItem('email_template_id', 'Необходимо установить шаблон письма')
        this.wKeyInItem('email_box_id', 'Необходимо установить отправителя')
        this.wKeyRecipients()
      }
    }
  },
  methods: {
    wKeyRecipients() {
      if (!this.item) return
      let msg = 'Необходимо добавить получателей'
      if (this.item.recipient_users_list && this.item.recipient_users_list.length < 1
          && !this.item.parse_recipient_customer && !this.item.parse_recipient_manager && !this.item.parse_recipient_user && !this.item.parse_recipient_from_event) {

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
    uploadItem: async function (showNotification = false) {

      this.loading = true
      await this.$api.emailNotification.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads: "EmailTemplate"
      })
          .then(resp => {
            if (resp['email_notification'] !== undefined) {
              this.item = resp['email_notification']
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
              title: 'Уведомление не найдено',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.loading = false)

    },

    updateItemData: async function () {
      let payload = this.item

      Object.assign(payload, {accountHashId: this.$store.state.account.hash_id,preloads: "EmailTemplate"})

      await this.$api.emailNotification.update(payload)
          .then((resp) => {
            if (resp['email_notification'] !== undefined) {
              this.item = resp['email_notification']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_notification]',
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
          let orderMemory = parseInt(this.emailQueueEmailTemplates[newIndex].order)
          console.log("orderMemory:", orderMemory)
          this.emailQueueEmailTemplates[newIndex].order = this.emailQueueEmailTemplates[oldIndex].order
          this.emailQueueEmailTemplates[oldIndex].order = orderMemory

          // console.log("element id:", element.id)
          // console.log("newIndex:", newIndex)
          // console.log("oldIndex:", oldIndex)
          break
      }

    },
    choiceEmailTemplate: async function (email_template) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_template_id: email_template.id,
        id: this.item.id,
      }
      this.$api.emailNotification.update(payload)
          .then(async (resp) => {
            if (resp['email_notification']) {
              this.item.email_template = resp['email_notification']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе email_notification',
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
    choiceEmailBox: async function (email_box) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_box_id: email_box.id,
        id: this.item.id
      }
      this.$api.emailNotification.update(payload)
          .then(async (resp) => {
            if (resp['email_notification'] !== undefined) {
              this.item.email_box_id = resp['email_notification'].email_box_id
              this.$notify({
                group: 'main',
                title: 'Отправитель успешно обновлен',
                type: 'success'
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Не хватает данных',
                type: 'warring'
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => {
            this.dialogChoiceEmailBox = false
          })

    },
    choiceUsers: async function (data) {

      let recipient_users_list = []
      if (data.append) {
        recipient_users_list = _.union(this.item['recipient_users_list'], data.users);
      } else {
        recipient_users_list = _.difference(this.item['recipient_users_list'], data.users);
      }

      // Обновляем данные
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        recipient_users_list: recipient_users_list
      }

      await this.$api.emailNotification.update(payload)
          .then((resp) => {
            if (resp['email_notification'] !== undefined) {
              this.item = resp['email_notification']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success'
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_notification]',
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
          .finally(() => this.dialogChoiceUsers = false)
    },

    removeUsersRecipient: async function () {

      this.dialogAddUser = false

      let notificationId = this.addedUsersNotificationId
      let removeListUsers = this.userSelected.map(el => el.id) // список id

      let item = this.items.find(el => el.id === notificationId)
      if (!item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          type: 'error'
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
          type: 'error'
        });
        return
      }
      item['recipient_users_list'] = _.union(item['recipient_users_list'], newListUsers);
      await this.updateItemData(item)
      this.userSelected = []
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
            if (resp['web_sites'] !== undefined) {
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
          .finally(() => this.loading = false)

    },
    changeStatusStep: async function (item) {

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
      if (this.enabledToEdit) {
        msg = "Подтвердите остановку письма"
      } else {
        msg = "Подтвердите запуск письма"
      }

      const res = await this.$confirm(msg)
      if (!res) {
        return;
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        email_queue_id: this.item.id,
        id: item.id,
      }
      this.$api.emailQueueEmailTemplate.update(payload)
          .then(async (resp) => {
            if (resp['email_queue_email_template'] !== undefined) {
              this.enabledToEdit = !this.enabledToEdit
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success'
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
            });
          })
          .finally(() => this.dialogChoiceEmailBox = false)
    },
    setNewStatus: async function (status) {

      if (status === '') {
        return
      }
      this.loading = true
      await this.$api.emailNotification.changeStatus({
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        status: status,
        preloads:"EmailTemplate",
      })
          .then(resp => {
            if (resp['email_notification'] !== undefined) {
              this.item = resp['email_notification']
              this.$notify({
                group: 'main',
                title: 'Статус изменен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_notification]',
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
    openDialogSetDelayTime() {
      this.dialogDelayTime = true
    },
    async setDelayTime(durationTime) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        delay_time: durationTime
      }
      await this.$api.emailNotification.update(payload)
          .then(async (resp) => {
            if (resp['email_notification'] !== undefined) {
              this.item.delay_time = resp['email_notification'].delay_time
            }
            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success'
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка внесения изменений',
              text: err['message'],
              type: 'error'
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
    }
  },
  components: {
    draggable,
    ChoiceEmailTemplate: () => import('@/components/views/tpl/email-marketing/templates/ChoiceEmailTemplate'),
    ChoiceEmailBox: () => import('@/components/views/tpl/web-sites/layouts/ChoiceEmailBox'),
    ChoiceUsers: () => import('@/components/views/tpl/organizations/users/ChoiceUsers'),
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    WorkStatus: () => import('@/components/views/tpl/layouts/WorkStatus'),
  },
  computed: {
    colorEnabled() {
      return this.enabledToEdit ? 'green' : 'orange'
    }
  }

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
