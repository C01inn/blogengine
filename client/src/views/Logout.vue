<template>
    <div id="logout">

    </div>
</template>
<script>
export default {
    name: 'logout',
    created() {
        if (sessionStorage.getItem("userToken") !== null || sessionStorage.getItem("username") !== null) {
            sessionStorage.removeItem("userToken")
            sessionStorage.removeItem("username")
        }
        if (this.$route.query.path) {
            this.$router.push(this.$route.query.path)
        } else {
            this.$router.push('/')
        }

    }
}
</script>