<template>
  <v-form :class="{'mt-4': $vuetify.breakpoint.smAndUp}" @submit.prevent="submit" class="form-auth form-sign-up mt-8"
          ref="form" v-model="formValid">

    <app-logo/>

    <v-card-title class="mb-8 justify-center title text-center">
      <h1 v-html="header"></h1>
    </v-card-title>

    <!-- Блок вывода информации -->
    <v-expand-transition>
      <v-card class="blue lighten-5 text-center mb-6" flat v-if="notificationMsg">
        <v-card-text class="text-center text--darken-2  ">
          <span v-text="notificationMsg"></span>
        </v-card-text>
      </v-card>
    </v-expand-transition>

    <v-expand-transition>
      <v-card class="red lighten-5 text-center mb-6" flat v-if="hasErrorsFromServer">

        <v-card-text class="text-center text--darken-2 red--text ">
          <span v-text="headerError"></span>
        </v-card-text>

      </v-card>
    </v-expand-transition>

    <v-text-field :error-messages="errors.token || tokenErrors" @blur="$v.token.$touch()" @input="$v.token.$touch(); errors.token = '';"
                  autofocus
                  class="mt-2 body-1"
                  label="Код из письма"
                  v-model.trim="token"
    />

    <v-btn :class="{'disable-events': !formValid}" :disabled="!formValid || $v.$invalid" :loading="loading" block class="mt-5 text-none white--text "
           color="teal darken-1" type="submit" v-text="'Подтвердить'"/>

    <!-- Только для авторизованных пользователей -->
    <div class="text-center body-1 pa-0 mt-8" v-if="userIsAuth">
      <v-btn @click.prevent="resend" small type="button">Отправить повторно код?</v-btn>
    </div>
  </v-form>

</template>

<script>

const {required, minLength, maxLength} = require('vuelidate/lib/validators')

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
    notificationMsg: ''
  }),
  props: ['queryToken', 'header'],
  computed: {
    tokenErrors() {
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

    userIsAuth() {
      return this.$store.getters['user/isAuth']
    },
  },
  mounted() {
//			this.$testGoGo()
    if (this.queryToken && this.queryToken.length > 0) this.token = this.queryToken
  },
  methods: {
    submit: function () {

      this.$v.$touch()
      this.clearErrors()

      if (!this.validate() || this.$v.$invalid) return

      // фиксируем данные для отправки из переменной form
      const {token} = this

      this.loading = true

      // посылаем запрос на создание пользователя verifications
      this.$store.dispatch('user/EMAIL_VERIFICATION', {token})
          .then(async (resp) => {

            /**
             * Чтобы установить новый токен нужно осовободиться от текущих данных и загрузить новые.
             * Если токен уже есть (новая регистрация 1-2-3 шага, то нужно просто идти дальше)
             */
            if (resp['token'] !== undefined) {
              await this.$store.dispatch("AUTH_LOGOUT", {root: true})
              await this.$store.dispatch("app/SET_TOKEN", resp['token'])
            }

            // дальнейшие манипуляции требуют действительного ключа пользователя
            this.$store.dispatch('app/AUTH_CHECK_USER').catch(() => this.goToLogin())

            if (resp['user'] !== undefined)
              await this.$store.dispatch("user/SET_PROFILE", resp['user'])
            else
              this.$store.dispatch("user/UPDATE_PROFILE", resp['user']).catch(() => this.goToLogin())


            if (resp['accounts'] !== undefined)
              await this.$store.dispatch("user/SET_AVAILABLE_ACCOUNTS", resp['accounts'])
            else
              this.$store.dispatch("user/UPDATE_AVAILABLE_ACCOUNTS", resp['user']).catch(() => this.goToLogin())


            this.$emit('input', resp)
          })
          .catch((err) => {
            this.showErrors(err)
          })
          .finally(() => (this.loading = false))
    },

    resend: function () {

      this.$v.$touch()
      this.clearErrors()

//				if ( !this.validate() || this.$v.$invalid) return

      // посылаем запрос на создание пользователя verifications
      this.$store.dispatch('user/SEND_EMAIL_INVITE_VERIFICATION')
          .then(async (resp) => {

                if (resp['accounts'] !== undefined && resp['user'] !== undefined) {

                  if (resp['token'] !== undefined) {
                    await this.$store.dispatch("AUTH_LOGOUT", {root: true})
                    await this.$store.dispatch("app/SET_TOKEN", resp['token'])
                  }

                  // дальнейшие манипуляции требуют действительного ключа
                  this.$store.dispatch('app/AUTH_CHECK_USER').catch(() => this.goToLogin())

                  if (resp['user'] !== undefined)
                    await this.$store.dispatch("user/SET_PROFILE", resp['user'])
                  else
                    this.$store.dispatch("user/UPDATE_PROFILE", resp['user']).catch(() => this.goToLogin())


                  if (resp['accounts'] !== undefined)
                    await this.$store.dispatch("user/SET_AVAILABLE_ACCOUNTS", resp['accounts'])
                  else
                    this.$store.dispatch("user/UPDATE_AVAILABLE_ACCOUNTS", resp['user']).catch(() => this.goToLogin())

                  this.$emit('input', resp)
                } else {
                  this.notificationMsg = resp['message']
                }
              }
          )
          .catch((err) => this.showErrors(err))
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
  components: {
    AppLogo: () => import('@/components/auth/AppLogo.vue'),
  }
}
</script>
