<template>
  <auth-container>
    <form-forgot-password @done="(q) => $router.push({name: 'login.password-recovery', params: { msgNotification: q.msgNotification }})"/>
  </auth-container>
</template>

<script>
  export default {
    components: {
      AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
      FormForgotPassword: () => import ('@/components/auth/form/ForgotPassword.vue'),
//      FormAuthByToken: () => import ('@/components/auth/form/AuthByToken.vue'),
    }
  }
</script>