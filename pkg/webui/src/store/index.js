import Vue from 'vue'
import Vuex from 'vuex'

import createPersistedState from 'vuex-persistedstate'
import createMutationsSharer from "vuex-shared-mutations";

import app from './modules/app'
import user from './modules/user'
import account from './modules/account'
import users from './modules/users'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production';

/*let mutations = Object.keys( require("./actions/auth"))
    .concat(Object.keys(require("./actions/user"))
        .concat(Object.keys(require("./actions/account"))
    ));*/

//let mutations = Object.keys( require("./actions/auth")).concat(Object.keys(require("./actions/account")));
let mutations = [
    'app/SET_DEFAULT_STATE','app/SET_TOKEN', 'app/SET_SETTINGS',
    'user/SET_DEFAULT_STATE','user/SET_PROFILE','user/SET_AVAILABLE_ACCOUNTS', 'user/SET_ACCOUNT_USERS',
    'account/SET_DEFAULT_STATE','account/UPDATE','account/SET_PROFILE',
    'users/SET_DEFAULT_STATE', 'users/ADD','users/UPDATE','users/LIST_SET','users/DELETE',
];

export default new Vuex.Store({
    modules: {
        app,
        user,
        account,
        users,
    },
    plugins: [createPersistedState(), createMutationsSharer({ predicate: mutations })],
//    plugins: [createPersistedState(), createMutationsSharer({ predicate: mutations })],
//    plugins: [createPersistedState()],
    strict: debug,
});