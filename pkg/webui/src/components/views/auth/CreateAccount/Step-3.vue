<template>
    <auth-container>
        <form-create-account
                header='<h1>Регистрация аккаунта<br><small class="text&#45;&#45;disabled font-weight-medium">Шаг 3/3: создание аккаунта</small></h1>'
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
			AppLogo: () => import('@/components/auth/AppLogo.vue'),
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			FormCreateAccount: () => import ('@/components/auth/form/CreateAccount.vue'),
		}
	}
</script>