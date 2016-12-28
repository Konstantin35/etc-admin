import Vue from 'vue'
import App from './App'
import VueRouter from 'vue-router'
import config from '../config'
import login from './components/login'
import content from './components/content'
import home from './components/home'

Vue.use(VueRouter)

//路由map
//Define router map
const routes = [
  { path: '/', component: login},
  { 
  	path: '/manage', 
  	component: content,
  	children: [
  	  {path: '/', component: home}
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
  router,
  // beforeCreate(){
  //   fetch(config.BTCC.PM_APIHOST+'login',{ method: 'OPTIONS' })
    // .then(resp => {
      // resp.ok ? console.log('success!') : console.log('failed!')
    //   }
    // )
  	// //检查本地环境，是否已经登录
  	// if(!!localStorage.getItem(config.BTCC.JWT)){
  	// 	router.replace({path:'/manage'})
  	// }
  // }
});
