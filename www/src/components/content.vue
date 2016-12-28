<template>
<div class="content">
  <div class="header">
    <div class="container">
		  <img class="logo" src="../assets/logo.png" width="108" height="26" />
		  <ul class="nav">
		    <li><a href="">主页</a></li>
		    <li><a href="">用户管理</a></li>
		    <li><a href="">矿池管理</a></li>
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

</style>