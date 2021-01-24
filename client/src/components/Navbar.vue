<template>
    <div id="nav" class="bg-gray-800">
        <div id="left-buttons">
            <router-link to="/" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">Home</router-link>
            <router-link v-if="loggedIn === false" to="/signup" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">Sign Up</router-link>
            <router-link v-if="loggedIn === false" to="/login" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">Login</router-link>
            <router-link v-if="loggedIn === true" :to="`/logout?page=` + page" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">Logout</router-link>
            <!-- right aligning buttons -->
            <router-link v-if="loggedIn === true" to="/new" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">New Blog Post</router-link>
            <router-link v-if="loggedIn === true" to="/account" class="sm:px-3 py-1 text-white sm:text-lg nav-item" style="color: white;">Account</router-link>
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

@media screen and (max-width: 639px) {
    .nav-item {
        font-size: 1rem;
        line-height: 1.2rem;
    }
}

@media screen and (max-width: 529px) {
    .nav-item {
        font-size: .85rem;
        line-height: 1.3rem;
        padding: 0;
        margin: -.1rem;
    }
}

</style>