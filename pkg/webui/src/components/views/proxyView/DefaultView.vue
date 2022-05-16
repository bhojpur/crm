<template>
  <v-app v-if="isAuth">
    <template>
      <app-header />
      <left-menu />
    </template>

    <v-main style="background-color: #f5f8fa !important;">
      <router-view/>
    </v-main>

    <notifications group="main" position="bottom right" />
  </v-app>

</template>

<script>
    export default {
      data: () => ({

      }),
      components: {
        LeftMenu: () => import('@/components/navigation/LeftMenu.vue'),
        AppHeader: () => import('@/components/navigation/AppHeader.vue')
      },
      computed: {
        isAuth() {
          return this.$store.getters['user/isAuth'] && this.$store.getters['account/isAuth']
        }
      }


  };
</script>