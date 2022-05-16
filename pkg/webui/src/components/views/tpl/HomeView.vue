<template>
    <v-container fluid >
        <h1>Добро пожаловать!</h1>

    </v-container>
</template>

<script>

    import store from '../../../store'
    export default {
        data: () => ({

        }),
        beforeRouteEnter (to, from, next) {
            if (!store.getters['user/isAuth']) {
                next({name: 'login'})
                return
            }
            next()
        },
    }
</script>