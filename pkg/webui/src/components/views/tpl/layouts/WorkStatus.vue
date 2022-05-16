<template>
  <!-- В работе или нет-->
  <v-card class="pt-4 pb-2 pl-6" outlined>
    <h4>Статус работы</h4>
    <span
        :class="$api.helper.getCSSColorByStatus(status, ownerType)" class="status-table mr-2"
    >
      {{$api.helper.getHumanStatus(status, ownerType)}}
    </span>

    <v-divider class="mx-2 mr-2" inset vertical v-show="status !== 'completed' && status !== 'failed' && status !== 'cancelled'"></v-divider>

    <!-- Запланировать -->
    <v-btn
        v-if="showPlanned(ownerType, status)"
        @click="setNewStatus('planned')"
        :disabled="!canRun"
        small outlined class="teal--text mx-2"
    >
      <v-icon class="mr-2" small v-text="getPlannedIconName(ownerType)" />
      {{getPlannedStatus(ownerType, status)}}
    </v-btn>

    <!-- Запустить в работу -->
    <v-btn
        v-if="showActive(ownerType, status)"
        @click="setNewStatus('active')"
        :disabled="!canRun"
        small outlined class="indigo--text text--darken-2 mx-2"
    >
      <v-icon class="mr-2" small v-text="getActiveIconName(ownerType)" />
      {{getActiveStatus(ownerType, status)}}
    </v-btn>

    <!-- Вернуть на доработку -->
    <v-btn
        v-if="showPending(ownerType, status)"
        @click="setNewStatus('pending')"
        :disabled="(status !== 'planned' && ownerType !== 'shipment') && !(status === 'completed' && ownerType === 'email_campaign')"
        small outlined class="brown--text mx-2"
    >
      <v-icon class="mr-2" small>fal fa-undo</v-icon>
      {{getPendingStatus(ownerType, status)}}
    </v-btn>

    <!-- Приостановить -->
    <v-btn
        v-show="showPaused(ownerType, status)"
        @click="setNewStatus('paused')"
        :disabled="status !== 'active'"
        small outlined
        class="blue-grey--text mx-2"
    >
      <v-icon class="mr-2" small>fal fa-pause-circle</v-icon>
      Приостановить
    </v-btn>

    <!-- Возобновить отправку / работу -->
    <v-btn
        @click="setNewStatus('active')"
        v-show="status === 'paused'"
        :disabled="!canRun"
        small outlined
        class="teal--text mx-2"
    >
      <v-icon class="mr-2" small>fal fa-play-circle</v-icon>
      <template v-if="ownerType==='email_campaign'">Возобновить отправку</template>
      <template v-else>Возобновить работу</template>
    </v-btn>

    <!-- В процессе разгрузки -->
    <v-btn
        v-if="showPosting(ownerType, status)"
        @click="setNewStatus('posting')"
        small outlined
        class="blue--text mx-2"
    >
      <v-icon class="mr-2" small>{{getPostingIconName(ownerType)}}</v-icon>
      {{getPostingStatus(ownerType, status)}}
    </v-btn>

    <!-- Отменить -->
    <v-btn
        v-if="showCancelled(ownerType, status)"
        @click="setNewStatus('cancelled')"
        small outlined
        class="red--text mx-2"
    >
      <v-icon class="mr-2" small>{{getCancelledIconName(ownerType)}}</v-icon>
      {{getCancelledStatus(ownerType, status)}}
    </v-btn>

    <!-- Завершить -->
    <v-btn
        v-if="showCompleted(ownerType, status)"
        @click="setNewStatus('completed')"
        small outlined
        class="green--text mx-2"
    >
      <v-icon class="mr-2" small>{{getCompletedIconName(ownerType)}}</v-icon>
      {{getCompletedStatus(ownerType, status)}}
    </v-btn>

    <!-- Зафейлить -->
    <v-btn
        v-if="showFailed(ownerType, status)"
        @click="setNewStatus('failed')"
        small outlined
        class="red--text mx-2"
    >
      <v-icon class="mr-2" small>{{getFailedIconName(ownerType)}}</v-icon>
      {{getFailedStatus(ownerType, status)}}
    </v-btn>

    <v-divider class="mx-2 mr-2" inset vertical v-show="status !== 'completed' && status !== 'failed' && status !== 'cancelled'"></v-divider>

    <div class="deep-orange--text  darken-4 mt-2 body-2" v-for="(error,key) in errors">{{key+1}}. {{error}}</div>
  </v-card>
</template>

<script>

