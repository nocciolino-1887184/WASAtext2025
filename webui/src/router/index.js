// in qusto file definisco i percorsi e gestisco gli accessi
import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import ConversationsView from '../views/ConversationView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/',
		  redirect: '/session' },
		{path: '/dashboard', component: HomeView, meta: { requiresAuth: true }},
		{path: '/session', component: LoginView},	// qui non è necessaria autorizzazione
		{path: '/users/:id', component: ProfileView, meta: { requiresAuth: true }},
		{path: '/users/:id/conversations', component: ConversationsView, meta: { requiresAuth: true }},
	]
})
	router.beforeEach((to, from, next) => {
	// se salvo questa variabile posso monitorare il login
	const isLoggedIn = sessionStorage.getItem('loggedIn') === 'true'
  
	// se non ho l'autorizzazione mi fa rifare il login
	if (to.meta.requiresAuth && !isLoggedIn) {
	  next('/session')
	} else {
	  next()
	}
  })
export default router
