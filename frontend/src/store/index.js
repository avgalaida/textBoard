import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueNativeSock from 'vue-native-websocket';

const BACKEND_URL = 'http://localhost:8080';
const PUSHER_URL = 'ws://localhost:8080/pusher';

const SET_POSTS = 'SET_POSTS';
const CREATE_POST = 'CREATE_POST';
const SEARCH_SUCCESS = 'SEARCH_SUCCESS';
const SEARCH_ERROR = 'SEARCH_ERROR';

const MESSAGE_POST_CREATED = 1;

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    posts: [],
    searchResults: [],
  },
  mutations: {
    SOCKET_ONOPEN(state, event) {
    },
    SOCKET_ONCLOSE(state, event) {
    },
    SOCKET_ONERROR(state, event) {
      console.error(event);
    },
    SOCKET_ONMESSAGE(state, message) {
      switch (message.kind) {
        case MESSAGE_POST_CREATED:
          this.commit(CREATE_POST, { id: message.id, body: message.body });
      }
    },
    [SET_POSTS](state, posts) {
      state.posts = posts;
    },
    [CREATE_POST](state, post) {
      state.posts = [post, ...state.posts];
    },
    [SEARCH_SUCCESS](state, posts) {
      state.searchResults = posts;
    },
    [SEARCH_ERROR](state) {
      state.searchResults = [];
    },
  },
  actions: {
    getPosts({ commit }) {
      axios
          .get(`${BACKEND_URL}/posts`)
          .then(({ data }) => {
            commit(SET_POSTS, data);
          })
          .catch((err) => console.error(err));
    },
    async createPost({ commit }, post) {
      const { data } = await axios.post(`${BACKEND_URL}/posts`, null, {
        params: {
          body: post.body,
        },
      });
    },
    async searchPosts({ commit }, query) {
      if (query.length === 0) {
        commit(SEARCH_SUCCESS, []);
        return;
      }
      axios
          .get(`${BACKEND_URL}/search`, {
            params: { query },
          })
          .then(({ data }) => commit(SEARCH_SUCCESS, data))
          .catch((err) => {
            console.error(err);
            commit(SEARCH_ERROR);
          });
    },
  },
});

Vue.use(VueNativeSock, PUSHER_URL, { store, format: 'json' });

store.dispatch('getPosts');

export default store;