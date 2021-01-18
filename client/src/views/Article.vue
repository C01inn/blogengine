<template>
    <div id="article">
        <Navbar class="sticky top-0"/>

        <div id="main" class="mx-auto shadow-2xl">
            <h1 class="text-2xl sm:text-3xl md:text-3xl lg:text-3xl xl:text-4xl 2xl:text-4xl">{{articleTitle}}</h1>
            <h3 class="text-md sm:text-lg md:text-xl lg:text-xl xl:text-xl 2xl:text-xl">Written By: <router-link id="author-link" :to="'/@' + articleAuthor">@{{articleAuthor}}</router-link></h3>
            <div id="article-content" class="mx-auto article-div" v-html="articleHtml">

            </div>
        </div>

    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'
import $ from 'jquery'

export default {
    name: "Article",
    components: {
        Navbar
    },
    data() {
        return {
            apiUrl: this.$store.state.apiUrl,
            articleId: this.$route.params.id,
            articleHtml: null,
            articleTitle: null,
            articleAuthor: null,
        }
    },
    async created() {
        // make request
        let resp = await axios.get(this.apiUrl + '/get-post/' + this.articleId);
        let htmlData = resp.data.Html
        let imagesData = resp.data["Images"]

        
        // add correct image srcs using jquery
        htmlData = $("<div>").html(htmlData)
        htmlData.find("img").each(function(idx, elem) {
            // the following line will stop "no unused variables error from eslint"
            console.log(elem.toString().replace(/./g, ""))
            let imageName = imagesData[idx].Id + "." + imagesData[idx].Filetype
            let imageUrl = resp.data.ImageBaseUrl + imageName
            $(this).attr("src", imageUrl);
            // set margins also
            $(this).attr("style", "margin-left: auto; margin-right: auto; min-width: 25%;")
        })
        htmlData = htmlData.html()

        // set article data variables
        this.articleTitle = resp.data.Title
        this.articleAuthor = resp.data.Author
        this.articleHtml = htmlData
    },
    mounted() {
        // set minimum font size based on screen width
        let minFontSize;
        if ($(window).width() > 1000) {
            minFontSize = 19
        } else {
            minFontSize = 16
        }
        $("#article-content").each( function () {
            var $this = $(this);
            if (parseInt($this.css("fontSize")) < minFontSize) {
                $this.css({ "font-size": minFontSize+"px" });   
            }
        });
    }
}
</script>
<style scoped>
#main {
    margin-top: 2rem;
    border-radius: 1rem;
    min-width: 40vw;
    width: fit-content;
    max-width: 60vw;
    padding: 2rem;
    margin-bottom: 10rem;
}

#main h1 {
    text-align: center;
    margin-bottom: .5rem;
}

#article-content {
    margin-top: 3rem;
    margin-left: auto;
    margin-right: auto;
    display: block;
    text-align: center !important;
}

@media screen and (max-width: 1200px) {
    #main {
        margin-top: 2rem;
        border-radius: 1rem;
        min-width: 60vw;
        width: fit-content;
        max-width: 75vw;
        padding: 2rem;
        margin-bottom: 10rem;
    }
}

@media screen and (max-width: 750px) {
    #main {
        margin-top: 2rem;
        border-radius: 1rem;
        min-width: 70vw;
        width: fit-content;
        max-width: 90vw;
        padding: 1rem;
        margin-bottom: 10rem;
    }
    #article-content {
        margin-top: 1.5rem;
        margin-left: auto;
        margin-right: auto;
        display: block;
        text-align: center !important;
    }
}



#author-link {
    color: blue;
    cursor: pointer;
}
#author-link:hover {
    text-decoration: underline;
    color: darkblue;
}


</style>