<template>
    <v-form v-model="formValid" class="form-auth form-recover mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp, 'mt-6': $vuetify.breakpoint.mdAndUp}"
            @submit.prevent="submit" autocomplete="off" ref="form">

        <app-logo />

        <v-card text flat class="mb-8">

            <v-card-title class="mb-3 justify-center title text-center" >
                <h1>Восстановление имени пользователя</h1>
            </v-card-title>

            <v-card-text class="text-center pa-1">Укажите email-адрес, на который зарегистрирован ваш аккаунт, и мы вышлем на него ваше имя пользователя.</v-card-text>

        </v-card>

        <!-- Блок вывода информации -->
        <v-expand-transition>
            <v-card flat class="blue lighten-5 text-center mb-6" v-if="msgNotification">
                <v-card-text class="text-center text--darken-2  ">
                    <span v-text="msgNotification"></span>
                </v-card-text>
            </v-card>
        </v-expand-transition>

        <!-- Блок для вывода ошибок -->
        <v-expand-transition>
            <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">
                <v-card-text class="text-center text--darken-2 red--text ">
                    <span v-text="msgError"></span>
                </v-card-text>
            </v-card>
        </v-expand-transition>

        <v-text-field label="Email" autofocus class="mt-4 body-1"
                      v-model.trim="email"
                      :error-messages="errors.email || emailErrors"
                      @input="$v.email.$touch();errors.email = '';"
                      @blur="$v.email.$touch()" />

        <v-btn type="submit" class="mt-8 text-none white--text " block color="teal darken-1"
               :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" >Отправить имя пользователя</v-btn>

        <v-row class="justify-center mt-2">
            <v-col class="col-auto">
                <router-link :to="{name: 'login'}" class="headline a-link" v-text="$t('button.recover.back-to-login')" />
            </v-col>
        </v-row>

    </v-form>
</template>

<script>
	const { required,email} = require('vuelidate/lib/validators')

    export default {
	    data: () => ({
		    email: 'info@bhojpur-marketing.net',
		    errors: {
			    email: '',
		    },

		    formValid: false, // сохраняет состояние валидации формы ?
            msgNotification: '', // текст уведомления
		    msgError: '',     // текст ошибки в заголовке
		    httpError: false, // ошибка уровня http
		    loading: false,
	    }),
	    validations: {
		    email: {
			    required,email
		    },
	    },
	    methods:{
		    submit: function () {

			    this.$v.$touch()
			    this.clearErrors()

			    if ( !this.validate() || this.$v.$invalid) return

			    // фиксируем данные для отправки из переменной form
			    const { email } = this

			    this.loading = true

			    // посылаем запрос на создание пользователя
			    this.$store.dispatch('user/SEND_USERNAME_RECOVERY',{ email })
				    .then( (resp) => {
					    this.clearErrors()
					    this.email = ""
					    this.loading = false

                        if (resp['message'] !== undefined ) {
                            this.msgNotification = resp['message']
                        }
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
			    this.msgError = ''
			    this.msgNotification = ''
			    this.errors = {
				    email: '',
				    username: '',
				    password: '',
				    inviteToken: '',
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
		    emailErrors () {
			    const errors = []
			    if (!this.$v.email.$dirty) return errors
			    !this.$v.email.required && errors.push('Email-адрес обязателен')
			    !this.$v.email.email && errors.push('Должен быть корректный email-адрес')
			    if (this.errors.email.length > 0) errors.push(this.errors.email)
			    return errors
		    },
		    hasErrorsFromServer() {
			    for (let key in this.errors) {
				    if (this.errors[key].length > 0) return true
			    }
			    return this.httpError;

		    },
	    },
	    components: {
		    AppLogo: () => import('@/components/auth/AppLogo.vue')
	    }

    }
</script>