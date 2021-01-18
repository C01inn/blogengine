<template>
    <div id="account">
        <Navbar class="sticky top-0" style="z-index: 3;"/>

        <h1 class="header-1">Your Blog Posts</h1>

        <div id="articles" class="mx-auto">
            <div class="article-options border border-gray-300 rounded-xl" v-for="article in articles" :key="article" >
                <button class="bg-red-500 hover:bg-red-400 rounded-lg text-white px-4 py-1 delete-btn" @click="removeArticle(article.PostId)">Delete Article</button>
                <button class="bg-blue-500 hover:bg-blue-400 rounded-lg text-white px-4 py-1 read-btn" @click="articleClick(article.PostId)">Read Article</button>
                <div class="shadow-2xl rounded-xl p-1 cursor-pointer article-card" @click="articleClick(article.PostId)">
                    <h1 class="p-2 article-title">{{article.Title}}</h1>
                    <img :src="imageBaseUrl + article.ImageId + `.` + article.Filetype" class="p-2 rounded-lg article-image" v-if="article.ImageId">
                </div>
            </div>
        </div>

    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'

export default {
    name: "Account",
    components: {
        Navbar
    },
    data() {
        return {
            apiUrl: this.$store.state.apiUrl,
            imageBaseUrl: null,
            articles: null,
            username: null,
        }
    },
    async created() {
        this.username = sessionStorage.getItem('username')
        if (sessionStorage.getItem('userToken') === null) {
            this.$router.push('/')
        }

        let resp = await axios.get(this.apiUrl +'/user-articles/'+ this.username);
        this.imageBaseUrl = resp.data.ImageBaseUrl;
        this.articles = resp.data.Articles
    },
    methods: {
        articleClick(id) {
            this.$router.push('/article/' + id);
        },
        async removeArticle(id) {
            
            // confirm that the user wants to delete this article
            if (!confirm("Are you sure you want to remove this article? This cannot be undone.")) {
                return
            }
            // make delete request
            let deleteData = {
              UserId: sessionStorage.getItem("userToken"),
              PostId: id
            }

            fetch(this.apiUrl + '/remove-article', {
              method: "delete",
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify(deleteData)
            })
            .then(res => res.text())
            .then(response => {
              if (response === 'error') {
                alert("Internal Server Error. Please try again later.")
                return
              }
              this.$router.go()
            })


        }
    }
}
</script>
<style scoped>

.article-image {
  margin-right: auto;
  margin-left: auto;
  max-width: 90%;
}

.read-btn {
    display: inline-block;
    position: relative;
    top: 1rem;
    left: 2rem;
    margin: 0;
    z-index: 0;
}

.delete-btn {
    display: inline-block;
    position: relative;
    top: 1rem;
    left: 0;
    margin: 0;
    z-index: 0;
}

.article-options {
    display: block;
    height: calc(fit-content + 5rem);
}

@keyframes articleHover {
  from {

  }
  to {
    box-shadow: rgba(0, 0, 0, 0.25) 0px 54px 55px, rgba(0, 0, 0, 0.12) 0px -12px 30px, rgba(0, 0, 0, 0.12) 0px 4px 6px, rgba(0, 0, 0, 0.17) 0px 12px 13px, rgba(0, 0, 0, 0.09) 0px -3px 5px;
  }
}

.article-card {
    --tw-border-opacity: 1;
    border-top: 1px solid rgba(209, 213, 219, var(--tw-border-opacity));
}

.article-card:hover {
  animation: articleHover 500ms ease;
  animation-fill-mode: forwards;
}

.article-title {
  font-size: 1.7em;
}

#articles {
  margin-top: 0rem;
  margin-bottom: 5rem;
  display: grid;
  gap: 2rem;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  width: 80vw;
}

@media screen and (max-width: 1540px) {
  #articles {
    margin-top: 0rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    width: 80vw;
  }
}
@media screen and (max-width: 1200px) {
  #articles {
    margin-top: 0rem;
    margin-bottom: 5rem;
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    width: 90vw;
  }
}

#articles div {
    margin-top: 2rem;
    height: fit-content;
}



.header-1 {
  margin-top: 2rem;
  font-size: 2.8em;
}
</style>