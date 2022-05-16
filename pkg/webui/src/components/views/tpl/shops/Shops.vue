<template>
  <v-container fluid>
    <v-data-table
            :headers="headers"
            :items="$store.getters['shops/getShops']"
            class="elevation-1"
            item-key="id"
            expand-icon="fal fa-angle-down"
            :single-expand="false"
            :expanded.sync="expanded"
            show-expand
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title><h1 style="font-size: 26px;">Магазины</h1></v-toolbar-title>
          <v-spacer></v-spacer>
          <v-card class="d-flex align-center" flat>

            <v-btn @click="updateLists(true)"
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

            <v-btn
                    @click="createItem"
                    color="cyan darken-4"
                    elevation="1"
                    small
                    dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-plus</v-icon>
              Добавить магазин
            </v-btn>

          </v-card>

        </v-toolbar>


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
          <h3 class="mb-8">Редактируемые данные</h3>
          <v-row dense class="mt-4">
            <v-col cols="3">

              <v-text-field
                      label="Имя магазина" dense class="body-1"
                      :value="$store.getters['shops/getShop'](item['id'])['name']"
                      @input="$store.commit('shops/SHOP_UPDATE_STATE', {item, field: 'name', value: $event})"
              />
              <v-text-field
                      label="Email" dense class="body-1"
                      :value="$store.getters['shops/getShop'](item['id'])['email']"
                      @input="$store.commit('shops/SHOP_UPDATE_STATE', {item, field: 'email', value: $event})"
              />
              <v-text-field
                      label="Телефон" dense class="body-1"
                      :value="$store.getters['shops/getShop'](item['id'])['phone']"
                      @input="$store.commit('shops/SHOP_UPDATE_STATE', {item, field: 'phone', value: $event})"
              />

            </v-col>
            <v-col cols="5" offset="1">
              <v-textarea
                      label="Адрес магазина" dense class="body-1 ml-0 pl-0" rows="2" flat
                      :value="$store.getters['shops/getShop'](item['id'])['address']"
                      @input="$store.commit('shops/SHOP_UPDATE_STATE', {item, field: 'address', value: $event})"
              />
            </v-col>

          </v-row>

          <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
            <v-btn @click="openDeleteDialog(item)"
                   color="red darken-3"
                   elevation="1"
                   small
                   dark
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-user-slash</v-icon>
              Удалить
            </v-btn>

            <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
            ></v-divider>

            <v-btn
                    @click="updateItemData(item)"
                    dark
                    color="teal darken-3"
                    elevation="1"
                    small
            >
              <v-icon
                      x-small
                      class="ml-1 mr-3"
              >far fa-save</v-icon>
              Сохранить
            </v-btn>
          </v-card>
        </td>
      </template>

    </v-data-table>

    <!-- Подтверждение удаления -->
    <v-dialog
            v-model="dialog.open"
            @keydown.esc="dialog.open = false"
            max-width="380"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name">
              <v-list-item-title class="headline">Имя: &laquo;{{dialog.name}}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.address">
              <v-list-item-title class="headline">Address: &laquo;{{dialog.address}}&raquo;</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
                  color="green darken-1"
                  small
                  text
                  @click.stop="dialog.open = false"
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
                  color="red darken-1"
                  small
                  text
                  @click.stop="deleteItemDialog"
          >
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
        { text: 'Id', align: 'start', value: 'id'},
        { text: 'Имя', value: 'name'},
        { text: 'Адрес', value: 'address'},
        { text: 'Удалить', value: 'actions'},
      ],
      items: [],
      expanded: [],
      dialog:{
        id: null,
        name: '',
        address: '',
        open: false,
      },
    }),
    mounted() {
      this.updateLists(false)
    },
    methods: {
      updateLists: async function (showNotification) {
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
        }

        await this.$store.dispatch('shops/GET_SHOPS', payload)
                .then(data => {
                  if (showNotification) {
                    this.$notify({
                      group: 'main',
                      title: 'Данные обновлены',
                      text: 'Список магазинов обновлен',
                      type: 'success'
                    });
                  }
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: err['message'],
                    type: 'error'
                  });
                })
                .finally(()=>this.loading = false)
      },
      // Обновляет все данные магазина магазин на стороне сервера
      updateItemData: async function(item) {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
        }
        Object.assign(payload, item)

        await this.$store.dispatch('shops/UPDATE_SHOP', payload)
                .then((data) => {
                  this.$notify({
                    group: 'main',
                    title: 'Данные обновлены',
                    text: 'Магазин обновлен',
                    type: 'success'
                  });
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка обновления',
                    text: 'Не удалось обновить данные',
                    type: 'error'
                  });
                })
      },

      createItem: async function (){
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          name: "Новый магазин",
          address: "-",
        }
        await this.$store.dispatch('shops/CREATE_SHOP', payload)
                .then((data) => {
                  this.$notify({
                    group: 'main',
                    title: 'Магазин создан',
                    type: 'success'
                  });
                  if (data['shop'] !== undefined) {
                    console.log(data)
                    // let shop = this.$store.getters['shops/getShop'](data['shop']['id'])
                    // console.log(shop)
                    // this.expanded.push(shop)
                  }
                })
                .catch( (err) => {
                  this.$notify({
                    group: 'main',
                    title: 'Ошибка создания',
                    text: err['message'],
                    type: 'error'
                  });
                })
      },

      // ниже переписать в общее диалоговое оконо
      openDeleteDialog(item) {
        this.dialog.id = item['id']
        this.dialog.name = item['name']
        this.dialog.address = item['address']
        this.dialog.open = true
      },
      deleteItemDialog() {
        this.dialog.open = false
        const payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.dialog.id
        }
        this.$store.dispatch('shops/DELETE_SHOP', payload)
                .then(async ()=> {
                  let index = this.items.findIndex(key => key['user'].hash_id === payload.hash_id, payload)
                  if (index !== -1) Vue.delete(this.items, index)

                  this.$notify({
                    group: 'main',
                    title: 'Магазин удален',
                    type: 'success'
                  });
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
                  this.dialog.token = ''
        })

      },
    },


  }
</script>