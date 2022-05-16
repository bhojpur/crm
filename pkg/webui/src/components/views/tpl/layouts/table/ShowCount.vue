<template>
   <span
       :class="{ 'deep-orange--text': !(items && items.length > 0), 'mr-3' : marginRight }"
       v-text="(items && items.length > 0) ? items.length : '0'"></span>
</template>

<script>

    export default {
        data: () => ({}),
        props: {
          items: {
            type: [Object,Array],
            default: 0
          },
          marginRight: {
            type: Boolean,
            default: true
          },
        },
    }
</script>