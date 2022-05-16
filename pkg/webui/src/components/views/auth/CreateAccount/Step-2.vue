<template>
  <auth-container>
    <email-verification-form
            header='Регистрация аккаунта<br /><small class="text--disabled font-weight-medium">Шаг 2/3: активация учетной записи</small>'
            :queryToken="queryToken" @input="validated"></email-verification-form>
  </auth-container>
</template>

<script>
  export default {

    data: () => ({
    }),
    props: ['queryToken'],
    mounted(){
//      console.log("Query: ", this.queryToken)
    },
    methods: {
      validated: function () {
//        console.log(this.goToLogin)
        // Если вдруг аккаунтов больше чем 0 => даем вариант выбрать один из аккаунтов
        let route_name = 'login.sign-up.step-3'
        if (this.$store.getters['user/listOfAccounts'].length > 0) route_name = 'login.accounts'

        return this.$router.push({ name: route_name }).catch(() => this.goToLogin())

      },
    },
    components: {
      AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
      EmailVerificationForm: ()=> import('@/components/auth/form/EmailVerification.vue')
    }
  }


</script>