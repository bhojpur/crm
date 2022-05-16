<template>
    <v-app-bar app dense elevate-on-scroll
               color="#1b2c2f"
               fixed
               height="46px"
               :dance="true"
               style="border-bottom: 1px solid rgb(79, 86, 93);border-left: none !important;left:0;"
    >
        <h3 class="white--text">
            <v-icon color="white" class="mr-2" small>
                fal fa-badge-check
            </v-icon>
            {{$store.state.account.name}}
        </h3>

        <v-spacer></v-spacer>

        <v-menu offset-y
                :open-on-hover="false"
                style="border-radius: 0 !important;"
                class="mr-2"
        >
            <template v-slot:activator="{ on }">
                <v-btn
                        color="primary"
                        v-on="on"
                        dense
                        dark
                        depressed
                        text
                        class="text-capitalize"
                        style="color: white !important"
                >
                    <v-icon class="mr-3" color="green lighten-1">far fa-user-shield</v-icon>
                    {{ userProfile.name}} {{ userProfile['surname']}}
                    <v-icon class="ml-4" dark>fal fa-angle-down</v-icon>
                </v-btn>

            </template>

            <!-- Выпадающее меню -->
            <v-list
                    :nav="false"
                    :dense="true"
                    style="border-radius: 0 !important;"
            >

                <!-- Настройки аккаунта -->
                <v-list-item
                        v-for="(item, index) in items"
                        :key="index"
                        :to="{name: item.route_name}"
                        link
                        :exact="true"
                        :dence="true"
                >
                    <v-list-item-icon style="margin-right: 12px !important;">
                        <v-icon small>{{ item.icon }}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item-content>

                </v-list-item>

                <!-- Смена аккауна -->
                <v-list-item
                        link
                        color="primary"
                        style="height: 10px;"
                        @click="goToLoginInAccount()"
                >
                    <v-list-item-icon style="margin-right: 12px !important;">
                        <v-icon>fal fa-exchange</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>Сменить аккаунт</v-list-item-title>
                    </v-list-item-content>

                </v-list-item>

                <v-divider class="mb-2"></v-divider>

                <!-- Выход из приложения -->
                <v-list-item
                        link
                        color="primary"
                        style="height: 10px;"
                        @click="goToAuthView()"
                >
                    <v-list-item-icon style="margin-right: 12px !important;">
                        <v-icon>fal fa-sign-out</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>Выход</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>

            </v-list>
        </v-menu>

    </v-app-bar>
</template>

<script>
    export default {
        data: () => ({
            langs: ['en', 'ru'],
            drawer: false,
            fullHeader: false,
            /* выпадающее меню */
	        items: [
		        { title: 'Настройки', route_name: 'account', icon: 'fal fa-user-cog', },
	        ],
        }),
        props: {
            msg: String
        },
        methods: {

        },
        computed: {
            navItems: function() {
                return this.$store.state.nav.mainItems
            },
	        userProfile: function() {
        		return this.$store.state.user.profile
            },
            account: function () {
                return this.$store.state.account.data
            }
        }
    }
</script>