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

import Vue from 'vue';
import axios from 'axios'
import VueAxios from 'vue-axios'

window.axios = require('axios');
// window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
// window.axios.defaults.headers.common['Accept-Language'] = process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'en'
// window.axios.defaults.headers.common['Accept-Language'] = 'en'
// window.axios.defaults.withCredentials = true

axios.defaults.credentials = false
axios.defaults.withCredentials = false; // надо ли жестко ограничивать CORS
axios.defaults.headers.common['Accept-Language'] = 'en'
axios.defaults.timeout = 20000 // 20s
axios.defaults.headers.common['Accept'] = 'application/json';
// axios.defaults.headers.post['Content-Type'] = 'application/json';
// axios.defaults.headers.post['Content-Type'] = 'multipart/boundary';

// обработка неавторизованного доступа, продление токена и т.д.
axios.interceptors.response.use(undefined, function (err) {
    let res = err.response;
    if (res.status === 401 && res.config && !res.config.__isRetryRequest) {
        window.$broadcastTab.message({command: "logout"})
        window.location.href = '/login'
    } else {
        throw err; // !!!! запускаем дальше необрабатываемые ошибки
    }
    throw err

    /*return new Promise(function (resolve, reject) {

        console.log('interceptors', err.config)
        // if (err.status === 401 && err.config && !err.config.__isRetryRequest) {
        if (err.response.status === 401 && err.config && !err.config.__isRetryRequest) {
            // if you ever get an unauthorized, logout the user
            // this.$store.dispatch(AUTH_LOGOUT)
            console.log('this.$store.dispatch(AUTH_LOGOUT)')
            // you can also redirect to /login if needed !
        }
        throw err;
    });*/
});

Vue.use(VueAxios, axios.create({
    // withCredentials: false,
    // credentials: false,
    /*headers: {
        'Content-Type': 'application/json',
        'Accept-Language' : process.env.VUE_APP_I18N_LOCALE || 'en'
    }*/
}))
