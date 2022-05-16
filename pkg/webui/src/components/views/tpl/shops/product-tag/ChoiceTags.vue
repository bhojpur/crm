<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card class="rounded-0" flat>

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор тегов для товара</h2>
        <v-divider
            class="ml-6 mr-4"
            dark
            vertical
        ></v-divider>
        <v-text-field
            class="mx-2 my-0 white--text"
            dark
            dense
            hide-details placeholder="поиск..."
            prepend-inner-icon="fal fa-search"
            single-line
            style="max-width: 380px;"
            v-model.trim="search">
        </v-text-field>
        <v-spacer></v-spacer>
        <span class="mr-4 body-1 grey--text">Выбрано: {{ selected.length }}</span>
        <v-btn :disabled="selected.length< 1" @click="emitItems(false)" class="orange--text mr-3" dark elevation="0"
               outlined small>
          <v-icon class="mr-4">far fa-trash-alt</v-icon>
          Убрать теги
        </v-btn>

        <v-btn :disabled="selected.length< 1" @click="emitItems(true)" class="teal--text" dark elevation="0"
               outlined small>
          <v-icon class="mr-4">far fa-save</v-icon>
          Добавить теги
        </v-btn>

        <v-divider
            class="mx-5"
            dark
            vertical
        ></v-divider>

        <v-btn :block="false" @click="getDataFromApi(true)" class="mt-0 blue--text text--lighten-4 mr-3" elevation="0"
               outlined small>
          <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
          Обновить
        </v-btn>

        <v-btn @click="$emit('close')" class="deep-orange--text" elevation="0" outlined small>
          <v-icon class="mr-3">far fa-times</v-icon>
          Закрыть
        </v-btn>
      </v-card-title>

      <v-data-table
          :headers="[
                  { text: 'Id', align: 'start', value: 'public_id'},
                  { text: 'Имя категории', value: 'label'},
                  { text: 'Группа тегов',  value: 'product_tag_group_id'},
                  { text: 'Код',  align: 'end', value: 'code'},
                  { text: 'Числится товаров',  align: 'center', value: '_product_count'},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :items-per-page="15"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          item-key="id"
          selectable-key="id"
          show-select
          v-model="selected"
      >
        <template v-slot:item.label="{ item }">
          {{item['label'] ? item['label'] : ''}}
        </template>
        <template v-slot:item.product_tag_group_id="{ item }">
          {{item['product_tag_group_id'] && item['product_tag_group'] ? item['product_tag_group'].label : ''}}
        </template>
        <template v-slot:item.code="{ item }">
          <span class="mr-3">{{item['code'] ? item['code'] : '-'}}</span>
        </template>1857
        <template v-slot:item.product_count="{ item }">
        <span :class="{ 'red--text': !(item['product_cards'] && item['product_cards'].length > 0) }">
          {{item['product_cards'] && item['product_cards'].length > 0 ? item['product_cards'].length : '0'}}
        </span>
        </template>

      </v-data-table>

    </v-card>
  </v-dialog>

</template>

<script>

export default {
  data: () => ({
    items: [],
    totalItems: 0,
    selected: [],
    options: {
        itemsPerPage: 15,
      },
    search: '',
    loading: false,
    dialog: {open: false},
  }),
  props: {
    open: {
      type: Boolean,
      default: false
    },
    catIds: {
      type: Array,
      default: []
    }
  },
  async mounted() {
    // await this.getDataFromApi(false).then(()=>this.loading = false)

  },
  watch: {
    search() {
      this.options.page = 1
      this.searchDataFromApi();
    },
    open() {
      this.getDataFromApi(false)
    }
  },
  methods: {
    getDataFromApi: async function (showNotification) {

      console.log("dsa")
      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page-1)*(this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        preloads: "ProductTagGroup"
      }

      await this.$api.productTag.getListPagination(payload)
          .then(resp => {
            if (resp['product_tags'] !== undefined) {
              this.items = resp['product_tags']
              this.totalItems = resp.total


              if (this.catIds && this.catIds.length > 0) {
                this.selected = _.intersectionWith(this.items, this.catIds, function (obj, other) {
                  let v1, v2
                  if (typeof obj === 'object')
                    v1 = parseInt(obj.id)
                  else
                    v1 = parseInt(obj)

                  if (typeof other === 'object')
                    v2 = parseInt(other.id)
                  else
                    v2 = parseInt(other)

                  return v1 === v2;
                })
              }

              // console.log(this.selected)
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  text: 'Список категорий обновлен',
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
    emitItems: function (status) {
      // console.log(this.selected)
      // console.log(this.selected.map(el => el.id))
      this.$emit('choice', {
        categories: this.selected.map(el => el.id),
        append: status
      })
      this.selected = []
    },
    isAdded(item) {
      if (!item) return false
      return !!_.find(this.catIds, (el => el === item.id))
    }
  }
}
</script>
