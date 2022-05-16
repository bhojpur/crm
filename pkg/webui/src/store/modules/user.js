import axios from 'axios'

import api from '../../api/utils'
import Vue from 'vue'

const state = {
    profile: {
    	id: 0
    },
//	accounts: [],
//	roles: [],
	aUsers: [], // совмещает в себе роли и аккаунты
}

const getters = {
    getProfile: state => state.profile,
	isAuth: state => state.profile.id !== 0, /* ничего не говорит нам о наличии токена */
	emailVerified: state => !!state.profile['email_verified_at'],
	passwordReset: state => !!state.profile['password_reset'],
    isProfileLoaded: state => state.profile.id !== 0,
	listOfAccounts: state => state.accounts,
	listAccountUsers: state => state.aUsers,  /* возвращает список аккаунтов с ролями */
	defaultAccountId: state => state.profile.id > 0 ? state.profile['default_account_id'] : 0,
}

const actions = {
	SET_DEFAULT_STATE: (context) => context.commit('SET_DEFAULT_STATE'),
	AUTH_LOGOUT: {
		root: true,
		handler: (context) => context.dispatch("SET_DEFAULT_STATE"),
	},
	/**
	 * Авторизация пользователя по логину и паролю
	 * @param context
	 * @param payload: {username:<string>,password:<string>,onceLogin:<bool>}
	 * @returns {Promise<any>}
	 * @constructor
	 */
	AUTH_BY_USERNAME: (context, payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + `/app/auth/username`, method: 'post', data: payload})
				.then(resp => {

					if (resp['token'] === undefined) reject(new Error('Не получен токен авторизации'))
					context.dispatch('app/SET_TOKEN', resp['token'], { root: true })

					if (resp['user'] !== undefined) context.dispatch('SET_PROFILE', resp['user']); // хз хз
					else context.dispatch('UPDATE_PROFILE')

					if (resp['aUsers'] !== undefined) context.dispatch('SET_ACCOUNT_USERS', resp['aUsers']); // хз хз
					//else context.dispatch('UPDATE_AVAILABLE_ACCOUNTS')

					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},

	// перенести в accounts ([GET] /accounts/${payload['account_id']}/auth)
	LOGIN_IN_ACCOUNT: (context, payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/user/auth`, method: 'get'})
				.then(resp => {
					if (resp['token'] !== undefined && resp['token'] !== "") context.dispatch('app/SET_TOKEN', resp['token'], { root: true })

					if (resp['account'] !== undefined) context.dispatch('account/SET_PROFILE', resp['account'], {root:true}); // хз хз
					else context.dispatch('account/GET', nil, {root:true})

					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},
	CREATE: (context, payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/accounts/1/users', method: 'post',data: payload}) // всегда account_id = 1
				.then(resp => resolve(resp))
				.catch((err) => reject(err))
		})
	},
	SET_PROFILE: ({commit}, data) => commit('SET_PROFILE', data),
	UPDATE_PROFILE: (context) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/user', method: 'get'})
				.then(resp => {
					if (resp['user'] !== undefined) context.dispatch('SET_PROFILE', resp["user"]);
					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},

	/**
	 * Отправляет пользователю email со ссылкой для сброса пароля по {"username":"..."}
	 * @param context
	 * @param payload
	 * @returns {Promise<any>}
	 * @constructor
	 */
	PASSWORD_RECOVERY_SEND_EMAIL: (context,payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/password/reset/send-email', data: payload, method: 'post'})
				.then(resp => resolve(resp))
				.catch((err) => reject(err))
		})
	},

	/**
	 * !!! NEW Сбрасывает пароль сообщая токен и и получает ключ авторизации
	 * @param context
	 * @param payload: {"token":"<...>"}
	 * @returns {Promise<any>}
	 * @constructor
	 */
	PASSWORD_RESET_CONFIRM: (context,payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/password/reset/confirm', data: payload, method: 'post'})
				.then((resp) => {

					if (resp['token'] === undefined) reject(new Error('Не получен токен авторизации'))
					context.dispatch('app/SET_TOKEN', resp['token'], { root: true })

					if (resp['user'] !== undefined) context.dispatch('SET_PROFILE', resp['user']); // хз хз
					else context.dispatch('UPDATE_PROFILE')

					if (resp['accounts'] !== undefined) context.dispatch('SET_AVAILABLE_ACCOUNTS', resp['accounts']); // хз хз
					else context.dispatch('UPDATE_AVAILABLE_ACCOUNTS')

					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},

	SET_PASSWORD: (context,payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/password', data: payload, method: 'post'})
				.then(resp => {
					if (resp['user'] !== undefined) context.dispatch('SET_PROFILE', resp['user']);
					else context.dispatch('UPDATE_PROFILE')
					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},
	SEND_USERNAME_RECOVERY: (context,payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/recovery/username', data: payload, method: 'post'})
				.then(resp => resolve(resp))
				.catch((err) => reject(err))
		})
	},
	// отправка кода для подтверждения
	EMAIL_VERIFICATION: (context, payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/email-verification', data: payload, method: 'post'})
				.then(resp => resolve(resp))
				.catch((err) => reject(err))
		})
	},
	// Повторная отправка кода верификации для подтвержения емейла
	SEND_EMAIL_INVITE_VERIFICATION: (context, payload) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/email-verification/invite', data: payload, method: 'get'})
				.then(resp => resolve(resp))
				.catch((err) => reject(err))
		})
	},
	SET_AVAILABLE_ACCOUNTS: ({commit}, accounts) => commit('SET_AVAILABLE_ACCOUNTS', accounts),
	SET_ACCOUNT_USERS: ({commit}, aUsers) => commit('SET_ACCOUNT_USERS', aUsers),
	UPDATE_AVAILABLE_ACCOUNTS: (context) => {
		return new Promise((resolve, reject) => {
			api({url: window.VUE_APP_SERVER_GUI_URL + '/users/accounts', method: 'get'})
				.then(resp => {
					if (resp['aUsers'] !== undefined) context.dispatch('SET_ACCOUNT_USERS', resp['aUsers']);
					// else // todo some code there...
					resolve(resp)
				})
				.catch((err) => reject(err))
		})
	},
};

const mutations = {
	SET_DEFAULT_STATE: (state) => Object.assign(state, getDefaultState()),
	SET_PROFILE: (state, data) => state.profile = data,
	SET_AVAILABLE_ACCOUNTS: (state, data) => state.accounts = data,
	SET_ACCOUNT_USERS: (state, data) => state.aUsers = data,
};

const getDefaultState = () => {
	return {
		profile: {
			id: 0
		},
		accounts: [],
		roles:[],
		aUsers: [],
	}
}

export default {
    state,
    getters,
    actions,
    mutations,
     namespaced: true
}