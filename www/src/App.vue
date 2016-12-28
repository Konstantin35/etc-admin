<template>
  <div id="app">
    <router-view></router-view>
    <div class="footer">
	    <div class="container">
  	    <p class="text">&copy; BTCC Pool Manage System</p>
  	  </div>
	  </div>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
html,body,div,p,ul,li,span,label,input,a{
  margin: 0;
  padding: 0;
}
body{
	background-color: #f3f3f3;
}

.container{
	position: relative;
	width: 1080px;
	margin: 0 auto;
}
/*. footer
---------------------*/
.footer{
	background-color: #fff;
	height: 100px;
	min-width: 1080px;
}

.footer .text{
	text-align: center;
	font-size: 18px;

	padding-top: 39px;
}
</style>
