<template>
    <div id='signup'>
        <Navbar class="sticky top-0"/>

        <div class="mx-auto rounded-xl shadow-2xl bg-white border border-gray-500" id="signup-card">
            <h1 class="text-3xl my-4">Sign Up Below</h1>
            <form class="mx-auto w-3/5" @submit.prevent="signup()" id="signup-form">
                <h4 style="text-align: left;" class="text-lg">Email: </h4>
                <input type="email" id="email-input" class="form-input mt-1 block mx-auto p-3 bg-gray-100 rounded-lg w-full" placeholder="Email" required>
                <br>

                <h4 style="text-align: left;" class="text-lg">Username: </h4>
                <input type="text" id="username-input" class="form-input mt-1 block mx-auto p-3 bg-gray-100 rounded-lg w-full" placeholder="Username" required>
                <br>

                <h4 style="text-align: left;" class="text-lg">Password: </h4>
                <input id="pass-input" class="form-input mt-1 block mx-auto p-3 bg-gray-100 rounded-lg w-full" type="password" placeholder="Password" required>
                <br>

                <button type="submit" class="bg-blue-500 hover:bg-blue-400 rounded-lg text-white text-lg px-6 py-1 mb-5">Sign Up</button>
            </form>
        </div>
    </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'
import axios from 'axios'
import $ from 'jquery'

export default {
    name: "Signup",
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
        async signup() {
            let username = $("#username-input").val()
            let email = $("#email-input").val()
            let password = $("#pass-input").val()
            if (email === null || password === null || username === null) {
                return
            }

            let resp = await axios.post(this.apiUrl +'/signup', {
                Username: username,
                Email: email,
                Password: password
            })
            if (resp.data === 'error') {
                alert("Internal Server Error")
                return
            }
            sessionStorage.setItem("userToken", resp.data.id)
            sessionStorage.setItem("username", resp.data.username)
            alert("Successfully Signed Up 😇")
            this.$router.push('/')
        }
    }
}
</script>
<style scoped>
#signup-card {
    position: relative;
    top: 5rem;
    width: 35vw;
    --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
}
@media screen and (max-width: 1200px) {
    #signup-card {
        position: relative;
        top: 5rem;
        width: 59vw;
        --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
    }
}
@media screen and (max-width: 550px) {
    #signup-card {
        position: relative;
        top: 5rem;
        width: 75vw;
        --tw-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow);
    }
}
</style>