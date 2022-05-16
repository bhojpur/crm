import Vue from 'vue'

const ShopMixin = {

	install(Vue, options) {

		Vue.myGlobalMethod = function () {
//			console.log(this.$router.push({name: 'login'}))
		}

		Vue.mixin({
			methods: {
				getFieldFromShops(obItem, fieldName) {
					let shops = this.webSites
					let index = shops.findIndex(item => item.id === obItem.web_site_id, obItem)
					let shop = shops[index]

					if (shop) return shop[fieldName]

					return "-"
				},
			}
		})
	}
}

Vue.use(ShopMixin)

export {ShopMixin}

