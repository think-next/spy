import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Announcement from '../views/Announcement.vue'
import Leave from '../views/Leave.vue'
import Contacts from '../views/Contacts.vue'

const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/announcement', name: 'Announcement', component: Announcement },
  { path: '/leave', name: 'Leave', component: Leave },
  { path: '/contacts', name: 'Contacts', component: Contacts },
]

const router = createRouter({
  history: createWebHistory('/oa/'),
  routes,
})

export default router
