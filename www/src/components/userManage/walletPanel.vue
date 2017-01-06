<template>
<div class="user-wallet">
	<h3>钱包地址: {{wallet.address}}</h3>
	<div class="wallet-info" >
		<ul class="wallet stats" >
			<li v-if="!edit">费率: <span>{{cacheFee}}</span> %</li>
			<li v-else>费率: <input type="text" name="walletFee" v-model="cacheFee"> %</li>
			<li>上次收益: <span>{{wallet.lastBanefit}}</span></li>
			<li>总收益: <span>{{wallet.totalBanefit}}</span></li>
			<li>状态: <span>{{wallet.stats}}</span></li>
			<li v-show="wallet.stats === '离线'">离线时间: <span>{{wallet.offLineTime}}</span></li>
			<a class="edit" v-if="!edit" @click="editClick">编辑</a>
			<template v-else>
		  	<a class="cansel" @click="canselClick">取消</a>
		  	<a class="save" @click="saveClick">保存</a>
		  	<p class="unsave-warning" >请确认保存，否则不会保存</p>
	  	</template>
		</ul>
	</div>
</div>
</template>

<script type="text/javascript">
	export default{
		props:{
			wallet: Object,
			uindex: Number,
			windex: Number,
			updateData: Function
		},
		data(){return{
			edit: false,
			cacheFee: this.wallet.fee
		}},
		methods:{
			editClick(){
				this.edit = true
			},
			canselClick(){
				this.cacheFee = this.wallet.fee
				this.edit = false
			},
			saveClick(){
				this.wallet.fee = parseFloat(this.cacheFee)
				console.log(this.uindex,this.windex)
				this.updateData(this.uindex,this.windex)
				this.edit = false
			}
		},
		watch: {
			wallet(newVal){
				this.cacheFee = this.wallet.fee
			}
		}
	}
</script>