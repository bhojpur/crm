<template>
    <v-container fluid>
        <v-card style="border-top-left-radius: 2px;border-top-right-radius: 2px;">

            <v-toolbar color="white" flat>
                <v-toolbar-title class="d-flex flex-column pt-6">
                    <h1 style="font-size: 26px;">Редактировать шаблон</h1>
                    <v-btn @click="$router.push({name: 'email-marketing.templates'})" class="ml-0 pl-0" outlined small text>
                        <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
                        <span style="padding-top: 2px">Назад к списку шаблонов</span>
                    </v-btn>
                </v-toolbar-title>
                <v-spacer></v-spacer>
                <v-card class="d-flex align-center" flat v-if="!loadingTemplate">

                    <!-- Сохранение и создание шаблона -->
                    <v-btn @click="sendTemplate" class="mr-2" color="blue lighten-5 text--primary"  elevation="1" small disabled>
                        <v-icon class="ml-1 mr-3" small>far fa-paper-plane</v-icon>
                        Отправить тест
                    </v-btn>

                  <v-divider class="mx-2 mt-0" inset vertical></v-divider>

                    <v-btn @click="copyTemplate" color="secondary" dark small>
                        <v-icon class="ml-1 mr-3" small>fal fa-clone</v-icon>
                        Клонировать
                    </v-btn>

                    <v-divider class="mx-2 mt-0" inset vertical></v-divider>

                    <v-btn @click="uploadItem(true)" color="secondary" small>
                        <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
                        Обновить данные
                    </v-btn>

                    <v-divider class="mx-2 mt-0" inset vertical></v-divider>

                    <v-btn @click="updateItemData" color="cyan darken-4" dark small>
                        <v-icon class="mr-3" small>fal fa-save</v-icon>
                        Сохранить
                    </v-btn>

                </v-card>

            </v-toolbar>
            <v-row class="mt-9 ml-1">
                <v-col cols="12" sm="6">
                    <v-text-field label="Имя шаблона" placeholder="Введите имя шаблона" v-model.trim="item.name"></v-text-field>
                    <v-textarea label="Описание (назначение) шаблона" placeholder="Укажите назначние шаблона" v-model.trim="item.description" rows="3"></v-textarea>
                </v-col>
                <v-col sm="4" сols="11" offset="1">
                    <v-card class="d-flex flex-column" flat>
                        <v-checkbox
                                :label="'Предпросмотр | id: ' + item.hash_id"
                                class="mb-0 pb-0"
                                v-model="item.public"
                        ></v-checkbox>
                        <div v-if="item.public">
<!--                            <a :href="$store.getters['storage/getCDN_EMAILS_RAW_URL'] + '/'+item.hash_id"-->
                            <a :href="$api.storage.getCDN_EMAILS_RAW_URL() + '/'+item.hash_id"
                               class="d-inline-flex"
                               style="text-decoration: none;color:#3c6eb6 !important;border-bottom: 1px dotted #3c6eb6;"
                               target="_blank">
                                link raw
                            </a>
                        </div>
                        <div v-if="item.public">
                            <a
                                    :href="$api.storage.getCDN_EMAILS_COMPILE_URL() + '/'+item.hash_id"
                                    style="text-decoration: none;color:#3c6eb6 !important;border-bottom: 1px dotted #3c6eb6;"
                                    target="_blank"
                            >
                                link compile
                            </a>
                        </div>
                        <span class="deep-orange--text darken-3" v-if="!item.public">скрыт</span>
                    </v-card>

                </v-col>
            </v-row>

            <!-- Редактор HTML -->
            <section class="px-2" v-show="!loadingTemplate">

                <v-card class="d-flex justify-space-between mb-1 mx-1" flat>
                    <v-card-title>
                        <h2 class="mb-1">Код шаблона</h2>
                    </v-card-title>

                    <v-card-actions>
                        <v-btn
                                @click="item.html_data = $api.emailTemplate.getBaseDataTemplate()"
                                outlined small text>Set HTML tpl
                        </v-btn>
                    </v-card-actions>

                </v-card>

                <prism-editor class="myPrismEditor mb-2 pb-2" language="html" v-model="item['html_data']"></prism-editor>
            </section>

            <!-- Прогресс загрузки -->
            <section class="d-flex align-center text-center justify-center pa-4 flex-column"
                     v-if="loadingTemplate"
            >
                <v-progress-circular
                        :size="40"
                        color="teal"
                        indeterminate
                ></v-progress-circular>

                <span class="mt-2 body-2">Подождите, шаблон загружается...</span>

            </section>

        </v-card>
    </v-container>

