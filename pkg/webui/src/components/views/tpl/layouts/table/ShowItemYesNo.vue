<template>
  <!--    <span :class="[show ? 'green&#45;&#45;text' : ({'deep-orange&#45;&#45;text text&#45;&#45;darken-1': bright}), {'grey&#45;&#45;text text&#45;&#45;darken-1': disabledView}]"-->
  <span :class="[show ? ClassYes : ({'deep-orange--text text--darken-1': bright}), disabled ? 'grey--text text--darken-1': '']"
        v-text="show ? 'да' : 'нет'">
    </span>
</template>

<script>

export default {
  data: () => ({}),

  props: {
    show: {
      type: Boolean,
      default: false,
    },
    bright: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    ClassYes: {
      type: String,
      default: 'green--text',
    },
    ClassNo: {
      type: String,
      default: '',
    },
  },
}
</script>