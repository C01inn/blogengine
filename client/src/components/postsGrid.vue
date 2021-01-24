<template>
    <div id="articleGrid">

        <div id="articles" class="mx-auto">
            <!-- old card -->
            <!-- 
            <div class="shadow-2xl rounded-xl p-1 border border-gray-300 cursor-pointer article-card" v-for="article in articles" :key="article" @click="articleClick(article.PostId)">
                <h1 class="p-2 article-title" v-if="article.ImageId">{{article.Title}}</h1>
                <h1 class="p-2 article-title-2" v-else-if="!article.ImageId">{{article.Title}}</h1>
                <img :src="imageBaseUrl + article.ImageId + `.` + article.Filetype" class="p-2 rounded-lg article-image" v-if="article.ImageId">
            </div>  -->
            <div class="shadow-2xl border border-gray-300 cursor-pointer article-card" v-for="article in articles" :key="article" @click="articleClick(article.PostId)">
                <h1 class="p-2 article-title" v-if="article.ImageId">{{article.Title}}</h1>
                <h1 class="p-2 article-title-2" v-else-if="!article.ImageId">{{article.Title}}</h1>
                <img :src="imageBaseUrl + article.ImageId + `.` + article.Filetype" class="article-image" v-if="article.ImageId">
            </div>
        </div>

    </div>
</template>

<script>
export default {
    name: "postsGrid",
    components: {

    },
    props: {
        articles: Array,
        imageBaseUrl: String
    },
    methods: {
        articleClick(id) {
            this.$router.push(`/article/${id}`);
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
  min-width: 80%;
  max-width: 100%;
}

.article-card {
  font-family: 'Lato', sans-serif;
  color: #0D0D0D;
  font-weight: 600;
}

.article-card:hover {
  animation: articleHover 500ms ease;
  animation-fill-mode: forwards;
}

.article-title-2 {
  font-size: 2em;
}

.article-title {
  font-size: 1.5em;
}

#articles {
  margin-top: 2rem;
  margin-bottom: 5rem;
  display: grid;
  gap: 2em;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  width: 80vw;
}

@media screen and (max-width: 1540px) {
  #articles {
    margin-top: 2rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1em;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    width: 80vw;
  }
}
@media screen and (max-width: 1200px) {
  #articles {
    margin-top: 2rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1em;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    width: 90vw;
  }
  .article-title {
    font-size: 1.3em;
  }
}

#articles div {
  height: fit-content;
}
</style>