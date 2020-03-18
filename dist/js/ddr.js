function updateProfile() {
    let btn = document.getElementById('ddr-player-details-update');
    btn.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'ddr-update-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    btn.parentNode.insertBefore(processing, btn);

    let update = $.ajax({url: "/external/bst_api/ddr_update",
        type: "PATCH"});

    update.done(function() {
        processing.innerText = "Done!"
    })
}

function refreshProfile() {
    let btn = document.getElementById('ddr-player-details-refresh');
    btn.style.display = 'none';

    let processing = document.createElement('span');
    processing.id = 'ddr-refresh-processing'
    processing.appendChild(document.createTextNode('processing ...'));
    btn.parentNode.insertBefore(processing, btn);

    let update = $.ajax({url: "/external/bst_api/ddr_refresh",
        type: "PATCH",
        timeout: 60000});

    update.done(function() {
        processing.innerText = "Done!"
    })
}

function loadStatsTable() {
    $("#ddr-stats-loc").html = "<a>Loading table - please wait</a>";
    let table = $.ajax({url: "/external/bst_api/ddr_stats"});

    table.done(function(data) {
        $("#ddr-stats-loc").html = data;

        $('#stats').DataTable();
    })
}