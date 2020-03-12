function eagateLogin() {
    let u = document.getElementById("eagate-username");
    let p = document.getElementById("eagate-password");
    let loginForm = document.getElementById('eagate-login-state');

    let xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            loginForm.parentNode.textContent = '';
            location.reload()
        } else if (this.readyState == 4) {
            let processing = document.getElementById('eagate-login-processing');
            processing.parentNode.removeChild(processing);

            u.value = '';
            p.value = '';

            loginForm.style.display = 'initial';
        }
    };

    loginForm.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'eagate-login-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    loginForm.parentNode.insertBefore(processing, loginForm);

    xhttp.open("POST", "/external/bst_api/eagate_login", true);
    xhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    let reqBody = JSON.stringify(JSON.parse(`{ "username": "${u.value}", "password": "${p.value}" }`));
    xhttp.send(reqBody);
}