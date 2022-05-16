<template>
    <div>
        <v-toolbar color="white" flat>
            <v-toolbar-title><h3 style="font-weight: 500;">Основные параметры товара</h3></v-toolbar-title>

            <v-spacer></v-spacer>
            <v-card class="d-flex justify-end mb-3" flat style="background-color: transparent;">
                <v-btn @click="openDeleteDialog(item)" color="red darken-3" dark small outlined>
                    <v-icon class="ml-1 mr-3" x-small>far fa-user-slash</v-icon>
                    Удалить
                </v-btn>
                <v-divider class="mx-2 mt-0" inset vertical></v-divider>
                <v-btn @click="updateProductData(item, true)" color="teal darken-3" dark small outlined>
                    <v-icon
                            class="ml-1 mr-3"
                            x-small
                    >far fa-save</v-icon>
                    Сохранить
                </v-btn>
            </v-card>
        </v-toolbar>
        <v-row class="pa-2">
            <v-col cols="4">
                <v-text-field
                        :value="item.id" class="body-1" dense
                        disabled
                        label="Id товара"
                />
                <v-switch
                        :label="item['enabled'] ? 'В продаже' : 'Снят с продаж'"
                        :success="item['enabled']"
                        class="mb-6 mt-0"
                        dense
                        hide-details
                        light
                        v-model="item['enabled']"
                ></v-switch>

                <v-text-field
                        class="body-1"  label="Артикул"
                        v-model="item['article']"
                />
                <v-text-field
                        class="body-1"  label="SKU"
                        v-model="item['sku']"
                />
                <v-text-field
                        class="body-1"  label="Модель"
                        v-model="item['model']"
                />

            </v-col>
            <v-col cols="4">
                <v-text-field
                        class="body-1"  label="Name"
                        v-model="item['name']"
                />
                <v-text-field
                        class="body-1"  label="Short name"
                        v-model="item['short_name']"
                />

                <v-textarea
                        class="body-1"  label="Краткое описание товара" rows="2"
                        v-model="item['short_description']"
                />

                <v-select v-model="item['payment_subject_id']" :items="payment_subjects" item-text="name" item-value="id" label="Признак предмета расчета" class="mt-4"></v-select>

                <v-select v-model="item['vat_code_id']" :items="vatCodes" item-text="name" item-value="id" label="Код ставки НДС" class="mt-4"></v-select>
            </v-col>
            <v-col cols="4">
                <v-text-field
                        class="body-1"  label="Розничная цена"
                        v-model="item['retail_price']"
                />
                <v-text-field
                        class="body-1"  label="Оптовая цена"
                        v-model="item['wholesale_price']"
                />
                <v-text-field
                        class="body-1"  label="Закупочная цена"
                        v-model="item['purchase_price']"
                />
                <v-text-field
                        class="body-1"  label="Розничная скидка (факт., в руб.)"
                        v-model="item['retail_discount']"
                />
            </v-col>
            <v-col cols="12">
                <h3 class="mb-3">Описание товара</h3>
                <prism-editor autoStyleLineNumbers class="myPrismEditor mb-3 pb-2" language="html"  v-model="item['description']"></prism-editor>
            </v-col>
        </v-row>

        <!-- Подтверждение удаления -->
        <v-dialog
                @keydown.esc="dialog.open = false"
                max-width="380"
                v-model="dialog.open"
        >
            <v-card>
                <v-card-title style="font-size: 21px;">Подтвердите удаление</v-card-title>

                <v-card-text class="d-flex flex-column">
                    <v-list>
                        <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.id">
                            <v-list-item-title class="headline">Id: &laquo;{{dialog.id}}&raquo;</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.name">
                            <v-list-item-title class="headline">Имя: &laquo;{{dialog.name}}&raquo;</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.model">
                            <v-list-item-title class="headline">Модель: &laquo;{{dialog.model}}&raquo;</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialog.sku">
                            <v-list-item-title class="headline">SKU: &laquo;{{dialog.sku}}&raquo;</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-card-text>

                <v-card-actions>
                    <v-spacer></v-spacer>

                    <!-- Отменить -->
                    <v-btn
                            @click.stop="dialog.open = false"
                            color="green darken-1"
                            small
                            text
                    >
                        Отменить
                    </v-btn>

                    <!-- Удалить -->
                    <v-btn
                            @click.stop="deleteItemDialog"
                            color="red darken-1"
                            small
                            text
                    >
                        Удалить
                    </v-btn>

                </v-card-actions>
            </v-card>
        </v-dialog>
    </div>

</template>

<script>

    export default {
        data: () => ({
            dialog:{
                id: null,
                name: '',
                model: null,
                sku: null,
                open: false,
            },
        }),

        props: ['item', 'options' ,'payment_subjects','vatCodes'],
        methods: {
            openDeleteDialog(item) {
                this.dialog.id = item['id']
                this.dialog.name = item['name']
                this.dialog.model = item['model']
                this.dialog.sku = item['sku']
                this.dialog.open = true
            },
            deleteItemDialog() {
                this.dialog.open = false
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    id: this.dialog.id
                }
                // this.$store.dispatch('shops/DELETE_PRODUCT', payload)
                this.$api.product.delete(payload)
                    .then(async ()=> {
                        this.$emit('removeItem', payload)

                        // let index = this.items.findIndex(item => item.id === payload.id, payload)
                        // if (index !== -1) Vue.delete(this.items, index)

                        this.$notify({
                            group: 'main',
                            title: 'Товар удален',
                            type: 'success'
                        });

                        // обновляем текущий список
                        await this.getDataFromApi(false)
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
                    })

            },
            getDataFromApi: async function(showNotification) {

                this.loading = true
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    offset: (this.options.page-1)*(this.options.itemsPerPage),
                    limit: this.options.itemsPerPage,
                  preloads: 'ProductCards,PaymentSubject'
                }

                // await this.$store.dispatch('shops/GET_PRODUCT_PAGINATION_LIST', payload)
                await this.$api.product.getListPagination(payload)
                    .then(data => {
                        // this.items = data.users
                        this.$set(this, 'items', data.products)
                        this.totalItems = data.total
                        if (showNotification) {
                            this.$notify({
                                group: 'main',
                                title: 'Данные обновлены',
                                text: 'Список товаров обновлен',
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
            getPaymentSubjectsFromApi: async function(showNotification = false) {

                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                }

                await this.$api.paymentSubject.getList(payload)
                    .then(resp => {
                        if (resp['payment_subjects'] !== undefined) {
                            this.payment_subject = resp['payment_subjects']
                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Ошибка получения списка payment Subjects',
                                type: 'warring'
                            });
                        }

                    })
                    .catch( (err) => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка получения списка Payment Subjects',
                            text: err['message'],
                            type: 'error'
                        });
                    })
            },
        },
        name: "ProductImage"
    }
</script>