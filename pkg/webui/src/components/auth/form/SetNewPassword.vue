<!-- Форма авторизации по токену -->
<template>
    <v-form v-model="formValid" class="form-auth form-sign-up mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp}" ref="form" @submit.prevent="submit">

        <app-logo />

        <v-card text flat class="mb-8">

            <v-card-title class="mb-3 justify-center title text-center">
                <h1>Создание нового пароля</h1>
            </v-card-title>

            <v-card-text class="text-center pa-1">
                Придумайте надежный пароль от вашей учетной записи.
            </v-card-text>

        </v-card>


        <!-- Блок вывода информации -->
        <v-expand-transition>
            <v-card flat class="blue lighten-5 text-center mb-6" v-if="notificationMsg">
                <v-card-text class="text-center text--darken-2  ">
                    <span v-text="notificationMsg"></span>
                </v-card-text>
            </v-card>
        </v-expand-transition>

        <v-expand-transition>
            <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">

                <v-card-text class="text-center text--darken-2 red--text ">
                    <span v-text="headerError"></span>
                </v-card-text>

            </v-card>
        </v-expand-transition>

        <!-- Поле ввода нового пароля -->
        <v-text-field   class="mt-4" label="Новый пароль" v-model="password"
                        :type="showPassword ? 'text' :'password'"
                        :error-messages="errors.password || passwordErrors"
                        @input="$v.password.$touch(); errors.password = ''"
                        @blur="$v.password.$touch()"
                        @focusin="showPwdRule = true"
                        @focusout="showPwdRule = false"
                        required>
            <v-icon small slot="append" class="mt-1 mr-0" @click="showPassword = !showPassword">{{showPassword ? 'far fa-eye' : 'far fa-eye-slash'}}</v-icon>
        </v-text-field>
        <div class=" text-left body-2">
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

        <v-btn type="submit" class="mt-5 text-none white--text " block color="teal darken-1" v-text="'Установить новый пароль'"
               :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" />

    </v-form>
</template>

<script>

	const { required, minLength, maxLength} = require('vuelidate/lib/validators')

	export default {
		validations: {
			password: {
				required,
				minLength: minLength(8),
				maxLength: maxLength(32),
				oneLowercase: (password) => /[a-z]/.test(password),
				oneUppercase: (password) => /[A-Z]/.test(password),
				oneNumber: (password) => /[0-9]/.test(password),
				oneSpecial: (password) => /[!@#$%^&*(),.?":{}|<>_+=]/g.test(password),
			},
		},
		data: () => ({
			password: 'qwerty109#QW',
			errors: {
				password: '',
			},

			showPassword: false, // скрытие поля пароля (type = text / password)
			showPwdRule: false, // скрытие блока с правилами заполнения пароля

            formValid: false, // сохраняет состояние валидации формы ?
			httpError: false, // ошибка уровня http
			loading: false,
			notificationMsg: ''
		}),
        computed: {
	        passwordErrors () {
		        const errors = []
		        if (!this.$v.password.$dirty) return errors
		        !this.$v.password.maxLength && errors.push(`Максимум ${this.$v.username.$params.maxLength.max} символов`)
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
        },
		mounted() {
			// дальнейшие манипуляции требуют действительного ключа пользователя

//			this.$store.dispatch('app/AUTH_CHECK_USER').catch(() => this.logout());
			/*this.$store.dispatch('app/AUTH_CHECK_USER').catch();

			if (!this.$store.state.user.profile && !this.$store.state.user.profile['password_reset']) {
				this.$router.push({name:'home', message:"Пароль уже установлен"})
            }*/

		},
		methods: {
			submit: function () {

				this.$v.$touch()
				this.clearErrors()

				if ( !this.validate() || this.$v.$invalid) return

				// фиксируем данные для отправки из переменной form
//				const { password } = this

				this.loading = true

				// посылаем запрос на создание пользователя verifications
				this.$store.dispatch('user/SET_PASSWORD',{ password_new: this.password })
					.then( (resp) => {

						this.notificationMsg = resp.message
						this.loading = true
						setTimeout(()=>this.$emit('done', resp), 800)

					})
					.catch( (err) => this.showErrors(err))
					.finally(() => (this.loading = false))
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
				this.headerError = ''
				this.notificationMsg = ''
				this.errors = {
					token: '',
				}
			},

			showErrors: function (err, transError) {

				this.httpError = true // ошибка уровня http

				if (transError === undefined) {
					transError = 'Неизвестная ошибка сети'
				}

				this.headerError = (err ? err.message : false) || transError

				// выводим конкретизированные по полям ошибки
				if (err.errors) {
					this.httpError = false // ошибка в полях данных
					for (let field in err.errors) {
						this.errors[field] = err.errors[field];
					}
				}
			},
		},
        components:{
	        AppLogo: () => import('@/components/auth/AppLogo.vue'),
        }
	}
</script>
