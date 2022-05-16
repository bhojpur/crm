<template>
    <v-form v-model="formValid" class="form-auth form-recover mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp, 'mt-6': $vuetify.breakpoint.mdAndUp}"
            @submit.prevent="submit" ref="form">

        <app-logo />

        <v-card text flat class="mb-8">

            <v-card-title class="mb-3 justify-center title" >
                <h1>Восстановление пароля</h1>
            </v-card-title>

            <v-card-text class="text-center pa-1">
                Укажите ваш email и мы вышлем инструкцию по сбросу пароля на ваш email-адрес.
            </v-card-text>

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

        <v-text-field label="Имя пользователя" class="mt-4 body-1" autofocus
                      :error-messages="errors.username || usernameErrors"
                      v-model.trim="username"
                      @input="$v.username.$touch();errors.username = '';"
                      @blur="$v.username.$touch()"
                      />

        <v-btn type="submit" class="mt-8 text-none white--text " block color="teal darken-1"
               :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" >Сбросить пароль</v-btn>

        <v-row class="justify-center mt-2">
            <v-col class="col-auto">
                <router-link :to="{name: 'login'}" class="headline a-link" v-text="$t('button.recover.back-to-login')" />
            </v-col>
        </v-row>

    </v-form>
</template>

<script>
	const { required, minLength, maxLength, helpers, email , numeric} = require('vuelidate/lib/validators')
    export default {
        data: () => ({
	        username: 'admin',
	        errors: {
		        username: '',
	        },

	        formValid: false, // сохраняет состояние валидации формы ?
	        msgNotification: '', // текст уведомления
	        msgError: '',     // текст ошибки в заголовке
	        httpError: false, // ошибка уровня http
	        loading: false,
        }),
	    validations: {
		    username: {
			    required,
			    minLength: minLength(3),
			    maxLength: maxLength(32),
			    allowedCharacters: (username) => /^[a-zA-Z0-9@.\-_]*$/.test(username),
		    },
	    },
        methods:{
	        submit: function () {

		        this.$v.$touch()
		        this.clearErrors()

		        if ( !this.validate() || this.$v.$invalid) return

		        // фиксируем данные для отправки из переменной form
		        const { username } = this

		        this.loading = true

		        // посылаем запрос на создание пользователя
		        this.$store.dispatch('user/PASSWORD_RECOVERY_SEND_EMAIL',{ username })
			        .then( (resp) => {
				        this.username = ""
				        this.loading = false

				        if (resp['message'] !== undefined ) {
					        this.msgNotification = resp['message']
				        }

				        if (resp['message'] !== undefined )
					        this.$emit('done', {msgNotification: this.msgNotification})
                        else
					        this.$emit('done', {msgNotification: "Код отправлен на вашу почту"})

				        // todo: какое именно сообщение?
//				        this.$emit('done', {msg: this.msgNotification})
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
		        this.msgNotification = ''
		        this.msgError = ''
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