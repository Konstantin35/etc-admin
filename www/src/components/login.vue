<template>
	<div class="login">
	  <div class="header">
	    <div class="container">
		    <img class="logo" src="../assets/logo.png" width="108" height="26" />
		    <label class="welcome">欢迎登录</label>
		  </div>
	  </div>
	  <!-- header over -->
	  <div class="login-bg">
		  <div class="container">
			  <div class="login-form">
			  	<div class="account">
			  		<input type="text" name="account" placeholder="账号" v-model="username" />
			  		<!-- <label for="account" class="account-label">账号：</label> -->
			  	</div>
			  	<div class="password">
			  		<input type="password" name="password" placeholder="密码" v-model="password" />
			  		<!-- <label for="password" class="password-label">密码：</label> -->
			  	</div>
			  	<div class="helper">
			  	  <a class="remember">
			  		<input class="remember_password" type="checkbox" name="remember_password" v-model="isremember" />
			  		<a @click="rememberchange">记住密码</a></a>
			  	</div>
			  	<div class="submit">
			  		<a class="submit-btn" @click="loginon">登录</a>
			  	</div>
			  </div>
			</div>
		</div>
		<!-- login-bg over -->

	</div>
</template>

<script type="text/javascript">
import config from '../../config'
export default{
	data(){
		return {
			//view state
			isremember: false,

			//data state
			username: "",
			password: ""
		}
	},
	created(){
		this.init();
	},
	methods:{
		rememberchange(){
			this.isremember = !this.isremember
		},
		remember(){
			if(this.isremember){
				localStorage.setItem(config.BTCC.PM_USERNAME,this.username)
				localStorage.setItem(config.BTCC.PM_PWD,this.password)
			}
		},
		loginon(){
			// console.log(this.$router.replace({path:'/manage'}))
			this.remember()
			
			var header = new Headers({
				'Content-Type' : 'application/x-www-form-urlencoded'
			})

			var formdata = new FormData()
			formdata.append('username',this.username)
			formdata.append('password',this.password)
			// fetch(config.BTCC.PM_APIHOST+'login',{
			// 	method: 'OPTIONS'
			// })
			// .then(oresp => {
			// 	console.log(oresp)
			// 	if(oresp.ok){
					fetch(config.BTCC.PM_APIHOST+'login',{
						method: 'POST',
						headers: header,
						body: formdata
						// body: "username="+this.username+"&password="+this.password
					})
					.then(resp => {
						if(resp.ok){
							console.log('accssce')
						}
					})
			// 	}else{
			// 		console.log('failed')
			// 	}
			// })

			// var xhr = new XMLHttpRequest();
			// xhr.onreadystatechange = function(){
			// 	// console.log(xhr.readyState,xhr.status)
			// 	console.log(document.cookie)
			// 	if(xhr.readyState === 4 && xhr.status==200){
			// 		console.log(xhr.responseText)
			// 	}
			// }
			// xhr.open('POST',config.BTCC.PM_APIHOST+'login',true)
			// xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
			// xhr.send("username="+this.username+"&password="+this.password)
		},
		init(){
			this.username = localStorage.getItem(config.BTCC.PM_USERNAME)
			this.password = localStorage.getItem(config.BTCC.PM_PWD)
			this.isremember = localStorage.getItem(config.BTCC.PM_ISREMEMBER) === 'true'
		}
	},
	watch:{
		isremember(newVal){
			if(newVal){
			  localStorage.setItem(config.BTCC.PM_ISREMEMBER,newVal)
			}else{
				localStorage.removeItem(config.BTCC.PM_USERNAME)
				localStorage.removeItem(config.BTCC.PM_PWD)
				localStorage.removeItem(config.BTCC.PM_ISREMEMBER)
			}
		}
	}
}
</script>

<style type="text/css" scoped>
.welcome{
	font-size: 18px;
	color: #4c8fdf;
	margin-left: 20px;
}
/*. content
---------------------*/
.login-bg{
	background-color: #171c30;
}
.login-bg .container{
	position: relative;
	height: 500px;
	background: url('../assets/minning.png') 0 0/100% auto no-repeat;
}
.login-form{
	position: absolute;
	right: 0;
	top: 30px;

	background-color: rgba(242, 242, 242, 0.1);
	width: 286px;
	padding: 20px 40px;
}

.account, .password{
  position: relative;
  /* padding-left: 48px; */
  height: 38px;
  line-height: 38px;
}

.password, .helper, .submit{
	margin-top: 20px;
}

/*. account & password
---------------------*/
.account-label, .password-label {
    position: absolute;
    z-index: 3;
    top: 0;
    left: 0;
    text-indent: -99999999px;
    width: 38px;
    height: 38px;
    background-image: url('../assets/pwdicons.png');
    background-repeat: no-repeat;
}

.password-label{
    background-position: -48px 0px;
}

.account input, .password input{
	height: 20px;
	width: 100%;
	border: none;
	outline: none;
	color: #fff;
	font-size: 14px;
	text-align: center;
	border-bottom: 1px solid #fff;
	background-color: transparent;
}

.account input:focus, .password input:focus{
	border-bottom: 1px solid #4c8fdf;
}

.account input:focus+label{
	background-position: 0 -48px;
}
.password input:focus+label{
	background-position: -48px -48px;
}
/*. helper
---------------------*/
.remember_password, .remember-label{
	cursor: pointer;
}
.remember{
	color: #fff;
	cursor: pointer;
	font-size: 12px;
}
.remember:hover{
	color: #4c8fdf;
}

/*. submit
---------------------*/
.submit{
	margin-top: 50px;
}
.submit-btn{
	display: block;
	width: 100%;
	text-align: center;
	height: 38px;
	line-height: 38px;
	color: #fff;
	cursor: pointer;
	text-decoration: none;
	background-color: #4c8fdf;
}
.submit-btn:hover{
	background-color: #4d9af5;
}
.submit-btn:focus{
	background-color: #4c8fdf;
}
</style>