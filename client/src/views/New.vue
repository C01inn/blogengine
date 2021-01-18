<template>
    <div id='new-post'>
        <Navbar class="sticky top-0"/>

        <h1 class="mt-6" style="font-size: 3em;">Create A New Post</h1>
        <input id="title-input" class="bg-white rounded-lg border border-solid border-gray-400 p-1 bg-gray-50 shadow" placeholder="Post Title">
        <div id="main">
            <div id="editor-1" class="shadow-2xl">
                <div id="editor-container" class="mx-auto"></div>
            </div>
            <button @click="submitPost()" class="mt-3 bg-blue-500 hover:bg-blue-400 text-white rounded-lg px-6 py-2">Submit Post</button>
        </div>

    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import * as Quill from '@/assets/quill.min.js'
require('@/assets/quill.css')
import { ImageResize } from 'quill-image-resize-module';
import $ from 'jquery'
import axios from 'axios'

Quill.register('modules/imageResize', ImageResize);

export default {
    name: "New",
    components: {
        Navbar,
    },
    data() {
        return {
            apiUrl: this.$store.state.apiUrl
        }
    },
    created() {
        if (sessionStorage.getItem("userToken") === null) {
            this.$router.push('/')
        }
    },
    mounted() {
        const toolbarOptions = [
            [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
            [{ 'header': [1, 2, 3, 4, 5, 6, false] }],

            [{ 'color': [] }],          // dropdown with defaults from theme
            [{ 'font': [] }],
            ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
            ["link", "image"],
            ['blockquote', 'code-block'],

            [{ 'list': 'ordered'}, { 'list': 'bullet' }],
            [{ 'script': 'sub'}, { 'script': 'super' }],      // superscript/subscript
            [{ 'indent': '-1'}, { 'indent': '+1' }],          // outdent/indent
            [{ 'direction': 'rtl' }],                         // text direction

            [{ 'align': 'left' }],
        ];
        var quill = new Quill('#editor-container', {
            theme: 'snow',
            modules: {
                imageResize: {
                    displaySize: true
                },
                toolbar: toolbarOptions
            }
        });
        quill.format("align", "left")
    },
    methods: {
        async submitPost() {
            let editor = document.querySelector("#editor-container");
            let htmlString = editor.children[0].innerHTML;
            let postTitle = $("#title-input").val()

            if (postTitle === null || postTitle === '') {
                alert("You must enter a title for your blog must.")
                return;
            }
            
            // get images as files
            let imageArr = []
            let images = $(htmlString).find('img')
            for (let i = 0;i<images.length;i++) {
                let newImg = this.dataURLtoFile(images[i].src, `file${i}`)
                imageArr.push(newImg)
            }

            // replace image srcs with nothing
            htmlString = $("<div>").html(htmlString)
            htmlString.find("img").attr("src", "")
            htmlString = htmlString.html()
            
            // create form data
            let formSize = 0;
            let formData = new FormData();
            formData.append('html', htmlString)
            formSize = formSize + htmlString.length
            formData.append('title', postTitle)
            formSize = formSize + postTitle.length
            formData.append("userToken", sessionStorage.getItem('userToken'))
            formSize = formSize + sessionStorage.getItem('userToken').length
            // add files to formdata
            for (let i=0; i<imageArr.length;i++) {
                formData.append(`files${i}`, imageArr[i])
                formSize = formSize + imageArr[i].size
            }
            // size of form in mb
            formSize = formSize / 1024/1024
            if (formSize > 4.5) {
                alert("Your article content size is too big. This is most likely due to large image file sizes. The max article size is 4.5 MB.")
                return;
            }


            // make request
            let resp = await axios.post(this.apiUrl + "/new-post", formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })

            if (resp.data === 'error') {
                alert('Internal Server Error.')
            } else {
                this.$router.push('/article/' + resp.data.id)
            }

            
            
        },
        dataURLtoFile(dataurl, filename) {
    
            var arr = dataurl.split(','),
                mime = arr[0].match(/:(.*?);/)[1],
                bstr = atob(arr[1]), 
                n = bstr.length, 
                u8arr = new Uint8Array(n);
                
            while(n--){
                u8arr[n] = bstr.charCodeAt(n);
            }
            
            return new File([u8arr], filename, {type:mime});
        },
    }

}
</script>
<style scoped>
#title-input {
    width: 60%;
    font-size: 2em;
}
#main {
    margin-top: 1rem;
    margin-right: auto;
    margin-left: auto;
    width: 60%;
    height: 80vh;
}

#editor-1 {
    margin-right: auto;
    margin-left: auto;
    width: 100%;
    height: 80%;
    padding-bottom: 40px;
}

</style>