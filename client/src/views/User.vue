<template>
    <div id="user">
        <Navbar class="sticky top-0"/>
        <h1 class="username mt-8">Posts from @{{username}}</h1>
        <postsGrid :articles="articles" :imageBaseUrl="imageBaseUrl"/>
    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'
import postsGrid from '@/components/postsGrid.vue'

export default {
    name: "User",
    components: {
        Navbar,
        postsGrid,
    },
    data() {
        return {
            apiUrl: this.$store.state.apiUrl,
            username: this.$route.params.id,
            imageBaseUrl: null,
            articles: null,
        }
    },
    async created() {
        let resp = await axios.get(this.apiUrl + '/user-articles/' + this.username);
        this.imageBaseUrl = resp.data.ImageBaseUrl;
        this.articles = resp.data.Articles;
    },
    methods: {
        articleClick(id) {
            this.$router.push('/article/' + id);
        }
    }
}
</script>

<style scoped>

.username {
    font-size: 2.7em;
}
</style>