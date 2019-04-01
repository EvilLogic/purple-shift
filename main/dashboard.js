/* globals Chart:false, feather:false */

function timer(){
    updateGraph()
}

setInterval(function(){
    timer()}, 30000)

function updateGraph(){
  var tnames = [];
  var bpoints = [];
  var rpoints = [];
  var teams;
	
  // Graphs
  var ctx = document.getElementById('myChart');
  

  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      teams = JSON.parse(this.responseText);
    }
  };
  xmlhttp.open("GET", "/api/info/teams", false);
  xmlhttp.send();

  for (t in teams) {
    tnames.push(teams[t]["name"]);
	bpoints.push(teams[t]["iscore"]);
	rpoints.push(teams[t]["rpoints"]);
  }
  
  
  // eslint-disable-next-line no-unused-vars
  var myChart = new Chart(ctx, {
    type: 'bar',
    data: {
      labels: tnames,
      datasets: [{
		label: 'Inject Points',
        data: bpoints,
		backgroundColor: ["rgb(236, 204, 104)", "rgb(255, 165, 2)", "rgb(123, 237, 159)", "rgb(46, 213, 115)", "rgb(255, 127, 80)", "rgb(255, 99, 72)", "rgb(112, 161, 255)", "rgb(30, 144, 255)"],
		borderColor: window.chartColors.green,
		borderWidth: 1
      },
	  //{
		//label: 'Red Team Points',
        //data: rpoints,
		//backgroundColor: window.chartColors.red,
		//borderColor: window.chartColors.red,
		//borderWidth: 1
	  //}
	  ]
    },
    options: {
      responsive: true,
      legend: {
        position: 'top',
      },
	  scales: {
            yAxes: [{
                ticks: {
                    beginAtZero:true
                }
            }]
        }
    }
  })
  
  var standings = document.getElementById("standings");
  //var head = standings.thead
  
  var table = document.createElement("table");

  var col = [];
  for (t in teams) {
    for (key in teams[t]) {
      if (col.indexOf(key) === -1) {
        col.push(key);
      }
    }
  }

    // ADD JSON DATA TO THE TABLE AS ROWS.
  for (t in teams) {

    var tr = table.insertRow(-1);

    for (j in col) {
      var tabCell = tr.insertCell(-1);
        tabCell.innerHTML = teams[t][col[j]];
    }
  }
  
  standings.tBodies[0].innerHTML = table.tBodies[0].innerHTML;
}

(function () {
  
  updateGraph();

  
}())

