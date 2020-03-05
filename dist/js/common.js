function GetApiStatus() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var status = this.responseText;

            $(".footer a")[0].setAttribute("id", "bst-api-status-" + status);
            $(".footer a").text("BST API STATUS: " + status);
        }
    };
    xhttp.open("GET", "/ajax/apistatus", true);
    xhttp.send();
}
GetApiStatus();

$(document).ready(function() {
    $('.btn-logout').click(function(e) {
        Cookies.remove('auth-session');
    });
});