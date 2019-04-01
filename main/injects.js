function inject(e){
  //Grab the other page and swap it out

  var id = e.id;
  var answer = document.getElementById(id + '-ans').value

  var title = "/inject/" + id;

  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      if (this.responseText == 'correct') {
        document.getElementById(id).className = "btn btn-success";
        document.getElementById(id).innerHTML = "Correct!";
      } else {
		  document.getElementById(id).innerHTML = this.responseText;
		  document.getElementById(id).className = "btn btn-outline-danger";
	  }
    }
  };
  xmlhttp.open("POST", title, false);
  
  var obj = new Object();
   obj.title = id;
   obj.answer = answer;
   obj.team = team;
   var jsonString= JSON.stringify(obj);
  
  xmlhttp.send(jsonString);
}

function updateQuestions(){
  var questions
  
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
   	if (this.readyState == 4 && this.status == 200) {
      questions = JSON.parse(this.responseText);
    }
  };
  xmlhttp.open("GET", "/api/info/injects", false);
  xmlhttp.send();
  
  var tmpl = $.templates("#injectTemplate");
  
  for(q in questions) {
    var data = {title: questions[q]["title"],question: questions[q]["question"],points: questions[q]["points"]}
	var html = tmpl.render(data);
	
	var d1 = document.getElementById('injects-go-here');
    d1.insertAdjacentHTML('beforeend', html);
  }
}

(function () {
  
  updateQuestions();
  
}())