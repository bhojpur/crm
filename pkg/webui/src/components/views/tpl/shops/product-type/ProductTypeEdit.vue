<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Тип товара</h1>
          <v-btn @click="$router.push({name: 'product-types.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="item">

          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" color="cyan darken-4" class="white--text" small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>
      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">
          <v-col cols="4">

            <!-- name -->
            <v-text-field
                class="mt-0"
                label="name"
                required
                v-model.trim="item['name']"
            ></v-text-field>

            <!-- Code -->
            <v-text-field
                class="mt-0"
                label="Code"
                v-model.trim="item['code']"
            ></v-text-field>

          </v-col>
          <v-col cols="4" offset="1">
            <!-- color -->
            <v-color-picker :value="item['color']" @update:color="el => item['color']=el.hexa" flat></v-color-picker>
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
    uploadedFile: null, // загружаемый файл изображения
    tab: null,
    enabledEdit: true,
    products: [],
    email_boxes: [],
    web_sites: [],
    loading: false, // абстрактный статус загрузки..
    carried: false,
  }),
  async mounted() {
    await Promise.all([
      this.uploadItem(),
    ])
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.productType.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
      })
          .then(resp => {
            if (resp['product_type'] !== undefined) {
              this.item = resp['product_type']
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
      })

      await this.$api.productType.update(payload)
          .then((resp) => {
            if (resp['product_type'] !== undefined) {
              this.item = resp['product_type']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_category]',
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
    setColor(color) {
      console.log(color['hexa'])
    }
  },
}
</script>

