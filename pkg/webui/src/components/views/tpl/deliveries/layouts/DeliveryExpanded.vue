<template>
    <td :colspan="headersLength">
        <v-row>
            <v-col cols="4">
                <h3 class="mt-4 mb-4">
                    <a v-if="item.code === 'indiaPost'" class="text-h3 a-link" target="_blank" href="https://otpravka.pochta.ru/specification#/nogroup-rate_calculate">API Почты России</a>
                </h3>

                <div class="d-flex my-0">
                    <v-text-field
                            style="max-width:200px"
                            label="Id доставки" dense class="body-1"
                            :value="item.id"
                            disabled
                    />
                    <v-switch
                            class="mb-2 mt-0 ml-6 mt-0" dense hide-details light
                            :success="item['enabled']"
                            :error="!item['enabled']"
                            v-model="item['enabled']"
                    >
                        <template v-slot:label>
                            <span class="input__label" :class="item['enabled'] ? 'green--text' : 'deep-orange--text darken-1'">{{item['enabled'] ? 'в работе' : 'недоступен'}}</span>
                        </template>
                    </v-switch>
                </div>

                <v-select
                        :items="$store.getters['shops/getShops']"
                        item-value="id" item-text="name"
                        return-object dense
                        :value="currentShopId"
                        @change="$emit('switchCurrentShop',$event)"
                        style="max-width: 350px"
                ></v-select>

                <v-text-field label="Name"  class="body-1" v-model.trim="item['name']" />

                <!-- Общие настройки любой доставки -->
                <div class="d-flex">
                    <v-checkbox v-model="item['address_required']" class="mt-0 py-0" label="Требовать адрес доставки" />
                    <v-checkbox v-model="item['postal_code_required']" class="ml-8 mt-0 py-0" label="Требовать индекс в адресе" />
                </div>
            </v-col>

            <!-- India Post -->
            <v-col cols="4" v-if="item.code === 'indiaPost'">
                <v-text-field label="access_token"  class="body-1" v-model.trim="item['access_token']" />
                <v-text-field label="x-User-Authorization (base64)"  class="body-1" v-model="item['x_user_authorization']" />
                <v-text-field label="Максимально разрешенный вес"  class="body-1" v-model.number.trim="item['max_weight']" />
            </v-col>
            <v-col cols="4" v-if="item.code === 'indiaPost'">
                <v-text-field label="Индекс пункта отправления"  class="body-1" v-model.trim="item['postal_codeFrom']"></v-text-field>
                <v-text-field class="body-1" v-model="item['mail_category']">
                    <template v-slot:label>
                        <a href="https://india.post.in/specification#/enums-base-mail-category" class="label-link" target="_blank">Категория РПО | mail_category</a>
                    </template>
                </v-text-field>
                <v-text-field class="body-1" v-model="item['mail_type']" >
                    <template v-slot:label>
                        <a href="https://india.post.in/specification#/enums-base-mail-type" class="label-link" target="_blank">Вид РПО | mail_type</a>
                    </template>
                </v-text-field>
            </v-col>
            <v-col cols="4" v-if="item.code === 'indiaPost'">
                <v-checkbox v-model="item['fragile']" class="mt-0" label="Осторожно/Хрупко" />
                <v-checkbox v-model="item['with_electronic_notice']" class="mt-0" label="С электронным уведомлением" />
                <v-checkbox v-model="item['with_order_of_notice']" class="mt-0" label="С заказным уведомлением" />
                <v-checkbox v-model="item['with_simple_notice']" class="mt-0" label="С простым уведомлением" />
            </v-col>

            <!-- courier -->
            <v-col cols="4" v-if="item.code === 'courier'">
                <v-text-field label="Стоимость (руб.)"  class="body-1" v-model.number.trim="item['price']" />
                <v-text-field label="Максимальный вес (кг)"  class="body-1" v-model.number.trim="item['max_weight']" />
            </v-col>

            <!-- pickup -->
            <v-col cols="4" v-if="item.code === 'pickup'">
                <v-text-field label="Стоимость (руб.)"  class="body-1" v-model.number.trim="item['price']" />
                <v-text-field label="Максимальный вес (кг)"  class="body-1" v-model.number.trim="item['max_weight']" />
            </v-col>

        </v-row>
        <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
            <v-btn @click="$emit('removeItem', item)" color="red darken-3" elevation="1" small dark>
                <v-icon x-small class="ml-1 mr-3">far fa-user-slash</v-icon>
                Удалить
            </v-btn>

            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="updateItemData(item)" dark color="teal darken-3" elevation="1" small>
                <v-icon x-small class="ml-1 mr-3">far fa-save</v-icon>
                Сохранить
            </v-btn>
        </v-card>
    </td>
</template>

<script>
    export default {
        data: () => ({
            item: {},
        }),

        props: ['itemProp', 'headersLength', 'currentShopId'],
        mounted() {
            this.item = this.itemProp
        },
        methods: {
            updateItemData: async function (item) {

                let payload = item
                Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

                await this.$api.deliveries.update(payload)
                    .then((resp) => {
                        if (resp['delivery'] !== undefined) {
                            // this.item = resp['delivery']
                            this.$set(this, 'item', resp['delivery'])
                            this.$notify({
                                group: 'main',
                                title: 'Данные обновлены',
                                type: 'success'
                            });
                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Ошибка обновления',
                                text: 'Не верный ответ сервера',
                                type: 'warring'
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

            },

        },
    }
</script>

<style lang="scss">

</style>