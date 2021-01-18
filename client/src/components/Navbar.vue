<template>
    <div id="nav" class="bg-gray-800">
        <div id="left-buttons">
            <router-link to="/" class="px-3 py-1 text-white text-lg" style="color: white;">Home</router-link>
            <router-link v-if="loggedIn === false" to="/signup" class="px-3 py-1 text-white text-lg" style="color: white;">Sign Up</router-link>
            <router-link v-if="loggedIn === false" to="/login" class="px-3 py-1 text-white text-lg" style="color: white;">Login</router-link>
            <router-link v-if="loggedIn === true" :to="`/logout?page=` + page" class="px-3 py-1 text-white text-lg" style="color: white;">Logout</router-link>
            <!-- right aligning buttons -->
            <router-link v-if="loggedIn === true" to="/new" class="px-3 py-1 text-white text-lg" style="color: white;">New Blog Post</router-link>
            <router-link v-if="loggedIn === true" to="/account" class="px-3 py-1 text-white text-lg" style="color: white;">Account</router-link>
        </div>
    </div>
</template>
<script>
export default {
    name: "Navbar",
    data() {
        return {
            loggedIn: Boolean,
            page: window.location.pathname
        }
    },
    created() {
        if (sessionStorage.getItem('userToken') !== null) {
            this.loggedIn = true
        } else {
            this.loggedIn = false
        }
    }
}
</script>

<style scoped>
#nav {
    padding: .5rem 1rem .5rem 1rem;
}

#left-buttons {
    width: 100%;
    display: inline-flex;
    flex-direction: row;
    gap: 1.5rem;
    align-items: flex-start;
}

</style>