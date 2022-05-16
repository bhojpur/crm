<template>
  <auth-container>
    <form-forgot-username />
  </auth-container>
</template>

<script>
  export default {
    components: {
      AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
      FormForgotUsername: () => import ('@/components/auth/form/ForgotUsername.vue'),
    }
  }
</script>