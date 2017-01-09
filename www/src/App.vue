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
html,body,div,p,ul,li,span,label,input,a,h1,h2,h3,h4,h5{
  margin: 0;
  padding: 0;
}
body{
	background-color: #f3f3f3;
}
h1{ font-size: 24px; }
h2{ font-size: 20px; }
h3{ font-size: 18px; }
h4{ font-size: 16px; }
h5{ font-size: 14px; }
h1,h2,h3,h4,h5{
	line-height: 24px;
}

.container{
	position: relative;
	width: 1080px;
	margin: 0 auto;
}
.pool-panel{
    border: 1px solid #ddd;
    margin-top: 10px;
    background-color: #fff;
    padding: 10px 30px;
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
