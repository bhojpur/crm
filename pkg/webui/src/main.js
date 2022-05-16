// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import Vue from 'vue'
import App from './components/views/App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import './registerServiceWorker'
import './plugins/axios';
import './plugins/broadcast-tabs';
import './plugins/global-mixin';
import './plugins/vuevalidate';
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@fortawesome/fontawesome-free/css/all.css'
import VueMoment from 'vue-moment'
import './store/plugins/shop-mixin';
import _ from 'lodash'
import "prismjs";
import VuePrismEditor from "vue-prism-editor";
import 'prismjs/themes/prism-coy.css'
import "vue-prism-editor/dist/prismeditor.min.css"; // import the styles
import Notifications from 'vue-notification'
import VuetifyConfirm from './plugins/vuetify-confirm/src'
import DatetimePicker from 'vuetify-datetime-picker'

import { TreeViewComponent, TreeViewPlugin } from '@syncfusion/ej2-vue-navigations';
import VCurrencyField from 'v-currency-field'
import { VTextField } from 'vuetify/lib'  //Globally import VTextField
import VueDadata from 'vue-dadata'
import DatePicker from 'vue2-datepicker';
import 'vue2-datepicker/index.css';
import 'vue2-datepicker/locale/ru';

Vue.use(DatePicker)


Vue.use(TreeViewPlugin);
Vue.component(TreeViewPlugin.name, TreeViewComponent);
Vue.use(TreeViewPlugin);

Vue.component('v-text-field', VTextField)
Vue.use(VCurrencyField, {
  locale: 'ru-RU',
  decimalLength: 2,
  autoDecimalMode: true,
  min: null,
  max: null,
  defaultValue: 0,
  valueAsInteger: false,
  allowNegative: true
})

Vue.use(VueDadata)
// Vue.use(SlVueTree);
// Vue.component(SlVueTree.name, SlVueTree);

const moment = require('moment')
require('moment/locale/ru')

Vue.use(require('vue-moment'), { moment })
Vue.component("prism-editor", VuePrismEditor);
Vue.use(DatetimePicker)

// Vue.use(VuetifyConfirm, {vuetify})
Vue.use(VuetifyConfirm, {
  vuetify,
  buttonTrueText: 'Подтвердить',
  buttonFalseText: 'Отмена',
  buttonTrueColor: 'primary',
  buttonFalseColor: 'grey',
  buttonTrueFlat: false,
  buttonFalseFlat: true,
  color: 'warning',
  icon: 'far fa-user',
  title: 'Внимание!',
  width: 350,
  property: '$confirm'
})

Vue.use(Notifications)

// window.i18n = i18n
Vue.config.productionTip = true

// url-подключения к gui-server
window.VUE_APP_SERVER_GUI_URL = (process.env.NODE_ENV === 'production') ? "https://crm.bhojpur.net/ui-api" : "http://app.crm.local/ui-api";
window.VUE_CDN_SERVER_URL = (process.env.NODE_ENV === 'production') ? "https://cdn.crm.bhojpur.net" : "http://cdn.crm.local";

new Vue({
  router,
  store,
  vuetify,
  // i18n,
  render: h => h(App)
}).$mount('#app')

