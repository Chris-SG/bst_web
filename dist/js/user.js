function GetEagateState() {
    $.ajax({url: "/external/ajax/eagate_login_status", success: function(result){
        console.log(result);
        $("#eagate-container").html(result);
    }});
}
GetEagateState();

function eagateLogin() {
    let u = document.getElementById("eagate-username");
    let p = document.getElementById("eagate-password");
    let loginForm = document.getElementById('eagate-login-state');

    loginForm.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'eagate-login-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    loginForm.parentNode.insertBefore(processing, loginForm);

    let loginReq = $.ajax({url: "/external/bst_api/eagate_login",
        type: "POST",
        data: JSON.stringify({username: u.value, password: p.value}),
        contentType: "application/json; charset=utf-8",
        dataType   : "json"
        });

    loginReq.done(function() {
        console.log("reload")
        location.reload()
    })
}

function eagateLogout(user) {
    $.ajax({url: "/external/bst_api/eagate_logout",
        type: "POST",
        data: JSON.stringify({username: user}),
        contentType: "application/json; charset=utf-8",
        dataType   : "json",
        success: function(result){
            location.reload()
        }});
}