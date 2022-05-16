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

// var MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = {
    css: {
      loaderOptions: {
        // new
        sass: {
          prependData: `@import "@/assets/sass/main.scss"`, // change to *.sass; remove semi-colon
        },
        scss: { // separate one for scss
          prependData: `@import "@/assets/sass/main.scss";`, // this is .scss
        },
      },
    },
  
    chainWebpack: config => {
      config.optimization.delete('splitChunks');
      config.optimization
          .splitChunks().clear();
      // config.plugins.delete('html');
      config.plugins.delete('preload');
      config.plugins.delete('prefetch');
  
    },
    pwa: {
      name: 'Bhojpur CRM',
      themeColor: '#6392f1',
      msTileColor: '#000000',    appleMobileWebAppCapable: 'yes',
      appleMobileWebAppStatusarStyle: 'black',
      // настройки манифеста
      manifestOptions: {
        display: 'landscape',
        background_color: '#4277b8'
        // ...другие настройки манифеста...
      }
    }
  }