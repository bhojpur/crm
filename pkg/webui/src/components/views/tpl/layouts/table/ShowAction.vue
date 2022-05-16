<template>
  <div style="width: 100%;display: flex;justify-content: center; margin-left: 5px;">
    <!--{{editOnly}}-->
    <v-btn :to="route" depressed text class="mx-0" small :class="editOnly ? 'px-8':'px-0'">
      <v-icon class="blue--text text--lighten-2" small>
        far fa-edit
      </v-icon>
    </v-btn>
    <v-btn v-if="!editOnly" @click="$emit('delete')" depressed text class="px-0 mx-0" small :disabled="disabledRemove" >
      <v-icon class="red--text text--lighten-2" small>
        {{removeIcon ? 'far fa-unlink' : 'fas fa-trash'}}
      </v-icon>
    </v-btn>
  </div>
</template>

<script>

export default {
  data: () => ({}),
  props: {
    route: {
      type: [Object, Array],
      default: 0,
    },
    removeIcon: {
      type: Boolean,
    },
    editOnly: {
      type: Boolean,
      default: false,
    },
    marginRight: {
      type: Boolean,
      default: true,
    },
    disabledRemove: {
      type: Boolean,
      default: false,
    },
    disabledEdit: {
      type: Boolean,
      default: false,
    },
  },
}
</script>