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
	  	<form class="query-form">
	  		<input type="text" name="queryPara" placeholder="用户名/钱包地址/电话号码/邮箱" v-model="queryPara" />
  			<a class="history" :class="{'open':openHistory}" @click="open"><span class="cart"></span>
	  			<ul>
	  				<li v-for="item in queryHistory">{{item}}<a class="delete"></a></li>
	  			</ul>
  			</a>
	  		<input type="submit" value="queryUser" method="get" :action="config.BTCC.PM_APIHOST + '/'" />
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
			edit: false,
			openHistory: false,
			cacheFee: 0,
			cacheVipFee: 0,

			//localStorage
			queryHistory:[],
			
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
			editClick(){
				this.edit = true
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
			},
			open(){
				this.openHistory = !this.openHistory
			}
		},
		created(){
			this.init()
			setInterval(()=>{console.log(this.users[0].tel,this.users[0].wallet[0].fee)},1000)
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
	width: 300px;
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
    padding-left: 10px;
    padding-right: 10px;
}
.history ul li{
	font-size: 12px;
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