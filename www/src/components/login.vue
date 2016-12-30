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
			  <div class="login-form" :class="{'login-err':loginErr}">
			  	<div class="account">
			  		<input type="text" name="name" placeholder="账号" v-model="username" @focus="inputFocused" />
			  	</div>
			  	<div class="password">
			  		<input type="password" name="password" placeholder="密码" v-model="password" @focus="inputFocused" />
			  	</div>
			  	<div class="helper">
			  	  <a class="remember">
			  		<input class="remember_password" type="checkbox" name="remember_password" v-model="isremember" />
			  		<a @click="rememberchange">记住我</a></a>
			  	</div>
			  	<div class="submit">
			  	  <div class="err-panel" > 账号或密码错误 </div>
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
			loginErr: false,

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
		inputFocused(){
			this.loginErr = false;
		},
		loginon(){
			var header = new Headers({
				'Content-Type' : 'application/x-www-form-urlencoded'
			})

			// var formdata = new FormData()
			// formdata.append('username',this.username)
			// formdata.append('password',this.password)
			fetch(config.BTCC.PM_APIHOST+'login',{
				method: 'POST',
				headers: header,
				// body: formdata
				body: "username="+this.username+"&password="+this.password
			})
			.then(resp => {
				if(resp.ok){
					if(this.isremember){
			      localStorage.setItem(config.BTCC.PM_USERNAME,this.username)
				    localStorage.setItem(config.BTCC.PM_JWT,resp.headers.get('Json-Web-Token'))
					}
					this.$router.replace('/manage')
				}else{
					this.loginErr = true
				}
			})
		},
		init(){
			this.isremember = localStorage.getItem(config.BTCC.PM_ISREMEMBER) === 'true'
			if(this.isremember){
				this.username = localStorage.getItem(config.BTCC.PM_USERNAME)
			}
		}
	},
	watch:{
		isremember(newVal){
			if(newVal){
			  localStorage.setItem(config.BTCC.PM_ISREMEMBER,newVal)
			}else{
				localStorage.removeItem(config.BTCC.PM_USERNAME)
				localStorage.removeItem(config.BTCC.PM_JWT)
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
}

.password, .helper, .submit{
	margin-top: 20px;
}

/*. account & password
---------------------*/
/* .account-label, .password-label {
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
 */
.account input, .password input{
	height: 38px;
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
	position: relative;
	margin-top: 50px;
}
.err-panel{
	display: none;
  position: absolute;
  width: 286px;
  line-height: 30px;
  top: -40px;
  color: #f00;
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

/*. VUe Bind Classes
---------------------*/
.login-err input{
	border-bottom-color: #f00;
}
.login-err .err-panel{
	display: block;
}
</style>