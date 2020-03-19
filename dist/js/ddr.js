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
    $("#ddr-stats-loc").html("<a>Loading table - please wait</a>");
    let table = $.ajax({url: "/external/bst_api/ddr_stats"});

    table.done(function(data) {
        $("#ddr-stats-loc").html(data);

        let statsDataTable = $('#stats').DataTable();
        addStatsTableFiltering(statsDataTable)
    })

}

function addStatsTableFiltering(statsDataTable) {

    $("div.level-filter").on("click", function() {
        if($(this).hasClass("enabled")) {
            $(this).addClass("disabled");
            $(this).removeClass("enabled");
        } else {
            $(this).addClass("enabled");
            $(this).removeClass("disabled");
        }
        statsDataTable.draw();
    });

    $("div#level-filter-all-enable").on("click", function() {
        $("div.level-filter").each(function() {
            $(this).addClass("enabled");
            $(this).removeClass("disabled");
        });
        statsDataTable.draw();
    });

    $("div#level-filter-all-disable").on("click", function() {
        $("div.level-filter").each(function() {
            $(this).addClass("disabled");
            $(this).removeClass("enabled");
        });
        statsDataTable.draw();
    });

    $.fn.dataTable.ext.search.push(
        function( settings, data, dataIndex ) {
            let level = parseInt(data[0]);
            return $('#level-filter-' + level).hasClass("enabled");
        }
    )

    $('#single-filter, #double-filter').change( function() {
        statsDataTable.draw();
    } );

    $.fn.dataTable.ext.search.push(
        function( settings, data, dataIndex ) {
            let mode = data[3].toLowerCase();
            return $('#' + mode + '-filter')[0].checked;
        }
    );

    $('#beginner-filter, #basic-filter, #difficult-filter, #expert-filter, #challenge-filter').change( function() {
        statsDataTable.draw();
    } );

    $.fn.dataTable.ext.search.push(
        function( settings, data, dataIndex ) {
            let diff = data[4].toLowerCase();
            return $('#' + diff + '-filter')[0].checked;
        }
    );

    $('#fail-filter, #clear-filter, #good-filter, #great-filter, #perfect-filter, #marvellous-filter, #unplayed-filter').change( function() {
        statsDataTable.draw();
    } );

    let lampMap = new Map();
    lampMap.set("マーベラスフルコンボ", "marvellous");
    lampMap.set("パーフェクトフルコンボ", "perfect");
    lampMap.set("グレートフルコンボ", "great");
    lampMap.set("グッドフルコンボ", "good");

    $.fn.dataTable.ext.search.push(
        function( settings, data, dataIndex ) {
            let lamp = data[5].toLowerCase();
            if(lampMap.has(lamp)) {
                return $('#' + lampMap.get(lamp) + '-filter')[0].checked;
            }
            if(lamp === '---') {
                if(data[6].toLowerCase() === 'e') {
                    return $('#fail-filter')[0].checked;
                } else {
                    return $('#clear-filter')[0].checked;
                }
            }
            return $('#unplayed-filter')[0].checked;
        }
    );
}