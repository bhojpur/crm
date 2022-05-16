<template>
    <v-form v-model="formValid" class="form-auth form-sign-up mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp}" ref="form" @submit.prevent="submit">
        <app-logo />

        <template v-if="false">
            <v-card text flat class="mb-8">

                <v-card-title class="mb-3 justify-center title title text-center" >
                    <h1>Регистрация аккаунта<br>
                        <small class="text--disabled font-weight-medium">Шаг 1/3: создание пользователя</small>
                    </h1>
                </v-card-title>

                <v-card-text class="text-center pa-1">
                    Регистрация возможна только по приглашениям от других участников системы. Уже есть аккаунт? <router-link slot="action" :to="{name: 'login'}" class="a-link">Войти</router-link>
                </v-card-text>

            </v-card>

            <!-- Блок для вывода ошибок -->
            <v-expand-transition>
                <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">
                    <v-card-text class="text-center text--darken-2 red--text ">
                        <span v-text="msgError"></span>
                    </v-card-text>
                </v-card>
            </v-expand-transition>

            <v-text-field label="Имя пользователя" class="mt-4 body-1"
                          hint="Используйте латинские буквы, цифры. Можно использовать email-адрес."
                          autofocus
                          :error-messages="errors.username || usernameErrors"
                          v-model.trim="username"
                          @input="$v.username.$touch(); errors.username = '';"
                          @blur="$v.username.$touch()"
            />

            <v-text-field label="Телефон" class="mt-4 body-1"
                          hint="Укажите номер в международном формате: +7 xxx xxx-xx-xx"
                          :error-messages="errors.phone || phoneErrors"
                          v-model.trim="phone"
                          @input="$v.phone.$touch(); errors.phone = '';"
                          @blur="$v.phone.$touch()"
            />

            <v-text-field label="Email" class="mt-4 body-1"
                          v-model.trim="email"
                          :error-messages="errors.email || emailErrors"
                          @input="$v.email.$touch();errors.email = '';"
                          @blur="$v.email.$touch()"
            />

            <v-text-field   class="mt-4" label="Пароль" v-model="password"
                            :type="showPassword ? 'text' :'password'"
                            :error-messages="errors.password || passwordErrors"
                            @input="$v.password.$touch(); errors.password = ''"
                            @blur="$v.password.$touch()"
                            @focusin="showPwdRule = true"
                            @focusout="showPwdRule = false"
                            required>
                <v-icon small slot="append" class="mt-1 mr-0" @click="showPassword = !showPassword">{{showPassword ? 'far fa-eye' : 'far fa-eye-slash'}}</v-icon>
            </v-text-field>
            <v-expand-transition>
                <div v-if="showPwdRule" class=" text-left body-2">
                    <!--                <p>Используйте латинские буквы [a-Z], цифры [0-9] и спецсимволы. </p>-->
                    <div class="text--darken-2">Используйте латинские буквы, цифры и спецсимволы.</div>
                    <v-row class="pa-4" >

                        <li class="mb-1 pa-0 col-6 "
                            :class="!$v.password.minLength || !$v.password.required ? 'error--text' : 'text--disabled'">
                            Минимум {{$v.password.$params.minLength.min}} символов
                        </li>

                        <!--One lowercase character-->
                        <li class="mb-1 pa-0 col-6 "
                            :class="!$v.password.oneLowercase || !$v.password.required ? 'error--text' : 'text--disabled'">
                            Одна буква в нижнем регистре
                        </li>

                        <!--One uppercase character-->
                        <li class="mb-1 pa-0 col-6 "
                            :class="!$v.password.oneUppercase || !$v.password.required ? 'error--text' : 'text--disabled'">
                            Одна заглавная буква
                        </li>

                        <!--One number-->
                        <li class="mb-1 pa-0 col-6 "
                            :class="!$v.password.oneNumber || !$v.password.required ? 'error--text' : 'text--disabled'">
                            Одна цифра
                        </li>

                        <!--One Special-->
                        <li class="mb-1 pa-0 col-6 "
                            :class="!$v.password.oneSpecial || !$v.password.required ? 'error--text' : 'text--disabled'">
                            Один спецсимвол
                        </li>

                    </v-row>
                </div>
            </v-expand-transition>

            <v-text-field label="Код приглашения" class="mt-4 body-1"
                          v-if="userRegistrationInviteOnly"
                          v-model="inviteToken"
                          :error-messages="errors.inviteToken || inviteTokenErrors"
                          @input="$v.inviteToken.$touch(); errors.inviteToken = ''"
                          @blur="$v.inviteToken.$touch()"
            />

            <div class="body-2 pa-0 mt-8">
                Нажав на кнопку "Приступить к работе" вы создаете учетную запись в BhojpurCRM и принимаете <router-link slot="terms-of-use" :to="{name: 'legal.terms-of-use'}" class="a-link" v-text="'Правила пользования сервисом'" /> и <router-link slot="privacy-policy" :to="{name: 'legal.privacy-policy'}" class="a-link" v-text="'Политику конфидициальности'" />.
            </div>

            <v-btn type="submit" class="mt-5 text-none white--text " block color="teal darken-1" v-text="'Приступить к работе'"
                   :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" />
        </template>
        <template v-else>
            <v-card text flat class="mb-8">

                <v-card-title class="mb-3 justify-center title title text-center" >
                    <h1>Регистрация аккаунта<br>
                        <small class="text--disabled font-weight-medium">Шаг 1/3: создание пользователя</small>
                    </h1>
                </v-card-title>

                <v-card-text class="text-center pa-1">
                    К сожалению, регистрация новых пользователей доступна только через администратора аккаунта. Приносим свои извенения за предоставленыне неудобства.
                </v-card-text>

                <v-card-actions class="justify-center">
                    <router-link slot="action" :to="{name: 'login'}" class="a-link">Вернуться на страницу входа</router-link>
                </v-card-actions>


            </v-card>

        </template>


    </v-form>

