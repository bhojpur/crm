import apiCall from '../../api/utils'
import Vue from 'vue'
import axios from 'axios'

import api from '../../api/utils'

const state = {
    id: '',
    hash_id: '',
}

const getters = {
    isAuth: state => state.hash_id !== '',
    hash_id: state => state.hash_id,
}

const actions = {
    SET_DEFAULT_STATE: (context) => context.commit('SET_DEFAULT_STATE'),
    SET_PROFILE: (context, profile) =>  context.commit('SET_PROFILE', profile),
    AUTH_LOGOUT: {
        root: true,
        handler: (context) => context.dispatch("SET_DEFAULT_STATE"),
    },
    ACCOUNT_LOGOUT: {
        root: true,
        handler: (context) => context.dispatch("SET_DEFAULT_STATE"),
    },
    // перенести в accounts ([GET] /accounts/${payload['account_id']}/auth)
    LOGIN_IN_ACCOUNT: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['accountHashId']}/auth`, data: payload, method: 'post'})
                .then(resp => {
                    if (resp['token'] !== undefined && resp['token'] !== "") context.dispatch('app/SET_TOKEN', resp['token'], { root: true })
                    else reject()

                    if (resp['account'] !== undefined) context.dispatch('account/SET_PROFILE', resp['account'], {root:true}); // хз хз
                    else context.dispatch('account/GET', payload, {root:true})

                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },

    CREATE: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + '/accounts', data: payload, method: 'post'})
                .then(resp => resolve(resp))
                .catch((err) => reject(err))
        })
    },
    UPDATE: (context, payload) => {
        let postData = {}
        Object.assign(postData,payload)
        delete postData['hash_id']
        delete postData['accountHashId']

        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['hash_id']}`, method: 'patch', data: postData})
                .then(resp => {
                    // todo: его бы в базе обновить
                    context.commit('UPDATE',resp['account'])
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },
    GET: (context, payload) => {
        return new Promise((resolve, reject) => {
            api({url: window.VUE_APP_SERVER_GUI_URL + `/accounts/${payload['hash_id']}`, method: 'get'})
                .then(resp => {
                    if (resp['account'] !== undefined) context.commit('UPDATE', resp["account"]); // хз хз
                    resolve(resp)
                })
                .catch((err) => reject(err))
        })
    },

}

const mutations = {
    SET_DEFAULT_STATE: (state) => Object.assign(state, getDefaultState()),
    UPDATE: (state, data) => Object.keys(data).forEach((v)=> Vue.set(state, v, data[v])),
    SET_PROFILE: (state, data) => Object.keys(data).forEach((v)=> Vue.set(state, v, data[v])),
}

const getDefaultState = () => {
    return {
        id: '',
        hash_id: ''
    }
}

export default {
    state,
    getters,
    actions,
    mutations,
     namespaced: true
}