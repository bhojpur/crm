<template>
  <div>
    <span v-for="(tag, index) in tags"><small style="color: #045e54;padding-right: 2px;">#</small>{{tag.label}}
      <template v-if="index < tags.length-1">,&nbsp;</template>
  </span>
  </div>
</template>

<script>

export default {
  data: () => ({}),

  props: {
    tags: {
      type: Array, Object,
      default: [],
    },
  },
}
</script>