</template>

<script>

	const { required, minLength, maxLength, email} = require('vuelidate/lib/validators')

	export default {

		validations () {
			let arr = {
				password: {
					required,
					minLength: minLength(8),
					maxLength: maxLength(32),
					oneLowercase: (password) => /[a-z]/.test(password),
					oneUppercase: (password) => /[A-Z]/.test(password),
					oneNumber: (password) => /[0-9]/.test(password),
					oneSpecial: (password) => /[!@#$%^&*(),.?":{}|<>_+=]/g.test(password),
				},
				email: {
					required,email
				},
                phone: {
	                required,
	                minLength: minLength(11),
	                maxLength: maxLength(32),
	                isNumber: (number) => /^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$/g.test(number),
                },
				username: {
					required,
					minLength: minLength(3),
					maxLength: maxLength(32),
					allowedCharacters: (username) => /^[a-zA-Z0-9@.\-_]*$/.test(username),
				},
            };
			if (this.userRegistrationInviteOnly) {
				arr['inviteToken'] = {
					required,
					minLength: minLength(9),
					maxLength: maxLength(60)
				}
            }
			return arr


		},
		data: () => ({
            phone: (process.env.NODE_ENV === 'production') ? '' : '+79251952295',
//            phone: '+19251952295',
			email: (process.env.NODE_ENV === 'production') ? '' : 'info@bhojpur-marketing.net',
//	        email: '',
			username: (process.env.NODE_ENV === 'production') ? '' : 'superAdmin',
//	        username: '',
			password: (process.env.NODE_ENV === 'production') ? '' : 'qwerty109#QW',
//	        password: '',
			inviteToken: (process.env.NODE_ENV === 'production') ? '' : 'fdslkfdslfsjdfnsfdnsdf',
			errors: {
				email: '',
				username: '',
				password: '',
				inviteToken: '',
                phone: '',
			},
			showPassword: false, // скрытие поля пароля (type = text / password)
			showPwdRule: false, // скрытие блока с правилами заполнения пароля

			formValid: false, // сохраняет состояние валидации формы ?
			msgError: '',     // текст ошибки в заголовке
			httpError: false, // ошибка уровня http
			loading: false

		}),
		methods: {
			submit: function () {

				this.$v.$touch()
				this.clearErrors()

				if ( !this.validate() || this.$v.$invalid) return

				// фиксируем данные для отправки из переменной form
				const { email, username, password, inviteToken } = this

				this.loading = true

				// посылаем запрос на создание пользователя
				this.$store.dispatch('user/CREATE',{ email, username, password, inviteToken })

					/*.then( (data) => {
						this.$store.dispatch('app/SET_TOKEN', data.token).catch()
						this.$store.dispatch('user/SET_PROFILE', data.user).catch()
						this.$router.push({ name: 'login.sign-up.step-2' }).catch(err => console.error(err))
					})
					.catch( (err) => this.showErrors(err))
					.finally(() => (this.loading = false))*/
			},

			/**
			 * Валидирует форму и выводит ошибки (встроенный метод)
			 * @returns Object
			 */
			validate: function () {
				this.$refs.form.validate()
				return this.formValid
			},

			// очищает ошибки
			clearErrors: function () {
				this.$v.$reset()
				this.httpError = false
				this.msgError = ''
				this.errors = {
					email: '',
					username: '',
					password: '',
					inviteToken: '',
          phone: '',
				}
			},

			showErrors: function (err, transError) {

				this.httpError = true // ошибка уровня http

				if (transError === undefined) {
					transError = 'Неизвестная ошибка сети'
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
		computed: {
			passwordErrors () {
				const errors = []
				if (!this.$v.password.$dirty) return errors
				!this.$v.password.maxLength && errors.push(`Максимум ${this.$v.username.$params.maxLength.max} символов`)
				return errors
			},
			emailErrors () {
				const errors = []
				if (!this.$v.email.$dirty) return errors
				!this.$v.email.required && errors.push('Email-адрес обязателен')
				!this.$v.email.email && errors.push('Должен быть корректный email-адрес')
				if (this.errors.email.length > 0) errors.push(this.errors.email)
				return errors
			},
			usernameErrors () {
				const errors = []
				if (!this.$v.username.$dirty) return errors
				!this.$v.username.allowedCharacters && errors.push('Используйте только разрешенные символы')
				!this.$v.username.required && errors.push('Имя пользователя обязательно')
				!this.$v.username.minLength && errors.push(`Минимум ${this.$v.username.$params.minLength.min} символа`)
				!this.$v.username.maxLength && errors.push(`Максимум ${this.$v.username.$params.maxLength.max} символов`)
				if (this.errors.username.length > 0) errors.push(this.errors.username)
				return errors
			},
			inviteTokenErrors () {
				const errors = [];
//				if (!this.userRegistrationInviteOnly) return errors

				if (!this.$v.inviteToken.$dirty) return errors
				!this.$v.inviteToken.required && errors.push('Код приглашения обязательно для заполнения')
				!this.$v.inviteToken.minLength && errors.push(`Слишком короткий код: минимум ${this.$v.inviteToken.$params.minLength.min} символов`)
				!this.$v.inviteToken.maxLength && errors.push(`Слишком длинный код: максимум ${this.$v.inviteToken.$params.maxLength.max} символов`)

				if (this.errors.inviteToken.length > 0) errors.push(this.errors.inviteToken)
				return errors
			},
			phoneErrors () {
				const errors = [];

				if (!this.$v.phone.$dirty) return errors;

				!this.$v.phone.required && errors.push('Телефон обязателен к заполнению')
				!this.$v.phone.isNumber && errors.push('Проверьте корректность ввода')

				!this.$v.phone.minLength && errors.push(`Слишком короткий номер: минимум ${this.$v.phone.$params.minLength.min} символов`)
				!this.$v.phone.maxLength && errors.push(`Слишком длинный номер: максимум ${this.$v.phone.$params.maxLength.max} символов`)

				if (this.errors.phone.length > 0) errors.push(this.errors.phone)
				return errors
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
			userRegistrationInviteOnly() {
				return !this.$store.state.app.settings['this.$store.state.app.settings']!== undefined && this.$store.state.app.settings.user_registration_invite_only
			}
		},
		components: {
			AppLogo: () => import('@/components/auth/AppLogo.vue'),
		}
	}
</script>

<style lang="scss">


</style>
