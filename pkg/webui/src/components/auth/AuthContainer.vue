<template>
    <v-container class="d-flex fill-height pa-0" fluid >
        <v-row class="pa-0 ma-0 fill-height">

            <v-col class="d-flex flex-column pa-0" :class="{'col-lg-5 col-xl-4' : oneThird}">

                <!-- Тут основное действие: форма авторизации, выбор аккаунта и т.д.  -->
                <div class="d-flex fill-height justify-center">
                    <slot />
                </div>

                <!-- Футер с информацией о правах -->
                <div class="d-flex body-2 align-end mb-5 disabled grey--text justify-center">
                    <copyright/>
                </div>
            </v-col>

            <!-- Новости и обновления -->
<!--            <v-col class="pa-0 hidden-md-and-down" v-if="dual || oneThird">-->
            <v-col class="col justify-center hidden-md-and-down" :class="{'col-lg-7 col-xl-8' : oneThird}" v-if="dual || oneThird">
                <slot name="right" />
            </v-col>

        </v-row>
    </v-container>
</template>

<script>
	export default {
		data: () => ({
		}),
		props: {
			dual: Boolean,
			oneThird: Boolean
		},
		components: {
			Copyright: () => import ('@/components/legal/Copyright.vue')
		}
	}
</script>