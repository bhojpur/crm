<template>
    <auth-container one-third>
        <form-auth-login :msg="msg"/>
        <template #right>
            <whats-new/>
        </template>
    </auth-container>
</template>

<script>
	export default {
		components: {
			AuthContainer: () => import ('@/components/auth/AuthContainer.vue'),
			FormAuthLogin: () => import ('@/components/auth/form/Login.vue'),
			WhatsNew: () => import ('@/components/auth/WhatsNew.vue'),
		},
        props: ["msg"]
	}
</script>