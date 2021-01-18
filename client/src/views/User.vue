<template>
    <div id="user">
        <Navbar class="sticky top-0"/>
        <h1 class="username mt-8">Posts from @{{username}}</h1>
        <div id="articles" class="mx-auto">
            <div class="shadow-2xl rounded-xl p-1 border border-gray-300 cursor-pointer article-card" v-for="article in articles" :key="article" @click="articleClick(article.PostId)">
                <h1 class="p-2 article-title">{{article.Title}}</h1>
                <img :src="imageBaseUrl + article.ImageId + `.` + article.Filetype" class="p-2 rounded-lg article-image" v-if="article.ImageId">
            </div>
        </div>
    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'

export default {
    name: "User",
    components: {
        Navbar
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
@keyframes articleHover {
  from {

  }
  to {
    box-shadow: rgba(0, 0, 0, 0.25) 0px 54px 55px, rgba(0, 0, 0, 0.12) 0px -12px 30px, rgba(0, 0, 0, 0.12) 0px 4px 6px, rgba(0, 0, 0, 0.17) 0px 12px 13px, rgba(0, 0, 0, 0.09) 0px -3px 5px;
  }
}

.article-image {
  margin-right: auto;
  margin-left: auto;
  max-width: 90%;
}

.article-card:hover {
  animation: articleHover 500ms ease;
  animation-fill-mode: forwards;
}

.article-title {
  font-size: 1.7em;
}

#articles {
  margin-top: 2rem;
  margin-bottom: 5rem;
  display: grid;
  gap: 2rem;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  width: 80vw;
}

@media screen and (max-width: 1540px) {
  #articles {
    margin-top: 2rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    width: 80vw;
  }
}
@media screen and (max-width: 1200px) {
  #articles {
    margin-top: 2rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    width: 90vw;
  }
}

#articles div {
  height: fit-content;
}



.username {
    font-size: 2.7em;
}
</style>