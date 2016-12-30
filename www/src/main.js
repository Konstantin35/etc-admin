import Vue from 'vue'
import App from './App'
import VueRouter from 'vue-router'
import config from '../config'
import login from './components/login'
import content from './components/content'
import home from './components/home'
import userManage from './components/userManage'
import poolManage from './components/poolManage'

Vue.use(VueRouter)

//路由map
//Define router map
const routes = [
  { path: '/', component: login},
  { 
  	path: '/manage', 
  	component: content,
  	children: [
      {path: '/', component: home},
      {path: '/userManage', component: userManage},
      {path: '/poolManage', component: poolManage},
  	]
  },
]

//使用路由map生成路由对象
//Create the router instance with router map
const router = new VueRouter({
  routes
})

//将路由关系关联到app上
//Inject the router with the router
new Vue({
  el: '#app',
  template: '<App/>',
  components: { App },
  router
});