</template>

<script>
    export default {
        data: () => ({
            item: {
                id: null,
                hash_id: '',
                name: '',
                code: '',
                public: false,
            },
            // newTpl: true, // работаем с новым шаблоном или нет
            loadingTemplate: false, // абстрактный статус загрузки..

        }),
        async mounted() {
            await this.uploadItem()
        },
        methods: {
            copyTemplate: function () {

                let payload = {
                    accountHashId: this.$store.state.account.hash_id,
                }
                Object.assign(payload, this.item)

                this.$api.emailTemplate.create(payload)
                    .then(resp => {
                        this.$notify({
                            group: 'main',
                            title: 'Шаблон создан!',
                            text: 'Шаблон успешно скопирован',
                            type: 'success'
                        });
                        if (resp['email_template'] !== undefined || resp['email_template'] !== null) this.$router.push({
                            name: 'email-marketing.templates.edit',
                            params: {public_id: resp['email_template']['public_id']}
                        })
                        else {
                            this.$set(this, 'tpl', resp['email_template'])
                        }
                    })
                    .catch((err) => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка создания',
                            text: 'Не удалось скопировать шаблон',
                            type: 'error'
                        });
                    })


            },
            sendTemplate: function () {
                // 1. open dialog test
                const payload = {
                    accountHashId: this.$store.state.account.hash_id,
                    hash_id: this.tpl.hash_id,
                    UserId: 2,
                    EmailBoxId: 4,
                    Subject: "Очередь работает!"
                }
                this.$store.dispatch('marketing/EMAIL_TEMPLATE_SEND_TO_USER', payload)
                    .then(resp => {
                        this.$notify({
                            group: 'main',
                            title: 'Шаблон отправлен!',
                            text: 'Шаблон успешно отправлен',
                            type: 'success'
                        });
                    })
                    .catch((err) => {
                        this.$notify({
                            group: 'main',
                            title: 'Ошибка отправки',
                            text: 'Не удалось отправить шаблон',
                            type: 'error'
                        });
                    })


            },
            uploadItem: async function (showNotification) {
                this.loading = true
                return this.$api.emailTemplate.getByPublicId({
                    accountHashId: this.$store.state.account.hash_id,
                  public_id: this.$route.params.public_id
                })
                    .then(resp => {
                        if (resp['email_template'] !== undefined) {
                            this.item = resp['email_template']
                          if (showNotification) {
                            this.$notify({
                              group: 'main',
                              title: 'Данные обновлены',
                              type: 'success'
                            });
                          }
                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Данные от сервера получены не полностью',
                                type: 'warring'
                            });
                        }
                    })
                    .catch(err => {
                        this.$notify({
                            group: 'main',
                            title: 'Шаблон не найден',
                            text: err['message'],
                            type: 'error'
                        });
                    })
                    .finally(() => this.loading = false)

            },
            updateItemData: async function () {

                let payload = this.item

                Object.assign(payload, {accountHashId: this.$store.state.account.hash_id})

                await this.$api.emailTemplate.update(payload)
                    .then((resp) => {
                        if (resp['email_template'] !== undefined) {
                            this.item = resp['email_template']
                            this.$notify({
                                group: 'main',
                                title: 'Шаблон обновлен',
                                type: 'success'
                            });

                        } else {
                            this.$notify({
                                group: 'main',
                                title: 'Данные от сервера не полные',
                                text: 'не хватает объекта [email_template]',
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
