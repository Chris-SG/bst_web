function addLine(text) {
    document.writeln(text);
}

function addToA() {
    $("a#abc").append("<p>abcd</p>")
}

function clearVals() {
    console.log($("a#abc").find("p").length)
    $("a#abc").find("p").remove()
}