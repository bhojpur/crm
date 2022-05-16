<template>
  <span :class="color">
    {{ price ? price.toLocaleString('ru-RU',{currency:"RUB"}) : '-' }}{{price && showCurrency ? ' руб.': ''}}</span>
</template>

<script>

    export default {
        data: () => ({}),
        props: {
          price: {
            type: Number,
            default: 0
          },
          showCurrency: {
            type: Boolean,
            default: true
          },
          color: {
            type: String,
            default: ''
          },
        }
    }
</script>