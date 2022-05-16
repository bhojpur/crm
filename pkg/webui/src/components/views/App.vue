<template>
  <router-view />
</template>

<script>
  import axios from 'axios'
export default {

  name: 'App',
  data: () => ({}),

  async created() {

    this.$vuetify.theme.dark = false

    // 1. Устанавливаем для Axios token или стираем, если его нет
    if (this.$store.getters['app/getToken'] !== undefined && this.$store.getters['app/getToken'].length > 3) {
      axios.defaults.headers.common['Authorization'] = 'Bearer ' + this.$store.getters['app/getToken']
    } else {
      axios.defaults.headers.common['Authorization'] = ''
    }

    // 2. Проверяем есть ли токен
    if ( this.$store.getters['user/isAuth'] && this.$store.getters['account/isAuth'] ) {

      // 3. Чекаем авторизацию.
      // await this.checkAuthIfNotLogout()

      const payload = {
        hash_id: this.$store.state.account.hash_id
      }
      await this.$store.dispatch('account/GET', payload)
              .then(()=>{})
              .catch(async (err)=>{
                // window.$broadcastTab.message({command: "logout"})
                await this.logout()
              })

    } else {
      let item = localStorage.getItem('message');
      let message = JSON.parse(item);
      // console.log("MSG: ", message)
      // await this.logout()
      localStorage.removeItem("vuex")
      await this.$store.dispatch("AUTH_LOGOUT", {root: true})
      this.$router.push({ name: 'login',  params: { msg: 'Авторизуйтесь в системе' } }).catch()
    }

    // 3. Проверяем авторизацию


    // 2. Загружаем первичные настройки приложения
    await this.$store.dispatch('app/GET_SETTINGS')
            .then(() => {
              // ...
            })
            .catch(() => {
              /*await this.$store.dispatch('app/AUTH_LOGOUT')
                      .then()
                      .catch()
                      .finally(()=> this.$router.push({name: 'login', params: { msg: 'Неудалось загрузить настройки клиента' }}))*/

            })


  },
}
</script>
