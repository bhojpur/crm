<template>
    <v-card flat class="pa-0 ma-0">
        <v-data-table
            :headers="headers"
            :items-per-page="10"
            :items.sync="images"
            class="elevation-1"
            show-expand
            expand-icon="fal fa-angle-down"
            sort-by="id"

        >

            <template v-slot:top>
                <v-toolbar color="white" flat>
                    <v-toolbar-title><h1 style="font-size: 26px;">Видео товара</h1></v-toolbar-title>

                    <v-spacer></v-spacer>
                    <v-card class="d-flex align-center" flat>

                        <v-btn @click="getDataFromApi(true)"
                               color="secondary"
                               dark
                               elevation="1"
                               small
                        >
                            <v-icon
                                    class="ml-1 mr-3"
                                    x-small
                            >far fa-sync</v-icon>
                            Обновить
                        </v-btn>

                        <v-divider
                                class="mx-2 mt-0"
                                inset
                                vertical
                        ></v-divider>


                        <v-card flat class="pa-0 ma-0">
                            <v-file-input @change.native="uploadImage(item)" ref="inputUpload" label="Выберите файл" show-size dense hide-details v-model="uploadedFile"></v-file-input>
                        </v-card>
                    </v-card>

                    <!--<div class="mt-2 d-flex">
                        <v-card flat max-width="300px" style="">
                            <v-card-title class="headline">Выберите файл (max 32mb)</v-card-title>

                            <v-card-text>
                                <v-file-input @change.native="uploadImage(item)" ref="inputUpload" label="Выберите файл" show-size v-model="uploadedFile"></v-file-input>
                            </v-card-text>
                        </v-card>
                    </div>-->

                </v-toolbar>

            </template>

            <template v-slot:item.img="{ item }">
                <v-img :src="item['url']" aspect-ratio="1" max-width="160px"></v-img>
            </template>

            <template v-slot:item.name="{ item }">
                <a :href="item['url']" target="_blank">{{item['name']}}</a>
            </template>
            <template v-slot:item.size="{ item }">
                {{formatBytes(item['size'])}}
            </template>

            <template v-slot:item.enabled="{ item }">
                <span v-text="item['enabled'] ? 'виден' : 'скрыт'"></span>
            </template>

            <template v-slot:item.actions="{ item }">

                <v-btn
                        @click="openDeleteDialog(item)" depressed small
                        style="background-color: transparent">
                    <v-icon
                            class="red--text text--lighten-2"
                            small
                    >
                        fas fa-trash
                    </v-icon>
                </v-btn>

            </template>

            <template v-slot:expanded-item="{ headers, item }">
                <td :colspan="headers.length">
                    <h3 class="mt-4 mb-6">Параметры карточки</h3>
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
                                @click="updateFileData(item)"
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
                            <v-list-item-title class="headline" style="font-size: 14px;">Описание: &laquo;{{dialog.model}}&raquo;</v-list-item-title>
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
    </v-card>
</template>

<script>
    import Vue from "vue";

    export default {
        data: () => ({
            headers: [
                { text: '', align: 'start', value: 'img'},
                { text: 'Имя файла', align: 'start', value: 'name'},
                { text: 'Порядок', value: 'priority'},
                { text: 'Размер', value: 'size'},
                { text: 'enabled', value: 'enabled'},
                { text: 'Удалить', value: 'actions'},
            ],
            uploadedFile: null, // загружаемый файл изображения
            images: [],
            upImg: false,
            dialog:{
                id: null,
                name: '',
                short_description: null,
                open: false,
            },
        }),
        props: ['item'],
        async mounted() {
            // console.log("item: ", this.item)
            await this.getDataFromApi(false)
        },
        methods: {
            openDeleteDialog(item) {
                this.dialog.id = item['id']
                this.dialog.name = item['name']
                this.dialog.short_description = item['short_description']
                this.dialog.open = true
            },
            deleteItemDialog() {

                this.dialog.open = false
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    id: this.dialog.id
                }
                this.$store.dispatch('storage/DELETE_FILE', payload)
                    .then(async ()=> {
                        let index = this.images.findIndex(item => item.id === payload.id, payload)
                        if (index !== -1) Vue.delete(this.images, index)

                        this.$notify({
                            group: 'main',
                            title: 'Изображение удалено удален',
                            type: 'success'
                        });

                        // обновляем текущий список
                        // await this.getDataFromApi(false)
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
                        this.dialog.short_description = ''
                    })

            },
            async uploadImage(item) {
                await this.uploadFile(item['id']).then((resp)=>{
                    // console.log("File: ", resp)
                    this.images.push(resp['file'])
                })
            },

            async getDataFromApi(showNotification) {
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    product_id: this.item['id'],
                    offset: 0,
                    limit: 20,
                }
                await this.getFilePaginationList(showNotification, payload).
                    then((resp) => {
                        if (resp['files'] !== undefined) this.images = resp['files']

                    })
            },

        },
        name: "ProductImage"
    }
</script>

<style scoped>

</style>