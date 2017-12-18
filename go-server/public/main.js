function showBlock(id) {
	$elem = document.getElementById(id)
	if ($elem.style.display === 'none') {
		$elem.style.display = 'block';
	} else {
		$elem.style.display = 'none';
	}
}

function checkOt() {
    var xhttp = new XMLHttpRequest();
    xhttp.open("GET", "check", true);
    xhttp.setRequestHeader("Content-type", "text/html");
    xhttp.send();

}

function checkTech(id) {
	var xhttp = new XMLHttpRequest();
    xhttp.open("GET", "check_tech", false);
    xhttp.setRequestHeader("Content-type", "text/html");
    xhttp.send();
    data = JSON.parse(xhttp.response);
    if (data['Status'] === true) {
    	$elem = document.getElementById(id)
		if ($elem.style.display === 'none') {
			$elem.style.display = 'block';
		} else {
			$elem.style.display = 'none';
		}
    	//var xhttp2 = new XMLHttpRequest();
    	//xhttp2.open("GET", "insertmarks", false);
    	//xhttp2.setRequestHeader("Content-type", "text/html");
    	//xhttp2.send();
    	//document.location.href = "/inmarks"
    } else {
    	alert("У вас нет доступа для выполнения данных функций");
    }
}