/* globals Chart:false, feather:false */
var tnames = [];
var bpoints = [];
var rpoints = [];
var teams;

(function () {

  feather.replace()

  // Graphs
  var ctx = document.getElementById('myChart')
  

  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      teams = JSON.parse(this.responseText);
    }
  };
  xmlhttp.open("GET", "/api/json_demo.txt", false);
  xmlhttp.send();

  for (t in teams) {
    tnames.push(teams[t]["name"]);
	bpoints.push(teams[t]["bpoints"]);
	rpoints.push(teams[t]["rpoints"]);
  }
  
  
  // eslint-disable-next-line no-unused-vars
  var myChart = new Chart(ctx, {
    type: 'bar',
    data: {
      labels: tnames,
      datasets: [{
		label: 'Blue Team Points',
        data: bpoints,
		backgroundColor: window.chartColors.blue,
		borderColor: window.chartColors.blue,
		borderWidth: 1
      },
	  {
		label: 'Red Team Points',
        data: rpoints,
		backgroundColor: window.chartColors.red,
		borderColor: window.chartColors.red,
		borderWidth: 1
	  }]
    },
    options: {
      responsive: true,
      legend: {
        position: 'top',
      }
    }
  })
}())
