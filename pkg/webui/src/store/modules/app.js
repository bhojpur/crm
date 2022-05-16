import axios from 'axios'

// Модуль отвечающий за организационные вопрос, связанные с авторизацией: авторизация, аутентификация, регистрация, восстановление пароля, почты, выход из аккаунта.

import api from '../../api/utils'

const state = {
    token: '', // jwt-token
    settings: {},
    hasLoadedOnce: false,
}

const getters = {
//    isAuthenticated: state => !(state.token === "" || state.token === undefined),
    /**
     * Уровень доступа GuestAccess - пользователь неавторизован
     * @param state
     */
    isAuthenticated: state => !!state.token, // мутная хрень
    hasToken: state => !!state.token, // токен может быть полного доступа и не полного
    getToken: state => state.token,

}

const actions = {
    /**
     * Важная функция - выставляет системные настройки на базовое значение.
     * @param context
     * @constructor
     */
    SET_DEFAULT_STATE: (context) => {
        context.commit('SET_DEFAULT_STATE')
        delete axios.defaults.headers.common['Authorization']
        // window.$broadcastTab.message({command: "logout"})
        //context.dispatch("GET_SETTINGS") // тут хз... скорее в момент логина надо
        // this.$router.push({ name: 'login',  params: { msg: 'Авторизуйтесь в системе' } }).catch()
    },
    AUTH_LOGOUT: {
        root: true,
        handler: (context) => context.dispatch("SET_DEFAULT_STATE"),
    },
    /**
     * Проверяем есть ли полная авторизация User + Account, чтобы продолжить работу. Бэкенд также проверяет имеет ли пользовател доступ к аккаунту.
     * @param dispatch
     * @param commit
     * @param getters
     * @returns {Promise<any>}
     * @constructor
     */
    AUTH_CHECK: ({dispatch, commit, getters}) => {
        return new Promise((resolve, reject) => {
            if (getters.isAuthenticated) {
                // устанавливаем токен
                //axios.defaults.headers.common['Authorization'] = 'Bearer ' + getters.getToken
                api({url: window.VUE_APP_SERVER_GUI_URL + '/app/auth/check/full', method: 'GET'})
                    .then(() => resolve() )
                    .catch(() => {
//                        dispatch('auth/AUTH_LOGOUT', {root: true})
                        reject()
                    })

            } else {
                reject()
            }
        })
    },
    AUTH_CHECK_USER: () => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + '/app/auth/user', method: 'GET'})
                .then(() => resolve())
                .catch(() => reject())
        })
    },
    AUTH_CHECK_ACCOUNT: () => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + '/app/auth/account', method: 'GET'})
                .then(() => resolve())
                .catch(() => reject())
        })
    },
    SET_TOKEN: (context,token) => {
        return new Promise((resolve) => {
            axios.defaults.headers.common['Authorization'] = 'Bearer ' + token
            context.commit("SET_TOKEN", token)
            resolve()
        })

    },
    SET_SETTINGS: (context,settings) => {
        context.commit("SET_SETTINGS", settings)
    },
    GET_SETTINGS: (context) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + '/app/settings', method: 'get'})
                .then(resp => {
                    if (resp['settings'] !== undefined) context.dispatch('SET_SETTINGS', resp['settings']);
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    }
}

const mutations = {
    SET_DEFAULT_STATE: (state) => Object.assign(state, getDefaultState()),
    SET_TOKEN: (state, token) => state.token = token,
    SET_SETTINGS: (state, settings) => state.settings = settings,
}

const getDefaultState = () => state

export default {
    state,
    getters,
    actions,
    mutations,
     namespaced: true,
}