function GetApiStatus() {
    var xhttp = new XMLHttpRequest();
    xhttp.responseType = 'json';
    xhttp.open('GET', '/external/bst_api/status', true);

    xhttp.onload = function() {
        if (this.status == 200) {
            let status = xhttp.response;
            console.log(status)

            $(".footer a")[0].setAttribute("class", "bst-api-status-" + status.api);
            $(".footer a").text("BST API STATUS: " + status.api);
        }
    };
    xhttp.send();

}
GetApiStatus();

$(document).ready(function() {
    $('.btn-logout').click(function(e) {
        Cookies.remove('auth-session');
    });
});

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