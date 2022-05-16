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
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates`,['preloads'], payload), method: 'post', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     get (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates/${payload['id']}`,['preloads'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     getByPublicId (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates/${payload['public_id']}`,['preloads'], payload, true), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     getListPagination (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates`,['offset','limit','sortBy','sortDesc','search','preloads'], payload), method: 'get'})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     update (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates/${payload['id']}`,['preloads'], payload), method: 'patch', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
     delete (payload) {
         return new Promise((resolve, reject) => {
             api({url: getUrlByKeys(`/accounts/${payload['accountHashId']}/email-templates/${payload['id']}`,[], payload), method: 'delete', data: this.parseDataFromPayload(payload)})
                 .then(resp => resolve(resp))
                 .catch((err) => reject(err))
 
         })
     },
 
     // ####### Helpers ########
 
 
     parseDataFromPayload(payload) {
         let postData = {}
         Object.assign(postData,payload)
         delete postData['accountHashId']
 
         return postData
     },
     getBaseTemplateItem() {
         return {
             name: "Новый шаблон",
             description: "...",
             public: false,
             code: this.getBaseDataTemplate()
         }
     },
     getBaseDataTemplate() {
         return `<!DOCTYPE html>
 <html lang="en">
 <head>
     <meta charset="UTF-8">
     <title>Bhojpur CRM - {{.Subject}}</title>
 </head>
 <body style="background-color: #F4F4F4;" leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">
 <span class="mcnPreviewText" style="display:none; font-size:0px; line-height:0px; max-height:0px; max-width:0px; opacity:0; overflow:hidden; visibility:hidden; mso-hide:all;">{{.PreviewText}}</span>
 <table align="center" border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" style="background-color: #F4F4F4;">
     <tbody>
     <tr>
         <td align="center" valign="top">
             <table align="center" bgcolor="#F4F4F4" border="0" cellpadding="0" cellspacing="0" width="540" style="color:#222222; font:15px/1.4 Roboto,arial,sans-serif;">
                 <tr>
                     <td align="center" valign="top">
                         <p>Paste code there...</p>
                         <p>{{.UnsubscribeURL}}</p>
                         <img style="width: 1px;height: 1px;opacity: 0;" src="{{.PixelURL}}" />
                     </td>                
                 </tr>
             </table>
         </td>
     </tr>
     </tbody>
 </table>
 </body>
 </html>`
     }
 
 }
 