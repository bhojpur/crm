<template>
  <v-dialog @keydown.esc="open = false" fullscreen transition="dialog-bottom-transition" v-model="open">
    <v-card class="rounded-0" flat>

      <v-card-title class="d-flex align-center blue-grey darken-3 white--text">
        <h2 style="font-weight: 500;">Выбор категорий</h2>
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
          Убрать категории
        </v-btn>

        <v-btn :disabled="selected.length< 1" @click="emitItems(true)" class="teal--text" dark elevation="0"
               outlined small>
          <v-icon class="mr-4">far fa-save</v-icon>
          Добавить категории
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
                  { text: 'Наименование',  value: 'label'},
                  { text: 'Наименование во множ. числе',  value: 'label_plural'},
                  { text: 'Отображать у товара', align:'center', value: 'show_property'},
                  { text: 'Code', align:'start',value: 'code'},
                ]"
          :items.sync="items"
          :loading="loading"
          :options.sync="options"
          :search="search"
          :server-items-length="totalItems"
          @update:search="searchDataFromApi"
          @update:options="getDataFromApi(false)"
          class="elevation-1"
          item-key="id"
          selectable-key="id"
          :items-per-page="15"
          show-select
          v-model="selected"
      >
        <template v-slot:item.data-table-select="{ isSelected, select }">
          <v-simple-checkbox :value="isSelected" @input="select($event)" :ripple="false"></v-simple-checkbox>
        </template>

        <template v-slot:item.public_id="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{ item.public_id }}</span>
        </template>

        <template v-slot:item.label="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{item.label}}</span>
        </template>

        <template v-slot:item.label_plural="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{item.label_plural}}</span>
        </template>

        <template v-slot:item.code="{ item }">
          <span :class="{'green--text': isAdded(item)}">{{item.code}}</span>
        </template>

        <template v-slot:item.show_property="{ item }">
          <show-item-yes-no :show="item['show_property']" :bright="false"></show-item-yes-no>
        </template>

      </v-data-table>

    </v-card>
  </v-dialog>

</template>

<script>
export default {
  data: () => ({
    items: [],
    roles: [],
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

      this.loading = true
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        offset: (this.options.page-1)*(this.options.itemsPerPage),
        limit: this.options.itemsPerPage,
        sortBy: (this.options.sortBy !== undefined && typeof this.options.sortBy[0] === 'string') ? this.options.sortBy[0] : "public_id",
        sortDesc: this.options.sortDesc[0] !== undefined ? this.options.sortDesc[0] : true,
        search: this.search,
        preloads: 'Products',
      }

      await this.$api.productCategory.getListPagination(payload)
          .then(resp => {
            if (resp['product_categories'] !== undefined) {
              this.items = resp['product_categories']
              this.totalItems = resp.total

              // Не нужно..
              /*if (this.catIds && this.catIds.length > 0) {
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
              }*/

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
  },
  components: {
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
  },
}
</script>
