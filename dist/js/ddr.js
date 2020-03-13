function updateProfile() {
    let xhttp = new XMLHttpRequest();
    let details = document.getElementById('ddr-player-details');
    let btn = details.getElementById('ddr-player-details-update');
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (this.status == 200) {
                processing.innerText = "Done!"
            } else if (this.status == 500) {
                processing.innerText = this.response
            }
        }
    };

    btn.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'ddr-update-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    btn.parentNode.insertBefore(processing, btn);

    xhttp.open("PATCH", "/external/bst_api/ddr_update", true);
    xhttp.send();
}