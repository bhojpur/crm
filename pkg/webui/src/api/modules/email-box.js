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
 * Mocking client-server processing
 */
 import {api} from '@/api/utils'
 import {getUrlByKeys} from '@/api/utils'
 
 export default {
 
     create (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-boxes`,['preloads'], payload), method: 'post', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     get (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-boxes/${payload['id']}`,['preloads'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     getList (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-boxes`,['preloads'], payload, false,true), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     getListPagination (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-boxes`,
                     ['offset','limit','sortBy','sortDesc','search','preloads','webSiteId'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     update (payload) {
         delete payload['_recipient_users']
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-boxes/${payload['id']}`,['preloads'], payload), method: 'patch', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     delete (payload) {
         return new Promise((resolve, reject) => {
             api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/email-boxes/${payload['id']}`, method: 'delete', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
 
     ////////////////////////////
 
     parseDataFromPayload(payload) {
         let postData = {}
         Object.assign(postData,payload)
         delete postData['accountHashId']
 
         return postData
     },
 
 }