<template>
    <v-form v-model="formValid" class="form-auth form-sign-up mt-8" :class="{'mt-4': $vuetify.breakpoint.smAndUp}" ref="form" @submit.prevent="submit">

        <app-logo/>

        <template v-if="false">
            <v-card-title class="mb-3 justify-center title title text-center" v-html="header" ></v-card-title>

            <v-card-text class="text-center pa-1 mb-4">
                Придумайте имя аккаунта и заполните пустые поля.<br>Вернуться к <router-link slot="sign-up" :to='{name: "login.accounts"}' class="a-link">списку аккаунтов</router-link>.
            </v-card-text>

            <!-- Блок для вывода ошибок -->
            <v-expand-transition>
                <v-card flat class="red lighten-5 text-center mb-6" v-if="hasErrorsFromServer">
                    <v-card-text class="text-center text--darken-2 red--text ">
                        <span v-text="headerError"></span>
                    </v-card-text>
                </v-card>
            </v-expand-transition>

            <v-text-field label="Имя аккаунта" autofocus class="mt-2 body-1"
                          v-model.trim="name"
                          :error-messages="errors.name || nameErrors"
                          @input="$v.name.$touch();errors.name = '';"
                          @blur="$v.name.$touch()"
            />

            <v-text-field label="Основной сайт компании" class="mt-2 body-1"
                          v-model.trim="website"
                          :error-messages="errors.website || websiteErrors"
                          @input="$v.website.$touch();errors.website = '';"
                          @blur="$v.website.$touch()"
            />

            <v-select :items="items" v-model.trim="type" class="mt-4 body-1" label="Выберите отрасль" />

            <v-btn type="submit" class="mt-5 text-none white--text " block color="teal darken-1" v-text="'Создать аккаунт'"
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
	const { required, minLength, maxLength} = require('vuelidate/lib/validators')
	import Vue from 'vue'
	export default {
		validations: {
			name: {
				required,
				minLength: minLength(2),
				maxLength: maxLength(64)
			},
			website: {
				maxLength: maxLength(255)
			},
		},
		data: () => ({
			name: (process.env.NODE_ENV === 'production') ? '' : 'company name',
			website: (process.env.NODE_ENV === 'production') ? '' : 'web site',
			type: "",
			items: ['IT технологии', 'Торговля', 'Производство', 'Архитектура'],
			errors: {
				name: '',
				website: '',
				type: '',
                phone: '',
			},
			headerError: '',     // текст ошибки в заголовке
			formValid: false, // сохраняет состояние валидации формы ?
			httpError: false, // ошибка уровня http
			loading: false,
		}),
        props: ['header'],
		methods: {
			submit: function () {

				this.$v.$touch()
				this.clearErrors()

				if ( !this.validate() || this.$v.$invalid) return

				// фиксируем данные для отправки из переменной form
				const { name, website, type } = this

				this.loading = true

				// посылаем запрос на создание пользователя verifications
				this.$store.dispatch('account/CREATE',{ name, website, type })
					.then( (resp) => {this.$emit('input', resp)} )
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
				this.errors = {
					name: '',
					website: '',
					type: '',
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
		computed: {
			nameErrors () {
				const errors = []
				if (!this.$v.name.$dirty) return errors
				!this.$v.name.required && errors.push('Введите  подтверждения')
				!this.$v.name.minLength && errors.push(`Минимум ${this.$v.name.$params.minLength.min} символа`)
				!this.$v.name.maxLength && errors.push(`Максимум ${this.$v.name.$params.maxLength.max} символов`)
				if (this.errors.name.length > 0) errors.push(this.errors.name)
				return errors
			},

			websiteErrors () {
				const errors = []
				if (!this.$v.name.$dirty) return errors
				!this.$v.website.maxLength && errors.push(`Максимум ${this.$v.website.$params.maxLength.max} символов`)
				if (this.errors.website.length > 0) errors.push(this.errors.website)
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
		components: {
			AppLogo: () => import('@/components/auth/AppLogo.vue'),
		}
	}
</script>
