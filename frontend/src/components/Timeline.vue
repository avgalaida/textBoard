<template>
  <div>
    <form v-on:submit.prevent="createPost">
      <div class="input-group">
        <input v-model.trim="postBody" type="text" class="form-control" placeholder="Что нового?">
        <div class="input-group-append">
          <button class="btn btn-primary" type="submit">Отправить</button>
        </div>
      </div>
    </form>

    <div class="mt-4">
      <Post v-for="post in posts" :key="post.id" :post="post" />
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import Post from '@/components/Post';
export default {
  data() {
    return {
      postBody: '',
    };
  },
  computed: mapState({
    posts: (state) => state.posts,
  }),
  methods: {
    createPost() {
      if (this.postBody.length != 0) {
        this.$store.dispatch('createPost', { body: this.postBody });
        this.postBody = '';
      }
    },
  },
  components: {
    Post,
  },
};
</script>