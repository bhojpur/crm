<template>
    <auth-container>
        <form-create-account
                header='<h1>Регистрация аккаунта<br /><small class="text--disabled font-weight-medium">Завершите создание нового аккаунта</small></h1>'
                @input="afterCreated" />

    </auth-container>
</template>

<script>

	export default {

		data: () => ({

		}),

        methods: {
			afterCreated: function (resp) {
				this.$store.dispatch('app/SET_TOKEN', resp.token).catch()
				this.$store.dispatch('account/SET_PROFILE', resp.account).catch()
				return this.$router.push({ name: 'home',  query: { msg: 'Добро пожаловать!' } }).catch()
			},
		},
		components: {
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			FormCreateAccount: () => import ('@/components/auth/form/CreateAccount.vue'),
		}
	}
</script>
