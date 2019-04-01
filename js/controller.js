var main = document.getElementsByTagName("main")[0];
var team;

function setTeamInfo(){
  var teaminfo;
	
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      teaminfo = JSON.parse(this.responseText);
    }
  };
  xmlhttp.open("GET", "/api/info/teams", false);
  xmlhttp.send();
  
  var tmpl = $.templates("#teamTemplate");
  
  for(t in teaminfo) {
    var data = {teamname: teaminfo[t]["name"]}
	var html = tmpl.render(data);
	
	var d1 = document.getElementById('teamSelect');
    d1.insertAdjacentHTML('beforeend', html);
  }
}

function selectTeam(){
  var e = document.getElementById("teamSelect");
  team = e.options[e.selectedIndex].value;
  if (team == "Choose...") {
    team = "Select Team"
  }
  
  document.getElementById("teamSelectButton").innerHTML = team;
}

function setPage(title){
  //Grab the other page and swap it out

  title = "main/" + title;

  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      main.innerHTML = this.responseText;
	  
	  var script = document.createElement('script');
      script.onload = function () {};
      script.src = title + ".js";

      document.head.appendChild(script);
    }
  };
  xmlhttp.open("GET", title + ".html", false);
  xmlhttp.send();
}

(function () {
	
  setPage("dashboard");
  setTeamInfo();

  feather.replace();
}())