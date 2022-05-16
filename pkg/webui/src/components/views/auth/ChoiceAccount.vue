<template>
    <auth-container>
        <v-card flat class="mt-8" >

            <app-logo />

            <v-card-title class="mb-0 justify-center title title text-center" >
                <h1>Добро пожаловать!</h1>
            </v-card-title>

            <div class="text-center pa-1 mb-8">
                <template v-if="hasAccounts">
                    Выберите аккаунт или <router-link slot="sign-up" :to="{name: 'login.create-account'}" class="a-link">создайте новый</router-link>
                </template>
                <template v-else>
                    Нет доступных аккаунтов. <router-link slot="sign-up" :to="{name: 'login.create-account'}" class="a-link" disabled>Создать новый?</router-link>
                </template>
            </div>

            <!-- Блок для вывода ошибок -->
            <v-expand-transition>
                <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">
                    <v-card-text class="text-center text--darken-2 red--text ">
                        <span v-text="msgError"></span>
                    </v-card-text>
                </v-card>
            </v-expand-transition>

            <!-- блок выбора аккаунта с опциями и прехедер текстом -->
            <v-list>

                <!-- Выводим через template список аккаунтов -->
                <template v-for="(item,index) in aUsers" v-if="item.role.tag !== 'client' || item.account['visible_to_clients']">

                    <!-- Выводим разделитель, если это не первый элемент -->
                    <v-divider inset v-if="index !== 0"></v-divider>

                    <!-- кликабельный блок с информацией об аккаунте -->
                    <v-list-item :key="item.account_id" @click="login(item['account']['hash_id'])" >

                        <!-- Аватарка компании, пока просто выводим иконку -->
                        <v-list-item-avatar>
                            <v-icon small class="mt-0 mr-0">{{'far fa-chart-network'}}</v-icon>
                        </v-list-item-avatar>

                        <!-- блок с основной информацией о компании: имя компании и роль пользователя в ней -->
                        <v-list-item-content>
                            <v-list-item-title v-text="aUsers[index].account.name" class="subtitle-2"></v-list-item-title>
                            <!--                            <v-list-item-subtitle v-text="index ? 'manager' : 'admin'"></v-list-item-subtitle>-->
                            <!--                            <v-list-item-subtitle v-text="aUsers[index].role.name" class="caption"></v-list-item-subtitle>-->
                            <v-list-item-subtitle v-text="aUsers[index].role.name" class="body-2"></v-list-item-subtitle>
                        </v-list-item-content>

                        <!-- псевдокликабельный блок с иконкой входа -->
                        <v-list-item-action >
                            <v-icon small class="">fas fa-sign-in</v-icon>
                        </v-list-item-action>

                    </v-list-item>

                    <v-progress-linear
                            v-if="accountHashId === aUsers[index].account.id"
                            :active="accountHashId.length > 0"
                            :indeterminate="true"
                            class="ma-0"
                            height="4"
                            style="top: -2px;"
                    ></v-progress-linear>

                    <!-- Выводим разделитель, если это не первый элемент -->
                    <v-divider inset v-if="index === aUsers.length"></v-divider>

                </template>

                <!-- Дополнительные элементы управления -->
                <v-list-item class="pa-0 mt-8 justify-space-between" >

                    <!-- Предлагаем запомнить выбор аккаунта -->
                    <v-checkbox class="justify-center caption"
                                v-if="hasAccounts"
                                :outlined="false"
                                color="teal darken-2"
                                v-model="rememberChoice"
                                label="Запомнить выбор?"
                    />

                    <!-- Выход из системы на оконо логина-->
                    <v-btn text outlined small @click="goToAuthView"  color="link">Выйти из системы</v-btn>

                </v-list-item>

            </v-list>
        </v-card>
    </auth-container>
</template>

<script>
	export default {
		data: () => ({
			accountHashId: '',
			rememberChoice: false,
			errorStatus: '', // ошибка сообщение

			formValid: false, // сохраняет состояние валидации формы ?
			msgError: '',     // текст ошибки в заголовке
			httpError: false, // ошибка уровня http
			loading: false
		}),
		components: {
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			AppLogo: () => import('@/components/auth/AppLogo.vue'),
		},
        mounted() {

        },
        methods: {
			login: function (accountHashId) {

			    // console.log(hash_id)
                // return
				this.errorStatus = ''
				const rememberChoice = this.rememberChoice

                let payload = {
                    accountHashId: accountHashId,
                    rememberChoice
                }

				this.$store.dispatch('account/LOGIN_IN_ACCOUNT', payload)
					.then(() => {
					    this.$router.push({ name: 'home' })
                    })
					.catch( (err) => this.showErrors(err))
					.finally(() => (this.loading = false))

			},
			showErrors: function (err, transError) {

				this.httpError = true // ошибка уровня http

				if (transError === undefined) {
					transError = "Неизвестная ошибка соединения"
				}

				this.msgError = (err ? err.message : false) || transError

				// выводим конкретизированные по полям ошибки
				if (err.errors) {
					this.httpError = false // ошибка в полях данных
					for (let field in err.errors) {
						this.errors[field] = err.errors[field];
					}
				}
			},
		},
		created: function () {
//            console.log(this.$store.state.auth.token)
		},

		computed: {
			username () {
				let profile = this.$store.state.user.profile
				return profile ? profile.name : ''
			},
			accounts () {
				return this.$store.state.user.accounts
//				return this.$store.state.user.aUsers
			},
			aUsers () {
				return this.$store.state.user.aUsers
			},
            phone() {
	            return this.$store.state.user.phone
            },
			/**
			 * Определяет есть ли ошибки со стороны сервера
			 * @returns {boolean}
			 */
			hasErrorsFromServer() {
				for (let key in this.errors) {
					if (this.errors[key].length > 0) return true
				}
				return this.httpError;

			},
            hasAccounts(){
			    return this.$store.getters['user/listAccountUsers'].length > 0
            },
		}
	}
</script>