function LoginAttempt() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 302) {
            var l = this.responseURL;
            window.location = l;
        }
    };
    xhttp.open("GET", "/login", true);
    xhttp.send();
}