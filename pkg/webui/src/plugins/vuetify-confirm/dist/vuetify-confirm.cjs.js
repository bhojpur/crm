'use strict';

var lib = require('vuetify/lib');

var Confirm = { render: function () {
    var _vm = this;var _h = _vm.$createElement;var _c = _vm._self._c || _h;return _c('v-dialog', { attrs: { "eager": "", "value": "true", "max-width": _vm.width, "persistent": _vm.persistent }, on: { "input": _vm.change, "keydown": function ($event) {
          if (!('button' in $event) && _vm._k($event.keyCode, "esc", 27, $event.key, "Escape")) {
            return null;
          }_vm.choose(false);
        } } }, [_c('v-card', { attrs: { "tile": "" } }, [Boolean(_vm.title) ? _c('v-toolbar', { attrs: { "dark": "", "color": _vm.color, "dense": "", "flat": "" } }, [Boolean(_vm.icon) ? _c('v-icon', { attrs: { "left": "" } }, [_vm._v(_vm._s(_vm.icon))]) : _vm._e(), _vm._v(" "), _c('v-toolbar-title', { staticClass: "white--text", domProps: { "textContent": _vm._s(_vm.title) } })], 1) : _vm._e(), _vm._v(" "), _c('v-card-text', { staticClass: "body-1 text-body-1 py-3", domProps: { "innerHTML": _vm._s(_vm.message) } }), _vm._v(" "), _c('v-card-actions', [_c('v-spacer'), _vm._v(" "), Boolean(_vm.buttonFalseText) ? _c('v-btn', { attrs: { "color": _vm.buttonFalseColor, "text": _vm.buttonFalseFlat }, on: { "click": function ($event) {
          _vm.choose(false);
        } } }, [_vm._v(" " + _vm._s(_vm.buttonFalseText) + " ")]) : _vm._e(), _vm._v(" "), Boolean(_vm.buttonTrueText) ? _c('v-btn', { attrs: { "color": _vm.buttonTrueColor, "text": _vm.buttonTrueFlat }, on: { "click": function ($event) {
          _vm.choose(true);
        } } }, [_vm._v(" " + _vm._s(_vm.buttonTrueText) + " ")]) : _vm._e()], 1)], 1)], 1);
  }, staticRenderFns: [],
  components: {
    VCard: lib.VCard,
    VCardActions: lib.VCardActions,
    VCardText: lib.VCardText,
    VDialog: lib.VDialog,
    VIcon: lib.VIcon,
    VToolbar: lib.VToolbar,
    VToolbarTitle: lib.VToolbarTitle,
    VSpacer: lib.VSpacer,
    VBtn: lib.VBtn
  },
  props: {
    buttonTrueText: {
      type: String,
      default: 'Yes'
    },
    buttonFalseText: {
      type: String,
      default: 'No'
    },
    buttonTrueColor: {
      type: String,
      default: 'primary'
    },
    buttonFalseColor: {
      type: String,
      default: 'grey'
    },
    buttonFalseFlat: {
      type: Boolean,
      default: true
    },
    buttonTrueFlat: {
      type: Boolean,
      default: true
    },
    color: {
      type: String,
      default: 'warning'
    },
    icon: {
      type: String,
      default: function default$1() {
        return this.$vuetify.icons.values.warning;
      }
    },
    message: {
      type: String,
      required: true
    },
    persistent: Boolean,
    title: {
      type: String
    },
    width: {
      type: Number,
      default: 450
    }
  },
  data: function data() {
    return {
      value: false
    };
  },
  mounted: function mounted() {
    document.addEventListener('keyup', this.onEnterPressed);
  },
  destroyed: function destroyed() {
    document.removeEventListener('keyup', this.onEnterPressed);
  },
  methods: {
    onEnterPressed: function onEnterPressed(e) {
      if (e.keyCode === 13) {
        e.stopPropagation();
        this.choose(true);
      }
    },
    choose: function choose(value) {
      this.$emit('result', value);
      this.value = value;
      this.$destroy();
    },
    change: function change(res) {
      this.$destroy();
    }
  }
};

function Install(Vue, options) {
  if ( options === void 0 ) options = {};

  var property = options.property || '$confirm';
  delete options.property;
  var vuetify = options.vuetify;
  delete options.vuetify;
  if (!vuetify) {
    console.warn('Module vuetify-confirm needs vuetify instance. Use Vue.use(VuetifyConfirm, { vuetify })');
  }
  var Ctor = Vue.extend(Object.assign({ vuetify: vuetify }, Confirm));
  function createDialogCmp(options) {
    var container = document.querySelector('[data-app=true]') || document.body;
    return new Promise(function (resolve) {
      var cmp = new Ctor(Object.assign({}, {
        propsData: Object.assign({}, Vue.prototype[property].options, options),
        destroyed: function () {
          container.removeChild(cmp.$el);
          resolve(cmp.value);
        }
      }));
      container.appendChild(cmp.$mount().$el);
    });
  }

  function show(message, options) {
    if ( options === void 0 ) options = {};

    options.message = message;
    return createDialogCmp(options);
  }

  Vue.prototype[property] = show;
  Vue.prototype[property].options = options || {};
}

if (typeof window !== 'undefined' && window.Vue) {
  window.Vue.use(Install);
}

module.exports = Install;