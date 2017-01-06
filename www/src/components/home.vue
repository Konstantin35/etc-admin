<template>
	<div class="home container">
	  <chart :type="'line'" :data="seriesData" :options="options"></chart>
    <div class="pool-stats pool-panel">
      <h1>矿池状态概览</h1>
      <ul>
        <li><p>矿池当前算力: <span>{{hashrate}} </span></p></li>
        <li><p>活跃用户数: <span>{{minersTotal}}</span></p></li>
        <li><p>矿池余额（包含未支付）: <span>{{poolbalance}}</span></p></li>
      </ul>
    </div>
	</div>
</template>

<script type="text/javascript">
import Chart from 'vue-bulma-chartjs'
import config from '../../config'

export default {
  components: {
    Chart
  },

  data () { return {
    hashrate: 0,
    minersTotal: 0,
    poolbalance: 0,
    // chart datas
    chartX: [],
    chartData: [],
    options:{
      title:{
        display: true,
        fontSize: 18,
        padding: 40,
        text: '7天算力波动图'
      },
      legend:{
        position: 'bottom'
      },
      animation: false,
      tooltips:{
        callbacks:{
          label(item,data){
            var y = item.yLabel
            var len = Math.round(y).toString().length;
            var Y = ''
            if(len > 12){
              Y = (y/Math.pow(10,12)).toFixed(2) + ' TH/s';
            }else if(len >9){
              Y = (y/Math.pow(10,9)).toFixed(2) + ' GH/s';
            }else if(len >6){
              Y = (y/Math.pow(10,6)).toFixed(2) + ' MH/s';
            }else if(len > 3){
              Y = (y/Math.pow(10,3)).toFixed(2) + ' KH/s';
            }else{
              Y = 0 + ' KH/s';
            }
            return data.datasets[0].label + ' : ' + Y
          }
        }
      },
      scales: {
        xAxes: [{
          type: 'time',
          position: 'bottom',
          ticks: {
            callback (value, index, values){
              if(values && values.length > 0){
                return new Date(values[index]._d).getMonth() + 1 +'-'+ value.substr(4)
              }
              return value;
            }
          }
        }],
        yAxes:[{
          ticks:{
            beginAtZero: true,
            callback(value,index,values){
              var len = Math.round(value).toString().length;
              var Y = ''
              if(len > 12){
                Y = (value/Math.pow(10,12)).toFixed(0) + ' TH';
              }else if(len >9){
                Y = (value/Math.pow(10,9)).toFixed(0) + ' GH';
              }else if(len >6){
                Y = (value/Math.pow(10,6)).toFixed(0) + ' MH';
              }else if(len > 3){
                Y = (value/Math.pow(10,3)).toFixed(0) + ' KH';
              }else{
                Y = 0 + ' KH';
              }
              return Y
            }
          }
        }]
      }
    } // options over
  }},
  // data over
  computed:{
    seriesData(){
      var seriesData = {
        labels: this.chartX,
        datasets: [{
          label: 'etc',
          fill: true,
          backgroundColor: 'rgba(11, 163, 247, 0.2)',
          borderColor: '#19a8f7',
          pointRadius: 0,
          pointHoverRadius: 5,
          pointHitRadius:5,
          data: this.chartData
        }]
      }
      return  seriesData
    }, //seriesData over
  }, //computed over
  beforeCreate(){
    var jwt = localStorage.getItem( config.BTCC.PM_JWT )
    var header = new Headers({ 'Json-Web-Token' : jwt })
    fetch(config.BTCC.PM_APIHOST + 'main/poolchart',{ headers : header })
    .then(resp => {
      if(resp.status === 403) this.$router.replace('/')
      return resp.json()
    })
    .then(json => {
      json.poolhashs.forEach((el)=>{
        this.chartX.push(el.tempstamp)
        this.chartData.push( el.value )
      })
    })
    .catch(err => {})

    fetch(config.BTCC.PM_APIHOST + 'main/statistic',{ headers : header })
    .then(resp => {
      if(resp.status === 403) this.$router.replace('/')
      return resp.json()
    })
    .then(json => {
      //need format
      this.hashrate = json.hashrate
      this.minersTotal = json.minersTotal
      this.poolbalance = json.poolbalance
    })
    .catch(err => {})
  }
}
</script>

<style type="text/css">
canvas.chartjs{
  margin-top: 20px;
  background-color: #fff;
}
.pool-stats {
  margin-bottom: 40px;
}
.pool-stats h1{
  text-align: center;
}
.pool-stats ul{
  margin: 0 auto;
  display: table;
  margin-top: 20px;
}
.pool-stats li{
  list-style: none;
  display: table-cell;
}
.pool-stats p{
  padding-left: 30px;
  padding-right: 30px;
}
</style>