export default {
  data: () => ({
  }),
  props: ['status','ownerType','canRun','errors'],
  async mounted() {
    // await this.getDataFromApi(false)
  },
  methods: {
    getDataFromApi: async function (showNotification, searchStr = '') {
      // if (!this.web_site_id)return
      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page - 1) * (this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: searchStr,
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
    setNewStatus: async function(status) {

      if (status === '') {
        return
      }

      this.$emit('change',status)
    },

    getPlannedStatus(owner, status) {
      switch (owner) {
        case 'email_campaign':
          if (status === 'completed') return 'Запланировать повторно'
          return 'Запланировать рассылку'
        case 'shipment' || 'question':
          return 'Запланировать'
      }
    },
    getPlannedIconName(owner) {
      switch (owner) {
        case 'email_campaign':
          return 'fal fa-calendar-alt'
        case 'shipment':
          return 'fal fa-calendar-alt'
      }
    },

    getActiveStatus(owner, status) {
      switch (owner) {
        case 'email_campaign':
        case  'email_queue':
          return 'Запустить в работу'
        case  'email_notification':
          return 'Запустить в работу'
        case 'shipment':
          return 'В процессе доставки'
      }
    },
    getActiveIconName(owner) {
      switch (owner) {
        case 'email_notification':
        case  'email_queue':
          return 'fal fa-mail-bulk'
        case 'shipment':
          return 'fal fa-dolly-flatbed-alt'
      }
    },

    getPendingStatus(owner, status) {
      switch (owner) {
        case 'email_campaign':
          if (status === 'completed') return 'Повторить кампанию'
          return 'На доработку'
        case 'shipment':
          return 'На доработку'
        default:
          return 'Вернуть на доработку'
      }
    },
    getPendingIconName(owner) {
      switch (owner) {
        default:
          return 'fal fa-ban'
      }
    },

    getPostingStatus(owner, status) {
      switch (owner) {
        case 'shipment':
          return 'В процессе разгрузки'
        default:
          return ''
      }
    },
    getPostingIconName(owner) {
      switch (owner) {
        default:
          return 'fal fa-dolly-flatbed'
      }
    },

    getCancelledStatus(owner, status) {
      switch (owner) {
        case 'email_campaign':
          return 'Отменить кампанию'
        /*case  'email_queue':
          return 'Отменить серию'*/
        default:
          return 'Отменить'
      }
    },
    getCancelledIconName(owner) {
      switch (owner) {
        default:
          return 'fal fa-ban'
      }
    },

    getCompletedStatus(owner, status) {
      switch (owner) {
        case 'shipment':
          return 'Оприходовать товары на склад'
        case 'question':
          return 'Закрыть обращение'
        default:
          return 'Завершить'
      }
    },
    getCompletedIconName(owner) {
      switch (owner) {
        case 'shipment':
          return 'fal fa-box-up'
        default:
          return 'fal fa-check-circle'
      }
    },

    getFailedStatus(owner, status) {
      switch (owner) {
        case 'shipment':
          return 'Ошибка поступления'
        default:
          return ''
      }
    },
    getFailedIconName(owner) {
      return 'fal fa-window-close'
    },

    showPending(owner, status) {
      if (status === 'completed' && (owner === 'email_campaign')) return true
      if (owner === 'shipment' && (status ==='active' || status === 'planned')) {
        return true
      }
      return status === 'planned'
    },
    showPlanned(owner, status) {
      // return true
      // if (status === 'completed' && (owner === 'email_campaign')) return true

      if (status !== 'pending' && status !== 'active' && status !== 'paused') {
        return false
      }

      return (owner === 'email_campaign') || (owner === 'shipment')
    },
    showActive(owner, status) {
      switch (owner) {
        case 'email_campaign':
          return false
        case 'shipment':
          return (status === 'planned') ||(status === 'paused')
        case 'question':
          return false
        default:
          return (status === 'pending')
      }
    },
    showPaused(owner, status) {
      return (owner !== 'shipment') && (status === 'active')
    },
    showCancelled(owner, status) {
      if (owner === 'email_campaign') {
        if (status !== 'active' && status !== 'paused' && status !== 'planned')
          return false
      } else {
        if (status !== 'active' && status !== 'paused')
          return false
      }


      return (owner==='email_campaign') || (owner==='email_notification') || (owner === 'email_queue') || (owner === 'shipment')
    },
    showCompleted(owner, status) {
      if (owner === 'email_campaign') {
          return false
      }
      if (owner === 'shipment') {
        return (status === 'posting')
      }
      if (owner === 'question') {
        return (status === 'pending')
      }

      return (status === 'active') || (status === 'paused') || (status === 'posting')
    },
    showPosting(owner, status) {
      if (owner !== 'shipment') {
        return false
      }
      return (status === 'active') || (status === 'paused')
    },
    showFailed(owner, status) {
      if (owner !== 'shipment') {
        return false
      }
      return (status === 'posting')
    },
  },

}
</script>
