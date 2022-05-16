<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Создание серии писем</h1>
          <v-btn @click="$router.push({name: 'email-marketing.queues'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку серий</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat>

          <!-- Создание заказа -->
          <v-btn @click="createItem" color="cyan darken-4" dark small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Создать серию писем
          </v-btn>
        </v-card>
      </v-toolbar>

      <v-form class="mx-4 mt-6 mb-4" ref="form" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-4">
          <v-col cols="5">
            <!-- Стоимость -->
            <v-text-field
                :rules="[v => !!v || 'Укажите имя серии']"
                class="mt-0" dense
                label="Имя серии писем"
                required
                v-model="item['name']"
            ></v-text-field>

            <!-- В работе или нет-->
            <v-switch
                :label="item['enabled'] ? 'В работе' : 'Остановлена'"
                :success="item['enabled']"
                class="mb-6 mt-0"
                dense
                hide-details
                light
                v-model="item['enabled']"
            ></v-switch>
          </v-col>
          <v-col cols="5" offset="2">

          </v-col>
        </v-row>

        <v-divider class="mb-8 mt-4"></v-divider>
      </v-form>

      <!-- Список писем -->
      <v-card class="d-flex justify-space-between mb-1 mx-1" flat>

        <v-card-title>
          <h2 class="text-h3 pb-3">Цепочка писем</h2>
        </v-card-title>

        <v-card-actions>
          <v-btn
              @click="item.htmlData = $api.email_template.getBaseDataTemplate()"
              disabled
              small color="secondary" elevation="1">
            <v-icon class="mr-3 ml-2" small>fal fa-plus-circle</v-icon>
            Добавить письмо
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-card>

  </v-container>

</template>

<script>
export default {
  data: () => ({
    item: {},
    formValid: true,
  }),
  methods: {
    createItem: async function () {

      if (!this.$refs.form.validate()) {
        this.$notify({
          group: 'main',
          title: 'Проверьте заполненные поля',
          type: 'warring'
        });
        return
      }

      let payload = {accountHashId: this.$store.state.account.hash_id}
      Object.assign(payload, this.item)


      await this.$api.emailQueue.create(payload)
          .then((resp) => {
            if (resp['email_queue'] !== undefined) {
              return this.$router.push({
                name: 'email-marketing.queues.edit',
                params: {public_id: resp['email_queue'].public_id}
              })
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_queue]',
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
  },

}
</script>
