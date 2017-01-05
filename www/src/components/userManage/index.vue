<template>
	<div class="userManage container">
	  <div class="fee-set pool-panel">
	  	<h1>矿池费率设置</h1>
	  	<ul>
	  	  <template v-if="!edit">
		  		<li>普通费率：<span v-if="!edit">{{cacheFee}} %</span></li>
	  		  <li>VIP费率：<span v-if="!edit">{{cacheVipFee}} %</span></li>
	  	    <a class="edit" v-if="!edit" @click="editClick">编辑</a>
	  	  </template>
	  	  <template v-else>
		  		<li>普通费率：<input type="text" name="defaultfee" v-model="cacheFee" /> %</li>
	  		  <li>VIP费率：<input type="text" name="defaultfee" v-model="cacheVipFee" /> %</li>
			  	<a class="cansel" @click="canselClick" >取消</a>
			  	<a class="save" @click="saveClick" >保存</a>
		  	  <p class="unsave-warning">请确认保存，否则不会保存</p>
	  	  </template>
	  	</ul>
	  </div>
	 
	  <div class="query-set pool-panel">
	  	<h1>查询&amp;设置</h1>
	  	<form class="query-form" name="queryUser" autocomplete="off" @submit.prevent="queryUser">
	  		<input type="text" name="queryPara" placeholder="用户名/钱包地址/电话号码/邮箱" v-model="queryPara" />
  			<a class="history" :class="{'open':openHistory}" @click="open" @mouseover="willopen" @mouseout="willclose" ><span class="cart"></span>
	  			<ul>
	  				<li v-for="(item,index) in queryHistory" @click="selectItem(index)">{{item}}
	  				  <a class="delete" @click.stop="deleteItem(index)">delete</a>
	  				</li>
	  			</ul>
  			</a>
	  		<input type="submit" value="查询" />
	  	</form>

	  	<div class="query-result">
	  	  <template v-for="(user,$index) in users">
	  		  <user-panel  :user="user" :index="$index" />
	  	  </template>
	  	</div>
	  </div>

	</div>
</template>

<script type="text/javascript">
import config from '../../../config'
import userPanel from './userPanel.vue'
	export default{
		components:{
			userPanel: userPanel
		},
		data(){return{
			//view
			edit: false,
			cacheFee: 0,
			cacheVipFee: 0,
			openHistory: false,
			closeTimeout: 0,


			//localStorage
			queryHistory:[],

			//api
			queryPara: '',
			
			//data
			Fee: 1.5,
			VipFee: 0.5,
			users:[
			  {
					name: 'world',
					tel: '1244235',
					email: '123141421',
					wallet:[
					{
						address: '1231453466236435345',
						fee: 1.2,
						lastBanefit:12452354,
						totalBanefit: 12414125,
						stats: '离线',
						offLineTime: '213124324215',
					},
					{
						address: 'dsfgelrjtwre;lk23',
						fee: 1.2,
						lastBanefit:12452354,
						totalBanefit: 12414125,
						stats: '在线',
						offLineTime: '213124324215',
					}
					]
				},
				{
					name: 'hello',
					tel: '1244235',
					email: '123141421',
					wallet:[
						{
							address: '1231453466236435345',
							fee: 1.2,
							lastBanefit:12452354,
							totalBanefit: 12414125,
							stats: '离线',
							offLineTime: '213124324215',
						},
						{
							address: 'dsfgelrjtwre;lk23',
							fee: 1.2,
							lastBanefit:12452354,
							totalBanefit: 12414125,
							stats: '在线',
							offLineTime: '213124324215',
						}
					]
				},
				{
					name: 'hello',
					tel: '1244235',
					email: '123141421',
					wallet:[
						{
							address: '1231453466236435345',
							fee: 1.2,
							lastBanefit:12452354,
							totalBanefit: 12414125,
							stats: '离线',
							offLineTime: '213124324215',
						},
						{
							address: 'dsfgelrjtwre;lk23',
							fee: 1.2,
							lastBanefit:12452354,
							totalBanefit: 12414125,
							stats: '在线',
							offLineTime: '213124324215',
						}
					]
				}
			]
		}},
		methods:{
			queryUser(){
		    var header = new Headers({ 'Json-Web-Token' : localStorage.getItem( config.BTCC.PM_JWT ) })
		    fetch(config.BTCC.PM_APIHOST + 'user/query/' + this.queryPara,{ headers : header })
		    .then(resp => {
		      if(resp.status === 403) this.$router.replace('/')
		      if(resp.ok){
		      	if(!this.queryHistory.some(el => el === this.queryPara)){
				      this.queryHistory.push(this.queryPara)
				      localStorage.setItem(config.BTCC.PM_QUERY_HISTORY,this.queryHistory)
		      	}
		        return resp.json()
		      }
		    })
		    .then(json => {
		    	console.log(json)
		    	this.userFormat(json)
		    })
			},
			userFormat(data){
				var users = []
				data.forEach(el => {
					var user = {
						name: el.BasicInfo.account,
						tel: el.BasicInfo.phone,
						email: el.BasicInfo.email,
						wallet: []
					}
					var wallet = {
						address: el.BasicInfo.walletAddress,
						fee: el.BasicInfo.fee,
						lastBanefit: el.LastRevenue,
						totalBanefit: el.AllRevenue,
						stats: !!el.OfflineTime ? '离线' : '在线',
						offLineTime: el.offLineTime
					}

					var index = 0

					if(!users.some( (el,i) => {
						if(el.name === user.name){
							users[i].wallet.push(wallet)
							return true
						}
						index=i
						return false
					})){
					  users.push(user)
					  users[index].wallet.push(wallet)
					}
				})

			  this.users = users
			},
			editClick(){
				this.edit = true
			},
			open(){
				this.openHistory = !this.openHistory
			},
			willopen(){
				clearTimeout(this.closeTimeout)
			},
			willclose(){
				clearTimeout(this.closeTimeout)
				this.closeTimeout = setTimeout(() => this.openHistory = false, 1000)
			},
			selectItem(index){
				this.queryPara = this.queryHistory[index]
			},
			deleteItem(index){
				this.queryHistory.splice(index,1)
				localStorage.setItem(config.BTCC.PM_QUERY_HISTORY,this.queryHistory)
			},
			canselClick(){
				this.edit = false
				this.cacheFee = this.Fee
				this.cacheVipFee = this.VipFee
			},
			saveClick(){
				this.edit = false
				this.Fee = this.cacheFee
				this.VipFee = this.cacheVipFee
			},
			init(){
				this.cacheFee = this.Fee
				this.cacheVipFee = this.VipFee
				this.edit = false
				this.queryHistory = localStorage.getItem(config.BTCC.PM_QUERY_HISTORY).split(',')
				this.queryPara = this.queryHistory[this.queryHistory.length-1]
			},
		},
		created(){
			this.init()
		}
	}
