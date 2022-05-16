<template>
    <v-form v-model="formValid" class="form-auth form-login align-self-center pa-3" ref="form" @submit.prevent="submit" :autocomplete="onceLogin ? 'off' : 'on'">

        <app-logo />

        <v-card text flat class="mb-5">

            <v-card-title class="mb-1 justify-center title title text-center" >
                <h1>Вход в Bhojpur CRM</h1>
            </v-card-title>

            <v-card-text class="text-center pa-1">
                Нет аккаунта? <router-link slot="sign-up" :to="{name: 'login.sign-up'}" class="a-link">Создать аккаунт</router-link>
            </v-card-text>

        </v-card>

        <!-- Блок вывода информации -->
        <v-expand-transition>
            <v-card flat class="blue lighten-5 text-center mb-6" v-if="message">
                <v-card-text class="text-center text--darken-2  ">
                    <span v-text="message"></span>
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

        <v-text-field label="Имя пользователя" autofocus class="mt-2 body-1"
                      v-model.trim="username"
                      :error-messages="errors.username || usernameErrors"
                      @input="$v.username.$touch();errors.username = '';"
                      @blur="$v.username.$touch()"
                      required
        />

        <v-text-field   class="mt-4 body-1" label="Пароль" v-model="password"
                        :type="showPassword ? 'text' :'password'"
                        :error-messages="errors.password || passwordErrors"
                        @input="$v.password.$touch(); errors.password = ''"
                        @blur="$v.password.$touch()"
                        required>
            <v-icon small slot="append" class="mt-1 mr-0" @click="showPassword = !showPassword">{{showPassword ? 'far fa-eye' : 'far fa-eye-slash'}}</v-icon>
        </v-text-field>

        <v-btn type="submit" class="mt-5 text-none white--text " block color="teal darken-1" v-text="'Войти'"
               :disabled="!formValid || $v.$invalid" :class="{'disable-events': !formValid}" :loading="loading" />

        <div class="d-flex justify-center mb-0">
            <v-checkbox class="mt-4 caption"
                        :outlined="false"
                        color="teal darken-1"
                        v-model="onceLogin"
                        label="Чужой компьютер"
            ></v-checkbox>
        </div>

        <v-row class="justify-space-between mt-0">

            <v-col class="col-auto text-left">
                <router-link :to="{name: 'login.forgot-username'}" class="headline a-link body-1" v-text="'Забыли имя пользователя?'" />
            </v-col>
            <v-col class="col-auto text-right">
                <router-link :to="{name: 'login.forgot-password'}" class="headline a-link body-1" v-text="'Забыли пароль?'" />
            </v-col>
        </v-row>

    </v-form>
</template>

<script>
    const { required, minLength, maxLength} = require('vuelidate/lib/validators')

    export default {
        data: () => ({
	        username: (process.env.NODE_ENV === 'production') ? '' : 'admin',
	        password: (process.env.NODE_ENV === 'production') ? '' : 'qwerty109#QW',
            onceLogin: true,
	        errors: {
		        username: '',
		        password: '',
	        },

            showPassword: false, // отображения пароля в ввода

	        formValid: false, // сохраняет состояние валидации формы ?
	        msgError: '',     // текст ошибки в заголовке
	        httpError: false, // ошибка уровня http
	        loading: false
        }),
	    validations: {
		    username: {
                required,
			    maxLength: maxLength(32),
			    allowedCharacters: (username) => /^[a-zA-Z0-9@.\-_]*$/.test(username),
		    },
		    password: {
                required,
			    minLength: minLength(8),
			    maxLength: maxLength(32),
		    }
	    },
        props: ["message"],
        methods: {
            validate () {
                this.$refs.form.validate()
                return this.formValid
            },

            // вызывается по кнопке login
            submit: function () {

            	// очищаемся от ошибок
	            this.clearErrors()

                // если форма не валидируется - выходим, ошибки сами выведутся
                if ( !this.validate() ) return

                // фиксируем данные для отправки из переменной form
                const { username, password, onceLogin } = this

                // посылаем запрос на авторизацию пользователя
	            this.$store.dispatch('user/AUTH_BY_USERNAME', {username, password, onceLogin})
		            .then(() => {
                        if (this.$store.getters['user/defaultAccountId'] > 0)
                        	return this.loginInAccount(this.$store.getters['user/defaultAccountId'], {name:'home'})

			            if (this.$store.getters['user/listOfAccounts'] !== undefined && this.$store.getters['user/listOfAccounts'].length === 1)
			            	return this.loginInAccount(this.$store.getters['user/listOfAccounts'][0].id, {name:'home'})


		            	this.goToLoginInAccount()
		            })
		            .catch( (err) => this.showErrors(err))
		            .finally(() => (this.loading = false))

            },

            // выводит ошибки поверх формы и в полях формы
	        showErrors: function (err, transError) {

		        this.httpError = true // ошибка уровня http

		        if (transError === undefined) {
			        transError = "Неизвестная ошибка соединения"
		        }

		        this.msgError = (err ? err.message : false) || transError

		        // выводим конкретизированные по полям ошибки
		        if (err['errors'] !== undefined && err.errors) {
			        this.httpError = false // ошибка в полях данных
			        for (let field in err.errors) {
				        this.errors[field] = err.errors[field];
			        }
		        }
	        },

            // очищает вывод ошибок
            clearErrors: function () {
		        this.$v.$reset()
		        this.httpError = false
		        this.msgError = ''
		        this.errors = {
			        username: '',
			        password: '',
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
	        usernameErrors () {
		        const errors = []
		        if (!this.$v.username.$dirty) return errors
		        !this.$v.username.allowedCharacters && errors.push('Используйте только разрешенные символы')
		        !this.$v.username.required && errors.push('Имя пользователя обязательно')
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
        },

    }
</script>
