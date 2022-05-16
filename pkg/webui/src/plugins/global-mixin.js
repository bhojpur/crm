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

/**
 * Всякие полезные миксины и плюшки для клиентаы
 */

 import Vue from 'vue'

 import deliveries from '@/api/modules/deliveries';
 import emailTemplate from '@/api/modules/email-template';
 import event from '@/api/modules/event';
 import eventHandlerItem from '@/api/modules/event-handler-item';
 import eventListener from '@/api/modules/event-listener';
 import userRole from '@/api/modules/user-role';
 import user from '@/api/modules/user';
 import webHook from '@/api/modules/web-hook';
 import emailNotification from '@/api/modules/email-notification';
 import webSite from '@/api/modules/web-site';
 import emailBox from '@/api/modules/email-box';
 import productCard from '@/api/modules/product-card';
 import productGroup from '@/api/modules/product-group';
 import storage from '@/api/modules/storage';
 import order from '@/api/modules/order';
 import orderChannel from '@/api/modules/order-channel';
 import paymentSubject from '@/api/modules/payment-subject';
 import payment from '@/api/modules/payment';
 import vatCode from '@/api/modules/vat-code';
 import deliveryOrder from '@/api/modules/delivery-order';
 import paymentMethod from '@/api/modules/payment-method';
 import paymentMode from '@/api/modules/payment-mode';
 import deliveryStatuses from '@/api/modules/delivery-status';
 import orderStatus from '@/api/modules/order-status';
 import emailQueue from '@/api/modules/email-queue';
 import emailQueueEmailTemplate from '@/api/modules/email-queue-email-template';
 import product from '@/api/modules/product';
 import article from '@/api/modules/article';
 import apiKey from '@/api/modules/api-key';
 import emailCampaign from '@/api/modules/email-campaign';
 import usersSegment from '@/api/modules/user-segment';
 import webPage from '@/api/modules/web-pages';
 import warehouse from '@/api/modules/warehouse';
 import shipment from '@/api/modules/shipment';
 import inventory from '@/api/modules/inventory';
 import helper from '@/api/modules/helper';
 import productCategory from '@/api/modules/product-category';
 import measurementUnit from '@/api/modules/measurement-unit';
 import manufacturer from '@/api/modules/manufacturer';
 import company from '@/api/modules/company';
 import productType from '@/api/modules/product-type';
 import productTag from '@/api/modules/product-tag';
 import productTagGroup from '@/api/modules/product-tag-group';
 import mtaHistory from '@/api/modules/mta-history';
 import cartItem from '@/api/modules/cart-item';
 import warehouseItem from '@/api/modules/warehouse-item';
 import question from '@/api/modules/question';
 
 
 import moment from 'moment';
 
 const GlobalMixin = {
 
     /**
      * В options можно добавить процессор обработки колбеков
      * @param Vue
      * @param options
      */
     install(Vue, options) {
         Vue.myGlobalMethod = function () {}
 
 
 
         Vue.prototype.moment = moment;
         Vue.prototype.$api = {
             emailTemplate,
             deliveries,
             event,
             eventHandlerItem,
             eventListener,
             user, userRole,
             webHook,
             emailNotification,
             webSite,
             emailBox,
             productCard,
             productGroup,
             order,
             orderChannel,
             storage,
             payment,
             paymentSubject,
             vatCode,
             deliveryOrder,
             paymentMethod,
             paymentMode,
             deliveryStatuses,
             orderStatus,
             emailQueue,
             emailQueueEmailTemplate,
             product,
             article,
             apiKey,
             emailCampaign,
             usersSegment,
             webPage,
             warehouse,
             shipment,
             inventory,
             helper,
             productCategory,
             measurementUnit,
             manufacturer,
             company,
             productType,
             productTag,
             productTagGroup,
             mtaHistory,
             cartItem,
             warehouseItem,
             question,
         };
 
         Vue.mixin({
             computed: {
                 isMainAccount() {
                     return this.$store.state.account.id === 1
                 },
             },
             methods: {
                 getMomentJSData(data){
                     return moment(data).format("HH:MM DD.MM.YYYY")
                 },
                 logout: async function () {
                     localStorage.removeItem("vuex")
                     window.$broadcastTab.message({command: "logout"})
                     // localStorage.removeItem('message');
                     await this.$store.dispatch("AUTH_LOGOUT", {root: true})
                     this.$router.push({ name: 'login',  params: { msg: 'Авторизуйтесь в системе' } }).catch()
                 },
                 sleep(ms) {
                     return new Promise(resolve => setTimeout(resolve, ms));
                 },
                 goToAuthView: function () {
                     this.$store.dispatch("AUTH_LOGOUT", {root: true})
                         .catch()
                         .finally(()=> this.$router.push({ name: 'login',  params: { msg: 'Авторизуйтесь в системе' } }).catch())
 
                 },
                 /* Выход в список доступных аккаунтов */
 
                 goToLoginInAccount: function () {
                     this.$router.push({ name: 'login.accounts' })
                         .then()
                         .catch(() => this.logout)
                 },
                 loginInAccount: function (accountHashId, rPayload) {
                     const payload = {
                         accountHashId: accountHashId,
                     }
                     return this.$store.dispatch('account/LOGIN_IN_ACCOUNT', payload)
                         .then(() => this.$router.push(rPayload))
                         .catch( () => this.goToLoginInAccount)
                 },
                 validateEmail(email) {
                     let regexEmail = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
                     return !!email.match(regexEmail);
                 },
                 parseDate(inData) {
                     return (new Date(inData)).toLocaleDateString()
                 },
                 formatBytes (bytes, decimals = 2) {
                     if (bytes === 0) return '0 Bytes';
 
                     const k = 1024;
                     const dm = decimals < 0 ? 0 : decimals;
                     const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
 
                     const i = Math.floor(Math.log(bytes) / Math.log(k));
 
                     return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
                 },
                 copyBuffer(buf, showNotification = true) {
                     navigator.clipboard.writeText(buf).then(()=> {
                         if (showNotification) {
                             this.$notify({
                                 group: 'main',
                                 title: 'Данные скопированы в память',
                                 type: 'success'
                             });
                         }
 
                     })
                 },
 
                 uploadFile(owner_id, owner_type) {
                     const formData = new FormData()
                     formData.append("file", this.uploadedFile)
 
                     let payload = {
                         accountHashId: this.$store.state.account.hash_id,
                         file: formData,
                         owner_id: owner_id,
                         owner_type: owner_type,
                     }
 
                     return this.$api.storage.create(payload)
                 },
 
                 // helpers
                 getEmailBoxName(item) {
 
                     if (item === undefined) return 'some errors'
                     if (this.web_sites === undefined) return 'some errors'
                     if (this.email_boxes === undefined) return 'some errors'
 
                     let email_box = this.email_boxes.find(el=>el.id === item.email_box_id,item)
                     if (email_box) {
                         let web_site = this.web_sites.find(el=>el.id === email_box.web_site_id, email_box)
                         if (web_site) return email_box.name + " / " + email_box.box + "@" + web_site.hostname
                     }
 
                     return 'n/a'
                 },
                 getEmailTemplateName(item) {
 
                     // console.log(item.email_template_id, item.email_template)
                     if (item.email_template_id && item.email_template ) {
                         return '№' + item.email_template.public_id + " " + item.email_template.name
                     }
                     return 'n/a'
                 },
                 getManufacturerName(item) {
                     if (item.manufacturer_id && item.manufacturer ) {
                         return '№' + item.manufacturer.public_id + " " + item.manufacturer.name
                     }
                     return 'n/a'
                 },
                 getRoleByRoleId(role_id) {
                     if (role_id === undefined || !this.roles) return null
 
                     let role = this.roles.find(el=>el['id'] === role_id)
                     if (!role) return null
                     return role
                 },
                 getDelay(time) {
                     if (!time || parseInt(time) < 10) return 'без задержки'
                     return '~ ' + this.$moment.duration(time/1000000).humanize(false)
                 },
                 unflatten(arr) {
                     let tree = [],
                         mappedArr = {},
                         arrElem,
                         mappedElem;
 
                     // First map the nodes of the array to an object -> create a hash table.
                     for(let i = 0, len = arr.length; i < len; i++) {
                         arrElem = arr[i];
                         mappedArr[arrElem.id] = arrElem;
                         mappedArr[arrElem.id]['children'] = [];
                     }
 
 
                     for (let id in mappedArr) {
                         if (mappedArr.hasOwnProperty(id)) {
                             mappedElem = mappedArr[id];
                             // If the element is not at the root level, add it to its parent array of children.
                             if (mappedElem.parent_id) {
                                 mappedArr[mappedElem['parent_id']]['children'].push(mappedElem);
                             }
                             // If the element is at the root level, add it to first level elements array.
                             else {
                                 tree.push(mappedElem);
                             }
                         }
                     }
                     return tree;
                 },
                 sortByKey(items, key = 'order') {
                     // if (items.length < 1) return []
                     return items.sort(function (a, b) {
                         if (a[key] > b[key]) {
                             return 1;
                         }
                         if (a[key] < b[key]) {
                             return -1;
                         }
                         // a должно быть равным b
                         return 0;
                     })
                 },
 
                 loadWebSites: async function(showNotification) {
 
                     this.loading = true
 
                     const payload = {
                         accountHashId: this.$store.state.account.hash_id,
                     }
 
                     return this.$api.webSite.getList(payload)
                         .then(resp => {
                             if (resp['web_sites'] !== undefined ) {
                                 let _webSites = resp['web_sites']
                                 _webSites.push({name:'Все сайты',id: null})
                                 this.webSites = _webSites
                                 // this.webSites = resp['web_sites']
 
                                 if (showNotification) {
                                     this.$notify({
                                         group: 'main',
                                         title: 'Данные обновлены',
                                         type: 'success'
                                     });
                                 }
                             } else {
                                 if (showNotification) {
                                     this.$notify({
                                         group: 'main',
                                         title: 'Ошибка обновления',
                                         text: 'Ошибка в ответе сервера: web_sites - not found',
                                         type: 'warring'
                                     });
                                 }
                             }
                         })
                         .catch((err) => {
                             this.$notify({
                                 group: 'main',
                                 title: 'Ошибка обновления',
                                 text: err['message'],
                                 type: 'error'
                             });
                         })
                         .finally(()=>this.loading = false)
                 },
 
                 // возвращает +, если число положительное
                 MathGetSign(num) {
                     if (parseFloat(num) === 0) return ''
                     return num > 0 ? `+` : '-'
                 },
                 DeclOfNum(number, words) {
                     return words[(number % 100 > 4 && number % 100 < 20) ? 2 : [2, 0, 1, 1, 1, 2][(number % 10 < 5) ? number % 10 : 5]];
                 }
 
             },
             filters: {
                 // formatDate: d => moment(String(d)).format('MM/DD/YYYY hh:mm')
             },
         })
 
         // Vue.prototype.$testGoGo = function (methodOptions) {}
     }
 }
 
 Vue.use(GlobalMixin)
 
 export {GlobalMixin}
 
 