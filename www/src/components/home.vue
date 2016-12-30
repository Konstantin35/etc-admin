<template>
	<div class="home container">
	  <chart :type="'line'" :data="data" :options="options"></chart>
	</div>
</template>

<script type="text/javascript">
import Chart from 'vue-bulma-chartjs'

export default {
  components: {
    Chart
  },

  data () {
    return {
      data: {
        datasets: [{
          label: 'etc',
          pointRadius: 1,
          pointHoverRadius: 5,
          pointHitRadius:100,
          data: [{
            x: 0,
            y: 3
          }, {
            x: 1,
            y: 10
          }, {
            x: 2,
            y: 5
          }]
        }]
      },
      options: {
        title:{
          display: true,
          fontSize: 18,
          padding: 40,
          text: '矿池算力波动图'
        },
        legend:{
          position: 'bottom'
        },
        tooltips:{
          // intersect: false
          template: "sd"
        },
        scales: {
          xAxes: [{
            type: 'linear',
            position: 'bottom'
          }]
        }
      }
    } // return over
  },
  // data over

  created(){
    // console.log(this.data)
    // window.setInterval(()=>{
    //   this.data.datasets[0].data.push({x:this.data.datasets[0].data.length,y:Math.random()*20+5})
    //   // console.log(this.data.datasets[0].data)
    // },1000)
  }

}
</script>

