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

import axios from 'axios'
// import router from '../router'

/**
 * Основная функция обращения к API (Auth, Gui, Data servers)
 * @param url
 * @param method
 * @param args
 * @returns {Promise<any>}
 */
let api = ({url, method, ...args}) => new Promise((resolve, reject) => {
    axios({url, method, ...args})
        .then( (response) => {
            if (response['data'] === undefined || response['data']['status'] === undefined) {
                reject(new Error("Неудается обработать ответ сервера"));
            }
            if ( response.data && !response.data.status ) {
                reject(response.data);
            } else {
                resolve(response.data);
            }
        })
        .catch( (error) => {
            if (error.response) {

                // The request was made and the server responded with a status code
                // that falls out of the range of 2xx
                // ответы перехыватываются в axios.interceptors (401)
                // console.log("error.response!")
                reject(new Error("Ошибка в обработке запроса"));

            } else if (error.request) {
                // console.log("error.request!")
                // The request was made but no response was received (Неудалось получить ответ от сервера / Сервер не отвечает)
                // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
                reject(new Error("Не удалось получить ответ от сервера")); // сервер не смог ответить по каким-то причинам

            } else {

                // Something happened in setting up the request that triggered an Error
                // что-то сгенерировало ошибку (сервер вообще не отвечает)
                // console.log("error Something happened!")
                reject(new Error("Не удаётся установить соединение с сервером"));
            }

            reject(new Error("Неизвестная ошибка соединения"));
        });

});
function getUrlByKeys(url, keys, payload, byPublicId = false, showAll = false){

    url = window.VUE_APP_SERVER_GUI_URL + url
    if (byPublicId) {
        Object.assign(payload,{public_id:'true'})
        keys.push("public_id")
    }
    if (showAll) {
        Object.assign(payload,{all:'true'})
        keys.push("all")
    }

    let _first = true
    keys.forEach(key=>{
        let v = payload[key]
        if (v !== undefined && v !== null && !(typeof v === 'string' && v.length < 1)) {
            url += `${_first ? '?' : '&'}` + key + `=${v}`
            if (_first) _first = false
        }
    })
    return url
}

export default api
export {api, getUrlByKeys}