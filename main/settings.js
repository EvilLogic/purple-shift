function download(){
    //Get those dank secrets

    pages = [
        "teams",
        "injects",
        "scans"
    ]
  
    var password = document.getElementById("password").value
  
    pages.forEach(element => {
        var title = "/settings/" + element;
  
        var xmlhttp = new XMLHttpRequest();
        xmlhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                document.getElementById(element + "json").innerHTML = this.responseText;
            }
        };
        xmlhttp.open("GET", title, false);
        xmlhttp.setRequestHeader("Password", password)
        xmlhttp.send();
    });
  }

  function upload(e){
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