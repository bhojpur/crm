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
import Router from 'vue-router'

// # Views
import HomeView from './components/views/tpl/HomeView'
import DefaultView from './components/views/proxyView/DefaultView.vue'
import AuthView from './components/views/proxyView/AuthView.vue'
import store from './store'


Vue.use(Router)

/**
 * Пропускает, только если нет никакой авторизации. Если есть токен - пробует перенаправить на домашнюю страницу
 * @param to
 * @param from
 * @param next
 */
const ifNotGuestGoHome = (to, from, next) => {

	if (store.getters['user/isAuth'] && store.getters['account/isAuth']) {
		next({name: "home"})
		return
	}
//	console.log("ifNotGuestGoHome ... not auth..")
//	console.log("user/isAuth: ", store.getters['user/isAuth'])
//	console.log("account/isAuth: ", store.getters['account/isAuth'])
	next()
}

/**
 * Если нет авторизации (нет токена), перенаправляет на страницу авторизации
 * @param to
 * @param from
 * @param next
 */
const ifNotAuthenticatedGoLogin = (to, from, next) => {

	// если нет токена - на страницу авторизации
	if (!store.getters['user/isAuth']) {
		next({name: 'login'})
		return
	}
	next()
}

/**
 * Перед загрузкой аккаунта проверяет пользователя.
 *
 * Перенаправляет пользователя на выбор аккаунта, если аккаунт не загружен, но авторизация есть. Если авторизации нет, то отправляет на страницу авторизации.
 * @param to
 * @param from
 * @param next
 */
const ifNotLoadedAccountGoChoiceAccounts = (to, from, next) => {

	// если нет токена авторизации - на страницу логина
	if (!store.getters['user/isAuth']) {
		next({name: 'login'})
		return
	}

	// подтверждаем аккаунт
	if (!store.getters['user/emailVerified']) {
		next({name: 'login.email-verification'})
		return
	}

	// подтверждаем аккаунт
	if (store.getters['user/passwordReset']) {
		next({name: 'login.password-recovery'})
		return
	}

	// если аккаунт не загружен - на страницу выбора аккаунта
	if (!store.getters['account/isAuth']) {
		next({name: 'login.accounts'})
		return
	}

	next()
}

