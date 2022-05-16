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
             api({
                 url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage`,['preloads','owner_id','owner_type'], payload), method: 'post', data: payload['file'] })
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
         })
     },
     get (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage/${payload['id']}`,['preloads'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
         })
     },
     getListPagination (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage`,['owner_id','owner_type','offset','limit','sortBy','sortDesc','search','preloads'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
         })
     },
     update (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage/${payload['id']}`,['preloads'], payload), method: 'patch', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
         })
     },
     massUpdates (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage`,['preloads','owner_id','owner_type'], payload), method: 'patch', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     delete (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/storage/${payload['id']}`,[], payload), method: 'delete', data: this.parseDataFromPayload(payload)})
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
 
 
     getCDN_URL: () => window.VUE_CDN_SERVER_URL,
     getCDN_PUBLIC_URL: () => window.VUE_CDN_SERVER_URL + '/public',
     getCDN_PRODUCTS_IMAGES_URL: () => window.VUE_CDN_SERVER_URL + '/products/images',
     getCDN_EMAILS_IMAGES_URL: () => window.VUE_CDN_SERVER_URL + '/emails/images',
     getCDN_EMAILS_RAW_URL: () => window.VUE_CDN_SERVER_URL  + '/emails/preview/raw',
     getCDN_EMAILS_COMPILE_URL: () => window.VUE_CDN_SERVER_URL + '/emails/preview/compile',
     getCDN_ARTICLES_RAW_URL: () => window.VUE_CDN_SERVER_URL + '/articles/preview/raw',
     getCDN_ARTICLES_COMPILE_URL: () => window.VUE_CDN_SERVER_URL + '/articles/preview/compile'
 }
 