import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import("@/views/Home.vue")
  },
  {
    path: '/signup',
    name: "Signup",
    component: () => import("@/views/Signup.vue")
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue")
  },
  {
    path: '/logout',
    name: "Logout",
    component: () => import("@/views/Logout.vue")
  },
  {
    path: "/new",
    name: "New",
    component: () => import("@/views/New.vue")
  },
  {
    path: '/article/:id',
    name: "Article",
    component: () => import("@/views/Article.vue")
  },
  {
    path: '/@:id',
    name: "User",
    component: () => import("@/views/User.vue")
  },
  {
    path: '/account',
    name: "Account",
    component: () => import("@/views/Account.vue")
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
