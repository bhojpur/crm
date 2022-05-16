<template>
  <auth-container>
    <form-create-user />
  </auth-container>
</template>

<script>
  export default {
    components: {
      AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
      FormCreateUser: () => import ('@/components/auth/form/CreateUser.vue'),
    },
    mounted() {
      // тут надо подгрузить настройки главного аккаунта
      // console.log(this.$store.state.app.settings)
    },
  }
</script>