import Vue from 'vue'
import axios from 'axios'

import api from '../../api/utils'

const state = {
    users: [],
    roles: [],
}

const getters = {
    getUserList: state => state.users,
    getRoleList: state => state.roles,
    getRole: state => id => state.roles.find(item => item.id === id),
    /*getRole2: state => id => {
        console.log("Ищем id: ", id)
        let index = state.roles.findIndex(key=> key.id === id, id)
        console.log("Role id:", index)
        console.log("Role:", state.roles[index])
        return ""
        // let index = state.roles.findIndex(key, id => key.id === id)

        if (index !== -1) {
            console.log("Role id -1")
            console.log( state.roles)
            return null
        }
        else {
            console.log(state.roles[index])
            return state.roles[index]
        }
    },*/
    countUsers: state => state.users.length,
}

const actions = {
    SET_DEFAULT_STATE: (context) => context.commit('SET_DEFAULT_STATE'),
    AUTH_LOGOUT: {
        root: true,
        handler: (context) => context.dispatch("SET_DEFAULT_STATE"),
    },
    CREATE: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({ url: window.VUE_APP_SERVER_GUI_URL + `/users`, method: 'post', data: payload })
                .then(async resp => {
                    if (resp['api_key'] !== undefined) context.commit('API_KEY_ADD', resp['api_key'])
                    else await context.dispatch('API_KEY_LIST')
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },
    GET: (context,payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/storage/${payload['hash_id']}`, method: 'get'})
                .then(resp => resolve(resp))
                .catch((err) => reject(err))
        })
    },
    GET_ROLE_LIST: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/roles`, method: 'get'})
                .then(resp => {
                    if (resp['roles'] !== undefined) context.commit('SET_ROLE_LIST', resp["roles"]);
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },
    GET_USERS_DATA: (context, payload) => {
        return new Promise((resolve, reject) => {
            // api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/users?offset=${payload['offset']}&limit=${payload['limit']}&types=name1,name2`, method: 'get'})
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/users?offset=${payload['offset']}&limit=${payload['limit']}`, method: 'get'})
                .then(resp => resolve(resp))
                .catch((err) => reject(err))
        })
    },
    SEARCH_USERS_DATA: (context, payload) => {
        if (payload['search'] === undefined) {
            console.error("не верные данные для поиска: ", payload['search'])
            return
        }
        return new Promise((resolve, reject) => {

            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/users?offset=${payload['offset']}&limit=${payload['limit']}&search=${payload['search']}`, method: 'get'})
                .then(resp => resolve(resp))
                .catch((err) => reject(err))
        })
    },

    UPDATE_USER: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/users/${payload['hash_id']}`, method: 'patch', data: payload})
                .then(resp => {
                    // context.commit('UPDATE_USER', resp['user'])
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },
    REMOVE_USER_FROM_ACCOUNT: (context,payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/users/${payload['hash_id']}`, method: 'delete'})
                .then(resp => {
                    context.commit('DELETE_USER', payload['hash_id']) // не очень нужная то штука
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },
}

const mutations = {
    SET_DEFAULT_STATE: (state) => Object.assign(state, getDefaultState()),
    ADD: (state, user) => state.users.push(user),
    UPDATE_USER: (state, user) => {
        if (user !== undefined && user !== null) {
            let id = state.users['id']
            let index = state.users.findIndex(item => item.id === id, id )
            state.users[index] = user
        }

        // else Object.assign(state.api_keys, api_keys)
    },
    SET_USER_LIST: (state, users) => {
        if (users !== undefined) state.users = users
    },
    SET_ROLE_LIST: (state, roles) => {
        if (roles !== undefined) state.roles = roles

    },
    DELETE_USER: (state, id) => {
        let index = state.users.findIndex(key => key.hash_id === id, id)
        if (index !== -1) Vue.delete(state.users, index)
    },
}

const getDefaultState = () => {
    return {
        users: [],
        roles: [],
    }
}


export default {
    state,
    getters,
    actions,
    mutations,
     namespaced: true
}