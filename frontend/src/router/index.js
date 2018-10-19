import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/pages/Home'
import Staff from '@/pages/Staff'
import Booking from '@/pages/Booking'
import Service from '@/pages/Service'
import Product from '@/pages/Product'
import Customer from '@/pages/Customer'
import Setting from '@/pages/Setting'
import Billing from '@/pages/Billing'
import GiftCard from '@/pages/GiftCard'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/home', name: 'home', component: Home },
    { path: '/booking', name: 'booking', component: Booking },
    { path: '/service', name: 'service', component: Service },
    { path: '/product', name: 'product', component: Product },
    { path: '/customer', name: 'customer', component: Customer },
    { path: '/staff', name: 'staff', component: Staff },
    { path: '/setting', name: 'setting', component: Setting },
    { path: '/billing', name: 'billing', component: Billing },
    { path: '/giftcard', name: 'giftcard', component: GiftCard }
  ]
})
