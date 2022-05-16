<template>
    <v-card flat class="pa-0 ma-0">
        <v-data-table
            :headers="headers"
            :options.sync="options"
            :items.sync="images"
            class="elevation-1"
            show-expand
            expand-icon="fal fa-angle-down"
            sort-by="id"
        >

            <template v-slot:top>
                <v-toolbar color="white" flat>
                    <v-toolbar-title><h3 style="font-weight: 500;">Изображения товара</h3></v-toolbar-title>

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

                </v-toolbar>

            </template>

            <template v-slot:item.img="{ item }">
                <v-img :src="item['_url']" aspect-ratio="1" max-width="160px"></v-img>
            </template>

            <template v-slot:item.name="{ item }">
                <a :href="item['_url']" target="_blank">{{item['name']}}</a>
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
                <image-expand-template :file.sync="item" :headersLength="headers.length" @removeItem="removeItemFromItems"></image-expand-template>
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
            options: {},
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
            removeItemFromItems(payload) {
                let index = this.images.findIndex(item => item.id === payload.id, payload)
                if (index !== -1) Vue.delete(this.images, index)
                else console.log("Image not deleted")
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
                            title: 'Изображение удалено',
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
                await this.uploadFile(this.item['id'],"products")
                    .then((resp)=> {
                        if (resp['file'] !== undefined && resp['file']!== null) {
                            this.uploadedFile = null
                            this.$notify({
                                group: 'main',
                                title: 'Файл создан',
                                type: 'success'
                            });
                            this.images.unshift(resp['file'])

                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Не полные данные от сервера',
                                type: 'warring'
                            });
                        }
                    })
                    .catch((err) => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка при создании файла',
                            text: err['message'],
                            type: 'error'
                        });
                    })
                    .finally(()=> {
                        this.uploadedFile = null
                    });
            },

            async getDataFromApi(showNotification) {
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    owner_id: this.item['id'],
                    owner_type: "products",
                    offset: 0,
                    limit: this.options['itemsPerPage'],
                }
                await this.$api.storage.getListPagination(payload).
                    then((resp) => {
                        if (resp['files'] !== undefined) this.images = resp['files']

                    })
            },

        },
        name: "ProductImage",
        components: {
            ImageExpandTemplate: () => import('@/components/views/tpl/content/layouts/ImagesExpandTemplate'),
        },
    }
</script>