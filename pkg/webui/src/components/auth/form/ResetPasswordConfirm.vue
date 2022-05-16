<!-- Форма авторизации по токену -->
<template>
    <v-form v-model="formValid" class="form-auth form-sign-up mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp}" ref="form" @submit.prevent="submit">

        <app-logo />

        <v-card text flat class="mb-8">

            <v-card-title class="mb-3 justify-center title text-center">
                <h1>Сброс старого пароля</h1>
            </v-card-title>

            <v-card-text class="text-center pa-1">
                Введите проверочный код для сброса вашего пароля
            </v-card-text>

        </v-card>


        <!-- Блок вывода информации -->
        <v-expand-transition>
            <v-card flat class="blue lighten-5 text-center mb-6" v-if="notification">
                <v-card-text class="text-center text--darken-2  ">
                    <span v-text="notification"></span>
                </v-card-text>
            </v-card>
        </v-expand-transition>

        <!-- Блок вывода ошибок -->
        <v-expand-transition>
            <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">

                <v-card-text class="text-center text--darken-2 red--text ">
                    <span v-text="headerError"></span>
                </v-card-text>

            </v-card>
        </v-expand-transition>

        <v-text-field label="Код из письма" autofocus class="mt-2 body-1"
                      v-model.trim="token"
                      :error-messages="errors.token || tokenErrors"
                      @input="$v.token.$touch(); errors.token = '';"
                      @blur="$v.token.$touch()"
        />

        <v-btn type="submit" class="mt-5 text-none white--text " block color="teal darken-1" v-text="'Подтвердить'"
               :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" />

        <!-- Только для авторизованных пользователей -->
        <div class="text-center body-1 pa-0 mt-8" v-if="userIsAuth">
            <v-btn small type="button" @click.prevent="resend">Отправить повторно код?</v-btn>
        </div>
    </v-form>

</template>

<script>

	const { required, minLength, maxLength} = require('vuelidate/lib/validators')

	export default {
		validations: {
			token: {
				required,
				minLength: minLength(12)
			},
		},
		data: () => ({
			token: '',
//			token: '1ukyryxpfprxpy17i4ldlrz9kg3',
			errors: {
				token: '',
			},

			headerError: '',     // текст ошибки в заголовке
            formValid: false, // сохраняет состояние валидации формы ?
			httpError: false, // ошибка уровня http
			loading: false,
			notification: '', // текст уведомления: ''
		}),
		props: ['queryToken','msgNotification'],
        computed: {
	        tokenErrors () {
		        const errors = []
		        if (!this.$v.token.$dirty) return errors
		        !this.$v.token.required && errors.push('Введите код подтверждения')
		        !this.$v.token.minLength && errors.push(`Минимум ${this.$v.token.$params.minLength.min} символов`)
		        if (this.errors.token.length > 0) errors.push(this.errors.token)
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

            userIsAuth(){
	            return this.$store.getters['user/isAuth']
            },
        },
		mounted(){

			// сбарсываем авторизацию
			//this.logout()

			if ( (this.msgNotification !== undefined && this.msgNotification !== '') ) {
				console.log("msgNotification: " + this.msgNotification)
				this.notification = this.msgNotification
			}

			if (this.queryToken && this.queryToken.length > 0) this.token = this.queryToken

			// если есть токен, пробуем сразу по нему зайти и сбросить пароль
			if (this.queryToken !== undefined && this.queryToken.length > this.$v.token.$params.minLength.min) {
				console.log("submit!")
				this.submit()
			}

		},
		methods: {
			submit: function () {

				this.$v.$touch()
				this.clearErrors()

				if ( !this.validate() || this.$v.$invalid) return

				// фиксируем данные для отправки из переменной form
				const { token } = this

				this.loading = true

				// посылаем запрос на создание пользователя verifications
				this.$store.dispatch('user/PASSWORD_RESET_CONFIRM',{ token })
					.then( (resp) => {this.$emit('done', resp)})
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
				this.notification = ''
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
