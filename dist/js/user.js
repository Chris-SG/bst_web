function eagateLogin() {
    let u = document.getElementById("eagate_username");
    let p = document.getElementById("eagate_password");
    let loginForm = document.getElementById('login_form');

    let xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            loginForm.parentNode.textContent = '';
        } else if (this.readyState == 4) {
            let processing = document.getElementById('eagate_login_processing');
            processing.parentNode.removeChild(processing);

            u.value = '';
            p.value = '';

            loginForm.style.display = 'initial';
        }
    };

    loginForm.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'eagate_login_processing'
    processing.appendChild(document.createTextNode('processing ...'));
    loginForm.parentNode.insertBefore(processing, loginForm);

    xhttp.open("POST", "/api_bridge/eagate_login", true);
    xmlhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    JSON.stringify(JSON.parse(`{ "username": "${u.value}", "password": "${p.value} }`))
    xhttp.send();
}