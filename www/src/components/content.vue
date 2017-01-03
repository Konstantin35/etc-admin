<template>
<div class="content">
  <div class="header">
    <div class="container">
		  <img class="logo" src="../assets/logo.png" width="108" height="26" />
		  <ul class="nav">
		    <li><router-link to="/manage">主页</router-link></li>
		    <li><router-link to="/userManage">用户管理</router-link></li>
		    <li><router-link to="/poolManage">矿池管理</router-link></li>
		    <li><a href="javascript:void(0);">退出</a></li>
		  </ul>
		</div>
	</div>
	<!-- hearder over -->
	<router-view></router-view>
</div>
</template>

<style type="text/css">
/*. header
---------------------*/
.header{
	height: 74px;
	min-width: 1080px;
	background-color: #fff;
}
.logo{
	margin-top: 24px;
}
.nav{
	float: right;
}
.nav li{
	float: left;
	list-style: none;
}
.nav a{
	display: block;
	text-decoration: none;
	line-height: 74px;
	padding-left: 25px;
	padding-right: 25px;
}
.nav a:hover{
	background-color: ;
}

a.router-link-active {
    color: #19a8f7;
}

</style>