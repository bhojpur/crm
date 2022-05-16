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
 * Тут описание, честно-честно.
 * Кратко: бродкастим важные сервисные сообщения между открытыми вкладками. Например - logout.
 * by Shashi Bhushan Rai
 */

 import Vue from 'vue'
 import router from '../router'
 
 const _f = {
     message: function (message) {
         let msg = Object.assign(message, {'uid': (new Date).getTime()+Math.random()})
         // console.log("Бродкастим!", msg)
         // localStorage.removeItem('message');
         localStorage.setItem('message',JSON.stringify(msg));
         localStorage.removeItem('message');
     },
 }
 
 const BroadcastTabs = {
 
     /**
      * В options можно добавить процессор обработки колбеков
      * @param Vue
      * @param options
      */
     install(Vue, options) {
 
         Vue.myGlobalMethod = function () {
 //			console.log(this.$router.push({name: 'login'}))
         }
 
         Vue.mixin({
             computed: {
                 // `users` так же будет доступен везде
                 /*users () {
                     return this.$store.getters.users
                 },*/
             },
 
         })
 
         Vue.prototype.$broadcast = _f
 
         async function receive  (e) {
 
             if (e.key!=='message') return; // ignore other keys
             let message = JSON.parse(e.newValue);
 
             if (!message) return; // ignore empty msg or msg reset
 
             if (message.command === 'logout') {
                 console.log("need to logout")
                 window.location.href = '/login'
                 // router.go(0)
 // 				await router.push('/login')
             }
         }
         addEventListener("storage", receive)
     }
 }
 
 Vue.use(BroadcastTabs)
 
 window.$broadcastTab = Vue.prototype.$broadcast
 
 export default _f
 
 export {BroadcastTabs}