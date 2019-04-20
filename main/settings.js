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

  function upload(){
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
        xmlhttp.open("POST", title, false);
        xmlhttp.setRequestHeader("Password", password)

        var jsonString = document.getElementById(element + "json").value;
       
        xmlhttp.send(jsonString);
    });
  }