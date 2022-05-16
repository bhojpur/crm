<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <div>
            <h1 style="font-size: 26px;display: inline-block">Event</h1>
            <v-divider class="ml-2 mt-0" inset vertical></v-divider>
            <span style="font-size: 14px;color: #de3327">{{isOwner ? 'Специальное' : 'Системное'}}</span>
          </div>
          <v-btn @click="$router.push({name: 'events.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="item">

          <!-- Обновление и сохранение серии -->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>
          <v-btn @click="updateItemData" color="cyan darken-4" class="white--text" small :disabled="!isOwner">
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>
      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">
          <v-col cols="5">

            <!-- Name -->
            <v-text-field
                class="mt-0"
                label="Имя события"
                counter="128"
                required :disabled="!isOwner"
                v-model.trim="item['name']"
            ></v-text-field>

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Код события"
                counter="128"
                required :disabled="!isOwner"
                v-model.trim="item['code']"
            ></v-text-field>

            <v-textarea
                label="Краткое описание" :disabled="!isOwner"
                counter="255" rows="2" class="mt-0"
                v-model.trim="item['description']"
            ></v-textarea>

          </v-col>
          <v-col cols="5" offset="1">
            <!-- External Call Available -->
            <v-checkbox
                v-model="item['external_call_available']"
                class="my-0 py-0" :disabled="!isOwner"
                label="Разрешить запуск по API"
            ></v-checkbox>

            <!-- Parsing Recipient List -->
            <v-checkbox
                v-model="item['parsing_recipient_list']"
                class="my-0 py-0" :disabled="!isOwner"
                label="Сбор Ids контактов из контекста вызова по API"
            ></v-checkbox>

            <!-- Parsing Payload -->
            <v-checkbox
                v-model="item['parsing_payload']"
                class="my-0 py-0" :disabled="!isOwner"
                label="Загрузка данных [payload] из контекста вызова по API"
            ></v-checkbox>
          </v-col>
        </v-row>

        <!-- Описание -->
        <v-divider class="mb-8 mt-4"></v-divider>

      </v-form>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data: () => ({
    formValid: true, // form
    item: null,
    email_boxes: [],
    web_sites: [],
    loading: false, // абстрактный статус загрузки..
    carried: false,

    errors: [],
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
    ])
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.event.get({
        accountHashId: this.$store.state.account.hash_id,
        id: this.$route.params.id,
        // preloads:"Products", // может надо отдельно подгружать карточки товаров через filter[]
      })
          .then(resp => {
            if (resp['event'] !== undefined) {
              this.item = resp['event']

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

      if (!this.isOwner) return
      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        id: this.$route.params.id,
        // preloads:"Products",
      })

      await this.$api.event.update(payload)
          .then((resp) => {
            if (resp['event'] !== undefined) {
              this.item = resp['event']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [event]',
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

  },
  computed: {
    isOwner() {
      if (!this.item || !this.$store.state.account) return false

      return this.item['account_id'] === this.$store.state.account.id;
    },
  },
  components: {
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
  },


}
</script>
