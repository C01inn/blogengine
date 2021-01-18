<template>
    <div id="login">
        <Navbar class="sticky top-0"/>

        <div class="mx-auto rounded-xl shadow-2xl bg-white border border-gray-500" id="login-card">
            <h1 class="text-3xl my-4">Login Below</h1>
            <form class="mx-auto w-3/5" @submit.prevent="login()" id="login-form">
                <h4 style="text-align: left;" class="text-lg">Email: </h4>
                <input type="email" id="email-input" class="form-input mt-1 block mx-auto p-3 bg-gray-100 rounded-lg w-full" placeholder="Email" required>
                <br>
                <h4 style="text-align: left;" class="text-lg">Password: </h4>
                <input id="pass-input" class="form-input mt-1 block mx-auto p-3 bg-gray-100 rounded-lg w-full" type="password" placeholder="Password" required>
                <br>
                <button type="submit" class="bg-blue-500 hover:bg-blue-400 rounded-lg text-white text-lg px-6 py-1 mb-5">Login</button>
            </form>
        </div>
    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'
import $ from 'jquery'

export default {
    name: "Login",
    components: {
        Navbar
    },
    data() {
        return {
            apiUrl: this.$store.state.apiUrl
        }
    },
    created() {
        if (sessionStorage.getItem("userToken") !== null) {
            this.$router.push('/')
        }
    },
    methods: {
        async login() {
            let email = $("#email-input").val()
            let password = $("#pass-input").val()
            if (email == null || password == null) {
                return
            }
            let resp = await axios.post(this.apiUrl + '/login', {
                Email: email,
                Password: password,
            })
            if (resp.data === 'error') {
                alert("Internal Server Error")
                return
            } else if (resp.data === 'bad email') {
                alert("Wrong Email")
                return
            } else if (resp.data === 'bad password') {
                alert("Wrong Password")
                return
            }
            sessionStorage.setItem("userToken", resp.data.id)
            sessionStorage.setItem("username", resp.data.username)


            alert("Logged In :)")
            this.$router.push('/')

        }
    }
}
</script>
<style scoped>
#login-card {
    position: relative;
    top: 5rem;
    width: 35vw;
    --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
}
@media screen and (max-width: 1200px) {
    #login-card {
        position: relative;
        top: 5rem;
        width: 59vw;
        --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
    }
}
@media screen and (max-width: 550px) {
    #login-card {
        position: relative;
        top: 5rem;
        width: 75vw;
        --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
    }
}

</style>