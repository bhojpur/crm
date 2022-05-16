<template>
    <auth-container>
        <form-set-new-password />
    </auth-container>
</template>

<script>
	export default {
		components:{
			AppLogo: () => import('@/components/auth/AppLogo.vue'),
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			FormSetNewPassword: () => import ('@/components/auth/form/SetNewPassword.vue'),
		}
	}
</script>