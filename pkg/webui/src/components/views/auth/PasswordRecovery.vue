<!-- Сбрасывает пароль и переадресует на страницу установки пароля -->
<template>
  <auth-container>
    <form-reset-password-confirm :queryToken="queryToken" :msgNotification="msgNotification" v-if="!$store.getters['user/isAuth']"/>
    <form-set-new-password v-else @done="goToLoginInAccount" />
  </auth-container>
</template>

<script>
  export default {
    data: () => ({
      loading: false
    }),
    props: ['queryToken','msgNotification'],
    components: {
      AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
      FormForgotPassword: () => import ('@/components/auth/form/ForgotPassword.vue'),
      FormResetPasswordConfirm: () => import ('@/components/auth/form/ResetPasswordConfirm.vue'),
      FormSetNewPassword: () => import ('@/components/auth/form/SetNewPassword.vue'),
    },
  }
</script>