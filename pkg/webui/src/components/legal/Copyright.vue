<template>
    <div class="text-center"  style="max-width: 440px;padding-top: 10px;padding-bottom: 10px;">
        <router-link slot="use-cookie" :to="{name: 'legal.use-cookie'}" class="grey-text a-link-grey" v-text="'Использование Cookie'" />
        <router-link slot="terms-of-use" :to="{name: 'legal.terms-of-use'}" class="grey-text a-link-grey" v-text="'Правила пользования сервисом'" />
        <router-link slot="privacy-policy" :to="{name: 'legal.privacy-policy'}" class="grey-text a-link-grey" v-text="'Политика конфиденциальности'" />
    </div>
</template>


<style scoped lang="scss">
    .a-link-grey {
        color: inherit;
        transition: .12s ease-in-out;
        &:hover {
            color: var(--v-anchor-base) !important;
        }
    }
</style>