<template>
  <div>
    <input @keyup="searchPosts" v-model.trim="query" class="form-control" placeholder="Search...">
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
      query: '',
    };
  },
  computed: mapState({
    posts: (state) => state.searchResults,
  }),
  methods: {
    searchPosts() {
      if (this.query != this.lastQuery) {
        this.$store.dispatch('searchPosts', this.query);
        this.lastQuery = this.query;
      }
    },
  },
  components: {
    Post,
  },
};
</script>
Footer
