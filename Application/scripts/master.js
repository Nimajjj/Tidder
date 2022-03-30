const popup = document.getElementById("create_sub_popup");

function hide_popup() {
  popup.style.display = "none"
}

function show_popup() {
  popup.style.display = "block"
}


if (document.getElementById("create_subtidder_error").innerHTML != "") {
  show_popup()
} else {
  hide_popup()
}