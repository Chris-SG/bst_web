function GetApiStatus() {
    $.ajax({url: "/external/bst_api/status", dataType: "json", success: function(result){
        console.log(result);
        $("#footer a")[0].setAttribute("class", "bst-api-status-" + result.api);
        $("#footer a").text("BST API STATUS: " + result.api);
    }});
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