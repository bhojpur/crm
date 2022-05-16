<template>
    <td :colspan="headersLength">
        <v-toolbar flat>
            <v-toolbar-title><h4 class="mt-4 mb-6" style="font-weight: 400;">Параметры изображения</h4></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-card flat style="background-color: transparent;" class="d-flex justify-end mb-3">
                <v-btn @click="openDeleteDialog(item)" color="red lighten-1" elevation="0" small dark outlined>
                    <v-icon x-small class="ml-1 mr-3">far fa-user-slash</v-icon>
                    Удалить
                </v-btn>
                <v-divider class="mx-2 mt-0" inset vertical></v-divider>
                <v-btn @click="updateItemData(item, true)" dark color="teal darken-1" elevation="0" small outlined>
                    <v-icon x-small class="ml-1 mr-3">far fa-save</v-icon>
                    Сохранить
                </v-btn>
            </v-card>
        </v-toolbar>


        <v-row>
            <v-col cols="6">
                <v-text-field
                        label="Id изображения" dense class="body-1"
                        :value="item.id"
                        disabled
                />
                <v-switch
                        class="mb-6 mt-0"
                        dense
                        :label="item['enabled'] ? 'Отображается' : 'Скрыт'"
                        hide-details
                        light
                        :success="item['enabled']"
                        v-model="item['enabled']"
                ></v-switch>

                <v-text-field
                        label="Порядок отображения"  class="body-1"
                        type="number"
                        v-model="item['priority']"
                />
            </v-col>
            <v-col cols="6">
                <v-text-field
                        label="Имя файла (не забудьте про расширение!)"  class="body-1"
                        v-model="item['name']"
                />
                <v-text-field
                        label="ALT-tag"  class="body-1"
                        v-model="item['alt']"
                />
                <v-textarea
                        label="Short Description" dense class="body-1"
                        rows="3"
                        v-model="item['short_description']"
                />
            </v-col>
            <v-col cols="12">
                <h3 class="mb-3">Описание изображения</h3>
                <prism-editor v-model="item['description']" class="myPrismEditor mb-3 pb-2" language="html"  autoStyleLineNumbers max-height></prism-editor>
            </v-col>
        </v-row>

        <!-- Подтверждение удаления -->
        <v-dialog
                v-model="dialog.open"
                @keydown.esc="dialog.open = false"
                max-width="290"
        >
            <v-card>
                <v-card-title class="headline">Подтвердите удаление</v-card-title>

                <v-card-text class="d-flex flex-column">
                    <span class="mb-1"><strong>Id</strong>: {{dialog.id}}</span>
                    <span><strong>Имя: </strong>{{dialog.name}}</span>
                </v-card-text>

                <v-card-actions>
                    <v-spacer></v-spacer>

                    <!-- Отменить -->
                    <v-btn
                            color="primary"
                            class="mr-2"
                            small
                            text
                            outlined
                            @click.stop="dialog.open = false"
                    >
                        Отменить
                    </v-btn>

                    <!-- Удалить -->
                    <v-btn
                            color="red darken-1"
                            small
                            text
                            outlined
                            @click.stop="deleteItemDialog"
                    >
                        Удалить
                    </v-btn>

                </v-card-actions>
            </v-card>
        </v-dialog>
    </td>
</template>

<script>
    export default {
        data: () => ({
            uploadedFile: null, // загружаемый файл изображения
            options: {},
            dialog:{
                id: null,
                name: '',
                short_description: null,
                open: false,
            },
            item: {},
        }),
        props: ['file', 'headersLength'],
        async mounted() {
            this.item = this.file
            await this.getDataFromApi(false)
        },
        methods: {
            openDeleteDialog(item) {
                this.dialog.id = item['id']
                this.dialog.name = item['name']
                this.dialog.short_description = item['short_description']
                this.dialog.open = true
            },
            updateItemData: async function (item, shopNotification = false) {

                let payload = item
                Object.assign(payload, {
                    accountHashId: this.$store.state.account.hash_id
                })

                await this.$api.storage.update(payload)
                    .then((resp) => {
                        if (resp['file'] !== undefined) {
                            this.item = resp['file']
                            if (shopNotification) {
                                this.$notify({
                                    group: 'main',
                                    title: 'Файл обновлен',
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
                    .catch( (err) => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка обновления',
                            text: err['message'],
                            type: 'error'
                        });
                    })
            },

            async deleteItemDialog() {
                let payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    id: this.dialog.id
                }

                await this.$api.storage.delete(payload)
                    .then(()=> {
                        this.$emit('removeItem', payload)

                        this.$notify({
                            group: 'main',
                            title: 'Файл удален',
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
                        this.dialog.open = false
                        this.dialog.hash_id = ''
                        this.dialog.name = ''
                    });

                // обновляем текущий список
                // await this.getDataFromApi(false)

            },
            async uploadImage(item) {
                await this.uploadFile(item['id'], "products").then((resp)=>{
                    this.images.push(resp['file'])
                })
            },

            getDataFromApi: async function (shopNotification) {

                this.loading = true;

                let payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    id: this.item.id,
                }
                await this.$api.storage.get(payload)
                    .then((resp) => {
                        if (resp['file'] !== undefined && resp['file']!== null) {
                            this.item = resp['file']
                            if (shopNotification) {
                                this.$notify({
                                    group: 'main',
                                    title: 'Файл загружен',
                                    type: 'success'
                                });
                            }

                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Не полные данные от сервера',
                                type: 'warring'
                            });
                        }

                    })
                    .catch(err => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка удаления',
                            text: err['message'],
                            type: 'error'
                        });
                    })
                    .finally(()=> this.loading = false)
            },


        },
        name: "ProductImage"
    }
</script>