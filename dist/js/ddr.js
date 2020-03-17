function updateProfile() {
    let btn = document.getElementById('ddr-player-details-update');
    btn.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'ddr-update-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    btn.parentNode.insertBefore(processing, btn);

    $.ajax({url: "/external/bst_api/ddr_update", dataType: "json", success: function(result, textStatus, xhr){
            if (xhr.status == 200) {
                processing.innerText = "Done!"
            } else {
                processing.innerText = this.response
            }
        }});
}

function refreshProfile() {
    let btn = document.getElementById('ddr-player-details-refresh');
    btn.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'ddr-refresh-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    btn.parentNode.insertBefore(processing, btn);

    $.ajax({url: "/external/bst_api/ddr_refresh", dataType: "json", success: function(result, textStatus, xhr){
        if (xhr.status == 200) {
            processing.innerText = "Done!"
        } else {
            processing.innerText = this.response
        }
        }});
}