//export default new Router({
let router = new Router({
	mode: 'history',
	base: process.env.BASE_URL,
	routes: [
		{
			path: '/',
			component: DefaultView, // with header, footer & Co
			/* Маршруты основного приложения */
			children: [
				{
					path: '',
					name: 'home',
					component: HomeView,
				},
				{
					path: 'account',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: '', // login/accounts
							name: 'account',
							component: () => import('./components/views/tpl/account/Index'),
						},
						{
							path: 'integrations', // login/accounts
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								/*{
									path: '',
									name: 'account.integrations',
									component: () => import('./components/views/tpl/integrations/Index'),
								},
								{
									path: 'api',
									name: 'account.integrations.api',
									component: () => import('./components/views/tpl/integrations/Api'),
								},*/

							]
						}
					]
				},
				{
					path: 'email-marketing',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'templates',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'email-marketing.templates',
									component: () => import('./components/views/tpl/email-marketing/templates/EmailTemplates'),
								},
								{
									path: ':public_id',
									name: 'email-marketing.templates.edit',
									component: () => import('./components/views/tpl/email-marketing/templates/EmailTemplateEdit'),
								},
							],
						},
						{
							path: 'notifications',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'email-marketing.notifications',
									component: () => import('./components/views/tpl/email-marketing/notification/Notifications'),
								},
								{
									path: ':public_id',
									name: 'email-marketing.notifications.edit',
									component: () => import('./components/views/tpl/email-marketing/notification/NotificationEdit'),
								},
							],
						},
						{
							path: 'queues',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'email-marketing.queues',
									component: () => import('./components/views/tpl/email-marketing/queue/Queues'),
								},
								{
									path: 'create',
									name: 'email-marketing.queues.create',
									component: () => import('./components/views/tpl/email-marketing/queue/QueueCreate'),
								},
								{
									path: ':public_id',
									name: 'email-marketing.queues.edit',
									component: () => import('./components/views/tpl/email-marketing/queue/QueueEdit'),
								},
							],
						},
						{
							path: 'campaigns',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'email-marketing.campaigns',
									component: () => import('./components/views/tpl/email-marketing/campaign/Campaigns'),
								},
								{
									path: 'create',
									name: 'email-marketing.campaign.create',
									component: () => import('./components/views/tpl/email-marketing/campaign/CampaignCreate'),
								},
								{
									path: ':public_id',
									name: 'email-marketing.campaign.edit',
									component: () => import('./components/views/tpl/email-marketing/campaign/CampaignEdit'),
								},
							],
						},
					]
				},
				{
					path: 'web',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'sites',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'web-sites.list',
									component: () => import('./components/views/tpl/web-sites/sites/WebSites'),
								},
								{
									path: ':public_id',
									name: 'web-sites.edit',
									component: () => import('./components/views/tpl/web-sites/sites/WebSitesEdit'),
								},
							],
						},
						{
							path: 'pages',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'web-pages.list',
									component: () => import('./components/views/tpl/web-sites/pages/WebPages'),
								},
								{
									path: ':public_id',
									name: 'web-pages.edit',
									component: () => import('./components/views/tpl/web-sites/pages/WebPageEdit'),
								},
							],
						},


					],
				},
				{
					path: 'content',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'storage',
							name: 'content.storage',
							component: () => import('./components/views/tpl/content/Storage'),
						},
						{
							path: 'articles',
							name: 'content.articles',
							component: () => import('./components/views/tpl/content/article/Article'),
						},
						{
							path: 'articles/:public_id',
							name: 'content.articles.edit',
							props: true,
							component: () => import('./components/views/tpl/content/article/ArticleEdit'),
						},
					]
				},
				{
					path: 'deliveries',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'orders',
							name: 'delivery.orders',
							component: () => import('./components/views/tpl/deliveries/DeliveryOrders'),
						},

						{
							path: 'orders/:public_id',
							name: 'delivery.order.edit',
							component: () => import('./components/views/tpl/deliveries/DeliveryOrderEdit'),
						},
					]
				},
				{
					path: 'integrations',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: '',
							name: 'integrations',
							component: () => import('./components/views/tpl/integrations/Index'),
						},
						{
							path: 'api',
							name: 'integrations.api',
							component: () => import('./components/views/tpl/integrations/Api'),
						},
						{
							path: 'web-hooks',
							name: 'integrations.webHooks',
							component: () => import('./components/views/tpl/integrations/WebHook'),
						},
						{
							path: 'deliveries',
							name: 'delivery.methods',
							component: () => import('./components/views/tpl/deliveries/DeliveryMethods'),
						},
					]
				},
				{
					path: 'observers',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'events',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'events.list',
									component: () => import('./components/views/tpl/observers/events/Events'),
								},
								{
									path: ':id',
									name: 'event.edit',
									component: () => import('./components/views/tpl/observers/events/EventEdit'),
								},
							],
						},
						{
							path: 'listeners',
							// name: 'events.listeners',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'listeners.list',
									component: () => import('./components/views/tpl/observers/listeners/Listeners'),
								},
								{
									path: ':id',
									name: 'listener.edit',
									component: () => import('./components/views/tpl/observers/listeners/ListenerEdit'),
								},
							],
						},
						{
							path: 'handlers',
							// name: 'events.handlers',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'handlers.list',
									component: () => import('./components/views/tpl/observers/handlers/Handlers'),
								},
								{
									path: ':id',
									name: 'handler.edit',
									component: () => import('./components/views/tpl/observers/handlers/HandlerEdit'),
								},
							],
						},

					]
				},

				{
					path: 'shops',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'products',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'products.list',
									component: () => import('./components/views/tpl/shops/products/Products'),
								},
								{
									path: ':public_id',
									name: 'product.edit',
									component: () => import('./components/views/tpl/shops/products/ProductEdit'),
								},
							],
						},
						{
							path: 'product-cards',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'product-cards.list',
									props: true,
									component: () => import('./components/views/tpl/shops/product-cards/ProductCards'),
								},
								{
									path: ':public_id',
									name: 'product-cards.edit',
									props: true,
									component: () => import('./components/views/tpl/shops/product-cards/ProductCardEdit'),
								},
							],
						},
						{
							path: 'product-categories',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'product-categories.list',
									props: true,
									component: () => import('./components/views/tpl/shops/product-category/ProductCategories'),
								},
								{
									path: ':public_id',
									name: 'product-category.edit',
									props: true,
									component: () => import('./components/views/tpl/shops/product-category/ProductCategoryEdit'),
								},
							],
						},
						{
							path: 'product-types',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'product-types.list',
									props: true,
									component: () => import('./components/views/tpl/shops/product-type/ProductTypes'),
								},
								{
									path: ':public_id',
									name: 'product-type.edit',
									props: true,
									component: () => import('./components/views/tpl/shops/product-type/ProductTypeEdit'),
								},
							],
						},
						{
							path: 'product-tags',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'product-tags.list',
									props: true,
									component: () => import('./components/views/tpl/shops/product-tag/ProductTags'),
								},
								{
									path: ':public_id',
									name: 'product-tag.edit',
									props: true,
									component: () => import('./components/views/tpl/shops/product-tag/ProductTagEdit'),
								},
							],
						},
						{
							path: 'product-tag-groups',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'product-tag-groups.list',
									props: true,
									component: () => import('./components/views/tpl/shops/product-tag-group/ProductTagGroups'),
								},
								{
									path: ':public_id',
									name: 'product-tag-group.edit',
									props: true,
									component: () => import('./components/views/tpl/shops/product-tag-group/ProductTagGroupEdit'),
								},
							],
						},
					]
				},
				{
					path: 'organizations',
					name: 'organizations',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'users',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'users.list',
									component: () => import('./components/views/tpl/organizations/users/Users'),
								},
								{
									path: ':public_id',
									name: 'user.edit',
									component: () => import('./components/views/tpl/organizations/users/UserEdit'),
								},
							]
						},
						{
							path: 'manufacturers',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'manufacturers.list',
									component: () => import('./components/views/tpl/organizations/manufacturer/Manufacturers'),
								},
								{
									path: ':public_id',
									name: 'manufacturer.edit',
									component: () => import('./components/views/tpl/organizations/manufacturer/ManufacturerEdit'),
								},
							],
						},
						{
							path: 'companies',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'companies.list',
									component: () => import('./components/views/tpl/organizations/companies/Companies'),
								},
								{
									path: ':public_id',
									name: 'company.edit',
									component: () => import('./components/views/tpl/organizations/companies/CompanyEdit'),
								},
							],
						},
					],
				},
				{
					path: 'inventory-control',
					name: 'inventory-control',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: 'shipments',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'shipments.list',
									component: () => import('./components/views/tpl/warehouse/shipment/Shipments'),
								},
								{
									path: ':public_id',
									name: 'shipment.edit',
									component: () => import('./components/views/tpl/warehouse/shipment/ShipmentEdit'),
								},
							],
						},
						{
							path: 'warehouse',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'warehouses.list',
									component: () => import('./components/views/tpl/warehouse/warehouse/Warehouses'),
								},
								{
									path: ':public_id',
									name: 'warehouse.edit',
									component: () => import('./components/views/tpl/warehouse/warehouse/WarehouseEdit'),
								},
							],
						},
						{
							path: 'inventories',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'inventories.list',
									component: () => import('./components/views/tpl/warehouse/inventories/Inventories'),
								},
								{
									path: ':public_id',
									name: 'inventory.edit',
									component: () => import('./components/views/tpl/warehouse/inventories/InventoryEdit'),
								},
							],
						},
					],
				},
				{
					path: 'orders',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: '',
							name: 'orders.list',
							component: () => import('./components/views/tpl/orders/Orders'),
						},
						{
							path: 'create',
							name: 'order.create',
							component: () => import('./components/views/tpl/orders/OrderEditOrCreate'),
						},
						{
							path: ':public_id',
							name: 'order.edit',
							component: () => import('./components/views/tpl/orders/OrderEditOrCreate'),
						},

					],
				},
				{
					path: 'questions',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: '',
							name: 'questions.list',
							component: () => import('./components/views/tpl/orders/Questions'),
						},
					],
				},
				{
					path: 'finance',
					component: () => import('./components/views/tpl/ViewProxy'),
					children:[
						{
							path: 'payments',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'finance.payments.list',
									component: () => import('./components/views/tpl/finance/payment/Payments'),
								},
								{
									path: ':id',
									name: 'finance.payment.edit',
									component: () => import('./components/views/tpl/finance/payment/PaymentEdit'),
								},

							],
						},
						{
							path: 'methods',
							component: () => import('./components/views/tpl/ViewProxy'),
							children: [
								{
									path: '',
									name: 'finance.payment-method.list',
									component: () => import('./components/views/tpl/finance/method/PaymentMethods'),
								},
								{
									path: ':id',
									name: 'finance.payment-method.edit',
									component: () => import('./components/views/tpl/finance/method/PaymentMethodEdit'),
								},

							],
						}
					],
				},
				{
					path: '/404',
					name: '404',
					meta: {
						title: '404 - страница не найдена',
						label: '404',
						breadcrumb: '404',
					},
					component: import(/* webpackChunkName: "about" */ './components/views/tpl/404.vue'),
				},
			],
			beforeEnter: ifNotLoadedAccountGoChoiceAccounts, // только при загруженном аккаунте
		},
		{
			path: '/login',
			// route level code-splitting
			// this generates a separate chunk (about.[hash].js) for this route
			// which is lazy-loaded when the route is visited.
			component: AuthView,
			children: [
				{
					path: '',
					name: 'login',
					component: () => import('./components/views/auth/Login.vue'),
					props: (route) => ({ msgNotification: route.params['message'] !== undefined ? route.params.message : ''}),
					beforeEnter: ifNotGuestGoHome,

				},
				{
					path: 'accounts', // login/accounts/create
					name: 'login.accounts',
					component: () => import('./components/views/auth/ChoiceAccount.vue'),
					beforeEnter: (to, from, next) => {

						if (!store.getters['user/isAuth']) {
							console.log("!store.getters['user/isAuth']")
							next({name: 'login'});
							return
						}

						// Todo доделать проверку верификации email'а
						/*if (!store.getters['user/emailVerified']) {
							console.log("!store.getters['user/emailVerified']")
							next({name: 'login.accounts.step-2'});
							return
						}*/

						if (store.getters['account/isAuth']) {
							// тут по идее должно быть предупреждение о выходе
							store.dispatch('ACCOUNT_LOGOUT') // root
								.then()
								.catch()
								.finally(()=> next());
//									next()
							return
						}

						next()
					}, //
				},
				{
					path: 'create-account',
					name: 'login.create-account',
					component: () => import('./components/views/auth/CreateAccount.vue'),
					props: (route) => ({ queryToken: route.query.t }),
					beforeEnter: (to, from, next) => {

						if (!store.getters['user/isAuth']) {
							next({name: 'login'});
							return
						}

						if (!store.getters['user/emailVerified']) {
							next({name: 'login.sign-up.step-2'});
							return
						}
						next()
					}, // не должен загружаться, если аккаунт уже был выбран

				},
				{
					path: 'sign-up',
					component: () => import('./components/views/tpl/ViewProxy'),
					children: [
						{
							path: '', // login/accounts
							name: 'login.sign-up',
							component: () => import('./components/views/auth/CreateAccount/Step-1.vue'),
							beforeEnter: (to, from, next) => {

								if (store.getters['account/isAuth']) {
									// тут по идее должно быть предупреждение о выходе
									store.dispatch('ACCOUNT_LOGOUT') // root
										.then()
										.catch()
										.finally(()=> next());
									return
								}
								next()
							},
						},
						{
							path: 'email-verification', // email-verification
							name: 'login.sign-up.step-2',
							component: () => import('./components/views/auth/CreateAccount/Step-2.vue'),
							props: (route) => ({ queryToken: route.query.t })
						},
						{
							path: 'create',
							name: 'login.sign-up.step-3',
							component: () => import('./components/views/auth/CreateAccount/Step-3.vue'),
							beforeEnter: (to, from, next) => {

								if (!store.getters['user/isAuth']) {
									next({name: 'login'});
									return
								}

								if (!store.getters['user/emailVerified']) {
									next({name: 'login.sign-up.step-2'});
									return
								}
								next()
							}, // не должен загружаться, если аккаунт уже был выбран
						},
					],
				},
				{
					path: 'email-verification', // email-verification
					name: 'login.email-verification',
					component: () => import('./components/views/auth/EmailVerification.vue'),
					props: (route) => ({ queryToken: route.query.t })

				},
				{
					path: 'forgot-password',
					name: 'login.forgot-password',
					component: () => import('./components/views/auth/ForgotPassword.vue'),
					beforeEnter: ifNotGuestGoHome,
				},
				{
					path: 'password-recovery',
					name: 'login.password-recovery',
					component: () => import('./components/views/auth/PasswordRecovery.vue'),
					props: (route) => ({ queryToken: route.query.t, msgNotification: route.params['msgNotification']}),
//					beforeEnter: ifNotGuestGoHome,
				},
				{
					path: 'forgot-username',
					name: 'login.forgot-username',
					component: () => import('./components/views/auth/ForgotUsername.vue'),
					beforeEnter: ifNotGuestGoHome,
				},
			],
		},
		{
			path: '/legal',
			component: () => import(/* webpackChunkName: "about" */ './components/views/legal/LegalViewProxy.vue'),
			/* Маршруты для правил использования сервисом */
			children: [
				{
					path: '',
					name: 'legal',
					component: () => import(/* webpackChunkName: "about" */ './components/views/legal/Legal.vue'),
				},
				{
					path: 'use-cookie',
					name: 'legal.use-cookie',
					component: () => import(/* webpackChunkName: "about" */ './components/views/legal/UseCookie.vue')
				},
				{
					path: 'privacy-policy',
					name: 'legal.privacy-policy',
					component: () => import(/* webpackChunkName: "about" */ './components/views/legal/PrivacyPolicy.vue'),
				},
				{
					path: 'terms-of-use',
					name: 'legal.terms-of-use',
					component: () => import(/* webpackChunkName: "about" */ './components/views/legal/TermsOfUse.vue'),
				},
			],
		},
		{
			path: '*',
			redirect: '/404',
			meta: {
				title: '404 - страница не найдена',
				label: '404',
			},
			// beforeEnter: ifUserAccess,
			// beforeEnter: ifAuthenticated,
		},
	],
})

/*router.beforeEach((to, from, next) => {
	next()
})*/

export default router
