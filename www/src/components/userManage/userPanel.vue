<template>
<div class="user-panel">
	<ul>
		<li><h3>用户名:{{user.name}}</h3></li>
		<template v-if="!edit">
			<li>电话: <span>{{tel}}</span></li>
			<li>Email: <span>{{email}}</span></li>
	    <a class="edit" @click="editClick" >编辑</a>
		</template>
		<template v-else>
			<li>电话: <input type="text" name="tel" v-model="tel" /></li>
			<li>Email: <input type="text" name="email" v-model="email" /></li>
			<a class="cansel" @click="canselClick">取消</a>
			<a class="save" @click="saveClick">保存</a>
			<p class="unsave-warning">请确认保存，否则不会保存</p>
		</template>
	</ul>
	<template v-for="(wallet,windex) in user.wallet">
		<wallet-panel :wallet="wallet" :index="windex" />
	</template>
</div>
</template>

<script type="text/javascript">
import walletPanel from './walletPanel.vue'
	export default{
		components:{
			walletPanel: walletPanel
		},
		props:{
			user: Object,
			index: Number,
			updateData: Function
		},
		data(){return {
			tel: this.user.tel,
			email: this.user.email,
			edit: false,
			warning: false
		}},
		methods:{
			editClick(){
				this.edit = true
			},
			canselClick(){
				this.edit = false
			  this.tel = this.user.tel
			  this.email = this.user.email
			},
			saveClick(){
				this.edit = false
				this.user.tel = this.tel
				this.user.eamil = this.eamil
				this.updateData()
			}
		}
	}
</script>