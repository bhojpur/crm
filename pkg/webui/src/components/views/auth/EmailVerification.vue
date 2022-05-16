<template>
    <auth-container>
        <email-verification-form
                header='Подтвердите Email<br /><small class="text--disabled font-weight-medium">Завершите создание учетной записи</small>'
                :queryToken="queryToken" @input="validated"></email-verification-form>
    </auth-container>
</template>

<script>
	export default {

		data: () => ({
		}),
		props: ['queryToken'],
		methods: {
			validated: function () {
				// Если все хорошо, делаем переход на выбор аккаунтов
				this.$router.push({ name: 'login.accounts' }).catch(() => this.goToLogin())

			},
		},
		components: {
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			EmailVerificationForm: ()=> import('@/components/auth/form/EmailVerification.vue')
		}
	}
</script>
