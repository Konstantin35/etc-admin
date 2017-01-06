export default function(num){
	var len = Math.round(num).toString().length;
  var Y = ''
  if(len > 12){
    Y = (num/Math.pow(10,12)).toFixed(2) + ' TH';
  }else if(len >9){
    Y = (num/Math.pow(10,9)).toFixed(2) + ' GH';
  }else if(len >6){
    Y = (num/Math.pow(10,6)).toFixed(2) + ' MH';
  }else if(len > 3){
    Y = (num/Math.pow(10,3)).toFixed(2) + ' KH';
  }else{
    Y = 0 + ' KH';
  }

  return Y
}