</script>

<style type="text/css">
.user-panel,.user-wallet,.fee-set {
    position: relative;
}
a.edit, a.cansel, a.save, p.unsave-warning{
	  display: none;
    position: absolute;
    top: 10px;
    right: 40px;
    line-height: 24px;
    font-size: 14px;
    color: #19a8f7;
    cursor: pointer;
}
a.cansel{
	right: 80px;
}
.fee-set a, p.unsave-warning{
	display: block;
}
p.unsave-warning{
	right: 120px;
	color: #f12345;
}
.userManage input[type="text"]{
    font-size: 16px;
    border: none;
    text-align: center;
    border-bottom: 1px solid #19a8f7;
    outline: none;
    width: 100px;
    background-color: transparent;
}
.userManage input[name="queryPara"]{
  width: 400px;
  font-size: 14px;
  border: 1px solid #ccc;
  float: left;
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px;
}
.userManage input[name="queryPara"]:focus{
	outline: -webkit-focus-ring-color auto 5px;
}
.userManage input[name="tel"],.userManage input[name="email"]{ width: 150px; }

.fee-set ul {
    margin-top: 20px;
}
.fee-set li {
    list-style: none;
    margin-top: 10px;
}
.query-set {
    margin-bottom: 40px;
}
form.query-form {
    margin-top: 20px;
}
.query-form input{
	height: 30px;
}
.query-form .history{
    position: relative;
    width: 30px;
    height: 30px;
    text-align: center;
    display: inline-block;
    line-height: 30px;
    cursor: pointer;
    border: 1px solid #ccc;
    border-left: none;
    float: left;
}
.query-form .cart{
    width: 0;
    height: 0;
    vertical-align: middle;
    border-top: 4px dashed;
    border-right: 4px solid transparent;
    border-left: 4px solid transparent;
    display: inline-block;
}
.history ul{
	display: none;
    position: absolute;
    right: 0;
    background: #fff;
    border: 1px solid #ccc;
    margin-top: 2px;
    border-radius: 5px;
    z-index: 1;
}
.history ul li{
	font-size: 12px;
	min-width: 300px;
	text-align: left;
  padding-left: 10px;
  padding-right: 10px;
}
.history ul li:hover{
	background-color: #f1f2f3;
}
a.delete {
  float: right;
  color: #19a8f7;
  text-decoration: underline;
}
a.delete:hover{
	font-weight: 600;
}
.open ul{
	display: block;
}
.query-form input[type="submit"]{
    background: none;
    font-size: 16px;
    line-height: 30px;
    display: inline-block;
    padding-left: 20px;
    padding-right: 20px;
    border: 1px solid #ccc;
    border-left: none;
    border-top-right-radius: 5px;
    border-bottom-right-radius: 5px;
    box-sizing: content-box;
    background-color: #19a8f7;
    color: #fff;
    cursor: pointer;
}
.query-result {
    margin-top: 10px;
}
.user-panel {
	padding: 10px;
}
div.user-panel:hover{
	background-color: #e4f5fe;
}
/* 悬浮时示edit等标签 */
div.user-panel:hover>ul>a, div.user-panel:hover>div>div>ul>a{
	display: block;
}
.user-panel ul {
    margin-top: 10px;
    height: 30px;
    line-height: 30px;
}
.user-panel li {
    float: left;
    list-style: none;
    margin-right: 50px;
}
.user-panel:nth-child(odd){
	background-color: #f4f4f4;
}
.user-panel:nth-child(even){
	background-color: #fff;
}
.user-wallet {
    display: table;
    width: 100%;
    box-sizing: border-box;
    clear: both;
    margin-top: 5px;
    border: 1px solid #f0f0f0;
    border-radius: 10px;
    background-color: #fff;
    padding: 5px 10px;
}
</style>